# Base Package

The `base` package provides a generic 3-layer architecture foundation (Repository, Usecase, Handler) that eliminates boilerplate code for standard CRUD operations, with built-in Redis caching, Prometheus metrics, and transaction support.

## Features

- Generic CRUD operations with type safety (Go generics)
- Entity/Model separation for Clean Architecture
- Automatic Entity ↔ Model conversion via `copier`
- Redis caching with automatic cache invalidation
- Prometheus metrics for all database operations
- Context-based transaction injection
- Composable query scopes (filters, preloads, ordering)
- Pagination with safe limits
- Batch operations (CreateBatch, DeleteBatch)
- Soft delete with Restore and ForceDelete

## Installation

This package is part of `github.com/budimanlai/go-pkg`. Import it as:

```go
import "github.com/budimanlai/go-pkg/base"
```

### Dependencies

```bash
go get gorm.io/gorm
go get github.com/redis/go-redis/v9
go get github.com/prometheus/client_golang
go get github.com/jinzhu/copier
go get github.com/gofiber/fiber/v2
```

---

## Architecture

The package implements a layered decorator pattern:

```
HTTP Request
     │
┌────▼───────────────────────────────────┐
│  BaseHandler[E, C, U]                  │  ← Parse, Validate, Map DTO → Entity
└────┬───────────────────────────────────┘
     │
┌────▼───────────────────────────────────┐
│  BaseUsecase[E]                        │  ← Business logic, transactions
└────┬───────────────────────────────────┘
     │
┌────▼───────────────────────────────────┐
│  Prometheus Decorator                  │  ← Metrics (optional)
└────┬───────────────────────────────────┘
     │
┌────▼───────────────────────────────────┐
│  Redis Decorator                       │  ← Caching (optional)
└────┬───────────────────────────────────┘
     │
┌────▼───────────────────────────────────┐
│  BaseRepositoryImpl[E, M]              │  ← Entity ↔ Model conversion + GORM
└────┬───────────────────────────────────┘
     │
  Database
```

### Entity vs Model

| | Entity (E) | Model (M) |
|---|---|---|
| Layer | Domain | Persistence |
| Purpose | Business logic | Database mapping |
| GORM tags | No | Yes |
| Used by | Handler, Usecase | Repository (internal) |
| Cached | Yes | No |

---

## Quick Start (5 minutes)

### Step 1: Define Domain (Entity + Model + DTOs)

```go
// domain/user.go
package domain

import (
    "time"
    "gorm.io/gorm"
)

// Entity - business domain object
type UserEntity struct {
    ID        uint      `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}

// Model - database representation (GORM tags here, not in Entity)
type UserModel struct {
    ID        uint           `gorm:"primaryKey"`
    Email     string         `gorm:"uniqueIndex;not null"`
    Name      string         `gorm:"not null"`
    Status    string         `gorm:"default:'active'"`
    CreatedAt time.Time      `gorm:"autoCreateTime"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserModel) TableName() string { return "users" }

// DTOs for HTTP input validation
type CreateUserRequest struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required,min=3"`
}

type UpdateUserRequest struct {
    Name   string `json:"name" validate:"required,min=3"`
    Status string `json:"status" validate:"omitempty,oneof=active inactive"`
}
```

### Step 2: Create Repository

```go
// repository/user_repository.go
package repository

import (
    "github.com/budimanlai/go-pkg/base"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
    "yourapp/domain"
)

type UserRepository interface {
    base.DomainRepository[domain.UserEntity]
    // Add custom methods here if needed
}

type userRepositoryImpl struct {
    base.BaseRepository[domain.UserEntity, domain.UserModel]
}

func NewUserRepository(db *gorm.DB, rdb *redis.Client) UserRepository {
    factory := base.NewFactory(db, base.RepoConfig{
        EnableCache:      true,
        EnablePrometheus: true,
        RedisClient:      rdb,
    })

    return &userRepositoryImpl{
        BaseRepository: base.NewRepository[domain.UserEntity, domain.UserModel](factory),
    }
}
```

That's it — 13 CRUD methods are available automatically.

### Step 3: Create Usecase

```go
// usecase/user_usecase.go
package usecase

