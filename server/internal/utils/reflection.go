package utils

import "reflect"

// DiscoverRelations inspects the given struct type and returns a map of relation field names to their types.
//
// It detects relations by examining exported, non-anonymous struct fields that are either:
//   - Structs from external packages (excluding gorm/time types)
//   - Slices of structs (i.e., one-to-many relations)
//
// Parameters:
//   - t: the reflect.Type of the model struct (typically from reflect.TypeOf(model).Elem())
//
// Returns:
//   - a map where keys are relation field names and values are their reflect.Type
func DiscoverRelations(t reflect.Type) map[string]reflect.Type {
	relations := map[string]reflect.Type{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() || field.Anonymous {
			continue
		}

		// Skip internal GORM/time fields
		if field.Type.String() == "gorm.DeletedAt" || field.Type.PkgPath() == "gorm.io/gorm" || field.Type.PkgPath() == "time" {
			continue
		}

		switch kind := field.Type.Kind(); kind {
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.Struct {
				relations[field.Name] = field.Type
			}
		case reflect.Struct:
			if field.Type.Name() != "" && field.Type.PkgPath() != "" {
				relations[field.Name] = field.Type
			}
		}
	}

	return relations
}
