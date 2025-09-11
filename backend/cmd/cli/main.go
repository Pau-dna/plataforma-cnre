package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	modelName := "Answer"
	templateContent := renderTemplate(modelName)
	if err := createFile(modelName, templateContent); err != nil {
		log.Fatalln(err.Error())
	}
}

func createFile(model string, content string) error {

	folder := "internal/repositories/"

	tmpFile, err := os.Create(fmt.Sprintf("%s/%s.go", folder, strings.ToLower(model)))
	if err != nil {
		return err
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		return err
	}

	tmpFile.Close()
	return nil
}

func renderTemplate(modelName string) string {

	template := `
package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm/clause"
)

type CourseRepository interface {
	Get(id uint) (*models.Course, error)
	Create(course *models.Course) error
	Update(course *models.Course) error
	Delete(id uint) error
}

type courseRepository struct {
	*Repository
}

func NewCourseRepository(r *Repository) CourseRepository {
	return &courseRepository{
		Repository: r,
	}
}

func (r *courseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) Get(id uint) (*models.Course, error) {
	var course models.Course
	if err := r.db.First(&course, id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) Update(course *models.Course) error {
	return r.db.Model(course).Clauses(clause.Returning{}).Updates(course).Error
}

func (r *courseRepository) Patch(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.Course{}).Where("id = ?", id).Updates(data).Error
}

func (r *courseRepository) Delete(id uint) error {
	var course models.Course
	course.ID = id
	return r.db.Delete(&course).Error
}

func (r *courseRepository) GetAll() ([]*models.Course, error) {
	var courses []*models.Course
	if err := r.db.Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}
`

	template = strings.ReplaceAll(template, "Course", modelName)
	template = strings.ReplaceAll(template, "course", strings.ToLower(modelName))

	return template

}