import (
    "context"
    "errors"
    "github.com/budimanlai/go-pkg/base"
    "gorm.io/gorm"
    "yourapp/domain"
    "yourapp/repository"
)

type UserUsecase interface {
    base.BaseUsecase[domain.UserEntity]
    FindByEmail(ctx context.Context, email string) (*domain.UserEntity, error)
}

type userUsecaseImpl struct {
    base.BaseUsecase[domain.UserEntity]
    repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository, db *gorm.DB) UserUsecase {
    return &userUsecaseImpl{
        BaseUsecase: base.NewBaseUsecase[domain.UserEntity](repo, db),
        repo:        repo,
    }
}

// Override Create to add business logic
func (s *userUsecaseImpl) Create(ctx context.Context, entity *domain.UserEntity) error {
    existing, _ := s.FindByEmail(ctx, entity.Email)
    if existing != nil {
        return errors.New("email already exists")
    }
    return s.BaseUsecase.Create(ctx, entity)
}

func (s *userUsecaseImpl) FindByEmail(ctx context.Context, email string) (*domain.UserEntity, error) {
    return s.FindOne(ctx, func(db *gorm.DB) *gorm.DB {
        return db.Where("email = ?", email)
    })
}
```

### Step 4: Create Handler

```go
// handler/user_handler.go
package handler

import (
    "github.com/budimanlai/go-pkg/base"
    "yourapp/domain"
    "yourapp/usecase"
)

type UserHandler struct {
    *base.BaseHandler[domain.UserEntity, domain.CreateUserRequest, domain.UpdateUserRequest]
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
    return &UserHandler{
        BaseHandler: base.NewBaseHandler[domain.UserEntity, domain.CreateUserRequest, domain.UpdateUserRequest](uc),
    }
}
```

5 RESTful endpoints ready: `Index`, `View`, `Create`, `Update`, `Delete`.

### Step 5: Wire Up in main.go

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/redis/go-redis/v9"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "yourapp/domain"
    "yourapp/handler"
    "yourapp/repository"
    "yourapp/usecase"
)

func main() {
    db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/dbname"), &gorm.Config{})
    db.AutoMigrate(&domain.UserModel{})

    rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

    userRepo    := repository.NewUserRepository(db, rdb)
    userUsecase := usecase.NewUserUsecase(userRepo, db)
    userHandler := handler.NewUserHandler(userUsecase)

    app := fiber.New()
    users := app.Group("/api/users")
    users.Get("/", userHandler.Index)        // GET  /api/users?page=1&limit=10
    users.Get("/:id", userHandler.View)      // GET  /api/users/:id
    users.Post("/", userHandler.Create)      // POST /api/users
    users.Put("/:id", userHandler.Update)    // PUT  /api/users/:id
    users.Delete("/:id", userHandler.Delete) // DELETE /api/users/:id

    app.Listen(":3000")
}
```

---

## API Reference

### BaseRepository[E, M]

The core data access layer. `E` = Entity (domain), `M` = Model (database).

```go
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
    Restore(ctx context.Context, id any) error
    ForceDelete(ctx context.Context, id any) error
    CreateBatch(ctx context.Context, entities []*E) error
    DeleteBatch(ctx context.Context, ids []any) error
}
```

#### Create

```go
func Create(ctx context.Context, entity *E) error
```

Inserts a single record. Entity → Model conversion is automatic. The inserted ID is copied back into `entity`.

```go
user := &domain.UserEntity{Email: "john@example.com", Name: "John"}
err := repo.Create(ctx, user)
// user.ID is now populated
```

#### FindByID

```go
func FindByID(ctx context.Context, id any, scopes ...func(*gorm.DB) *gorm.DB) (*E, error)
```

