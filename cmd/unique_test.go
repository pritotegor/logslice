package cmd

import (
	"strings"
	"testing"
)

func TestExtractUniqueFieldValue_String(t *testing.T) {
	line := `{"level":"info","msg":"started"}`
	got := extractUniqueFieldValue(line, "level")
	if got != "info" {
		t.Errorf("expected 'info', got %q", got)
	}
}

func TestExtractUniqueFieldValue_Number(t *testing.T) {
	line := `{"status":200}`
	got := extractUniqueFieldValue(line, "status")
	if got != "200" {
		t.Errorf("expected '200', got %q", got)
	}
}

func TestExtractUniqueFieldValue_MissingField(t *testing.T) {
	line := `{"level":"info"}`
	got := extractUniqueFieldValue(line, "service")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestExtractUniqueFieldValue_NonJSON(t *testing.T) {
	got := extractUniqueFieldValue("plain text line", "level")
	if got != "" {
		t.Errorf("expected empty string for non-JSON, got %q", got)
	}
}

func TestCollectUniqueValues_Basic(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"info"}`,
		`{"level":"error"}`,
		`{"level":"info"}`,
		`{"level":"warn"}`,
	}, "\n")

	counts, order, err := collectUniqueValues(strings.NewReader(input), "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 3 {
		t.Errorf("expected 3 unique values, got %d", len(order))
	}
	if counts["info"] != 2 {
		t.Errorf("expected info count=2, got %d", counts["info"])
	}
	if counts["error"] != 1 {
		t.Errorf("expected error count=1, got %d", counts["error"])
	}
}

func TestCollectUniqueValues_SkipsEmptyLines(t *testing.T) {
	input := "\n" + `{"level":"info"}` + "\n\n" + `{"level":"info"}` + "\n"
	counts, order, err := collectUniqueValues(strings.NewReader(input), "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(order) != 1 {
		t.Errorf("expected 1 unique value, got %d", len(order))
	}
	if counts["info"] != 2 {
		t.Errorf("expected info count=2, got %d", counts["info"])
	}
}

func TestCollectUniqueValues_PreservesInsertionOrder(t *testing.T) {
	input := strings.Join([]string{
		`{"svc":"auth"}`,
		`{"svc":"api"}`,
		`{"svc":"auth"}`,
		`{"svc":"db"}`,
	}, "\n")
	_, order, err := collectUniqueValues(strings.NewReader(input), "svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"auth", "api", "db"}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("order[%d]: expected %q, got %q", i, v, order[i])
		}
	}
}
