# Types Package

The `types` package provides custom time types for consistent UTC time handling in JSON serialization and deserialization.

## Features

- ⏰ Automatic UTC conversion for time values
- 📝 Consistent RFC3339 format with 'Z' suffix
- 🔄 Seamless JSON marshaling and unmarshaling
- 🌍 Timezone-agnostic time handling
- ✅ Database-friendly UTC storage

## Installation

This package is part of `github.com/budimanlai/go-pkg`. Import it as:

```go
import "github.com/budimanlai/go-pkg/types"
```

## UTCTime

`UTCTime` is a custom time type that wraps `time.Time` and ensures all time values are stored and transmitted in UTC timezone.

### Why UTCTime?

- **Consistency**: All times are automatically converted to UTC
- **JSON Format**: Always outputs RFC3339 format with 'Z' suffix (e.g., "2025-11-15T04:56:56Z")
- **Timezone Safety**: Prevents timezone-related bugs in distributed systems
- **Database Compatibility**: Standardized UTC storage for all time values

## API Reference

### Type Definition

```go
type UTCTime time.Time
```

A custom time type that implements `json.Marshaler` and `json.Unmarshaler` interfaces.

### Methods

#### MarshalJSON
```go
func (t UTCTime) MarshalJSON() ([]byte, error)
```
Converts UTCTime to JSON string in RFC3339 format with UTC timezone (Z suffix).

**Output Format:** `"YYYY-MM-DDTHH:MM:SSZ"`

**Example:**
```go
t := types.UTCTime(time.Now())
json, _ := t.MarshalJSON()
// Output: []byte(`"2025-11-15T04:56:56Z"`)
```

#### UnmarshalJSON
```go
func (t *UTCTime) UnmarshalJSON(data []byte) error
```
Parses a JSON time string in RFC3339 format into UTCTime.

**Accepted Formats:**
- `"2025-11-15T04:56:56Z"` (UTC with Z)
- `"2025-11-15T04:56:56+07:00"` (with timezone offset)
- Any valid RFC3339 format

**Example:**
```go
var t types.UTCTime
err := t.UnmarshalJSON([]byte(`"2025-11-15T04:56:56Z"`))
```

#### String
```go
func (t UTCTime) String() string
```
Returns string representation in UTC RFC3339 format.

**Example:**
```go
t := types.UTCTime(time.Now())
fmt.Println(t.String())
// Output: 2025-11-15T04:56:56Z
```

#### ToTime
```go
func (t UTCTime) ToTime() time.Time
```
Converts UTCTime back to standard `time.Time`.

**Example:**
```go
utcTime := types.UTCTime(time.Now())
standardTime := utcTime.ToTime()
```

### Helper Functions

#### NewUTCTime
```go
func NewUTCTime(t time.Time) UTCTime
```
Creates a new UTCTime from `time.Time`, automatically converting to UTC.

**Example:**
```go
now := time.Now()
utcNow := types.NewUTCTime(now)
```

#### Now
```go
func Now() UTCTime
```
Returns the current time in UTC as UTCTime.

**Example:**
```go
currentTime := types.Now()
```

## Usage Examples

### Basic Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"
    "github.com/budimanlai/go-pkg/types"
)

type Event struct {
    ID        int            `json:"id"`
    Name      string         `json:"name"`
    CreatedAt types.UTCTime  `json:"created_at"`
    UpdatedAt types.UTCTime  `json:"updated_at"`
}

func main() {
    event := Event{
        ID:        1,
        Name:      "Conference",
        CreatedAt: types.UTCTime(time.Now()),
        UpdatedAt: types.UTCTime(time.Now()),
    }
    
    // Marshal to JSON
    jsonData, _ := json.Marshal(event)
    fmt.Println(string(jsonData))
    // Output: {"id":1,"name":"Conference","created_at":"2025-11-15T04:56:56Z","updated_at":"2025-11-15T04:56:56Z"}
}
```

### Database Models

```go
import (
    "gorm.io/gorm"
    "github.com/budimanlai/go-pkg/types"
)

type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Name      string         `json:"name"`
    Email     string         `json:"email"`
    CreatedAt types.UTCTime  `json:"created_at"`
    UpdatedAt types.UTCTime  `json:"updated_at"`
    DeletedAt *types.UTCTime `gorm:"index" json:"deleted_at,omitempty"`
}

func createUser(db *gorm.DB) {
    user := User{
        Name:      "John Doe",
        Email:     "john@example.com",
        CreatedAt: types.UTCTime(time.Now()),
        UpdatedAt: types.UTCTime(time.Now()),
    }
    
    db.Create(&user)
    // Timestamps are stored in UTC in the database
}
```

### API Responses

```go
import (
    "github.com/gofiber/fiber/v2"
    "github.com/budimanlai/go-pkg/types"
)