Returns `(entity, nil)` if found, `(nil, nil)` if not found, `(nil, err)` on error.
Results are cached automatically when no scopes are provided.

```go
// Simple lookup (cached)
user, err := repo.FindByID(ctx, 1)

// With preload (bypasses cache)
user, err := repo.FindByID(ctx, 1, func(db *gorm.DB) *gorm.DB {
    return db.Preload("Profile")
})
```

#### FindOne

```go
func FindOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*E, error)
```

Find the first record matching the given conditions. Not cached.

```go
user, err := repo.FindOne(ctx, func(db *gorm.DB) *gorm.DB {
    return db.Where("email = ?", "john@example.com")
})
```

#### FindAll

```go
func FindAll(ctx context.Context, page, limit int, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResult[E], error)
```

Paginated list. `page` is 1-indexed, `limit` is capped at 100.

```go
type PaginationResult[E any] struct {
    Data      []E   `json:"data"`
    Total     int64 `json:"total"`
    TotalPage int   `json:"total_page"`
    Page      int   `json:"page"`
    Limit     int   `json:"limit"`
}
```

```go
result, err := repo.FindAll(ctx, 1, 20, func(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", "active").Order("created_at DESC")
})
fmt.Printf("Page %d of %d, total %d records\n", result.Page, result.TotalPage, result.Total)
```

#### Update

```go
func Update(ctx context.Context, entity *E) error
```

Updates all fields using GORM `Save()`. Entity → Model conversion is automatic.

```go
user, _ := repo.FindByID(ctx, 1)
user.Name = "New Name"
err := repo.Update(ctx, user)
```

#### UpdateFields

```go
func UpdateFields(ctx context.Context, id any, fields map[string]interface{}) error
```

Partial update without loading the entity. Automatically invalidates cache.

```go
err := repo.UpdateFields(ctx, 1, map[string]interface{}{
    "status":     "inactive",
    "updated_at": time.Now(),
})
```

#### Delete / ForceDelete / Restore

```go
func Delete(ctx context.Context, id any) error      // Soft delete if DeletedAt exists
func ForceDelete(ctx context.Context, id any) error // Permanent delete
func Restore(ctx context.Context, id any) error     // Undo soft delete
```

```go
repo.Delete(ctx, 1)      // sets deleted_at
repo.Restore(ctx, 1)     // clears deleted_at
repo.ForceDelete(ctx, 1) // DELETE FROM table WHERE id = 1
```

All three automatically invalidate the cache for the given ID.

#### CreateBatch / DeleteBatch

```go
func CreateBatch(ctx context.Context, entities []*E) error
func DeleteBatch(ctx context.Context, ids []any) error
```

```go
users := []*domain.UserEntity{
    {Email: "a@example.com", Name: "User A"},
    {Email: "b@example.com", Name: "User B"},
}
err := repo.CreateBatch(ctx, users) // inserts 100 records per batch
// All IDs populated after insert

repo.DeleteBatch(ctx, []any{1, 2, 3})
```

#### Count

```go
func Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error)
```

More efficient than `FindAll(...).Total` when you only need a count.

```go
total, _ := repo.Count(ctx)

activeCount, _ := repo.Count(ctx, func(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", "active")
})
```

---

### DomainRepository[E]

The bridge interface used by Usecase. Same methods as `BaseRepository` but without the Model type parameter — the Usecase never sees `M`.

```go
type DomainRepository[E any] interface {
    Create(ctx context.Context, entity *E) error
    FindByID(ctx context.Context, id any, scopes ...func(*gorm.DB) *gorm.DB) (*E, error)
    // ... (same as BaseRepository minus GetDB)
}
```

---

### BaseUsecase[E]

Business logic layer.

