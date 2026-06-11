package base

import (
	"context"

	"gorm.io/gorm"
)

// -------------------------------------------------------------
// 1. Domain Repository Interface (The Bridge)
// -------------------------------------------------------------
// Interface ini "menyembunyikan" Model (M) dari Usecase.
// BaseRepositoryImpl[E, M] otomatis mengimplementasikan interface ini
// karena method signature-nya cocok (Duck Typing).
type DomainRepository[E any] interface {
	Create(ctx context.Context, entity *E) error
	FindByID(ctx context.Context, id any, scopes ...func(*gorm.DB) *gorm.DB) (*E, error)
	Update(ctx context.Context, entity *E) error
	UpdateFields(ctx context.Context, id any, fields map[string]interface{}) error
	Delete(ctx context.Context, id any) error
	FindAll(ctx context.Context, page, limit int, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResult[E], error)
	FindOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*E, error)

	Restore(ctx context.Context, id any) error
	ForceDelete(ctx context.Context, id any) error
	CreateBatch(ctx context.Context, entities []*E) error
	DeleteBatch(ctx context.Context, ids []any) error
	Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error)
}

// -------------------------------------------------------------
// 2. Base Usecase Interface (The Standard Contract)
// -------------------------------------------------------------
type BaseUsecase[E any] interface {
	// Standard CRUD
	Create(ctx context.Context, entity *E) error
	Update(ctx context.Context, entity *E) error
	Delete(ctx context.Context, id any) error
	FindByID(ctx context.Context, id any) (*E, error)
	FindAll(ctx context.Context, page, limit int, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResult[E], error)

	// Advanced Access (Optional, bisa diekspos jika perlu)
	FindOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*E, error)
	Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error)

	// Batch operations
	CreateBatch(ctx context.Context, entities []*E) error
	DeleteBatch(ctx context.Context, ids []any) error

	// Soft delete management
	Restore(ctx context.Context, id any) error
	ForceDelete(ctx context.Context, id any) error

	// Partial update
	UpdateFields(ctx context.Context, id any, fields map[string]interface{}) error

	// Akses DB untuk Custom Transaction
	GetDB() *gorm.DB

	// Helper untuk Transaction
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}
