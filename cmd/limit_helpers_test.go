package cmd

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// runLimitOnReader is a helper for table-driven limit tests.
func runLimitOnReader(t *testing.T, input string, offset, count int) []string {
	t.Helper()
	var out bytes.Buffer
	if err := RunLimit(strings.NewReader(input), &out, offset, count); err != nil {
		t.Fatalf("RunLimit: %v", err)
	}
	raw := strings.TrimSpace(out.String())
	if raw == "" {
		return []string{}
	}
	return strings.Split(raw, "\n")
}

// limitLineCount returns the number of output lines for given params.
func limitLineCount(t *testing.T, r io.Reader, offset, count int) int {
	t.Helper()
	var out bytes.Buffer
	if err := RunLimit(r, &out, offset, count); err != nil {
		t.Fatalf("RunLimit: %v", err)
	}
	raw := strings.TrimSpace(out.String())
	if raw == "" {
		return 0
	}
	return len(strings.Split(raw, "\n"))
}

func TestRunLimitHelper_TableDriven(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		offset int
		count  int
		want   int
	}{
		{"zero count", "a\nb\nc\n", 0, 0, 0},
		{"all lines", "a\nb\nc\n", 0, 10, 3},
		{"offset one", "a\nb\nc\n", 1, 10, 2},
		{"offset and count", "a\nb\nc\nd\n", 1, 2, 2},
		{"offset equals length", "a\nb\n", 2, 5, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := runLimitOnReader(t, tc.input, tc.offset, tc.count)
			if len(got) != tc.want {
				t.Errorf("expected %d lines, got %d: %v", tc.want, len(got), got)
			}
		})
	}
}
