package inspector2util

import (
	"strings"
)

func CanonicalSeverity(s string) string {
	m := map[string]string{
		"":          "None",
		"critical":  "Critical",
		"high":      "High",
		"important": "High",
		"medium":    "Medium",
		"moderate":  "Medium",
		"low":       "Low",
	}
	if v, ok := m[strings.ToLower(s)]; ok {
		return v
	} else {
		return s
	}
}
