package base

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type cachedRepository[E any, M any] struct {
	next BaseRepository[E, M]
	rdb  *redis.Client
	ttl  time.Duration
}

func (r *cachedRepository[E, M]) getKey(id any) string {
	return fmt.Sprintf("cache:entity:%T:%v", *new(E), id)
}

// Helper untuk mengambil ID dari Generic Struct menggunakan Reflection
func (r *cachedRepository[E, M]) getIDFromEntity(entity any) (any, bool) {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, false
	}

	// Cari field bernama "ID", "Id", atau "UID"
	possibleNames := []string{"ID", "Id", "UID"}
	for _, name := range possibleNames {
		field := val.FieldByName(name)
		if field.IsValid() && !field.IsZero() {
			return field.Interface(), true
		}
	}
	return nil, false
}

func (r *cachedRepository[E, M]) GetDB(ctx context.Context) *gorm.DB {
	return r.next.GetDB(ctx)
}

// PERBAIKAN: Tambahkan parameter scopes ...func
func (r *cachedRepository[E, M]) FindByID(ctx context.Context, id any, scopes ...func(*gorm.DB) *gorm.DB) (*E, error) {
	// 1. SAFETY CHECK: Jika ada scopes (filter/preload), JANGAN pakai cache.
	// Karena cache key kita cuma based on ID.
	// Nanti bahaya kalau data preload tercampur dengan data polos.
	if len(scopes) > 0 {
		return r.next.FindByID(ctx, id, scopes...)
	}

	// 2. Logic Cache Standar (Hanya jalan kalau query polos by ID)
	key := r.getKey(id)
	val, err := r.rdb.Get(ctx, key).Result()

	if err == nil {
		var entity E
		if err := json.Unmarshal([]byte(val), &entity); err == nil {
			return &entity, nil
		}
	}

	// 3. Cache MISS -> Panggil Repo Asli
	entity, err := r.next.FindByID(ctx, id) // scopes kosong
	if err != nil {
		return nil, err
	}

	// 4. Set Cache
	if entity != nil {
		go func() {
			data, _ := json.Marshal(entity)
			r.rdb.Set(context.Background(), key, data, r.ttl)
		}()
	}

	return entity, nil
}

// Method lain tetap sama (Pastikan signature-nya match interface)
func (r *cachedRepository[E, M]) Create(ctx context.Context, entity *E) error {
	return r.next.Create(ctx, entity)
}

func (r *cachedRepository[E, M]) UpdateFields(ctx context.Context, id any, fields map[string]interface{}) error {
	if err := r.next.UpdateFields(ctx, id, fields); err != nil {
		return err
	}
	r.rdb.Del(context.Background(), r.getKey(id))
	return nil
}

func (r *cachedRepository[E, M]) Update(ctx context.Context, entity *E) error {
	if err := r.next.Update(ctx, entity); err != nil {
		return err
	}

	if id, ok := r.getIDFromEntity(entity); ok {
		// Jika ketemu ID-nya, hapus cache!
		r.rdb.Del(context.Background(), r.getKey(id))
	}
	return nil
}

func (r *cachedRepository[E, M]) Delete(ctx context.Context, id any) error {
	if err := r.next.Delete(ctx, id); err != nil {
		return err
	}
	r.rdb.Del(context.Background(), r.getKey(id))
	return nil
}

func (r *cachedRepository[E, M]) FindAll(ctx context.Context, page, limit int, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResult[E], error) {
	return r.next.FindAll(ctx, page, limit, scopes...)
}

func (r *cachedRepository[E, M]) Restore(ctx context.Context, id any) error {
	if err := r.next.Restore(ctx, id); err != nil {
		return err
	}
	r.rdb.Del(context.Background(), r.getKey(id)) // Invalidate
	return nil
}

func (r *cachedRepository[E, M]) ForceDelete(ctx context.Context, id any) error {
	if err := r.next.ForceDelete(ctx, id); err != nil {
		return err
	}
	r.rdb.Del(context.Background(), r.getKey(id)) // Invalidate
	return nil
}

func (r *cachedRepository[E, M]) FindOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*E, error) {
	// No caching for FindOne (complex filters)
	return r.next.FindOne(ctx, scopes...)
}

func (r *cachedRepository[E, M]) CreateBatch(ctx context.Context, entities []*E) error {
	return r.next.CreateBatch(ctx, entities)
}

func (r *cachedRepository[E, M]) DeleteBatch(ctx context.Context, ids []any) error {
	if err := r.next.DeleteBatch(ctx, ids); err != nil {
		return err
	}
	// Invalidate cache for all deleted IDs
	go func() {
		if len(ids) == 0 {
			return
		}

		// Kumpulkan semua keys dulu
		keys := make([]string, len(ids))
		for i, id := range ids {
			keys[i] = r.getKey(id)
		}

		// Hapus sekaligus (pipeline/variadic)
		r.rdb.Del(context.Background(), keys...)
	}()
	return nil
}

func (r *cachedRepository[E, M]) Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error) {
	// No caching for count (filters may vary)
	return r.next.Count(ctx, scopes...)
}
