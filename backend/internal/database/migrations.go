package postgres

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {

	err := db.AutoMigrate(
		&models.User{},
		&models.Notification{},
		&models.PushNotificationSubscription{},
		&models.File{},
		&models.Answer{},
		&models.EvaluationAttempt{},
		&models.Content{},
		&models.Course{},
		&models.Enrollment{},
		&models.Evaluation{},
		&models.Module{},
		&models.UserProgress{},
		&models.Question{},
	)

	return err
}
