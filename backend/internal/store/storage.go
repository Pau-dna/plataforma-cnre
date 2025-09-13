package store

import (
	"github.com/imlargo/go-api-template/internal/repositories"
)

type Store struct {
	Files              repositories.FileRepository
	PushSubscriptions  repositories.PushNotificationSubscriptionRepository
	Notifications      repositories.NotificationRepository
	Users              repositories.UserRepository
	Answers            repositories.AnswerRepository
	EvaluationAttempts repositories.EvaluationAttemptRepository
	Contents           repositories.ContentRepository
	Courses            repositories.CourseRepository
	Enrollments        repositories.EnrollmentRepository
	Evaluations        repositories.EvaluationRepository
	Modules            repositories.ModuleRepository
	UserProgresss      repositories.UserProgressRepository
	Questions          repositories.QuestionRepository
}

func NewStorage(container *repositories.Repository) *Store {
	return &Store{
		Files:              repositories.NewFileRepository(container),
		Notifications:      repositories.NewNotificationRepository(container),
		PushSubscriptions:  repositories.NewPushSubscriptionRepository(container),
		Users:              repositories.NewUserRepository(container),
		Answers:            repositories.NewAnswerRepository(container),
		EvaluationAttempts: repositories.NewEvaluationAttemptRepository(container),
		Contents:           repositories.NewContentRepository(container),
		Courses:            repositories.NewCourseRepository(container),
		Enrollments:        repositories.NewEnrollmentRepository(container),
		Evaluations:        repositories.NewEvaluationRepository(container),
		Modules:            repositories.NewModuleRepository(container),
		UserProgresss:      repositories.NewUserProgressRepository(container),
		Questions:          repositories.NewQuestionRepository(container),
	}
}
