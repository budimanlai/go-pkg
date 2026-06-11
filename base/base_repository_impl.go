package base

import (
	"context"
	"errors"
	"math"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// Implementasi Struct
type BaseRepositoryImpl[E any, M any] struct {
	db *gorm.DB
}

// InjectTx: Memasukkan object Transaction ke dalam Context
func InjectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// ExtractTx: Mengambil object Transaction dari Context (jika ada)
func ExtractTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}

// Constructor Public
func NewGormRepository[E any, M any](db *gorm.DB) BaseRepository[E, M] {
	return &BaseRepositoryImpl[E, M]{
		db: db,
	}
}

// Helper internal untuk memilih DB mana yang dipakai
func (r *BaseRepositoryImpl[E, M]) GetDB(ctx context.Context) *gorm.DB {
	// 1. Cek apakah ada Transaksi "titipan" di context?
	tx := ExtractTx(ctx)
	if tx != nil {
		// Jika ada, gunakan TX tersebut, dan jangan lupa WithContext
		return tx.WithContext(ctx)
	}

	// 2. Jika tidak ada, gunakan DB default milik repo
	return r.db.WithContext(ctx)
}

func (r *BaseRepositoryImpl[E, M]) Create(ctx context.Context, entity *E) error {
	var model M
	if err := copier.Copy(&model, entity); err != nil {
		return err
	}

	if err := r.GetDB(ctx).Create(&model).Error; err != nil {
		return err
	}

	if err := copier.Copy(entity, &model); err != nil {
		return err
	}

	return nil
}

func (r *BaseRepositoryImpl[E, M]) FindByID(ctx context.Context, id any, scopes ...func(*gorm.DB) *gorm.DB) (*E, error) {
	var model M

	// 1. Ambil DB dasar
	db := r.GetDB(ctx)

	// 2. Apply Scopes (misal: Preload("Profile"))
	for _, scope := range scopes {
		db = scope(db)
	}

	// 3. Eksekusi
	err := db.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var entity E
	if err := copier.Copy(&entity, &model); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *BaseRepositoryImpl[E, M]) UpdateFields(ctx context.Context, id any, fields map[string]interface{}) error {
	return r.GetDB(ctx).Model(new(M)).Where("id = ?", id).Updates(fields).Error
}

func (r *BaseRepositoryImpl[E, M]) Update(ctx context.Context, entity *E) error {
	var model M
	if err := copier.Copy(&model, entity); err != nil {
		return err
	}
	return r.GetDB(ctx).Save(&model).Error
}

func (r *BaseRepositoryImpl[E, M]) Delete(ctx context.Context, id any) error {
	var model M
	return r.GetDB(ctx).Delete(&model, id).Error
}

// 5. LIST with Pagination
func (r *BaseRepositoryImpl[E, M]) FindAll(ctx context.Context, page, limit int, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResult[E], error) {
	var models []M
	var total int64

	// Default value safety
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Mulai build query
	db := r.GetDB(ctx).Model(new(M))

	// Apply Filter Dinamis (Scopes) SEBELUM count
	for _, scope := range scopes {
		db = scope(db)
	}

	if err := db.Session(&gorm.Session{}).
		Limit(-1).
		Offset(-1).
		Count(&total).Error; err != nil {
		return PaginationResult[E]{}, err
	}

	// Pagination Offset
	offset := (page - 1) * limit

	// Ambil Data
	if err := db.Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return PaginationResult[E]{}, err
	}

	var entities []E
	// Copier pintar, dia bisa copy slice ke slice otomatis
	if err := copier.Copy(&entities, &models); err != nil {
		return PaginationResult[E]{}, err
	}

	// Hitung Total Page
	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return PaginationResult[E]{
		Data:      entities,
		Total:     total,
		TotalPage: totalPage,
		Page:      page,
		Limit:     limit,
	}, nil
}

func (r *BaseRepositoryImpl[E, M]) Restore(ctx context.Context, id any) error {
	var models M
	return r.GetDB(ctx).Unscoped().Model(&models).
		Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *BaseRepositoryImpl[E, M]) ForceDelete(ctx context.Context, id any) error {
	var models M
	return r.GetDB(ctx).Unscoped().Delete(&models, id).Error
}

// FindOne: Find single entity by any condition using scopes
func (r *BaseRepositoryImpl[E, M]) FindOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*E, error) {
	var models M

	// Build query
	db := r.GetDB(ctx)

	// Apply scopes (e.g., Where conditions, Preload, etc.)
	for _, scope := range scopes {
		db = scope(db)
	}

	// Execute
	err := db.First(&models).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var entities E
	// Copier pintar, dia bisa copy slice ke slice otomatis
	if err := copier.Copy(&entities, &models); err != nil {
		return nil, err
	}
	return &entities, nil
}

// CreateBatch: Bulk insert entities (efficient batch processing)
func (r *BaseRepositoryImpl[E, M]) CreateBatch(ctx context.Context, entities []*E) error {
	if len(entities) == 0 {
		return nil
	}

	// Convert entities to models
	var models []M
	if err := copier.Copy(&models, &entities); err != nil {
		return err
	}

	// GORM's CreateInBatches automatically handles chunking
	// Default batch size: 100 records per INSERT
	if err := r.GetDB(ctx).CreateInBatches(&models, 100).Error; err != nil {
		return err
	}

	// Copy back IDs to entities
	if err := copier.Copy(&entities, &models); err != nil {
		return err
	}

	return nil
}

// DeleteBatch: Delete multiple entities by IDs
func (r *BaseRepositoryImpl[E, M]) DeleteBatch(ctx context.Context, ids []any) error {
	if len(ids) == 0 {
		return nil
	}

	var model M
	// DELETE FROM table WHERE id IN (?, ?, ?)
	return r.GetDB(ctx).Delete(&model, ids).Error
}

// Count: Count entities with optional filters
func (r *BaseRepositoryImpl[E, M]) Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var count int64

	// Build query
	db := r.GetDB(ctx).Model(new(M))

	// Apply scopes (filters)
	for _, scope := range scopes {
		db = scope(db)
	}

	// Execute count
	err := db.Session(&gorm.Session{}).Limit(-1).Offset(-1).Count(&count).Error
	return count, err
}
