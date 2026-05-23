# Tab & Filter System

This is the most sophisticated part of the API. It powers the Airtable/Notion-style data table in the frontend — per-user saved views with custom column configurations and recursive filters.

---

## Overview

```
User logs in
     │
     ▼
GET /api/tabs?model=Maid
     │
     ├── Returns FormField metadata (what columns exist, their types, options)
     └── Returns UserTab list (user's saved views)
             │
             ├── If no tabs exist → creates a default "All Maids" tab automatically
             └── Each tab has columns (visibility, order, width, locked state)

User changes filters / columns / search
     │
     ▼
POST /api/maids/index  (TabPayload)
     │
     ├── Persists filters + search term to UserTab
     ├── Persists column settings to UserTabColumn
     └── Executes query with filters + search + pagination
             │
             └── Returns { data: [...], meta: { totalRowCount: N } }
```

---

## FormField System

`FormField` records are seeded in the database and drive everything:

| Column | Purpose |
|--------|---------|
| `field_key` | snake_case — matches the DB column name |
| `model_name` | `"Maid"` or `"User"` |
| `form_field_type` | Input type (see types below) |
| `label` | Human-readable display name |
| `data_source` | `"Model:IDField:LabelField"` — e.g. `"Country:ID:Name"` |
| `options` | Static `[{ value, label }]` JSON — used when no data_source |
| `table_is_visible` | Default column visibility in table view |
| `table_order` | Default column position |
| `table_is_pinned` | Whether column is locked/pinned by default |
| `form_order` | Position in the create/edit form |
| `form_stage` | Multi-step form stage number |
| `form_width` | Column width in pixels |
| `form_is_required` | Whether the form field is required |

### FormField Types

| Type | Description |
|------|-------------|
| `string_field` | Plain text input |
| `number_field` | Numeric input |
| `date_field` | Date picker |
| `boolean_field` | Toggle / checkbox |
| `select_field` | Single select (static options) |
| `multi_select_field` | Multi select (static options) |
| `single_relation` | Belongs-to relationship — stored as `{field_key}_id` |
| `multi_relation` | Many-to-many relationship |
| `image_field` | File upload |

---

## Tab Structure

### UserTab

| Field | Description |
|-------|-------------|
| `id` | Primary key |
| `user_id` | Owner (composite indexed with `model_name`) |
| `model_name` | `"Maid"` or `"User"` |
| `tab_name` | Display name |
| `is_default` | The "All {Model}s" tab created automatically |
| `search_term` | Persisted search string |
| `filters` | Persisted FilterGroup JSON |
| `columns` | `[]UserTabColumn` |

### UserTabColumn

| Field | Description |
|-------|-------------|
| `user_tab_id` | Parent tab (composite unique indexed with `form_field_id`) |
| `form_field_id` | Which FormField this column represents |
| `field_key` | Denormalised for fast lookup |
| `visible` | Whether column is shown |
| `locked` | Whether column is pinned |
| `order` | Column position |
| `width` | Column width in pixels |

---

## Filter System

Filters are a **recursive tree** stored as JSON in `UserTab.Filters` and sent in every `POST /api/{model}/index` request.

### FilterGroup Shape

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
    },
    {
      "id": "group-2",
      "type": "GROUP",
      "conjunction": "or",
      "children": [ ... ]
    }
  ]
}
```

### How Filters Are Applied

```
ApplyFilters(db, rawJSON, fieldsByID)
     │
     └── applyGroupFilters(db, group, fields)
              │
              ├── conjunction = "and"
              │   └── applies each child sequentially with .Where()
              │
              └── conjunction = "or"
                  └── builds each child in an isolated session
                      then combines with .Or()
```

`fieldsByID` is a `map[string]FormField` keyed by `string(field.ID)`, loaded from `services.LoadFormFields()`.

### Column Name Resolution

| FormFieldType | Resolved column |
|--------------|----------------|
| `single_relation` | `{field_key}_id` (e.g. `nationality_id`) |
| everything else | `{field_key}` as-is |

### Supported Operators

| Operator value(s) | SQL generated |
|------------------|---------------|
| `is`, `=`, `==`, `isExactly` | `col = ?` |
| `isNot`, `!=` | `col != ?` |
| `contains` | `col ILIKE '%val%'` |
| `doesNotContain` | `col NOT ILIKE '%val%'` |
| `isEmpty` | `col IS NULL` |
| `isNotEmpty` | `col IS NOT NULL` |
| `<`, `isBefore` | `col < ?` |
| `<=`, `isOnOrBefore` | `col <= ?` |
| `>`, `isAfter` | `col > ?` |
| `>=`, `isOnOrAfter` | `col >= ?` |
| `isAnyOf`, `hasAnyOf` | `col IN (?)` |
| `hasNoneOf`, `isNoneOf` | `col NOT IN (?)` |
| `hasAllOf` | `col @> ?` (PostgreSQL array) |
| `isWithin` | date range (see below) |

### Date Range Operators (used with `isWithin`)

Set in `secondOperator.value`:

| Value | Range |
|-------|-------|
| `today` | Current day 00:00 – 23:59 |
| `yesterday` / `tomorrow` | Previous / next day |
| `oneWeekAgo` / `oneWeekFromNow` | ±7 days |
| `oneMonthAgo` / `oneMonthFromNow` | ±1 month |
| `numberOfDaysAgo` / `numberOfDaysFromNow` | ±N days (value = number) |
| `exactDate` | Specific date (value = `"YYYY-MM-DD"`) |
| `thePastWeek` / `thePastMonth` / `thePastYear` | Rolling past period |
| `theNextWeek` / `theNextMonth` / `theNextYear` | Rolling future period |
| `thisCalendarWeek` / `thisCalendarMonth` / `thisCalendarYear` | Current calendar period |

---

## TabPayload (sent on every index request)

```json
{
  "tab_id": 1,
  "page": 0,
  "size": 50,
  "search": "optional search term",
  "filters": {
    "id": "root",
    "type": "GROUP",
    "conjunction": "and",
    "children": []
  },
  "columns": [
    {
      "field_key": "first_name",
      "visible": true,
      "locked": false,
      "order": 1,
      "width": 200
    }
  ]
}
```

All column and filter state is **round-tripped through the server** — the frontend sends its current state, the server persists it, then executes the query using that state. This means the user's view is restored correctly on any device.

---

## Default Tab Creation

When `GET /api/tabs?model=Maid` is called and the user has no tabs for that model yet, the server automatically:

1. Loads all `FormField` records for the model
2. Creates a `UserTab` with `tab_name = "All Maids"`, `is_default = true`
3. Creates a `UserTabColumn` for each FormField using its default visibility, pinned state, and width
4. Returns the newly created tab in the response

This happens in a single database transaction (`internal/utils/tab_defaults.go`).
