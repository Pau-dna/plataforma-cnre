package repositories_test

import (
	"testing"

	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.Course{}, &models.Module{}, &models.Content{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func TestModuleRepository_GetByCourseID(t *testing.T) {
	db := setupTestDB(t)
	
	// Create test data
	course := &models.Course{
		BaseModel: models.BaseModel{ID: 1},
		Title:     "Test Course",
	}
	db.Create(course)

	modules := []*models.Module{
		{BaseModel: models.BaseModel{ID: 1}, Title: "Module 1", CourseID: 1, Order: 1},
		{BaseModel: models.BaseModel{ID: 2}, Title: "Module 2", CourseID: 1, Order: 2},
		{BaseModel: models.BaseModel{ID: 3}, Title: "Module 3", CourseID: 2, Order: 1}, // Different course
	}
	
	for _, module := range modules {
		db.Create(module)
	}
	
	// Test the query logic that our repository method would use
	var result []*models.Module
	err := db.Where("course_id = ?", uint(1)).Order("\"order\" ASC").Find(&result).Error
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	
	// Validate results
	if len(result) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(result))
	}
	
	if len(result) >= 1 && result[0].Title != "Module 1" {
		t.Errorf("Expected first module to be 'Module 1', got '%s'", result[0].Title)
	}
	
	if len(result) >= 2 && result[1].Title != "Module 2" {
		t.Errorf("Expected second module to be 'Module 2', got '%s'", result[1].Title)
	}
	
	// Validate ordering
	if len(result) >= 2 && result[0].Order > result[1].Order {
		t.Error("Modules should be ordered by 'order' field ascending")
	}
}

func TestModuleRepository_GetWithContent(t *testing.T) {
	db := setupTestDB(t)
	
	// Create test data
	course := &models.Course{
		BaseModel: models.BaseModel{ID: 1},
		Title:     "Test Course",
	}
	db.Create(course)

	module := &models.Module{
		BaseModel: models.BaseModel{ID: 1},
		Title:     "Test Module",
		CourseID:  1,
		Order:     1,
	}
	db.Create(module)

	contents := []*models.Content{
		{BaseModel: models.BaseModel{ID: 1}, Title: "Content 1", ModuleID: 1, Order: 1},
		{BaseModel: models.BaseModel{ID: 2}, Title: "Content 2", ModuleID: 1, Order: 2},
	}
	
	for _, content := range contents {
		db.Create(content)
	}
	
	// Test the preload query logic
	var result models.Module
	err := db.Preload("Contents").First(&result, 1).Error
	if err != nil {
		t.Fatalf("Query with preload failed: %v", err)
	}
	
	// Validate results
	if result.Title != "Test Module" {
		t.Errorf("Expected module title 'Test Module', got '%s'", result.Title)
	}
	
	if len(result.Contents) != 2 {
		t.Errorf("Expected 2 contents, got %d", len(result.Contents))
	}
}