```go
type BaseUsecase[E any] interface {
    Create(ctx context.Context, entity *E) error
    Update(ctx context.Context, entity *E) error
    Delete(ctx context.Context, id any) error
    FindByID(ctx context.Context, id any) (*E, error)
    FindAll(ctx context.Context, page, limit int, scopes ...func(*gorm.DB) *gorm.DB) (PaginationResult[E], error)
    FindOne(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (*E, error)
    Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error)
    CreateBatch(ctx context.Context, entities []*E) error
    DeleteBatch(ctx context.Context, ids []any) error
    Restore(ctx context.Context, id any) error
    ForceDelete(ctx context.Context, id any) error
    UpdateFields(ctx context.Context, id any, fields map[string]interface{}) error
    GetDB() *gorm.DB
    WithTransaction(ctx context.Context, fn func(context.Context) error) error
}
```

#### Constructor

```go
func NewBaseUsecase[E any](repo DomainRepository[E], db *gorm.DB) BaseUsecase[E]
```

#### WithTransaction

Wraps multiple operations in a single database transaction. Pass the `txCtx` to all repository calls inside the function.

```go
err := usecase.WithTransaction(ctx, func(txCtx context.Context) error {
    if err := orderRepo.Create(txCtx, order); err != nil {
        return err // triggers rollback
    }
    if err := inventoryRepo.UpdateFields(txCtx, productID, map[string]interface{}{
        "stock": gorm.Expr("stock - ?", order.Qty),
    }); err != nil {
        return err // triggers rollback
    }
    return nil // commits
})
```

---

### BaseHandler[E, C, U]

HTTP layer. `E` = Entity, `C` = Create DTO, `U` = Update DTO.

```go
type BaseHandler[E any, C any, U any] struct {
    Service BaseUsecase[E]
}

func NewBaseHandler[E any, C any, U any](service BaseUsecase[E]) *BaseHandler[E, C, U]
```

Provides 5 methods compatible with Fiber:

| Method | Route | Description |
|--------|-------|-------------|
| `Index` | `GET /` | List with pagination (`?page=1&limit=10`) |
| `View` | `GET /:id` | Find by ID |
| `Create` | `POST /` | Parse C DTO → validate → copy to E → create |
| `Update` | `PUT /:id` | Load existing → parse U DTO → validate → copy → update |
| `Delete` | `DELETE /:id` | Delete by ID |

To override a method while keeping the others:

```go
type UserHandler struct {
    *base.BaseHandler[domain.UserEntity, domain.CreateUserRequest, domain.UpdateUserRequest]
}

// Override only Create
func (h *UserHandler) Create(c *fiber.Ctx) error {
    // custom pre-logic
    return h.BaseHandler.Create(c) // call base
}
```

---

### Factory and RepoConfig

```go
type RepoConfig struct {
    EnableCache      bool
    EnablePrometheus bool
    RedisClient      *redis.Client
}

func NewFactory(db *gorm.DB, cfg RepoConfig) *Factory
func NewRepository[E any, M any](f *Factory) BaseRepository[E, M]
```

Decorators are applied in order: `GORM → Redis → Prometheus`.

```go
// Development (no overhead)
factory := base.NewFactory(db, base.RepoConfig{})

// Production
factory := base.NewFactory(db, base.RepoConfig{
    EnableCache:      true,
    EnablePrometheus: true,
    RedisClient:      rdb,
})

repo := base.NewRepository[UserEntity, UserModel](factory)
```

---

## Transaction Support

### Via InjectTx (manual)

```go
err := db.Transaction(func(tx *gorm.DB) error {
    txCtx := base.InjectTx(ctx, tx)

    if err := userRepo.Create(txCtx, user); err != nil {
        return err
    }
    if err := profileRepo.Create(txCtx, profile); err != nil {
        return err
    }
    return nil
})
```

### Via WithTransaction (usecase helper)

```go
err := svc.WithTransaction(ctx, func(txCtx context.Context) error {
    if err := svc.Create(txCtx, order); err != nil {
        return err
    }
    // more operations with txCtx...
    return nil
})
```

---

## Scopes Pattern

Scopes are `func(*gorm.DB) *gorm.DB` functions that compose filters, ordering, and preloads.

