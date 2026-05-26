package utils

import "strings"

// ToPascalCase converts a snake_case string to PascalCase.
// e.g. "languages" → "Languages", "first_name" → "FirstName"
func ToPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}