type OrderResponse struct {
    OrderID    string        `json:"order_id"`
    Status     string        `json:"status"`
    Total      float64       `json:"total"`
    OrderedAt  types.UTCTime `json:"ordered_at"`
    ShippedAt  *types.UTCTime `json:"shipped_at,omitempty"`
    DeliveredAt *types.UTCTime `json:"delivered_at,omitempty"`
}

app.Get("/orders/:id", func(c *fiber.Ctx) error {
    order := OrderResponse{
        OrderID:   "ORD-12345",
        Status:    "shipped",
        Total:     299.99,
        OrderedAt: types.UTCTime(time.Now().Add(-48 * time.Hour)),
        ShippedAt: &types.UTCTime(time.Now().Add(-24 * time.Hour)),
    }
    
    return c.JSON(order)
    // All timestamps are automatically in UTC format
})
```

### JSON Marshaling & Unmarshaling

```go
import (
    "encoding/json"
    "github.com/budimanlai/go-pkg/types"
)

type Task struct {
    Title     string        `json:"title"`
    DueDate   types.UTCTime `json:"due_date"`
    Completed bool          `json:"completed"`
}

func example() {
    // Marshal to JSON
    task := Task{
        Title:     "Complete documentation",
        DueDate:   types.UTCTime(time.Date(2025, 11, 15, 12, 0, 0, 0, time.UTC)),
        Completed: false,
    }
    
    jsonData, _ := json.Marshal(task)
    fmt.Println(string(jsonData))
    // Output: {"title":"Complete documentation","due_date":"2025-11-15T12:00:00Z","completed":false}
    
    // Unmarshal from JSON
    jsonStr := `{"title":"Review code","due_date":"2025-11-16T10:00:00Z","completed":false}`
    var newTask Task
    json.Unmarshal([]byte(jsonStr), &newTask)
    
    fmt.Println(newTask.DueDate.String())
    // Output: 2025-11-16T10:00:00Z
}
```

### Timezone Conversion

```go
import (
    "time"
    "github.com/budimanlai/go-pkg/types"
)

func timezoneExample() {
    // Time in Jakarta timezone (UTC+7)
    jakartaLoc, _ := time.LoadLocation("Asia/Jakarta")
    jakartaTime := time.Date(2025, 11, 15, 14, 30, 0, 0, jakartaLoc)
    
    // Convert to UTCTime - automatically converts to UTC
    utcTime := types.UTCTime(jakartaTime)
    
    fmt.Println("Jakarta time:", jakartaTime)
    // Output: 2025-11-15 14:30:00 +0700 WIB
    
    fmt.Println("UTC time:", utcTime.String())
    // Output: 2025-11-15T07:30:00Z
    
    // When marshaled to JSON, always in UTC
    jsonData, _ := json.Marshal(struct {
        Time types.UTCTime `json:"time"`
    }{Time: utcTime})
    
    fmt.Println(string(jsonData))
    // Output: {"time":"2025-11-15T07:30:00Z"}
}
```

### Comparison and Manipulation

```go
import (
    "time"
    "github.com/budimanlai/go-pkg/types"
)

func compareAndManipulate() {
    now := types.UTCTime(time.Now())
    future := types.UTCTime(time.Now().Add(24 * time.Hour))
    
    // Convert to time.Time for comparison
    nowTime := time.Time(now)
    futureTime := time.Time(future)
    
    if futureTime.After(nowTime) {
        fmt.Println("Future is after now")
    }
    
    // Add duration
    tomorrow := types.UTCTime(time.Time(now).Add(24 * time.Hour))
    
    // Format
    formatted := time.Time(now).Format("2006-01-02")
    fmt.Println(formatted)
}
```

### Nullable Timestamps

```go
type Article struct {
    ID          int            `json:"id"`
    Title       string         `json:"title"`
    PublishedAt *types.UTCTime `json:"published_at,omitempty"`
    ArchivedAt  *types.UTCTime `json:"archived_at,omitempty"`
}

func publishArticle(article *Article) {
    now := types.UTCTime(time.Now())
    article.PublishedAt = &now
    // ArchivedAt remains nil until archived
}

// JSON output for unpublished article:
// {"id":1,"title":"Draft Article"}

// JSON output for published article:
// {"id":1,"title":"Published Article","published_at":"2025-11-15T04:56:56Z"}
```

## Best Practices

1. **Always Use UTCTime**: For API responses and database models to ensure consistency
2. **Timezone Conversion**: Convert local times to UTCTime for storage and transmission
3. **Nullable Fields**: Use `*types.UTCTime` for optional timestamp fields
4. **Database Storage**: Store all times in UTC in the database
5. **Display**: Convert from UTC to user's local timezone only for display purposes
6. **Comparison**: Convert to `time.Time` for time operations and comparisons

## Testing

Run tests with:
```bash
go test ./types/ -v
```

## License

This package is part of the go-pkg project and follows the same license.
