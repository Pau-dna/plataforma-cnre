package repositories

import (
	"github.com/imlargo/go-api-template/internal/cache"
	"github.com/imlargo/go-api-template/pkg/kv"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	db        *gorm.DB
	cacheKeys *cache.CacheKeys
	cache     kv.KeyValueStore
	logger    *zap.SugaredLogger
}

func NewRepository(
	db *gorm.DB,
	cacheKeys *cache.CacheKeys,
	cache kv.KeyValueStore,
	logger *zap.SugaredLogger,
) *Repository {
	return &Repository{
		db,
		cacheKeys,
		cache,
		logger,
	}
}

// DB returns the database instance for transactions
func (r *Repository) DB() *gorm.DB {
	return r.db
}
