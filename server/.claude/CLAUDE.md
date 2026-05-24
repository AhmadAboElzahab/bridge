# Bridge API — Claude Project Reference

## Overview
UAE-based domestic service marketplace. The API manages **Maid** profiles and **User** accounts with advanced filtering, per-user customizable table views (tabs), and dynamic form field metadata. The frontend is a separate Next.js app on `http://localhost:3000`.

Module: `github.com/AhmadAboElzahab/bridge`
Go: 1.24 | Gin | GORM | PostgreSQL | JWT (HS256) | Swagger (swaggo)

---

## Directory Map
```
cmd/
  main.go              — server entrypoint, hardcoded :8080
  migrations/main.go   — GORM AutoMigrate runner
  seed/main.go         — seeder runner
internal/
  controllers/
    auth/              — Signup, Signin (public)
    base/              — BaseController (Index, Show, Delete; Update is empty)
    maid/              — embeds BaseController, overrides Store
    user/              — embeds BaseController, overrides Store
    tabs/              — GetTabs, AddNewTab (UpdateTab/DeleteTab exist but are NOT registered)
  middlewares/
    auth_middleware.go — JWT validation, sets "user" and "userID" in ctx
  models/              — GORM structs (see Model Map below)
  services/
    filter_service.go  — ApplyFilters() — pure query builder, no DB exec
    query_service.go   — QueryModelRecords() — executes on passed *gorm.DB
    tab_service.go     — BindTabPayload, UpdateTabSettings, LoadFormFields, UpsertTabColumns
  utils/
    responses.go       — ErrorResponse struct, ErrorJSON(), ErrorWithCodeJSON()
    search_builder.go  — multi-field ILIKE search with auto-JOIN
    reflection.go      — DiscoverRelations() for auto-Preload
    forms.go           — ResolveOptionsFromDataSource(), ResolveOptionsForField()
    model_mapper.go    — GetModelForDataSource() — string name → model instance
    image_utils.go     — ProcessImageUpload() — WebP + blurhash
    tab_defaults.go    — CreateDefaultTabsForUserModel()
  initializers/
    ConnectDataBase.go — initializers.DB (*gorm.DB global)
    LoadEnvVars.go     — loads .env via godotenv
  constants/
    contants.go        — FormFieldType constants
    module.go          — model name constants (MAID, USER, etc.)
    options.go         — predefined select options (gender, visa status, cities, etc.)
  routes/routes.go     — all route registrations
docs/                  — auto-generated Swagger files (DO NOT edit manually)
```

---

