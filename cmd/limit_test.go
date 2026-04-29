package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func buildLimitLines(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		sb.WriteString(strings.Repeat("a", i) + "\n")
	}
	return sb.String()
}

func TestRunLimit_NoOffset_LimitFive(t *testing.T) {
	input := buildLimitLines(10)
	var out bytes.Buffer
	if err := RunLimit(strings.NewReader(input), &out, 0, 5); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(lines))
	}
}

func TestRunLimit_WithOffset(t *testing.T) {
	input := "line1\nline2\nline3\nline4\nline5\n"
	var out bytes.Buffer
	if err := RunLimit(strings.NewReader(input), &out, 2, 2); err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(out.String())
	want := "line3\nline4"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestRunLimit_OffsetBeyondInput(t *testing.T) {
	input := "line1\nline2\n"
	var out bytes.Buffer
	if err := RunLimit(strings.NewReader(input), &out, 10, 5); err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Fatalf("expected empty output, got %q", out.String())
	}
}

func TestRunLimit_CountZero_NoOutput(t *testing.T) {
	input := "line1\nline2\nline3\n"
	var out bytes.Buffer
	if err := RunLimit(strings.NewReader(input), &out, 0, 0); err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Fatalf("expected empty output, got %q", out.String())
	}
}

func TestRunLimit_SkipsEmptyLines(t *testing.T) {
	input := "line1\n\nline2\n\nline3\n"
	var out bytes.Buffer
	if err := RunLimit(strings.NewReader(input), &out, 0, 10); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 non-empty lines, got %d: %v", len(lines), lines)
	}
}

func TestRunLimit_FewerLinesThanCount(t *testing.T) {
	input := "a\nb\nc\n"
	var out bytes.Buffer
	if err := RunLimit(strings.NewReader(input), &out, 0, 100); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}
