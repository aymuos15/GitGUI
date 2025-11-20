package utils

import "strings"

// Truncate truncates a string to maxLen characters
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// PadRight pads a string with spaces to reach the desired length
func PadRight(s string, length int) string {
	// Strip ANSI codes for length calculation
	visibleLen := len(StripAnsi(s))
	if visibleLen >= length {
		return s
	}
	return s + strings.Repeat(" ", length-visibleLen)
}

// StripAnsi removes ANSI escape codes from a string for length calculation
func StripAnsi(s string) string {
	result := ""
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		result += string(r)
	}
	return result
}
