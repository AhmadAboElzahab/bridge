# Skill: Add a Filter

This project uses a hierarchical filter system. Filters are sent in the `POST /api/{model}/index` body as nested `FilterGroup` and `FilterItem` nodes. Filters are processed by `services/filter_service.go`.

---

## How Filters Work End-to-End

```
Frontend sends TabPayload.Filters (JSON)
  → services.BindTabPayload() binds it as json.RawMessage
  → services.UpdateTabSettings() persists it to UserTab.Filters
  → services.ApplyFilters(db, raw, fieldsByID) builds WHERE clauses
      → applyGroupFilters() handles AND/OR recursion
          → applyItemFilter() resolves field and calls buildCondition()
              → buildCondition() returns (sqlString, value)
```

`fieldsByID` is `map[string]models.FormField` keyed by `string(field.ID)`.
It comes from `services.LoadFormFields(tx, modelName)`.

---

## FilterGroup JSON Shape

```json
{
  "id": "root",
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
      "children": [
        {
          "id": "item-2",
          "type": "FILTER",
          "fieldId": 15,
          "columnType": "single_select",
          "operator": { "label": "is", "value": "is" },
          "value": "Available"
        }
      ]
    }
  ]
}
```

`fieldId` must match `FormField.ID` in the database.

---

## Adding a New Filterable Field

To make a field filterable, it needs a `FormField` record. The field must already exist as a DB column on the model.

**Step 1** — Ensure FormField record exists in the database with the correct:
- `field_key` — must match the DB column name (or `column_name + "_id"` for single_relation)
- `model_name` — "Maid", "User", etc.
- `form_field_type` — determines operator behavior and column resolution

**Step 2** — The field is automatically available for filtering because `LoadFormFields` fetches all FormFields for the model and `buildCondition` handles all standard operators.

No code changes needed in `filter_service.go` for standard operators.

---

## Column Name Resolution (`resolveColumnName`)

Location: `services/filter_service.go:resolveColumnName()`

```go
func resolveColumnName(field models.FormField) string {
    switch field.FormFieldType {
    case "single_relation":
        return field.FieldKey + "_id"  // e.g. "nationality" → "nationality_id"
    default:
        return field.FieldKey           // e.g. "first_name"
    }
}
```

If filtering a `single_relation` field, the filter compares against the foreign key column, not the relation itself.

---

## Adding a New Operator

To add a new operator to `buildCondition()` in `services/filter_service.go`:

```go
func buildCondition(column, operator string, value interface{}, secondOperator *models.Operator) (string, interface{}) {
    switch operator {
    // ... existing cases ...
    case "myNewOperator":
        return fmt.Sprintf("%s MY_SQL_CLAUSE ?", column), value
    }
    return "", nil
}
```

The operator `value` string must match what the frontend sends in `operator.value`.

For date range operators, add to `handleDateRangeCondition()` instead — these are invoked when `operator.value == "isWithin"` and the actual range type is in `secondOperator.Value`.

---

## Known Limitation

OR conjunctions with **nested subgroups** are not applied — `applyGroupFilters` skips nested groups inside an OR group (`continue` statement at line ~44). Only flat FilterItems work inside OR groups. AND groups support full nesting.

---

## TabPayload Structure (reference)

```go
type TabPayload struct {
    TabID   uint              `json:"tabId"`
    Filters json.RawMessage   `json:"filters"`
    Columns []ColumnInput     `json:"columns"`
    Page    int               `json:"page"`    // 0-indexed
    Size    int               `json:"size"`    // default 50, max 100
    Search  string            `json:"search"`
}
```

---

## Testing a Filter

```bash
curl -X POST http://localhost:8080/api/maids/index \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "tabId": 1,
    "filters": {
      "id": "root",
      "type": "GROUP",
      "conjunction": "and",
      "children": [
        {
          "id": "f1",
          "type": "FILTER",
          "fieldId": 1,
          "columnType": "string_field",
          "operator": { "label": "contains", "value": "contains" },
          "value": "Ahmad"
        }
      ]
    },
    "columns": [],
    "page": 0,
    "size": 20,
    "search": ""
  }'
```
