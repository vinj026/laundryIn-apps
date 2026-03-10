package utils

import "strings"

// Sanitize removes leading/trailing whitespace and null bytes from input strings.
// Use this for all user input sanitization across the application.
func Sanitize(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\x00", "")
	return s
}
