package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureGrepOutput(t *testing.T, lines []string, pattern, field string, invert, ignoreCase bool) string {
	t.Helper()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	grepPattern = pattern
	grepField = field
	grepInvert = invert
	grepIgnoreCase = ignoreCase

	input := strings.NewReader(strings.Join(lines, "\n") + "\n")
	scanner := bufio.NewScannerFromReader(input)
	_ = scanner

	// Use RunGrep via args-based openInput substitute
	origOpen := openInputFn
	openInputFn = func(_ []string) (io.Reader, error) {
		return strings.NewReader(strings.Join(lines, "\n") + "\n"), nil
	}
	defer func() { openInputFn = origOpen }()

	err := RunGrep(nil, nil)
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("RunGrep error: %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestRunGrep_BasicMatch(t *testing.T) {
	lines := []string{
		`{"level":"error","msg":"disk full"}`,
		`{"level":"info","msg":"all good"}`,
		`{"level":"error","msg":"timeout"}`,
	}

	grepPattern = "error"
	grepField = ""
	grepInvert = false
	grepIgnoreCase = false

	input := strings.NewReader(strings.Join(lines, "\n") + "\n")
	var out bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	_ = runGrepOnReader(input)

	w.Close()
	os.Stdout = stdout
	io.Copy(&out, r)

	result := out.String()
	if !strings.Contains(result, "disk full") {
		t.Errorf("expected 'disk full' in output, got: %s", result)
	}
	if strings.Contains(result, "all good") {
		t.Errorf("expected 'all good' to be filtered out, got: %s", result)
	}
}

func TestRunGrep_InvertMatch(t *testing.T) {
	lines := []string{
		`{"level":"error","msg":"bad"}`,
		`{"level":"info","msg":"ok"}`,
	}

	grepPattern = "error"
	grepField = ""
	grepInvert = true
	grepIgnoreCase = false

	input := strings.NewReader(strings.Join(lines, "\n") + "\n")
	result := runGrepOnReader(input)

	if strings.Contains(result, "bad") {
		t.Errorf("inverted: 'bad' should be excluded")
	}
	if !strings.Contains(result, "ok") {
		t.Errorf("inverted: 'ok' should be included")
	}
}
