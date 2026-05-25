package services

import (
	"strings"
	"testing"

	"github.com/AhmadAboElzahab/bridge/internal/models"
)

// buildCondition

func TestBuildCondition_IsOperators(t *testing.T) {
	for _, op := range []string{"is", "=", "==", "isExactly"} {
		sql, val := buildCondition("name", op, "Ahmad")
		if sql != "name = ?" {
			t.Errorf("op=%q: got SQL %q, want %q", op, sql, "name = ?")
		}
		if val != "Ahmad" {
			t.Errorf("op=%q: got val %v, want Ahmad", op, val)
		}
	}
}

func TestBuildCondition_IsNotOperators(t *testing.T) {
	for _, op := range []string{"isNot", "!="} {
		sql, val := buildCondition("name", op, "Ahmad")
		if sql != "name != ?" {
			t.Errorf("op=%q: got SQL %q, want %q", op, sql, "name != ?")
		}
		if val != "Ahmad" {
			t.Errorf("op=%q: got val %v, want Ahmad", op, val)
		}
	}
}

func TestBuildCondition_Contains(t *testing.T) {
	sql, val := buildCondition("bio", "contains", "nurse")
	if sql != "bio ILIKE ?" {
		t.Errorf("got SQL %q, want %q", sql, "bio ILIKE ?")
	}
	if val != "%nurse%" {
		t.Errorf("got val %v, want %%nurse%%", val)
	}
}

func TestBuildCondition_DoesNotContain(t *testing.T) {
	sql, val := buildCondition("bio", "doesNotContain", "nurse")
	if sql != "bio NOT ILIKE ?" {
		t.Errorf("got SQL %q, want %q", sql, "bio NOT ILIKE ?")
	}
	if val != "%nurse%" {
		t.Errorf("got val %v, want %%nurse%%", val)
	}
}

func TestBuildCondition_IsEmpty(t *testing.T) {
	sql, val := buildCondition("bio", "isEmpty", nil)
	if sql != "bio IS NULL" {
		t.Errorf("got SQL %q, want %q", sql, "bio IS NULL")
	}
	if val != nil {
		t.Errorf("got val %v, want nil", val)
	}
}

func TestBuildCondition_IsNotEmpty(t *testing.T) {
	sql, val := buildCondition("bio", "isNotEmpty", nil)
	if sql != "bio IS NOT NULL" {
		t.Errorf("got SQL %q, want %q", sql, "bio IS NOT NULL")
	}
	if val != nil {
		t.Errorf("got val %v, want nil", val)
	}
}

func TestBuildCondition_ComparisonOperators(t *testing.T) {
	cases := []struct{ op, wantSQL string }{
		{"<", "age < ?"},
		{"isBefore", "age < ?"},
		{"<=", "age <= ?"},
		{"isOnOrBefore", "age <= ?"},
		{">", "age > ?"},
		{"isAfter", "age > ?"},
		{">=", "age >= ?"},
		{"isOnOrAfter", "age >= ?"},
	}
	for _, tc := range cases {
		sql, val := buildCondition("age", tc.op, 25)
		if sql != tc.wantSQL {
			t.Errorf("op=%q: got SQL %q, want %q", tc.op, sql, tc.wantSQL)
		}
		if val != 25 {
			t.Errorf("op=%q: got val %v, want 25", tc.op, val)
		}
	}
}

func TestBuildCondition_SetOperators(t *testing.T) {
	vals := []interface{}{"Available", "Hired"}

	for _, op := range []string{"isAnyOf", "hasAnyOf"} {
		sql, _ := buildCondition("status", op, vals)
		if sql != "status IN (?)" {
			t.Errorf("op=%q: got SQL %q, want %q", op, sql, "status IN (?)")
		}
	}
	for _, op := range []string{"hasNoneOf", "isNoneOf"} {
		sql, _ := buildCondition("status", op, vals)
		if sql != "status NOT IN (?)" {
			t.Errorf("op=%q: got SQL %q, want %q", op, sql, "status NOT IN (?)")
		}
	}

	sql, _ := buildCondition("tags", "hasAllOf", vals)
	if sql != "tags @> ?" {
		t.Errorf("got SQL %q, want %q", sql, "tags @> ?")
	}
}

func TestBuildCondition_UnknownOperator(t *testing.T) {
	sql, val := buildCondition("name", "doesNotExist", "x")
	if sql != "" || val != nil {
		t.Errorf("unknown operator should return empty, got SQL=%q val=%v", sql, val)
	}
}

// handleDateRangeCondition (isWithin logic)

