# Architecture

## Layer Responsibilities

```
HTTP Request
     │
     ▼
┌─────────────┐
│  Middleware  │  JWT validation, CORS, logging
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Controller  │  Bind input → call service/util → return JSON
│             │  No raw SQL. No business logic.
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Service   │  Business logic, query building
│             │  May call initializers.DB directly
└──────┬──────┘
       │
       ▼
┌─────────────┐
│    Utils    │  Stateless pure helpers
│             │  forms.go and model_mapper.go may read DB
└─────────────┘
       │
       ▼
┌─────────────┐
│  GORM / DB  │  PostgreSQL via gorm.io/driver/postgres
└─────────────┘
```

**Rule:** never put DB queries directly in route handlers. Delegate to a service function.

---

## Directory Map

```
cmd/
  main.go                   Server entrypoint — hardcodes :8080
  migrations/main.go        GORM AutoMigrate runner
  seed/main.go              Seeder runner

internal/
  controllers/
    auth/                   Signup, Signin (public routes)
    base/                   BaseController — Index, Show, Update, Delete
    maid/                   Embeds BaseController, overrides Store
    user/                   Embeds BaseController, overrides Store
    tabs/                   GetTabs, AddNewTab, UpdateTab, DeleteTab

  middlewares/
    auth_middleware.go      JWT validation — sets "user" and "userID" in ctx

  models/
    Users.go                User struct
    Maids.go                Maid struct (40+ fields)
    FormFields.go           Drives table columns, forms, search, filters
    UserTabs.go             Per-user tab config
    user_tab_column.go      Per-column visibility, order, width, locked state
    FilterGroup.go          Recursive filter tree (GROUP | FILTER nodes)
    Country.go / Skills.go / Language.go   Reference data

  services/
    filter_service.go       ApplyFilters() — pure query builder, no DB exec
    query_service.go        QueryModelRecords() — executes on passed *gorm.DB
    tab_service.go          BindTabPayload, UpdateTabSettings, LoadFormFields,
                            UpsertTabColumns

  utils/
    responses.go            ErrorResponse struct, ErrorJSON(), ErrorWithCodeJSON()
    search_builder.go       Multi-field ILIKE search with auto-JOIN
    reflection.go           DiscoverRelations() for auto-Preload
    forms.go                ResolveOptionsFromDataSource(), ResolveOptionsForField()
    model_mapper.go         GetModelForDataSource() — string name → model instance
    image_utils.go          ProcessImageUpload() — WebP conversion + BlurHash + upload
    tab_defaults.go         CreateDefaultTabsForUserModel()

  initializers/
    ConnectDataBase.go      initializers.DB (*gorm.DB global) with connection pool
    LoadEnvVars.go          Loads .env via godotenv
    storage.go              initializers.Storage (storage.Driver global)

  storage/
    driver.go               Driver interface
    local.go                Local disk implementation
    s3.go                   AWS S3 / Cloudflare R2 implementation

  constants/
    contants.go             FormFieldType constants
    module.go               Model name constants (MAID, USER, etc.)
    options.go              Predefined select options (gender, visa, cities…)

  routes/routes.go          All route registrations

docs/
  swagger.json / swagger.yaml / docs.go    Auto-generated — DO NOT edit manually
  *.md                      This documentation
```

---

## Route Map

| Method | Path | Controller | Auth |
|--------|------|-----------|------|
| POST | `/api/auth/signup` | auth.Signup | Public |
| POST | `/api/auth/signin` | auth.Signin | Public |
| POST | `/api/users/index` | base.Index | JWT |
| POST | `/api/users/` | user.Store | JWT |
| GET | `/api/users/:id` | base.Show | JWT |
| PUT | `/api/users/:id` | base.Update | JWT |
| DELETE | `/api/users/:id` | base.Delete | JWT |
| POST | `/api/maids/index` | base.Index | JWT |
| POST | `/api/maids/` | maid.Store | JWT |
| GET | `/api/maids/:id` | maid.Show | JWT |
| PUT | `/api/maids/:id` | maid.Update | JWT |
| DELETE | `/api/maids/:id` | maid.Delete | JWT |
| GET | `/api/tabs` | tabs.GetTabs | JWT |
| POST | `/api/tabs` | tabs.AddNewTab | JWT |
| PUT | `/api/tabs/:id` | tabs.UpdateTab | JWT |
| DELETE | `/api/tabs/:id` | tabs.DeleteTab | JWT |
| GET | `/ping` | inline | Public |
| GET | `/swagger/*any` | SwaggerUI | Public |
| GET | `/storage/*` | static files | Public (local driver only) |

---

## Global Singletons

| Variable | Package | Type | Description |
|----------|---------|------|-------------|
| `initializers.DB` | initializers | `*gorm.DB` | PostgreSQL connection (pool configured) |
| `initializers.Storage` | initializers | `storage.Driver` | Active storage backend |

---

## Migration Approach

This project uses **GORM AutoMigrate** — not SQL migration files.

To add a new table or column:
1. Add or modify the struct in `internal/models/`
2. Register the model in `cmd/migrations/main.go`
3. Run: `go run ./cmd/migrations/main.go`

Never alter a model struct without running AutoMigrate. Never create raw SQL migration files.

---

## Known Issues

| # | Issue | Location |
|---|-------|---------|
| 1 | `UpdateTab` and `DeleteTab` are implemented but **not registered** in routes | `routes/routes.go` |
| 2 | `BaseController.Store()` is empty — sub-controllers override it | `base/base_controller.go` |
| 3 | `PORT` env var is declared but the server hardcodes `:8080` | `cmd/main.go` |
| 4 | No test files (`*_test.go`) anywhere in the project | — |
| 5 | `models.Patient` referenced in old Swagger docs — model does not exist | `docs/` |
