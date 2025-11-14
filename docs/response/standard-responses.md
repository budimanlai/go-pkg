# Standard Response Functions

Standard response functions provide basic JSON response formatting without internationalization. These are useful for simple applications or when i18n is not required.

## Success

Returns a 200 OK response with a success message and optional data.

### Signature

```go
func Success(c *fiber.Ctx, message string, data interface{}) error
```

### Parameters

- `c` (*fiber.Ctx) - The Fiber context
- `message` (string) - Success message to include in response
- `data` (interface{}) - Response data (can be nil, struct, map, slice, etc.)

### Response Format

```json
{
  "meta": {
    "success": true,
    "message": "Success message"
  },
  "data": {
    // your data here
  }
}
```

### Examples

**Simple success response:**

```go
app.Get("/health", func(c *fiber.Ctx) error {
    return response.Success(c, "Service is healthy", nil)
})
```

**Success with data:**

```go
app.Get("/users/:id", func(c *fiber.Ctx) error {
    user := User{
        ID:    123,
        Name:  "John Doe",
        Email: "john@example.com",
    }
    return response.Success(c, "User retrieved successfully", user)
})
```

**Success with map data:**

```go
app.Post("/login", func(c *fiber.Ctx) error {
    token := generateToken()
    return response.Success(c, "Login successful", fiber.Map{
        "token": token,
        "expiresIn": 3600,
    })
})
```

## Error

Returns a JSON error response with a custom HTTP status code.

### Signature

```go
func Error(c *fiber.Ctx, code int, message string) error
```

### Parameters

- `c` (*fiber.Ctx) - The Fiber context
- `code` (int) - HTTP status code (e.g., 400, 404, 500)
- `message` (string) - Error message to include in response

### Response Format

```json
{
  "meta": {
    "success": false,
    "message": "Error message"
  },
  "data": null
}
```

### Examples

**Internal server error:**

```go
app.Get("/process", func(c *fiber.Ctx) error {
    if err := processData(); err != nil {
        return response.Error(c, 500, "Failed to process data")
    }
    return response.Success(c, "Data processed", nil)
})
```

**Forbidden error:**

```go
app.Delete("/users/:id", func(c *fiber.Ctx) error {
    if !hasPermission(c) {
        return response.Error(c, 403, "You don't have permission to delete this user")
    }
    deleteUser(c.Params("id"))
    return response.Success(c, "User deleted", nil)
})
```

## BadRequest

Returns a 400 Bad Request response with an error message.

### Signature

```go
func BadRequest(c *fiber.Ctx, message string) error
```

### Parameters

- `c` (*fiber.Ctx) - The Fiber context
- `message` (string) - Error message to include in response

### Response Format

```json
{
  "meta": {
    "success": false,
    "message": "Invalid request"
  },
  "data": null
}
```

### Examples

**Invalid input:**

```go
app.Post("/users", func(c *fiber.Ctx) error {
    var user User
    if err := c.BodyParser(&user); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }
    
    if user.Email == "" {
        return response.BadRequest(c, "Email is required")
    }
    
    return response.Success(c, "User created", user)
})
```

**Invalid parameter:**

```go
app.Get("/users/:id", func(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return response.BadRequest(c, "Invalid user ID format")
    }
    
    user := getUserByID(id)
    return response.Success(c, "User found", user)
})
```

## NotFound

Returns a 404 Not Found response with an error message.

### Signature

```go
func NotFound(c *fiber.Ctx, message string) error
```

### Parameters

- `c` (*fiber.Ctx) - The Fiber context
- `message` (string) - Error message to include in response

### Response Format

```json
{
  "meta": {
    "success": false,
    "message": "Resource not found"
  },
  "data": null
}
```

### Examples

**Resource not found:**

```go
app.Get("/users/:id", func(c *fiber.Ctx) error {
    user := getUserByID(c.Params("id"))
    if user == nil {
        return response.NotFound(c, "User not found")
    }
    return response.Success(c, "User found", user)
})
```

**Route not found:**

```go
app.Get("/api/*", func(c *fiber.Ctx) error {
    return response.NotFound(c, "API endpoint not found")
})
```

## Usage Patterns

### CRUD Operations

```go
// Create
app.Post("/users", func(c *fiber.Ctx) error {
    var user User
    if err := c.BodyParser(&user); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }
    
    if err := createUser(&user); err != nil {
        return response.Error(c, 500, "Failed to create user")
    }
    
    return response.Success(c, "User created successfully", user)
})

// Read
app.Get("/users/:id", func(c *fiber.Ctx) error {
    user := getUserByID(c.Params("id"))
    if user == nil {
        return response.NotFound(c, "User not found")
    }
    return response.Success(c, "User retrieved successfully", user)
})

// Update
app.Put("/users/:id", func(c *fiber.Ctx) error {
    var user User
    if err := c.BodyParser(&user); err != nil {
        return response.BadRequest(c, "Invalid request body")
    }
    
    existing := getUserByID(c.Params("id"))
    if existing == nil {
        return response.NotFound(c, "User not found")
    }
    
    updateUser(c.Params("id"), &user)
    return response.Success(c, "User updated successfully", user)
})

// Delete
app.Delete("/users/:id", func(c *fiber.Ctx) error {
    user := getUserByID(c.Params("id"))
    if user == nil {
        return response.NotFound(c, "User not found")
    }
    
    deleteUser(c.Params("id"))
    return response.Success(c, "User deleted successfully", nil)
})
```

### Error Handling Chain

```go
app.Post("/orders", func(c *fiber.Ctx) error {
    var order Order
    
    // Parse request
    if err := c.BodyParser(&order); err != nil {
        return response.BadRequest(c, "Invalid order data")
    }
    
    // Validate user
    user := getUserByID(order.UserID)
    if user == nil {
        return response.NotFound(c, "User not found")
    }
    
    // Validate product
    product := getProductByID(order.ProductID)
    if product == nil {
        return response.NotFound(c, "Product not found")
    }
    
    // Check stock
    if product.Stock < order.Quantity {
        return response.BadRequest(c, "Insufficient stock")
    }
    
    // Create order
    if err := createOrder(&order); err != nil {
        return response.Error(c, 500, "Failed to create order")
    }
    
    return response.Success(c, "Order created successfully", order)
})
```

## When to Use Standard Responses

Use standard responses when:

- Building simple applications without i18n requirements
- Prototyping or development
- Internal APIs where English is sufficient
- You want full control over error messages

For production applications with multilingual support, consider using [I18n Responses](i18n-responses.md).
