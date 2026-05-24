# Skill: Error Handling

All error responses in this project must go through `internal/utils/responses.go`. Never use `gin.H{"error": "..."}` directly in new code.

---

## The ErrorResponse Struct

```go
// internal/utils/responses.go
type ErrorResponse struct {
    Code    string `json:"code,omitempty"`    // machine-readable, for frontend switch
    Error   string `json:"error"`             // human-readable message
    Details string `json:"details,omitempty"` // debug/context info
}
```

---

## Helper Functions

### `utils.ErrorJSON` — standard error

```go
utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request payload", err.Error())
```

Sends:
```json
{ "error": "Invalid request payload", "details": "json: cannot unmarshal..." }
```

The `details` variadic argument is optional. Include it when it helps the frontend or debugging.

### `utils.ErrorWithCodeJSON` — error with machine-readable code

```go
utils.ErrorWithCodeJSON(ctx, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Email is required", "field: email")
```

Sends:
```json
{ "code": "VALIDATION_ERROR", "error": "Email is required", "details": "field: email" }
```

Use this when the frontend needs to branch on a specific error type.

---

## When to Use Which HTTP Status Code

| Situation | Code | Helper |
|-----------|------|--------|
| Missing or malformed request body | 400 | ErrorJSON |
| Binding/validation failed (`ShouldBindJSON` error) | 400 | ErrorJSON |
| Unauthenticated (no/invalid token) | 401 | ErrorJSON (middleware handles this) |
| Forbidden (authenticated but not allowed) | 403 | ErrorJSON |
| Resource not found | 404 | ErrorJSON |
| Business rule violation (duplicate email, etc.) | 422 | ErrorWithCodeJSON |
| Database or server error | 500 | ErrorJSON |

---

## Pattern for Each Handler Type

### Create (Store)
```go
func (c *MyController) Store(ctx *gin.Context) {
    var input MyInput
    if err := ctx.ShouldBindJSON(&input); err != nil {
        utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request", err.Error())
        return
    }
    record := models.MyModel{Field: input.Field}
    if err := initializers.DB.Create(&record).Error; err != nil {
        utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to create record", err.Error())
        return
    }
    ctx.JSON(http.StatusCreated, record)
}
```

### Read (Show)
```go
func (c *MyController) Show(ctx *gin.Context) {
    id := ctx.Param("id")
    var record models.MyModel
    if err := initializers.DB.First(&record, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            utils.ErrorJSON(ctx, http.StatusNotFound, "Record not found")
            return
        }
        utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to fetch record", err.Error())
        return
    }
    ctx.JSON(http.StatusOK, record)
}
```

### Update
```go
func (c *MyController) Update(ctx *gin.Context) {
    id := ctx.Param("id")
    var record models.MyModel
    if err := initializers.DB.First(&record, id).Error; err != nil {
        utils.ErrorJSON(ctx, http.StatusNotFound, "Record not found")
        return
    }
    var input MyInput
    if err := ctx.ShouldBindJSON(&input); err != nil {
        utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request", err.Error())
        return
    }
    record.Field = input.Field
    if err := initializers.DB.Save(&record).Error; err != nil {
        utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to update record", err.Error())
        return
    }
    ctx.JSON(http.StatusOK, record)
}
```

### Delete
```go
func (c *MyController) Delete(ctx *gin.Context) {
    id := ctx.Param("id")
    var record models.MyModel
    if err := initializers.DB.First(&record, id).Error; err != nil {
        utils.ErrorJSON(ctx, http.StatusNotFound, "Record not found")
        return
    }
    if err := initializers.DB.Delete(&record).Error; err != nil {
        utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to delete record", err.Error())
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
```

---

## How the JWT Middleware Returns Errors

The middleware (`internal/middlewares/auth_middleware.go`) uses `gin.H` directly — this is existing code, do not change:
```go
c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
```

This is an inconsistency from before the `utils.ErrorResponse` helper existed. New middleware should use `utils.ErrorJSON`.

---

## What NOT to Do

```go
// ❌ Do not use gin.H for errors in new code
ctx.JSON(http.StatusBadRequest, gin.H{"error": "something went wrong"})

// ❌ Do not swallow errors silently
initializers.DB.Create(&record)  // missing .Error check

// ❌ Do not return 200 for errors
ctx.JSON(http.StatusOK, gin.H{"success": false, "error": "..."})
```
