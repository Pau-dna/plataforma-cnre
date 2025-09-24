package repositories

import (
	"gorm.io/gorm"
)

// QueryOptimizer provides query optimization utilities for GORM
type QueryOptimizer struct {
	db *gorm.DB
}

// NewQueryOptimizer creates a new query optimizer
func NewQueryOptimizer(db *gorm.DB) *QueryOptimizer {
	return &QueryOptimizer{db: db}
}

// OptimizeForReads applies read optimizations to a query
func (qo *QueryOptimizer) OptimizeForReads(query *gorm.DB) *gorm.DB {
	// Set statement timeout for long queries
	return query.Session(&gorm.Session{
		QueryFields: true, // Only select required fields
	})
}

// OptimizeForBulkInsert applies bulk insert optimizations
func (qo *QueryOptimizer) OptimizeForBulkInsert(query *gorm.DB) *gorm.DB {
	return query.Session(&gorm.Session{
		SkipDefaultTransaction: true, // Skip transaction for better performance
		CreateBatchSize:        100,  // Batch size for bulk operations
	})
}

// OptimizeForUpdates applies update optimizations
func (qo *QueryOptimizer) OptimizeForUpdates(query *gorm.DB) *gorm.DB {
	return query.Session(&gorm.Session{
		SkipHooks: true, // Skip hooks for better performance when not needed
	})
}

// ApplyIndexHint adds index hints for PostgreSQL queries
func (qo *QueryOptimizer) ApplyIndexHint(query *gorm.DB, indexName string) *gorm.DB {
	// PostgreSQL doesn't support index hints like MySQL, but we can use 
	// query optimization techniques and ensure proper WHERE clause ordering
	return query
}