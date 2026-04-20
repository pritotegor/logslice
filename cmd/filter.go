package cmd

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

// FilterOptions holds all parameters for a filter run.
type FilterOptions struct {
	TimeField  string
	From       *time.Time
	To         *time.Time
	FieldMatch map[string]string // key=value pairs that must match
}

// RunFilter reads lines from r, applies opts, and writes matching lines to w.
func RunFilter(r io.Reader, w io.Writer, opts FilterOptions, writer *Writer) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		raw := scanner.Text()
		pl, err := ParseLine(raw, opts.TimeField)
		if err != nil {
			// skip unparseable lines silently
			continue
		}
		if !matchesTime(pl, opts) {
			continue
		}
		if !matchesAllFields(pl, opts.FieldMatch) {
			continue
		}
		if err := writer.WriteLine(raw, pl.Fields); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
	}
	return scanner.Err()
}

// matchesTime returns true when the line's timestamp falls within [From, To].
// Lines without a timestamp pass through unless a range is specified.
func matchesTime(pl *ParsedLine, opts FilterOptions) bool {
	if opts.From == nil && opts.To == nil {
		return true
	}
	if pl.Timestamp == nil {
		return false
	}
	if opts.From != nil && pl.Timestamp.Before(*opts.From) {
		return false
	}
	if opts.To != nil && pl.Timestamp.After(*opts.To) {
		return false
	}
	return true
}

// matchesField checks whether a single field in the parsed line matches the
// expected value (case-insensitive substring match).
func matchesField(fields map[string]interface{}, key, value string) bool {
	v, ok := fields[key]
	if !ok {
		return false
	}
	return strings.Contains(
		strings.ToLower(fmt.Sprintf("%v", v)),
		strings.ToLower(value),
	)
}

// matchesAllFields returns true when every entry in want matches the line.
func matchesAllFields(pl *ParsedLine, want map[string]string) bool {
	if len(want) == 0 {
		return true
	}
	if !pl.IsJSON || pl.Fields == nil {
		return false
	}
	for k, v := range want {
		if !matchesField(pl.Fields, k, v) {
			return false
		}
	}
	return true
}
