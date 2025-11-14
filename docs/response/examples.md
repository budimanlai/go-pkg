# Response Package Examples

Practical examples and common patterns for using the response package in real-world applications.

## Table of Contents

- [REST API CRUD](#rest-api-crud)
- [Authentication & Authorization](#authentication--authorization)
- [File Upload](#file-upload)
- [Pagination](#pagination)
- [Search & Filtering](#search--filtering)
- [Batch Operations](#batch-operations)
- [Multilingual API](#multilingual-api)

## REST API CRUD

Complete CRUD implementation with validation and i18n support:

```go
package main

import (
    "github.com/budimanlai/go-pkg/i18n"
    "github.com/budimanlai/go-pkg/response"
    "github.com/budimanlai/go-pkg/validator"
    "github.com/gofiber/fiber/v2"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name" validate:"required,min=3,max=50"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"required,gte=18"`
}

func main() {
    // Setup i18n
    i18nMgr, _ := i18n.NewI18nManager(i18n.Config{
        LocalesPath:     "./locales",
        DefaultLanguage: "en",
    })
    response.SetI18nManager(i18nMgr)
    validator.SetI18nManager(i18nMgr)

    // Create Fiber app
    app := fiber.New(fiber.Config{
        ErrorHandler: response.FiberErrorHandler,
    })

    app.Use(i18n.I18nMiddleware(i18nMgr))

    // Routes
    api := app.Group("/api/v1")
    
    // Create
    api.Post("/users", createUser)
    
    // Read all
    api.Get("/users", getUsers)
    
    // Read one
    api.Get("/users/:id", getUser)
    
    // Update
    api.Put("/users/:id", updateUser)
    
    // Delete
    api.Delete("/users/:id", deleteUser)

    app.Listen(":3000")
}

func createUser(c *fiber.Ctx) error {
    var user User
    
    // Parse request body
    if err := c.BodyParser(&user); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    // Validate
    if err := validator.ValidateStructWithContext(c, user); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Check if email already exists
    if emailExists(user.Email) {
        return response.BadRequestI18n(c, "email_already_exists", map[string]string{
            "Email": user.Email,
        })
    }
    
    // Create user in database
    user.ID = saveUser(&user)
    
    return response.SuccessI18n(c, "user_created", user)
}

func getUsers(c *fiber.Ctx) error {
    users := getAllUsers()
    return response.SuccessI18n(c, "users_retrieved", users)
}

func getUser(c *fiber.Ctx) error {
    id := c.Params("id")
    user := getUserByID(id)
    
    if user == nil {
        return response.NotFoundI18n(c, "user_not_found")
    }
    
    return response.SuccessI18n(c, "user_retrieved", user)
}

func updateUser(c *fiber.Ctx) error {
    id := c.Params("id")
    
    // Check if user exists
    existing := getUserByID(id)
    if existing == nil {
        return response.NotFoundI18n(c, "user_not_found")
    }
    
    var user User
    if err := c.BodyParser(&user); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    // Validate
    if err := validator.ValidateStructWithContext(c, user); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Update user
    user.ID = existing.ID
    updateUserInDB(&user)
    
    return response.SuccessI18n(c, "user_updated", user)
}

func deleteUser(c *fiber.Ctx) error {
    id := c.Params("id")
    
    // Check if user exists
    user := getUserByID(id)
    if user == nil {
        return response.NotFoundI18n(c, "user_not_found")
    }
    
    // Delete user
    deleteUserFromDB(id)
    
    return response.SuccessI18n(c, "user_deleted", nil)
}
```

## Authentication & Authorization

```go
type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
    Name     string `json:"name" validate:"required,min=3"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

func register(c *fiber.Ctx) error {
    var req RegisterRequest
    
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    if err := validator.ValidateStructWithContext(c, req); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Check if email exists
    if userExists(req.Email) {
        return response.BadRequestI18n(c, "email_already_registered", nil)
    }
    
    // Create user
    user := createUser(req.Name, req.Email, req.Password)
    
    // Generate token
    token := generateJWT(user.ID)
    
    return response.SuccessI18n(c, "registration_success", fiber.Map{
        "user":  user,
        "token": token,
    })
}

func login(c *fiber.Ctx) error {
    var req LoginRequest
    
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    if err := validator.ValidateStructWithContext(c, req); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Find user
    user := findUserByEmail(req.Email)
    if user == nil {
        return response.BadRequestI18n(c, "invalid_credentials", nil)
    }
    
    // Verify password
    if !verifyPassword(user.Password, req.Password) {
        return response.BadRequestI18n(c, "invalid_credentials", nil)
    }
    
    // Generate token
    token := generateJWT(user.ID)
    
    return response.SuccessI18n(c, "login_success", fiber.Map{
        "user":  user,
        "token": token,
    })
}

func logout(c *fiber.Ctx) error {
    // Invalidate token (implementation depends on your token strategy)
    token := c.Get("Authorization")
    invalidateToken(token)
    
    return response.SuccessI18n(c, "logout_success", nil)
}

// Middleware for protected routes
func authRequired(c *fiber.Ctx) error {
    token := c.Get("Authorization")
    
    if token == "" {
        return response.ErrorI18n(c, 401, "unauthorized", nil)
    }
    
    user, err := validateToken(token)
    if err != nil {
        return response.ErrorI18n(c, 401, "invalid_token", nil)
    }
    
    c.Locals("user", user)
    return c.Next()
}

// Usage
app.Post("/auth/register", register)
app.Post("/auth/login", login)
app.Post("/auth/logout", authRequired, logout)
app.Get("/profile", authRequired, getProfile)
```

## File Upload

```go
func uploadFile(c *fiber.Ctx) error {
    // Get file from form
    file, err := c.FormFile("file")
    if err != nil {
        return response.BadRequestI18n(c, "file_required", nil)
    }
    
    // Validate file size (max 10MB)
    maxSize := int64(10 * 1024 * 1024)
    if file.Size > maxSize {
        return response.BadRequestI18n(c, "file_too_large", map[string]string{
            "MaxSize": "10MB",
            "Size":    formatBytes(file.Size),
        })
    }
    
    // Validate file type
    allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
    contentType := file.Header.Get("Content-Type")
    if !contains(allowedTypes, contentType) {
        return response.BadRequestI18n(c, "invalid_file_type", map[string]string{
            "Type":    contentType,
            "Allowed": "JPEG, PNG, GIF",
        })
    }
    
    // Save file
    filename := generateFilename(file.Filename)
    path := filepath.Join("./uploads", filename)
    
    if err := c.SaveFile(file, path); err != nil {
        return response.ErrorI18n(c, 500, "file_upload_failed", nil)
    }
    
    return response.SuccessI18n(c, "file_uploaded", fiber.Map{
        "filename": filename,
        "size":     file.Size,
        "url":      "/uploads/" + filename,
    })
}

func uploadMultiple(c *fiber.Ctx) error {
    // Get form
    form, err := c.MultipartForm()
    if err != nil {
        return response.BadRequestI18n(c, "invalid_form_data", nil)
    }
    
    files := form.File["files"]
    if len(files) == 0 {
        return response.BadRequestI18n(c, "files_required", nil)
    }
    
    // Validate max files
    if len(files) > 10 {
        return response.BadRequestI18n(c, "too_many_files", map[string]string{
            "Max":   "10",
            "Count": fmt.Sprintf("%d", len(files)),
        })
    }
    
    uploadedFiles := []fiber.Map{}
    
    for _, file := range files {
        filename := generateFilename(file.Filename)
        path := filepath.Join("./uploads", filename)
        
        if err := c.SaveFile(file, path); err != nil {
            return response.ErrorI18n(c, 500, "file_upload_failed", nil)
        }
        
        uploadedFiles = append(uploadedFiles, fiber.Map{
            "filename": filename,
            "size":     file.Size,
            "url":      "/uploads/" + filename,
        })
    }
    
    return response.SuccessI18n(c, "files_uploaded", fiber.Map{
        "count": len(uploadedFiles),
        "files": uploadedFiles,
    })
}
```

## Pagination

```go
type PaginationQuery struct {
    Page     int `query:"page"`
    PageSize int `query:"page_size"`
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    TotalPages int         `json:"total_pages"`
}

func getUsers(c *fiber.Ctx) error {
    var query PaginationQuery
    
    if err := c.QueryParser(&query); err != nil {
        return response.BadRequestI18n(c, "invalid_query_params", nil)
    }
    
    // Default values
    if query.Page < 1 {
        query.Page = 1
    }
    if query.PageSize < 1 {
        query.PageSize = 10
    }
    if query.PageSize > 100 {
        query.PageSize = 100
    }
    
    // Get data
    offset := (query.Page - 1) * query.PageSize
    users := getUsersPaginated(offset, query.PageSize)
    total := getTotalUsers()
    
    totalPages := int(total) / query.PageSize
    if int(total)%query.PageSize > 0 {
        totalPages++
    }
    
    result := PaginatedResponse{
        Data:       users,
        Total:      total,
        Page:       query.Page,
        PageSize:   query.PageSize,
        TotalPages: totalPages,
    }
    
    return response.SuccessI18n(c, "users_retrieved", result)
}
```

## Search & Filtering

```go
type SearchQuery struct {
    Query    string `query:"q"`
    Status   string `query:"status"`
    Category string `query:"category"`
    SortBy   string `query:"sort_by"`
    Order    string `query:"order"`
    Page     int    `query:"page"`
    PageSize int    `query:"page_size"`
}

func searchUsers(c *fiber.Ctx) error {
    var query SearchQuery
    
    if err := c.QueryParser(&query); err != nil {
        return response.BadRequestI18n(c, "invalid_query_params", nil)
    }
    
    // Validate sort field
    allowedSorts := []string{"name", "email", "created_at"}
    if query.SortBy != "" && !contains(allowedSorts, query.SortBy) {
        return response.BadRequestI18n(c, "invalid_sort_field", map[string]string{
            "Field":   query.SortBy,
            "Allowed": "name, email, created_at",
        })
    }
    
    // Validate order
    if query.Order != "" && query.Order != "asc" && query.Order != "desc" {
        return response.BadRequestI18n(c, "invalid_sort_order", nil)
    }
    
    // Build filters
    filters := buildFilters(query)
    
    // Search
    users := searchUsersWithFilters(filters, query.Page, query.PageSize)
    total := countUsersWithFilters(filters)
    
    return response.SuccessI18n(c, "search_results", fiber.Map{
        "data":  users,
        "total": total,
        "query": query.Query,
    })
}
```

## Batch Operations

```go
type BatchDeleteRequest struct {
    IDs []int `json:"ids" validate:"required,min=1,max=100"`
}

func batchDelete(c *fiber.Ctx) error {
    var req BatchDeleteRequest
    
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequestI18n(c, "invalid_request_body", nil)
    }
    
    if err := validator.ValidateStructWithContext(c, req); err != nil {
        return response.ValidationErrorI18n(c, err)
    }
    
    // Delete users
    deleted := 0
    failed := []int{}
    
    for _, id := range req.IDs {
        if err := deleteUser(id); err != nil {
            failed = append(failed, id)
        } else {
            deleted++
        }
    }
    
    if len(failed) > 0 {
        return response.SuccessI18n(c, "batch_delete_partial", fiber.Map{
            "deleted": deleted,
            "failed":  failed,
        })
    }
    
    return response.SuccessI18n(c, "batch_delete_success", fiber.Map{
        "deleted": deleted,
    })
}
```

## Multilingual API

Complete example with multiple languages:

**locales/en.json:**
```json
{
  "user_created": "User created successfully",
  "user_updated": "User updated successfully",
  "user_deleted": "User deleted successfully",
  "user_retrieved": "User retrieved successfully",
  "users_retrieved": "Users retrieved successfully",
  "user_not_found": "User not found",
  "email_already_exists": "Email {{.Email}} is already registered",
  "invalid_request_body": "Invalid request body",
  "invalid_credentials": "Invalid email or password",
  "login_success": "Login successful",
  "registration_success": "Registration successful",
  "logout_success": "Logged out successfully",
  "unauthorized": "Unauthorized access",
  "invalid_token": "Invalid or expired token",
  "file_uploaded": "File uploaded successfully",
  "file_required": "File is required",
  "file_too_large": "File size {{.Size}} exceeds maximum {{.MaxSize}}",
  "invalid_file_type": "Invalid file type {{.Type}}. Allowed: {{.Allowed}}"
}
```

**locales/id.json:**
```json
{
  "user_created": "Pengguna berhasil dibuat",
  "user_updated": "Pengguna berhasil diperbarui",
  "user_deleted": "Pengguna berhasil dihapus",
  "user_retrieved": "Pengguna berhasil diambil",
  "users_retrieved": "Daftar pengguna berhasil diambil",
  "user_not_found": "Pengguna tidak ditemukan",
  "email_already_exists": "Email {{.Email}} sudah terdaftar",
  "invalid_request_body": "Body request tidak valid",
  "invalid_credentials": "Email atau password salah",
  "login_success": "Login berhasil",
  "registration_success": "Registrasi berhasil",
  "logout_success": "Logout berhasil",
  "unauthorized": "Akses tidak diizinkan",
  "invalid_token": "Token tidak valid atau kadaluarsa",
  "file_uploaded": "File berhasil diunggah",
  "file_required": "File harus diisi",
  "file_too_large": "Ukuran file {{.Size}} melebihi maksimum {{.MaxSize}}",
  "invalid_file_type": "Tipe file {{.Type}} tidak valid. Diizinkan: {{.Allowed}}"
}
```

**Testing with different languages:**

```bash
# English (default)
curl -X GET http://localhost:3000/api/v1/users/999

# Indonesian
curl -X GET http://localhost:3000/api/v1/users/999 \
  -H "Accept-Language: id"

# Chinese
curl -X GET http://localhost:3000/api/v1/users/999 \
  -H "Accept-Language: zh"

# Using query parameter
curl -X GET "http://localhost:3000/api/v1/users/999?lang=id"
```

## Related Documentation

- [README](README.md) - Package overview
- [Standard Responses](standard-responses.md) - Basic response functions
- [I18n Responses](i18n-responses.md) - Internationalized responses
- [Error Handler](error-handler.md) - Custom Fiber error handler
