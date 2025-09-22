package services

import (
	"errors"

	"github.com/imlargo/go-api-template/internal/dto"
	"github.com/imlargo/go-api-template/internal/models"
)

type UserService interface {
	CreateUser(user *dto.RegisterUser) (*models.User, error)
	DeleteUser(userID uint) error
	UpdateUser(userID uint, data *models.User) (*models.User, error)
	GetUserByID(userID uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type userService struct {
	*Service
}

func NewUserService(service *Service) UserService {
	return &userService{
		Service: service,
	}
}

func (s *userService) CreateUser(data *dto.RegisterUser) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (s *userService) DeleteUser(userID uint) error {
	return nil
}

func (s *userService) UpdateUser(userID uint, data *models.User) (*models.User, error) {
	return nil, nil
}

func (s *userService) GetUserByID(userID uint) (*models.User, error) {
	return nil, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return nil, nil
}
