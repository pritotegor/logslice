package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempLimitFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "limit_input.log")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func captureLimitOutput(t *testing.T, r *os.File, offset, count int) string {
	t.Helper()
	var out bytes.Buffer
	if err := RunLimit(r, &out, offset, count); err != nil {
		t.Fatalf("RunLimit error: %v", err)
	}
	return strings.TrimSpace(out.String())
}

func TestRunLimit_FromFile_BasicCount(t *testing.T) {
	lines := "alpha\nbeta\ngamma\ndelta\nepsilon\n"
	path := writeTempLimitFile(t, lines)

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	result := captureLimitOutput(t, f, 0, 3)
	got := strings.Split(result, "\n")
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if got[0] != "alpha" {
		t.Errorf("expected first line 'alpha', got %q", got[0])
	}
}

func TestRunLimit_FromFile_WithOffset(t *testing.T) {
	lines := "one\ntwo\nthree\nfour\nfive\n"
	path := writeTempLimitFile(t, lines)

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	result := captureLimitOutput(t, f, 2, 2)
	got := strings.Split(result, "\n")
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}
	if got[0] != "three" {
		t.Errorf("expected 'three', got %q", got[0])
	}
	if got[1] != "four" {
		t.Errorf("expected 'four', got %q", got[1])
	}
}

func TestRunLimit_EmptyFile(t *testing.T) {
	path := writeTempLimitFile(t, "")

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var out bytes.Buffer
	if err := RunLimit(f, &out, 0, 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Len() != 0 {
		t.Errorf("expected empty output for empty file")
	}
}
