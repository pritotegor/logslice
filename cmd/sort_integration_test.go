package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func writeTempSortFile(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp("", "sort_test_*.log")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

func TestRunSort_FromFile_TimestampAscending(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-03T00:00:00Z","msg":"third"}`,
		`{"ts":"2024-01-01T00:00:00Z","msg":"first"}`,
		`{"ts":"2024-01-02T00:00:00Z","msg":"second"}`,
	}
	path := writeTempSortFile(t, lines)
	defer os.Remove(path)

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("opening file: %v", err)
	}
	defer f.Close()

	var out bytes.Buffer
	if err := RunSort(f, &out, "", false); err != nil {
		t.Fatalf("RunSort error: %v", err)
	}

	result := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(result) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(result))
	}
	if !strings.Contains(result[0], "first") {
		t.Errorf("expected first msg first, got %q", result[0])
	}
	if !strings.Contains(result[2], "third") {
		t.Errorf("expected third msg last, got %q", result[2])
	}
}

func TestRunSort_FromFile_TimestampDescending(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-01T00:00:00Z","msg":"first"}`,
		`{"ts":"2024-01-03T00:00:00Z","msg":"third"}`,
		`{"ts":"2024-01-02T00:00:00Z","msg":"second"}`,
	}
	path := writeTempSortFile(t, lines)
	defer os.Remove(path)

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("opening file: %v", err)
	}
	defer f.Close()

	var out bytes.Buffer
	if err := RunSort(f, &out, "", true); err != nil {
		t.Fatalf("RunSort error: %v", err)
	}

	result := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if !strings.Contains(result[0], "third") {
		t.Errorf("expected third msg first in descending, got %q", result[0])
	}
}

func TestRunSort_MixedJSON_NonJSON(t *testing.T) {
	input := strings.Join([]string{
		`plain log line`,
		`{"ts":"2024-01-01T00:00:00Z","msg":"structured"}`,
	}, "\n") + "\n"

	var out bytes.Buffer
	if err := RunSort(strings.NewReader(input), &out, "", false); err != nil {
		t.Fatalf("RunSort error: %v", err)
	}
	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 output lines, got %d", len(lines))
	}
}
