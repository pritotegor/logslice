package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunDedupe_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	content := "alpha\nbeta\nalpha\ngamma\nbeta\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	var out bytes.Buffer
	err := dedupeLines(openFileReader(t, path), &out, "", 0)
	if err != nil {
		t.Fatalf("RunDedupe error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 unique lines, got %d: %v", len(lines), lines)
	}
}

func TestRunDedupe_JSONFieldDedupe(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"info","msg":"started"}`,
		`{"level":"warn","msg":"retrying"}`,
		`{"level":"info","msg":"done"}`,
		`{"level":"error","msg":"failed"}`,
	}, "\n") + "\n"

	var out bytes.Buffer
	err := dedupeLines(strings.NewReader(input), &out, "level", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := strings.TrimSpace(out.String())
	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (unique levels), got %d: %v", len(lines), lines)
	}

	for _, line := range lines {
		if !strings.HasPrefix(line, "{") {
			t.Errorf("expected JSON output, got: %s", line)
		}
	}
}

func TestRunDedupe_EmptyInput(t *testing.T) {
	var out bytes.Buffer
	err := dedupeLines(strings.NewReader(""), &out, "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Len() != 0 {
		t.Errorf("expected empty output, got: %q", out.String())
	}
}

func openFileReader(t *testing.T, path string) *os.File {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open file: %v", err)
	}
	t.Cleanup(func() { f.Close() })
	return f
}
