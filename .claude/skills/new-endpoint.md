# Skill: Add a New Endpoint

Use this checklist every time a new endpoint is added to the Bridge API. Follow in order — skipping steps causes inconsistency.

---

## Step 1 — Model (if new fields or new table needed)

- Add/modify the struct in `internal/models/`
- Use correct GORM tags: `gorm:"primarykey"`, `gorm:"uniqueIndex"`, `gorm:"not null"`, etc.
- Add to the AutoMigrate call in `cmd/migrations/main.go`
- Run migrations: `go run ./cmd/migrations/main.go`
- If new model needs FormFields, add a seeder in `internal/database/seeder/` and run: `go run ./cmd/seed/main.go`

## Step 2 — Service / Util (if business logic needed)

- Add logic to an existing service file or create a new one under `internal/services/`
- Services may call `initializers.DB` directly (current project pattern)
- Keep service functions pure where possible: accept `*gorm.DB` as parameter when the function is called inside a transaction
- Never put DB calls directly in the controller

## Step 3 — Controller

**If the new endpoint follows the standard CRUD pattern on an existing model:**
- The controller already embeds `base.BaseController` — `Index`, `Show`, `Delete` are inherited
- Override only what differs from the base (e.g. `Store`, `Update`)

**If it's a new model:**
```go
type MyController struct {
    base.BaseController
}

func NewMyController() *MyController {
    return &MyController{
        BaseController: base.BaseController{Model: &models.MyModel{}},
    }
}
```

**Input binding pattern:**
```go
var input MyInput
if err := ctx.ShouldBindJSON(&input); err != nil {
    utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request", err.Error())
    return
}
```

**Success response:**
```go
ctx.JSON(http.StatusCreated, result)         // 201 for create
ctx.JSON(http.StatusOK, result)              // 200 for read/update/delete
ctx.JSON(http.StatusOK, gin.H{"message": "..."})  // for delete confirmation
```

**Error responses — always use the helpers:**
```go
utils.ErrorJSON(ctx, http.StatusNotFound, "Resource not found")
utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to save", err.Error())
utils.ErrorWithCodeJSON(ctx, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "message", "detail")
```

## Step 4 — Route Registration

Open `internal/routes/routes.go`.

Public route:
```go
authRoutes.POST("/my-route", myCtrl.MyHandler)
```

Protected route (JWT required):
```go
myGroup := protected.Group("/my-resource")
{
    myGroup.GET("/:id", myCtrl.Show)
    myGroup.POST("/", myCtrl.Store)
    myGroup.PUT("/:id", myCtrl.Update)
    myGroup.DELETE("/:id", myCtrl.Delete)
    myGroup.POST("/index", myCtrl.Index)  // list with filters
}
```

## Step 5 — Swagger Annotations

Add above **every** handler function. Do not skip.

```go
// Store godoc
// @Summary     Create a new resource
// @Description Creates a resource with the provided data
// @Tags        resource-name
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body body    MyInput true "Request body"
// @Success     201  {object} models.MyModel
// @Failure     400  {object} utils.ErrorResponse
// @Failure     401  {object} utils.ErrorResponse
// @Failure     500  {object} utils.ErrorResponse
// @Router      /resource [post]
func (c *MyController) Store(ctx *gin.Context) {
```

Tags to use: `maids`, `users`, `tabs`, `auth` — or add a new tag for a new resource group.

## Step 6 — Regenerate Swagger Docs

```bash
swag init -g cmd/main.go
```

Run from the project root. This updates `docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`.

Verify the new endpoint appears at `http://localhost:8080/swagger/index.html`.

## Step 7 — Smoke Test

```bash
# Health check
curl http://localhost:8080/ping

# Auth (get token)
curl -X POST http://localhost:8080/api/auth/signin \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'

# Test new endpoint
curl -X POST http://localhost:8080/api/my-resource/ \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"field":"value"}'
```

---

## Checklist Summary

- [ ] Model added/updated in `internal/models/`
- [ ] AutoMigrate updated and run
- [ ] FormField seeder added if needed
- [ ] Service function added if business logic required
- [ ] Controller created or extended
- [ ] Route registered in `internal/routes/routes.go` (correct group — public vs protected)
- [ ] Swagger annotations on every handler
- [ ] `swag init` run and docs regenerated
- [ ] Smoke tested with curl or Postman
