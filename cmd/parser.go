package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// CommonTimeFormats lists timestamp formats tried during auto-detection.
var CommonTimeFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
}

// ParsedLine represents a decoded log line with its raw form preserved.
type ParsedLine struct {
	Raw       string
	Fields    map[string]interface{}
	Timestamp *time.Time
	IsJSON    bool
}

// ParseLine attempts to parse a raw log line as JSON and extract a timestamp.
func ParseLine(raw string, tsField string) (*ParsedLine, error) {
	pl := &ParsedLine{Raw: raw}

	trimmed := strings.TrimSpace(raw)
	if len(trimmed) == 0 {
		return pl, nil
	}

	if trimmed[0] == '{' {
		var fields map[string]interface{}
		if err := json.Unmarshal([]byte(trimmed), &fields); err != nil {
			return pl, fmt.Errorf("json parse error: %w", err)
		}
		pl.Fields = fields
		pl.IsJSON = true

		if tsField != "" {
			if v, ok := fields[tsField]; ok {
				if s, ok := v.(string); ok {
					if t, err := parseTimestamp(s); err == nil {
						pl.Timestamp = &t
					}
				}
			}
		}
	}

	return pl, nil
}

// parseTimestamp tries each CommonTimeFormat until one succeeds.
func parseTimestamp(s string) (time.Time, error) {
	for _, layout := range CommonTimeFormats {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised timestamp: %q", s)
}

// ParseTimestamp is the exported wrapper used by other packages.
func ParseTimestamp(s string) (time.Time, error) {
	return parseTimestamp(s)
}
