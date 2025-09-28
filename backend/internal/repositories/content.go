package repositories

import (
	"github.com/imlargo/go-api-template/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ContentProgressResult represents the result of content progress query
type ContentProgressResult struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type ContentRepository interface {
	Get(id uint) (*models.Content, error)
	Create(content *models.Content) error
	Update(content *models.Content) error
	Patch(id uint, data map[string]interface{}) error
	Delete(id uint) error
	GetAll() ([]*models.Content, error)
	GetByModuleID(moduleID uint) ([]*models.Content, error)
	GetContentProgressByModule(userID, moduleID uint) ([]*ContentProgressResult, error)
}

type contentRepository struct {
	*Repository
}

func NewContentRepository(r *Repository) ContentRepository {
	return &contentRepository{
		Repository: r,
	}
}

func (r *contentRepository) Create(content *models.Content) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get the maximum order for this module
		var maxOrder int
		err := tx.Model(&models.Content{}).
			Where("module_id = ?", content.ModuleID).
			Select("COALESCE(MAX(\"order\"), 0)").
			Scan(&maxOrder).Error
		if err != nil {
			return err
		}

		// Set the next order
		content.Order = maxOrder + 1

		// Create the content
		return tx.Create(content).Error
	})
}

func (r *contentRepository) Get(id uint) (*models.Content, error) {
	var content models.Content
	if err := r.db.First(&content, id).Error; err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *contentRepository) Update(content *models.Content) error {
	return r.db.Model(content).Clauses(clause.Returning{}).Updates(content).Error
}

func (r *contentRepository) Patch(id uint, data map[string]interface{}) error {
	return r.db.Model(&models.Content{}).Where("id = ?", id).Updates(data).Error
}

func (r *contentRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete user progress for this content
		if err := tx.Where(&models.UserProgress{ContentID: id}).Delete(&models.UserProgress{}).Error; err != nil {
			return err
		}

		// Delete the content itself
		return tx.Delete(&models.Content{ID: id}).Error
	})
}

func (r *contentRepository) GetAll() ([]*models.Content, error) {
	var contents []*models.Content
	if err := r.db.Find(&contents).Error; err != nil {
		return nil, err
	}
	return contents, nil
}

func (r *contentRepository) GetByModuleID(moduleID uint) ([]*models.Content, error) {
	var contents []*models.Content
	if err := r.db.Where("module_id = ?", moduleID).Order("\"order\" ASC").Find(&contents).Error; err != nil {
		return nil, err
	}
	return contents, nil
}

// GetContentProgressByModule gets all contents in a module with their completion status for a user
func (r *contentRepository) GetContentProgressByModule(userID, moduleID uint) ([]*ContentProgressResult, error) {
	var results []*ContentProgressResult

	query := `
		SELECT 
			c.id,
			c.title,
			CASE 
				WHEN up.completed_at IS NOT NULL THEN true 
				ELSE false 
			END as completed
		FROM contents c
		LEFT JOIN user_progress up ON (
			c.id = up.content_id 
			AND up.user_id = ? 
			AND up.completed_at IS NOT NULL
		)
		WHERE c.module_id = ?
		ORDER BY c."order" ASC
	`

	if err := r.db.Raw(query, userID, moduleID).Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
