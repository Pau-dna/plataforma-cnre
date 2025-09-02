package repositories

import (
	"github.com/imlargo/cnre/internal/cache"
	"github.com/imlargo/cnre/pkg/kv"
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
