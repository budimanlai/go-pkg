package base

import (
	"context"

	"gorm.io/gorm"
)

type baseUseaseImpl[E any] struct {
	Repo DomainRepository[E]
	db   *gorm.DB
}

// Constructor Service
// Parameter 'repo' bisa menerima 'BaseRepositoryImpl[E, M]' apapun M-nya.
func NewBaseUsecase[E any](repo DomainRepository[E], db *gorm.DB) BaseUsecase[E] {
	return &baseUseaseImpl[E]{
		Repo: repo,
		db:   db,
	}
}

// Helper untuk akses DB (buat transaction manual di custom service)
func (s *baseUseaseImpl[E]) GetDB() *gorm.DB {
	return s.db
}

// --- Transaction Helper ---
func (s *baseUseaseImpl[E]) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txCtx := InjectTx(ctx, tx)
		return fn(txCtx)
	})
}

// --- Standard CRUD (Simple Delegation) ---
func (s *baseUseaseImpl[E]) Create(ctx context.Context, entity *E) error {
	return s.Repo.Create(ctx, entity)
}

func (s *baseUseaseImpl[E]) Update(ctx context.Context, entity *E) error {
	return s.Repo.Update(ctx, entity)
}

func (s *baseUseaseImpl[E]) UpdateFields(ctx context.Context, id any, fields map[string]interface{}) error {
	return s.Repo.UpdateFields(ctx, id, fields)
}

func (s *baseUseaseImpl[E]) Delete(ctx context.Context, id any) error {
	return s.Repo.Delete(ctx, id)
}

// --- Query Methods ---
func (s *baseUseaseImpl[E]) FindByID(ctx context.Context, id any) (*E, error) {
	return s.Repo.FindByID(ctx, id)
}

func (s *baseUseaseImpl[E]) FindAll(ctx context.Context, page, limit int, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResult[E], error) {
	return s.Repo.FindAll(ctx, page, limit, scopes...)
}

func (s *baseUseaseImpl[E]) FindOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*E, error) {
	return s.Repo.FindOne(ctx, scopes...)
}

func (s *baseUseaseImpl[E]) Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error) {
	return s.Repo.Count(ctx, scopes...)
}

// --- Batch Operations ---
func (s *baseUseaseImpl[E]) CreateBatch(ctx context.Context, entities []*E) error {
	return s.Repo.CreateBatch(ctx, entities)
}

func (s *baseUseaseImpl[E]) DeleteBatch(ctx context.Context, ids []any) error {
	return s.Repo.DeleteBatch(ctx, ids)
}

// --- Soft Delete Management ---
func (s *baseUseaseImpl[E]) Restore(ctx context.Context, id any) error {
	return s.Repo.Restore(ctx, id)
}

func (s *baseUseaseImpl[E]) ForceDelete(ctx context.Context, id any) error {
	return s.Repo.ForceDelete(ctx, id)
}
