package cmd

import (
	"fmt"
	"strings"
	"time"
)

// SupportedFormats lists all available output formats.
var SupportedFormats = []string{"raw", "json", "csv"}

// SupportedTimeFormats lists human-friendly aliases for time formats.
var SupportedTimeFormats = map[string]string{
	"rfc3339": time.RFC3339,
	"rfc3339nano": time.RFC3339Nano,
	"unix": "2006-01-02T15:04:05Z07:00",
	"datetime": "2006-01-02 15:04:05",
	"date": "2006-01-02",
}

// ValidateFormat checks whether the given output format is supported.
func ValidateFormat(format string) error {
	for _, f := range SupportedFormats {
		if f == format {
			return nil
		}
	}
	return fmt.Errorf("unsupported format %q: must be one of [%s]", format, strings.Join(SupportedFormats, ", "))
}

// ResolveTimeFormat maps a human-friendly alias to a Go time layout string.
// If the alias is not found, the input is returned as-is (allowing raw layouts).
func ResolveTimeFormat(alias string) string {
	if layout, ok := SupportedTimeFormats[strings.ToLower(alias)]; ok {
		return layout
	}
	return alias
}

// FormatTimestamp formats a time.Time value using the given alias or layout.
func FormatTimestamp(t time.Time, alias string) string {
	layout := ResolveTimeFormat(alias)
	return t.Format(layout)
}

// ParseTimeArg parses a time string using a given alias or layout.
func ParseTimeArg(value, alias string) (time.Time, error) {
	layout := ResolveTimeFormat(alias)
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot parse %q with format %q: %w", value, alias, err)
	}
	return t, nil
}
