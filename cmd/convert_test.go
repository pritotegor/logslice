package cmd

import (
	"encoding/json"
	"testing"
)

func TestConvertParseLine_JSON(t *testing.T) {
	line := `{"level":"info","msg":"hello"}`
	m, err := convertParseLine(line, "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
	if m["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", m["msg"])
	}
}

func TestConvertParseLine_InvalidJSON(t *testing.T) {
	_, err := convertParseLine("not json", "json")
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestConvertParseLine_Raw(t *testing.T) {
	line := "some raw log line"
	m, err := convertParseLine(line, "raw")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["message"] != line {
		t.Errorf("expected message=%q, got %v", line, m["message"])
	}
}

func TestConvertParseLine_CSV(t *testing.T) {
	line := "foo,bar,baz"
	m, err := convertParseLine(line, "csv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["field1"] != "foo" {
		t.Errorf("expected field1=foo, got %v", m["field1"])
	}
	if m["field3"] != "baz" {
		t.Errorf("expected field3=baz, got %v", m["field3"])
	}
}

func TestParseCSVLine_QuotedFields(t *testing.T) {
	line := `"hello, world",42,true`
	m, err := parseCSVLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["field1"] != "hello, world" {
		t.Errorf("expected 'hello, world', got %v", m["field1"])
	}
}

func TestConvertParseLine_JSONRoundtrip(t *testing.T) {
	orig := map[string]interface{}{"a": "1", "b": float64(2)}
	bytes, _ := json.Marshal(orig)
	m, err := convertParseLine(string(bytes), "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["a"] != "1" {
		t.Errorf("expected a=1, got %v", m["a"])
	}
	if m["b"] != float64(2) {
		t.Errorf("expected b=2, got %v", m["b"])
	}
}
