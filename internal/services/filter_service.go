package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/AhmadAboElzahab/bridge/internal/models"
	"gorm.io/gorm"
)

func ApplyFilters(db *gorm.DB, raw json.RawMessage, fields map[string]models.FormField) *gorm.DB {
	var group models.FilterGroup

	if err := json.Unmarshal(raw, &group); err != nil {
		log.Println("Failed to unmarshal filter group:", err)
		return db
	}
	return applyGroupFilters(db, group, fields)
}

func applyGroupFilters(db *gorm.DB, group models.FilterGroup, fields map[string]models.FormField) *gorm.DB {
	if len(group.Children) == 0 {
		return db
	}

	if group.Conjunction == "or" {
		var orParts []*gorm.DB

		for _, raw := range group.Children {
			childJSON, err := json.Marshal(raw)
			if err != nil {
				log.Println("Failed to marshal OR child filter:", err)
				continue
			}
			// Each OR part gets a fresh statement so sub-conditions don't bleed into siblings.
			partDB := db.Session(&gorm.Session{NewDB: true})

			if rawMap, ok := raw.(map[string]interface{}); ok && rawMap["type"] == "GROUP" {
				var subgroup models.FilterGroup
				if err := json.Unmarshal(childJSON, &subgroup); err != nil {
					log.Println("Failed to unmarshal OR subgroup:", err)
					continue
				}
				partDB = applyGroupFilters(partDB, subgroup, fields)
			} else {
				var item models.FilterItem
				if err := json.Unmarshal(childJSON, &item); err != nil {
					log.Println("Failed to unmarshal OR filter item:", err)
					continue
				}
				partDB = applyItemFilter(partDB, item, fields)
			}

			orParts = append(orParts, partDB)
		}

		if len(orParts) > 0 {
			combined := orParts[0]
			for _, part := range orParts[1:] {
				combined = combined.Or(part)
			}
			db = db.Where(combined)
		}

		return db
	}

	// AND conjunction — apply each child's conditions sequentially
	for _, raw := range group.Children {
		childJSON, err := json.Marshal(raw)
		if err != nil {
			log.Println("Failed to marshal AND child filter:", err)
			continue
		}

		if rawMap, ok := raw.(map[string]interface{}); ok && rawMap["type"] == "GROUP" {
			var subgroup models.FilterGroup
			if err := json.Unmarshal(childJSON, &subgroup); err != nil {
				log.Println("Failed to unmarshal AND subgroup:", err)
				continue
			}
			db = applyGroupFilters(db, subgroup, fields)
		} else {
			var item models.FilterItem
			if err := json.Unmarshal(childJSON, &item); err != nil {
				log.Println("Failed to unmarshal AND filter item:", err)
				continue
			}
			db = applyItemFilter(db, item, fields)
		}
	}

	return db
}

func applyItemFilter(db *gorm.DB, item models.FilterItem, fields map[string]models.FormField) *gorm.DB {
	field, ok := fields[fmt.Sprintf("%d", item.FieldID)]
	if !ok {
		return db
	}

	column := resolveColumnName(field)
	value := normalizeValue(extractActualValue(item.Value))

	condition, conditionValue := buildCondition(column, item.Operator.Value, value, item.SecondOperator)
	if condition != "" && conditionValue != nil {
		if slice, ok := conditionValue.([]interface{}); ok {
			db = db.Where(condition, slice...)
		} else {
			db = db.Where(condition, conditionValue)
		}
	} else if condition != "" {
		db = db.Where(condition)
	}

	return db
}

func buildCondition(column, operator string, value interface{}, secondOperator *models.Operator) (string, interface{}) {
	switch operator {
	case "is", "=", "==", "isExactly":
		return fmt.Sprintf("%s = ?", column), value
	case "isNot", "!=":
		return fmt.Sprintf("%s != ?", column), value
	case "contains":
		return fmt.Sprintf("%s ILIKE ?", column), fmt.Sprintf("%%%s%%", toStr(value))
	case "doesNotContain":
		return fmt.Sprintf("%s NOT ILIKE ?", column), fmt.Sprintf("%%%s%%", toStr(value))
	case "isEmpty":
		return fmt.Sprintf("%s IS NULL", column), nil
	case "isNotEmpty":
		return fmt.Sprintf("%s IS NOT NULL", column), nil
	case "<", "isBefore":
		return fmt.Sprintf("%s < ?", column), value
	case "<=", "isOnOrBefore":
		return fmt.Sprintf("%s <= ?", column), value
	case ">", "isAfter":
		return fmt.Sprintf("%s > ?", column), value
	case ">=", "isOnOrAfter":
		return fmt.Sprintf("%s >= ?", column), value
	case "isAnyOf", "hasAnyOf":
		return fmt.Sprintf("%s IN (?)", column), value
	case "hasNoneOf", "isNoneOf":
		return fmt.Sprintf("%s NOT IN (?)", column), value
	case "hasAllOf":
		return fmt.Sprintf("%s @> ?", column), value
	case "isWithin":
		if secondOperator != nil {
			return handleDateRangeCondition(column, secondOperator.Value, toStr(value))
		}
	}
	return "", nil
}

