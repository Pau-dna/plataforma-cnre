package services

import (
	"github.com/imlargo/cnre/internal/cache"
	"github.com/imlargo/cnre/internal/config"
	"github.com/imlargo/cnre/internal/store"
	"github.com/imlargo/cnre/pkg/kv"
	"go.uber.org/zap"
)

type Service struct {
	store     *store.Store
	logger    *zap.SugaredLogger
	config    *config.AppConfig
	cacheKeys *cache.CacheKeys
	cache     kv.KeyValueStore
}

func NewService(
	store *store.Store,
	logger *zap.SugaredLogger,
	config *config.AppConfig,
	cacheKeys *cache.CacheKeys,
	cache kv.KeyValueStore,
) *Service {
	return &Service{
		store,
		logger,
		config,
		cacheKeys,
		cache,
	}
}