```go
// Define reusable scopes
func ActiveScope(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", "active")
}

func CreatedAfter(t time.Time) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("created_at > ?", t)
    }
}

func WithProfile(db *gorm.DB) *gorm.DB {
    return db.Preload("Profile")
}

// Compose at call site
result, err := repo.FindAll(ctx, 1, 20,
    ActiveScope,
    CreatedAfter(lastWeek),
    WithProfile,
)
```

---

## Caching Behavior

Redis caching is applied automatically by the factory when `EnableCache: true`.

| Operation | Cached? | Cache Invalidated? |
|-----------|---------|-------------------|
| `FindByID` (no scopes) | Yes (TTL 10min) | — |
| `FindByID` (with scopes) | No | — |
| `FindOne` | No | — |
| `FindAll` | No | — |
| `UpdateFields(id, ...)` | — | Yes |
| `Update(entity)` | — | Yes (by reflection) |
| `Delete(id)` | — | Yes |
| `DeleteBatch(ids)` | — | Yes (all IDs) |
| `Restore(id)` | — | Yes |
| `ForceDelete(id)` | — | Yes |

Cache key format: `cache:entity:<TypeName>:<id>`

---

## Prometheus Metrics

Metric name: `gocore_db_query_duration_seconds`  
Type: Histogram  
Labels: `entity`, `operation`, `status` (`success` or `error`)

```promql
# Average query latency
rate(gocore_db_query_duration_seconds_sum[5m])
  / rate(gocore_db_query_duration_seconds_count[5m])

# 95th percentile latency
histogram_quantile(0.95, rate(gocore_db_query_duration_seconds_bucket[5m]))

# Error rate
sum(rate(gocore_db_query_duration_seconds_count{status="error"}[5m]))
  / sum(rate(gocore_db_query_duration_seconds_count[5m]))
```

---

## Project Structure

Recommended layout when using the base package:

```
yourapp/
├── main.go
├── domain/
│   └── user.go          # UserEntity, UserModel, CreateUserRequest, UpdateUserRequest
├── repository/
│   └── user_repository.go
├── usecase/
│   └── user_usecase.go
└── handler/
    └── user_handler.go
```

---

## Best Practices

**Handle not-found correctly** — `FindByID` returns `(nil, nil)` when not found, not an error:

```go
user, err := repo.FindByID(ctx, id)
if err != nil {
    return err        // DB error
}
if user == nil {
    return ErrNotFound // record doesn't exist
}
```

**Prefer UpdateFields over Update for partial updates:**

```go
// Avoid: loads full entity then saves all fields
user, _ := repo.FindByID(ctx, 1)
user.Status = "inactive"
repo.Update(ctx, user)

// Prefer: targeted update, also invalidates cache
repo.UpdateFields(ctx, 1, map[string]interface{}{"status": "inactive"})
```

**Use Count instead of FindAll when you only need the total:**

```go
// Avoid
result, _ := repo.FindAll(ctx, 1, 1)
count := result.Total

// Prefer
count, _ := repo.Count(ctx)
```

**Use CreateBatch for bulk inserts:**

```go
// Avoid: N separate INSERT queries
for _, u := range users {
    repo.Create(ctx, u)
}

// Prefer: batched 100 records per INSERT
repo.CreateBatch(ctx, users)
```

---

## Troubleshooting

**Cache not working**  
Verify `EnableCache: true` and `RedisClient != nil`. Check Redis connectivity with `rdb.Ping(ctx)`. Note that `FindByID` with scopes intentionally skips the cache.

**Metrics not appearing**  
Verify `EnablePrometheus: true` and that `/metrics` is exposed. Look for `gocore_db_query_duration_seconds`.

**Transaction not rolling back**  
Ensure the error is returned from the `fn` function passed to `WithTransaction` or `db.Transaction`. Only returned errors trigger rollback.

**FindAll total count is wrong**  
Check if soft-deleted records are being excluded (expected behavior when model has `DeletedAt`). Use `db.Unscoped()` in a scope to include them.