func handleDateRangeCondition(column, operator, value string) (string, interface{}) {
	now := time.Now()

	switch operator {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 0, 1)
		return fmt.Sprintf("%s >= ? AND %s < ?", column, column), []interface{}{start, end}
	case "tomorrow":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
		end := start.AddDate(0, 0, 1)
		return fmt.Sprintf("%s >= ? AND %s < ?", column, column), []interface{}{start, end}
	case "yesterday":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -1)
		end := start.AddDate(0, 0, 1)
		return fmt.Sprintf("%s >= ? AND %s < ?", column, column), []interface{}{start, end}
	case "oneWeekAgo":
		return fmt.Sprintf("%s >= ?", column), now.AddDate(0, 0, -7)
	case "oneWeekFromNow":
		return fmt.Sprintf("%s <= ?", column), now.AddDate(0, 0, 7)
	case "oneMonthAgo":
		return fmt.Sprintf("%s >= ?", column), now.AddDate(0, -1, 0)
	case "oneMonthFromNow":
		return fmt.Sprintf("%s <= ?", column), now.AddDate(0, 1, 0)
	case "numberOfDaysAgo":
		if days, err := strconv.Atoi(value); err == nil {
			return fmt.Sprintf("%s >= ?", column), now.AddDate(0, 0, -days)
		}
	case "numberOfDaysFromNow":
		if days, err := strconv.Atoi(value); err == nil {
			return fmt.Sprintf("%s <= ?", column), now.AddDate(0, 0, days)
		}
	case "exactDate":
		if parsed, err := time.Parse("2006-01-02", value); err == nil {
			return fmt.Sprintf("%s = ?", column), parsed
		}
	case "thePastWeek":
		return fmt.Sprintf("%s >= ?", column), now.AddDate(0, 0, -7)
	case "thePastMonth":
		return fmt.Sprintf("%s >= ?", column), now.AddDate(0, -1, 0)
	case "thePastYear":
		return fmt.Sprintf("%s >= ?", column), now.AddDate(-1, 0, 0)
	case "theNextWeek":
		return fmt.Sprintf("%s <= ?", column), now.AddDate(0, 0, 7)
	case "theNextMonth":
		return fmt.Sprintf("%s <= ?", column), now.AddDate(0, 1, 0)
	case "theNextYear":
		return fmt.Sprintf("%s <= ?", column), now.AddDate(1, 0, 0)
	case "thisCalendarWeek":
		start := now.AddDate(0, 0, -int(now.Weekday()))
		end := start.AddDate(0, 0, 6)
		return fmt.Sprintf("%s BETWEEN ? AND ?", column), []interface{}{start, end}
	case "thisCalendarMonth":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, -1)
		return fmt.Sprintf("%s BETWEEN ? AND ?", column), []interface{}{start, end}
	case "thisCalendarYear":
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end := time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
		return fmt.Sprintf("%s BETWEEN ? AND ?", column), []interface{}{start, end}
	}

	return "", nil
}

func extractActualValue(val interface{}) interface{} {
	if m, ok := val.(map[string]interface{}); ok {
		if v, exists := m["value"]; exists {
			return v
		}
	}
	return val
}

func normalizeValue(val interface{}) interface{} {
	if floatVal, ok := val.(float64); ok && floatVal == float64(int(floatVal)) {
		return int(floatVal)
	}
	return val
}

func toStr(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func resolveColumnName(field models.FormField) string {
	switch field.FormFieldType {
	case "single_relation":
		return field.FieldKey + "_id"
	default:
		return field.FieldKey
	}
}
