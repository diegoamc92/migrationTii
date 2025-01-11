package utils

import "strings"

// cleanRUT limpia y formatea el RUT.
func CleanRUT(rut string) string {
	if rut == "" {
		return ""
	}
	cleaned := strings.ReplaceAll(rut, ".", "")
	cleaned = strings.TrimLeft(cleaned, "0")
	if len(cleaned) > 1 && strings.ToUpper(string(cleaned[len(cleaned)-1])) == "K" {
		cleaned = cleaned[:len(cleaned)-1] + "K"
	}
	return cleaned
}
