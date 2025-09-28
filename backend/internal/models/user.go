package models

import (
	"time"

	"github.com/imlargo/go-api-template/internal/enums"
)

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

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
