# API Reference

The live interactive docs are available via Swagger UI when the server is running:

```
http://localhost:8080/swagger/index.html
```

This document is a human-readable reference for all endpoints.

---

## Authentication

All protected routes require a JWT token in the `Authorization` header:

```
Authorization: Bearer <token>
```

Tokens are obtained from `POST /api/auth/signin`. They expire after **24 hours**.

---

## Auth Endpoints

### POST /api/auth/signup

Register a new user account.

**Content-Type:** `multipart/form-data`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `first_name` | string | Yes | |
| `last_name` | string | No | |
| `email` | string (email) | Yes | Must be unique |
| `password` | string | Yes | Will be bcrypt-hashed |
| `date_of_birth` | string (date) | No | |
| `avatar` | file | No | JPG, PNG, or WebP — converted to WebP on upload |

**Responses:**

| Code | Description |
|------|-------------|
| 201 | User created — returns the User object |
| 400 | Validation error |
| 409 | Email already exists (`code: EMAIL_EXISTS`) |
| 500 | Server error |

---

### POST /api/auth/signin

Authenticate and receive a JWT token.

**Content-Type:** `application/json`

```json
{
  "email": "user@example.com",
  "password": "secret"
}
```

**Responses:**

| Code | Body | Description |
|------|------|-------------|
| 200 | `{ "message": "Login successful", "token": "eyJ..." }` | Success |
| 400 | ErrorResponse | Invalid request body |
| 401 | ErrorResponse | Wrong email or password |
| 500 | ErrorResponse | Server error |

---

## Maid Endpoints (all require JWT)

### POST /api/maids/index

Query maids with filtering, search, and pagination. Also persists the tab state server-side.

**Content-Type:** `application/json`

```json
{
  "tab_id": 1,
  "page": 0,
  "size": 50,
  "search": "Ahmad",
  "filters": { "id": "root", "type": "GROUP", "conjunction": "and", "children": [] },
  "columns": [
    { "field_key": "first_name", "visible": true, "order": 1, "width": 200 }
  ]
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `tab_id` | uint | Yes | ID of the active UserTab |
| `page` | int | No | 0-based page index (default: 0) |
| `size` | int | No | Rows per page, 1–100 (default: 50) |
| `search` | string | No | Global search across searchable fields |
| `filters` | FilterGroup JSON | Yes | Recursive filter tree (can be `{}` for no filters) |
| `columns` | array | Yes | Column visibility/order settings to persist |

**Response:**
```json
{
  "data": [ { "id": 1, "first_name": "Fatima", ... } ],
  "meta": { "totalRowCount": 142 }
}
```

---

### POST /api/maids

Create a new maid profile.

```json
{
  "first_name": "Fatima",
  "email": "fatima@example.com"
}
```

| Code | Description |
|------|-------------|
| 201 | Maid created |
| 400 | Validation error |

---

### GET /api/maids/:id

Fetch a single maid by ID.

| Code | Description |
|------|-------------|
| 200 | Returns Maid object |
| 404 | Not found |

---

### PUT /api/maids/:id

Update maid fields. Accepts any valid maid field except `id`, `created_at`, `updated_at`, `deleted_at`, `password`.

```json
{
  "first_name": "Nour",
  "visa_status": "valid"
}
```

| Code | Description |
|------|-------------|
| 200 | Returns updated Maid |
| 400 | Invalid body |
| 404 | Not found |

---

### DELETE /api/maids/:id

Soft-delete a maid (sets `deleted_at`, record is hidden but not removed).

| Code | Description |
|------|-------------|
| 200 | `{ "message": "Resource deleted successfully" }` |
| 404 | Not found |

---

## User Endpoints (all require JWT)

Same shape as maid endpoints — replace `/api/maids` with `/api/users`.

- `POST /api/users/index` — list with filters
- `POST /api/users/` — create
- `GET /api/users/:id` — get by ID
- `PUT /api/users/:id` — update
- `DELETE /api/users/:id` — soft delete

---

## Tab Endpoints (all require JWT)

### GET /api/tabs?model=Maid

Returns form field metadata and the user's saved tab configurations for the given model.

**Query params:**

| Param | Required | Values | Description |
|-------|----------|--------|-------------|
| `model` | Yes | `Maid`, `User` | Which model's tabs to load |

**Response:**
```json
{
  "form_fields": [
    {
      "id": 1,
      "label": "First Name",
      "field_key": "first_name",
      "form_field_type": "string_field",
      "table_is_visible": true,
      "table_order": 1,
      "options": []
    }
  ],
  "tabs": [
    {
      "id": 1,
      "tab_name": "All Maids",
      "model_name": "Maid",
      "is_default": true,
      "search_term": "",
      "filters": {},
      "columns": [
        { "field_key": "first_name", "label": "First Name", "visible": true, "order": 1, "width": 200 }
      ]
    }
  ]
}
```

If no tabs exist for the user + model combination, a default "All {Model}s" tab is created automatically.

---

### POST /api/tabs

Create a new tab.

```json
{
  "model_name": "Maid",
  "tab_name": "Active Maids"
}
```

| Code | Description |
|------|-------------|
| 201 | `{ "message": "New tab created", "tab_id": 4 }` |
| 400 | Invalid payload or model name |

---

### PUT /api/tabs/:id

Rename a tab.

```json
{
  "tab_name": "My Custom View"
}
```

---

### DELETE /api/tabs/:id

Delete a tab and all its column configurations.

---

## Error Response Format

All errors use this consistent shape:

```json
{
  "error": "Human readable message",
  "details": "Optional detail string",
  "code": "MACHINE_READABLE_CODE"
}
```

`details` and `code` are omitted when empty.

**Common error codes:**

| Code | HTTP | Meaning |
|------|------|---------|
| `EMAIL_EXISTS` | 409 | Email already registered |
| `VALIDATION_ERROR` | 422 | Input failed validation |

---

## Pagination

Used on all `index` endpoints:

| Field | Type | Default | Max | Description |
|-------|------|---------|-----|-------------|
| `page` | int | 0 | — | 0-based page number |
| `size` | int | 50 | 100 | Rows per page |

Response always includes `meta.totalRowCount` for frontend pagination controls.