func TestHandleDateRangeCondition_Today(t *testing.T) {
	sql, val := handleDateRangeCondition("created_at", "today", "")
	if !strings.Contains(sql, "created_at >= ?") || !strings.Contains(sql, "created_at < ?") {
		t.Errorf("today: got SQL %q, want range condition with >= and <", sql)
	}
	vals, ok := val.([]interface{})
	if !ok || len(vals) != 2 {
		t.Errorf("today: val should be []interface{} with 2 elements, got %v", val)
	}
}

func TestHandleDateRangeCondition_ExactDate(t *testing.T) {
	sql, val := handleDateRangeCondition("date_of_birth", "exactDate", "2000-01-15")
	if sql != "date_of_birth = ?" {
		t.Errorf("exactDate: got SQL %q, want %q", sql, "date_of_birth = ?")
	}
	if val == nil {
		t.Error("exactDate: val should not be nil")
	}
}

func TestHandleDateRangeCondition_InvalidExactDate(t *testing.T) {
	sql, val := handleDateRangeCondition("date_of_birth", "exactDate", "not-a-date")
	if sql != "" || val != nil {
		t.Errorf("invalid date should return empty, got SQL=%q val=%v", sql, val)
	}
}

func TestHandleDateRangeCondition_NumberOfDaysAgo(t *testing.T) {
	sql, val := handleDateRangeCondition("created_at", "numberOfDaysAgo", "7")
	if sql != "created_at >= ?" {
		t.Errorf("numberOfDaysAgo: got SQL %q, want %q", sql, "created_at >= ?")
	}
	if val == nil {
		t.Error("numberOfDaysAgo: val should not be nil")
	}
}

func TestHandleDateRangeCondition_NumberOfDaysAgoInvalidValue(t *testing.T) {
	sql, val := handleDateRangeCondition("created_at", "numberOfDaysAgo", "notanumber")
	if sql != "" || val != nil {
		t.Errorf("invalid days value should return empty, got SQL=%q val=%v", sql, val)
	}
}

func TestHandleDateRangeCondition_CalendarRanges(t *testing.T) {
	for _, op := range []string{"thisCalendarWeek", "thisCalendarMonth", "thisCalendarYear"} {
		sql, val := handleDateRangeCondition("created_at", op, "")
		if !strings.Contains(sql, "BETWEEN ? AND ?") {
			t.Errorf("op=%q: got SQL %q, want BETWEEN condition", op, sql)
		}
		vals, ok := val.([]interface{})
		if !ok || len(vals) != 2 {
			t.Errorf("op=%q: val should be []interface{} with 2 elements", op)
		}
	}
}

func TestHandleDateRangeCondition_UnknownOperator(t *testing.T) {
	sql, val := handleDateRangeCondition("created_at", "unknownRange", "")
	if sql != "" || val != nil {
		t.Errorf("unknown range should return empty, got SQL=%q val=%v", sql, val)
	}
}

// resolveColumnName

func TestResolveColumnName_SingleRelation(t *testing.T) {
	field := models.FormField{FieldKey: "nationality", FormFieldType: "single_relation"}
	got := resolveColumnName(field)
	if got != "nationality_id" {
		t.Errorf("got %q, want %q", got, "nationality_id")
	}
}

func TestResolveColumnName_OtherTypes(t *testing.T) {
	for _, ft := range []string{"string_field", "integer_field", "single_select", "multi_relation", "date_field"} {
		field := models.FormField{FieldKey: "first_name", FormFieldType: ft}
		got := resolveColumnName(field)
		if got != "first_name" {
			t.Errorf("type=%q: got %q, want %q", ft, got, "first_name")
		}
	}
}

// normalizeValue

func TestNormalizeValue_WholeFloat(t *testing.T) {
	got := normalizeValue(float64(42))
	if got != 42 {
		t.Errorf("got %v (%T), want int 42", got, got)
	}
}

func TestNormalizeValue_FractionalFloat(t *testing.T) {
	got := normalizeValue(float64(3.14))
	if got != float64(3.14) {
		t.Errorf("got %v, want 3.14", got)
	}
}

func TestNormalizeValue_NonFloat(t *testing.T) {
	got := normalizeValue("hello")
	if got != "hello" {
		t.Errorf("got %v, want hello", got)
	}
}

// extractActualValue

func TestExtractActualValue_MapWithValueKey(t *testing.T) {
	m := map[string]interface{}{"value": "Ahmad", "label": "display"}
	got := extractActualValue(m)
	if got != "Ahmad" {
		t.Errorf("got %v, want Ahmad", got)
	}
}

func TestExtractActualValue_MapWithoutValueKey(t *testing.T) {
	m := map[string]interface{}{"label": "display"}
	got := extractActualValue(m)
	if got == nil {
		t.Error("expected map back, got nil")
	}
}

func TestExtractActualValue_PlainValue(t *testing.T) {
	got := extractActualValue("raw")
	if got != "raw" {
		t.Errorf("got %v, want raw", got)
	}
}
