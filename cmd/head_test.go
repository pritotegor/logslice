package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func buildLines(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&sb, "line%d\n", i)
	}
	return sb.String()
}

func headN(input string, n int) ([]string, error) {
	r := strings.NewReader(input)
	lines, err := readLastN(r, countLines(input))
	if err != nil {
		return nil, err
	}
	if n > len(lines) {
		n = len(lines)
	}
	return lines[:n], nil
}

func countLines(s string) int {
	return strings.Count(s, "\n")
}

func TestHeadN_FewerLinesThanN(t *testing.T) {
	input := buildLines(3)
	result, err := headN(input, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3, got %d", len(result))
	}
}

func TestHeadN_ExactN(t *testing.T) {
	input := buildLines(5)
	result, err := headN(input, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 5 {
		t.Errorf("expected 5, got %d", len(result))
	}
}

func TestHeadN_MoreLinesThanN(t *testing.T) {
	input := buildLines(20)
	result, err := headN(input, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 5 {
		t.Errorf("expected 5, got %d", len(result))
	}
	if result[0] != "line1" {
		t.Errorf("first line should be line1, got %q", result[0])
	}
}

func TestHeadN_WriteOutput(t *testing.T) {
	input := buildLines(5)
	result, _ := headN(input, 3)

	var buf bytes.Buffer
	for _, l := range result {
		buf.WriteString(l + "\n")
	}
	got := buf.String()
	if !strings.Contains(got, "line1") {
		t.Errorf("missing line1 in output: %q", got)
	}
	if strings.Contains(got, "line4") {
		t.Errorf("should not contain line4: %q", got)
	}
}
