package base

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

// Pagination Result Wrapper (Opsional tapi recommended)
type PaginationResult[T any] struct {
	Data      []T   `json:"data"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
}

type BaseRepository[E any, M any] interface {
	GetDB(ctx context.Context) *gorm.DB

	Create(ctx context.Context, entity *E) error
	Update(ctx context.Context, entity *E) error
	UpdateFields(ctx context.Context, id any, fields map[string]interface{}) error
	Delete(ctx context.Context, id any) error

	FindByID(ctx context.Context, id any, scopes ...func(*gorm.DB) *gorm.DB) (*E, error)
	FindAll(ctx context.Context, page, limit int, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResult[E], error)
	FindOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*E, error)
	Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error)

	// Soft delete specific
	Restore(ctx context.Context, id any) error
	ForceDelete(ctx context.Context, id any) error

	// Batch operations
	CreateBatch(ctx context.Context, entities []*E) error
	DeleteBatch(ctx context.Context, ids []any) error
}
