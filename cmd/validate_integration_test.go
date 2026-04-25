package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempValidateFile(t *testing.T, lines []string) string {
	t.Helper()
	tmp, err := os.CreateTemp(t.TempDir(), "validate_*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer tmp.Close()
	for _, l := range lines {
		tmp.WriteString(l + "\n")
	}
	return tmp.Name()
}

func captureValidateOutput(t *testing.T, args []string) (string, error) {
	t.Helper()
	cmd := validateCmd
	cmd.ResetFlags()
	cmd.Flags().BoolP("strict", "s", false, "")
	cmd.Flags().BoolP("quiet", "q", false, "")

	var errBuf bytes.Buffer
	cmd.SetErr(&errBuf)
	cmd.SetOut(&errBuf)

	err := cmd.RunE(cmd, args)
	return errBuf.String(), err
}

func TestRunValidate_AllValidFile(t *testing.T) {
	lines := []string{
		`{"ts":"2024-01-01T00:00:00Z","level":"info","msg":"boot"}`,
		`{"ts":"2024-01-01T00:00:01Z","level":"debug","msg":"ready"}`,
	}
	file := writeTempValidateFile(t, lines)

	reader, err := os.Open(file)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer reader.Close()

	var errBuf bytes.Buffer
	invalid, total, err := validateLines(reader, &errBuf, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if invalid != 0 {
		t.Errorf("expected 0 invalid, got %d", invalid)
	}
	if total != 2 {
		t.Errorf("expected 2 total, got %d", total)
	}
}

func TestRunValidate_MixedFile_ReportsInvalid(t *testing.T) {
	lines := []string{
		`{"level":"info"}`,
		`plain text log line`,
		`{"level":"warn"}`,
	}
	file := writeTempValidateFile(t, lines)
	_ = filepath.Base(file)

	reader, err := os.Open(file)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer reader.Close()

	var errBuf bytes.Buffer
	invalid, total, err := validateLines(reader, &errBuf, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if invalid != 1 {
		t.Errorf("expected 1 invalid, got %d", invalid)
	}
	if total != 3 {
		t.Errorf("expected 3 total, got %d", total)
	}
	if !strings.Contains(errBuf.String(), "plain text log line") {
		t.Errorf("expected error output to mention invalid line, got: %s", errBuf.String())
	}
}

func TestRunValidate_EmptyFile(t *testing.T) {
	file := writeTempValidateFile(t, []string{})

	reader, err := os.Open(file)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer reader.Close()

	var errBuf bytes.Buffer
	invalid, total, err := validateLines(reader, &errBuf, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if invalid != 0 || total != 0 {
		t.Errorf("expected 0/0 for empty file, got invalid=%d total=%d", invalid, total)
	}
}