## Route Map
| Method | Path | Controller | Auth | Swagger |
|--------|------|-----------|------|---------|
| POST | /api/auth/signup | auth.Signup | Public | ✅ |
| POST | /api/auth/signin | auth.Signin | Public | ✅ |
| POST | /api/users/index | base.Index | JWT | ❌ |
| POST | /api/users/ | user.Store | JWT | ❌ |
| GET | /api/users/:id | base.Show | JWT | ❌ |
| PUT | /api/users/:id | base.Update | JWT | ❌ |
| DELETE | /api/users/:id | base.Delete | JWT | ❌ |
| POST | /api/maids/index | base.Index | JWT | ✅ |
| POST | /api/maids/ | maid.Store | JWT | ✅ |
| GET | /api/maids/:id | maid.Show | JWT | ✅ |
| PUT | /api/maids/:id | maid.Update | JWT | ✅ |
| DELETE | /api/maids/:id | maid.Delete | JWT | ✅ |
| GET | /api/tabs | tabs.GetTabs | JWT | ❌ |
| POST | /api/tabs | tabs.AddNewTab | JWT | ❌ |
| GET | /ping | inline | Public | — |
| GET | /swagger/*any | SwaggerUI | Public | — |

---

## Layer Responsibilities

**Controller** — thin. Bind input, call service/util, return JSON. No raw SQL, no business logic.

**Service** — business logic and query building. May call `initializers.DB` directly (current pattern — no dedicated repo layer).

**Utils** — stateless helpers. Pure functions. No DB access except `forms.go` and `model_mapper.go` which read reference tables.

**Never** put DB queries directly in route handlers. Always delegate to a service function.

---

## Model Map

### User (`models/Users.go`)
```go
type User struct {
    ID          uint           `gorm:"primarykey"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
    FirstName   string         `json:"first_name"`
    LastName    string         `json:"last_name"`
    Email       string         `json:"email"`
    Password    string         `json:"-"`
    DateOfBirth string         `json:"date_of_birth"`
    Avatar      string         `json:"avatar"`
    Blurhash    string         `json:"blurhash"`
}
```

### Maid (`models/Maids.go`)
Key fields (40+ total): `Email` (unique), `FirstName`, `LastName`, `NationalityID` → `Nationality Country`, `Status`, `VisaStatus`, `ProfileCompletionPercentage`, `SubscriptionStatus`, `IsProfileVisible`.
Relations: `Languages []Language` (many2many: maid_languages), `Skills []Skill` (many2many: maid_skills).

### FormField (`models/FormFields.go`)
Drives everything — table columns, form layout, search, filters.
```go
FieldKey       string         // snake_case, matches DB column
ModelName      string         // "Maid", "User"
FormFieldType  string         // see constants/contants.go
DataSource     string         // "Model:IDField:LabelField" e.g. "Country:ID:Name"
Options        datatypes.JSON // static [{value, label}] — used if no DataSource
TableIsVisible bool
TableOrder     int
FormOrder      int            // gorm column:"field_order"
```

### UserTab / UserTabColumn (`models/UserTabs.go`, `models/user_tab_column.go`)
Per-user view config. `UserTab.Filters` is `datatypes.JSON` storing a `FilterGroup`. `UserTabColumn` stores per-column visibility, order, width, locked state.

### FilterGroup / FilterItem (`models/FilterGroup.go`)
Recursive filter tree. `FilterGroup.Conjunction` = "and" | "or". `FilterItem.Operator.Value` = operator string (see Filter Pattern section).

---

## Error Response Format

**Always** use these helpers from `internal/utils/responses.go`:

```go
// Simple error
utils.ErrorJSON(ctx, http.StatusBadRequest, "human message", "optional detail")

// With machine-readable code
utils.ErrorWithCodeJSON(ctx, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "human message", "optional detail")
```

Response shape:
```json
{
  "code": "VALIDATION_ERROR",   // omitempty — only when ErrorWithCodeJSON used
  "error": "human message",
  "details": "optional detail"  // omitempty
}
```

**Never** use `ctx.JSON(code, gin.H{"error": "..."})` directly — use the helpers.

---

## JWT Middleware

Middleware: `middlewares.AuthMiddleware()`
Application: wrap any protected route group with `group.Use(middlewares.AuthMiddleware())`

```go
protected := api.Group("/")
protected.Use(middlewares.AuthMiddleware())
```

Token format: `Authorization: Bearer <token>`
Claims: `UserID uint` + `jwt.StandardClaims`
Context keys: `"user"` (models.User), `"userID"` (uint)

Get user in handler:
```go
user := ctx.MustGet("user").(models.User)
userID := ctx.MustGet("userID").(uint)
```

**Known issue**: `ExpiresAt` is commented out in `generateJWT()` — tokens never expire.

---

## Swagger Annotation Template

Required on every handler. Run `swag init` from project root after any change.

```go
// MethodName godoc
// @Summary     Short action description
// @Description Longer description
// @Tags        resource-name
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id   path    int          true  "Resource ID"
// @Param       body body    InputStruct  true  "Request body"
// @Success     200  {object} models.ModelName
// @Failure     400  {object} utils.ErrorResponse
// @Failure     401  {object} utils.ErrorResponse
// @Failure     404  {object} utils.ErrorResponse
// @Failure     500  {object} utils.ErrorResponse
// @Router      /resource/{id} [get]
```

The `@Security BearerAuth` tag is required on all JWT-protected routes.

---

## Filter Pattern

Filters arrive as JSON (stored in `UserTab.Filters` and sent in `TabPayload.Filters`).

### FilterGroup shape
```json
{
  "id": "group-1",
  "type": "GROUP",
  "conjunction": "and",
  "children": [
    {
      "id": "item-1",
      "type": "FILTER",
      "fieldId": 42,
      "columnType": "string_field",
      "operator": { "label": "contains", "value": "contains" },
      "value": "Ahmad"
    }
  ]
}
```

### How filters are applied
`services.ApplyFilters(db, raw, fieldsByID)` → `applyGroupFilters` → `applyItemFilter` → `buildCondition`

`fieldsByID` is a `map[string]models.FormField` keyed by `string(field.ID)`. It comes from `services.LoadFormFields()`.

### Column name resolution (`resolveColumnName`)
- `single_relation` fields: `field.FieldKey + "_id"` (e.g. `nationality_id`)
- all others: `field.FieldKey` as-is

### Supported operators
| Operator values | SQL |
|----------------|-----|
| is, =, ==, isExactly | `col = ?` |
| isNot, != | `col != ?` |
| contains | `col ILIKE '%val%'` |
| doesNotContain | `col NOT ILIKE '%val%'` |
| isEmpty | `col IS NULL` |
| isNotEmpty | `col IS NOT NULL` |
| <, isBefore | `col < ?` |
| <=, isOnOrBefore | `col <= ?` |
| >, isAfter | `col > ?` |
| >=, isOnOrAfter | `col >= ?` |
| isAnyOf, hasAnyOf | `col IN (?)` |
| hasNoneOf, isNoneOf | `col NOT IN (?)` |
| hasAllOf | `col @> ?` (PostgreSQL array) |
| isWithin | date range via `secondOperator.Value` |

### Date range secondOperator values
`today`, `tomorrow`, `yesterday`, `oneWeekAgo`, `oneWeekFromNow`, `oneMonthAgo`, `oneMonthFromNow`, `numberOfDaysAgo`, `numberOfDaysFromNow`, `exactDate`, `thePastWeek`, `thePastMonth`, `thePastYear`, `theNextWeek`, `theNextMonth`, `theNextYear`, `thisCalendarWeek`, `thisCalendarMonth`, `thisCalendarYear`

---

## Migration Approach

This project uses **GORM AutoMigrate**, not SQL migration files.

To add a new table or column:
1. Add/modify the struct in `internal/models/`
2. Add the model to the AutoMigrate call in `cmd/migrations/main.go`
3. Run: `go run ./cmd/migrations/main.go`

**Never** alter a model struct without running AutoMigrate. **Never** create raw SQL migration files.

---

## FormField System

FormFields are seeded records (not hardcoded) that define:
- What columns appear in a model's table view (`TableIsVisible`, `TableOrder`, `TableIsPinned`)
- What fields appear in the creation/edit form (`FormOrder`, `FormStage`, `FormWidth`, `FormIsRequired`)
- What type of input to render (`FormFieldType` — see `constants/contants.go`)
- Where to load options from (`DataSource` = "Model:IDField:LabelField", or static `Options` JSON)
- Which fields are searchable (SearchBuilder checks FormFieldType)

Seeded in: `internal/database/seeder/maids_form_fields.go`, `user_form_fields.go`

---

## Tab System

Each user has `UserTab` records per model. Each tab stores:
- `Filters` (JSON) — current filter configuration
- `SearchTerm` — persisted search string
- `Columns []UserTabColumn` — per-column visibility/order/width/locked

The `POST /api/{model}/index` endpoint:
1. Binds `TabPayload` (TabID, Filters, Columns, Page, Size, Search)
2. Updates tab settings and column config in a transaction
3. Runs `QueryModelRecords` which applies filters + search + pagination
4. Returns `{ data: [...], meta: { totalRowCount: N } }`

---

## Known Issues (do not reproduce these patterns)

1. **JWT tokens never expire** — `ExpiresAt` commented out in `auth_controller.go:generateJWT()`
2. **`BaseController.Update()` is empty** — PUT endpoints do nothing
3. **`UpdateTab` and `DeleteTab` implemented but not registered** in `routes/routes.go`
4. **No tests** — zero `*_test.go` files in the project
5. **Missing Swagger** on all user routes and tab routes
6. **Port mismatch** — `.env` has `PORT=9000` but server hardcodes `:8080`
7. **`models.Patient`** referenced in Swagger docs — does not exist (dead reference)
8. **`BaseController.Show` and `Delete` use `gin.H` directly** instead of `utils.ErrorJSON`
9. **OR filter conjunction skips nested groups** — `applyGroupFilters` has a `continue` for OR+subgroup case
