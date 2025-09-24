package models

import "github.com/imlargo/go-api-template/internal/enums"

type User struct {
	BaseModel

	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Fullname  string         `json:"fullname" gorm:"not null"`
	AvatarUrl string         `json:"avatar_url" gorm:"not null"`
	Role      enums.UserRole `json:"role" gorm:"not null;default:'student'"`

	// Relaciones
	Enrollments []*Enrollment `json:"enrollments" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "users"
}

func (User) ValidateUserCreation() error {
	return nil
}

func (User) ValidatePassword() error {
	return nil
}
