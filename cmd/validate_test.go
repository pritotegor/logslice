package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestValidateLines_AllValid(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"info","msg":"started"}`,
		`{"level":"warn","msg":"slow query"}`,
		`{"level":"error","msg":"failed"}`,
	}, "\n") + "\n"

	var errBuf bytes.Buffer
	invalid, total, err := validateLines(strings.NewReader(input), &errBuf, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if invalid != 0 {
		t.Errorf("expected 0 invalid, got %d", invalid)
	}
	if total != 3 {
		t.Errorf("expected 3 total, got %d", total)
	}
	if errBuf.Len() != 0 {
		t.Errorf("expected no error output, got: %s", errBuf.String())
	}
}

func TestValidateLines_SomeInvalid(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"info"}`,
		`not json at all`,
		`{"level":"error"}`,
		`{broken`,
	}, "\n") + "\n"

	var errBuf bytes.Buffer
	invalid, total, err := validateLines(strings.NewReader(input), &errBuf, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if invalid != 2 {
		t.Errorf("expected 2 invalid, got %d", invalid)
	}
	if total != 4 {
		t.Errorf("expected 4 total, got %d", total)
	}
}

func TestValidateLines_StrictStopsEarly(t *testing.T) {
	input := strings.Join([]string{
		`{"ok":true}`,
		`bad line`,
		`also bad`,
	}, "\n") + "\n"

	var errBuf bytes.Buffer
	invalid, total, err := validateLines(strings.NewReader(input), &errBuf, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if invalid != 1 {
		t.Errorf("expected 1 invalid (strict), got %d", invalid)
	}
	if total != 2 {
		t.Errorf("expected 2 total lines read before stop, got %d", total)
	}
}

func TestValidateLines_QuietSuppressesOutput(t *testing.T) {
	input := `not json\n`

	var errBuf bytes.Buffer
	invalid, _, err := validateLines(strings.NewReader(input), &errBuf, false, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if invalid != 1 {
		t.Errorf("expected 1 invalid, got %d", invalid)
	}
	if errBuf.Len() != 0 {
		t.Errorf("expected no error output in quiet mode, got: %s", errBuf.String())
	}
}

func TestValidateLines_SkipsEmptyLines(t *testing.T) {
	input := "\n\n{\"ok\":true}\n\n"

	var errBuf bytes.Buffer
	invalid, total, err := validateLines(strings.NewReader(input), &errBuf, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if invalid != 0 {
		t.Errorf("expected 0 invalid, got %d", invalid)
	}
	if total != 1 {
		t.Errorf("expected 1 non-empty line, got %d", total)
	}
}
