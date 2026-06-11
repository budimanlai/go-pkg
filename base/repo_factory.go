package base

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RepoConfig struct {
	EnableCache      bool
	EnablePrometheus bool
	RedisClient      *redis.Client
}

// Factory Struct
type Factory struct {
	DB     *gorm.DB
	config RepoConfig
}

func NewFactory(db *gorm.DB, cfg RepoConfig) *Factory {
	return &Factory{
		DB:     db,
		config: cfg,
	}
}

func NewRepository[E any, M any](f *Factory) BaseRepository[E, M] {

	// 1. Layer Inti: Database (Gorm)
	// Akses f.DB (karena f sekarang parameter)
	var repo BaseRepository[E, M] = NewGormRepository[E, M](f.DB)

	// 2. Layer Wrapper: Redis (Jika enabled)
	if f.config.EnableCache && f.config.RedisClient != nil {
		repo = &cachedRepository[E, M]{
			next: repo,
			rdb:  f.config.RedisClient,
			ttl:  10 * time.Minute, // Default TTL
		}
	}

	// 3. Layer Wrapper: Prometheus (Jika enabled)
	if f.config.EnablePrometheus {
		repo = &prometheusRepository[E, M]{
			next: repo,
			name: fmt.Sprintf("%T", *new(E)),
		}
	}

	return repo
}
