package cmd

import (
	"strings"
	"testing"
)

func TestCountLines_NoField_Total(t *testing.T) {
	input := strings.NewReader("line1\nline2\nline3\n")
	_, total, err := countLines(input, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 3 {
		t.Errorf("expected 3, got %d", total)
	}
}

func TestCountLines_SkipsEmptyLines(t *testing.T) {
	input := strings.NewReader("line1\n\nline2\n\n")
	_, total, err := countLines(input, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Errorf("expected 2, got %d", total)
	}
}

func TestCountLines_ByField_GroupsCorrectly(t *testing.T) {
	lines := []string{
		`{"level":"info","msg":"a"}`,
		`{"level":"error","msg":"b"}`,
		`{"level":"info","msg":"c"}`,
		`{"level":"warn","msg":"d"}`,
	}
	input := strings.NewReader(strings.Join(lines, "\n"))
	counts, total, err := countLines(input, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 4 {
		t.Errorf("expected total 4, got %d", total)
	}
	if counts["info"] != 2 {
		t.Errorf("expected info=2, got %d", counts["info"])
	}
	if counts["error"] != 1 {
		t.Errorf("expected error=1, got %d", counts["error"])
	}
	if counts["warn"] != 1 {
		t.Errorf("expected warn=1, got %d", counts["warn"])
	}
}

func TestCountLines_ByField_MissingField(t *testing.T) {
	lines := []string{
		`{"level":"info"}`,
		`{"msg":"no level here"}`,
	}
	input := strings.NewReader(strings.Join(lines, "\n"))
	counts, _, err := countLines(input, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if counts["<missing>"] != 1 {
		t.Errorf("expected <missing>=1, got %d", counts["<missing>"])
	}
}

func TestCountLines_ByField_InvalidJSON(t *testing.T) {
	input := strings.NewReader("not json\n{\"level\":\"info\"}\n")
	counts, total, err := countLines(input, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if counts["<unparsed>"] != 1 {
		t.Errorf("expected <unparsed>=1, got %d", counts["<unparsed>"])
	}
}
