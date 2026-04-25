package cmd

import (
	"encoding/json"
	"testing"
)

func TestReformatTimestamp_NonJSON_Passthrough(t *testing.T) {
	line := "plain text log line"
	out, err := reformatTimestamp(line, "timestamp", "", "2006-01-02T15:04:05Z07:00", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != line {
		t.Errorf("expected passthrough, got %q", out)
	}
}

func TestReformatTimestamp_FieldMissing_NoAdd(t *testing.T) {
	line := `{"level":"info","msg":"hello"}`
	out, err := reformatTimestamp(line, "timestamp", "", "2006-01-02", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != line {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestReformatTimestamp_FieldMissing_WithAdd(t *testing.T) {
	line := `{"level":"info"}`
	out, err := reformatTimestamp(line, "timestamp", "", "2006-01-02", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := obj["timestamp"]; !ok {
		t.Error("expected timestamp field to be added")
	}
}

func TestReformatTimestamp_Reformats_RFC3339(t *testing.T) {
	line := `{"timestamp":"2024-03-15T10:00:00Z","msg":"ok"}`
	out, err := reformatTimestamp(line, "timestamp", "", "2006-01-02", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["timestamp"] != "2024-03-15" {
		t.Errorf("expected '2024-03-15', got %q", obj["timestamp"])
	}
}

func TestReformatTimestamp_WithExplicitInLayout(t *testing.T) {
	line := `{"ts":"15/03/2024","msg":"test"}`
	out, err := reformatTimestamp(line, "ts", "02/01/2006", "2006-01-02", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["ts"] != "2024-03-15" {
		t.Errorf("expected '2024-03-15', got %q", obj["ts"])
	}
}

func TestReformatTimestamp_NonStringField_Passthrough(t *testing.T) {
	line := `{"timestamp":1234567890}`
	out, err := reformatTimestamp(line, "timestamp", "", "2006-01-02", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != line {
		t.Errorf("expected unchanged line for non-string field, got %q", out)
	}
}

func TestReformatTimestamp_CustomField(t *testing.T) {
	line := `{"time":"2024-06-01T00:00:00Z","level":"debug"}`
	out, err := reformatTimestamp(line, "time", "", "01/02/2006", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["time"] != "06/01/2024" {
		t.Errorf("expected '06/01/2024', got %q", obj["time"])
	}
}
