package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/enums"
	"github.com/imlargo/go-api-template/internal/models"
	"github.com/imlargo/go-api-template/pkg/jwt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(email, password string) (*dto.UserAuthResponse, error)
	Register(user *dto.RegisterUser) (*dto.UserAuthResponse, error)
	Logout(userID uint) error
	RefreshToken(userID uint, refreshToken string) (*dto.AuthTokens, error)
	GetUser(userID uint) (*models.User, error)
	GoogleLogin(code string) (*dto.UserAuthResponse, error)
}

type authService struct {
	*Service
	userService       UserService
	jwtAuthenticator  *jwt.JWT
	googleOauthConfig *oauth2.Config
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func NewAuthService(service *Service, userService UserService, jwtAuthenticator *jwt.JWT, googleOauthConfig *oauth2.Config) AuthService {
	return &authService{
		service,
		userService,
		jwtAuthenticator,
		googleOauthConfig,
	}
}

func (s *authService) Login(email, password string) (*dto.UserAuthResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *authService) Register(user *dto.RegisterUser) (*dto.UserAuthResponse, error) {

	createdUser, err := s.userService.CreateUser(user)
	if err != nil {
		return nil, err
	}

	accessExpiration := time.Now().Add(s.config.Auth.TokenExpiration)
	refreshExpiration := time.Now().Add(s.config.Auth.RefreshExpiration)
	accessToken, err := s.jwtAuthenticator.GenerateToken(createdUser.ID, accessExpiration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtAuthenticator.GenerateToken(createdUser.ID, refreshExpiration)
	if err != nil {
		return nil, err
	}

	authResponse := &dto.UserAuthResponse{
		User: *createdUser,
		Tokens: dto.AuthTokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    refreshExpiration.Unix(),
		},
	}

	return authResponse, nil
}

func (s *authService) Logout(userID uint) error {
	return nil
}

func (s *authService) RefreshToken(userID uint, refreshToken string) (*dto.AuthTokens, error) {
	return nil, nil
}

func (s *authService) GetUser(userID uint) (*models.User, error) {
	if userID == 0 {
		return nil, errors.New("user ID cannot be zero")
	}

	user, err := s.store.Users.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *authService) GoogleLogin(code string) (*dto.UserAuthResponse, error) {

	token, err := s.exchange(code)
	if err != nil {
		return nil, err
	}

	googleUser, err := s.getUserInfo(token)
	if err != nil {
		return nil, err
	}

	// Check if existing already exists
	existing, err := s.store.Users.GetByEmail(googleUser.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user := existing
	if existing == nil {
		// Register new user
		newUser := &models.User{
			Email:     googleUser.Email,
			Fullname:  googleUser.Name,
			AvatarURL: googleUser.Picture,
			Role:      enums.UserRoleStudent,
		}

		if err := s.store.Users.Create(newUser); err != nil {
			return nil, err
		}

		user = newUser
	}

	accessExpiration := time.Now().Add(s.config.Auth.TokenExpiration)
	refreshExpiration := time.Now().Add(s.config.Auth.RefreshExpiration)
	accessToken, err := s.jwtAuthenticator.GenerateToken(user.ID, accessExpiration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtAuthenticator.GenerateToken(user.ID, refreshExpiration)
	if err != nil {
		return nil, err
	}

	authResponse := &dto.UserAuthResponse{
		User: *user,
		Tokens: dto.AuthTokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    refreshExpiration.Unix(),
		},
	}

	return authResponse, nil
}

func (s *authService) exchange(code string) (*oauth2.Token, error) {

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{})
	if s.googleOauthConfig.RedirectURL == "" {
		return nil, fmt.Errorf("redirect URL is empty")
	}
	token, err := s.googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("invalid code: %v", err)
	}

	// Verify that the token is valid
	if !token.Valid() {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func (s *authService) getUserInfo(token *oauth2.Token) (*GoogleUserInfo, error) {

	// Create an HTTP client with the access token
	client := s.googleOauthConfig.Client(context.Background(), token)

	// Retrieve user information using the token
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("error retrieving user information: %v", err)
	}

	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error in Google API response: %s", resp.Status)
	}

	// Decode the response into a generic map
	var googleUser GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("error decoding user information: %v", err)
	}

	return &googleUser, nil
}
