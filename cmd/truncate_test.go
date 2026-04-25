package cmd

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTruncateLine_AllStringFields(t *testing.T) {
	line := `{"msg":"hello world","level":"info"}`
	out, err := truncateLine(line, map[string]bool{}, 5, "...")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["msg"] != "hello..." {
		t.Errorf("expected 'hello...', got %q", obj["msg"])
	}
	if obj["level"] != "info" {
		t.Errorf("expected 'info' (not truncated), got %q", obj["level"])
	}
}

func TestTruncateLine_SpecificField(t *testing.T) {
	line := `{"msg":"this is a long message","level":"warning"}`
	out, err := truncateLine(line, map[string]bool{"msg": true}, 7, "~")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["msg"] != "this is~" {
		t.Errorf("expected 'this is~', got %q", obj["msg"])
	}
	if obj["level"] != "warning" {
		t.Errorf("level should be unchanged, got %q", obj["level"])
	}
}

func TestTruncateLine_NonStringFieldUnchanged(t *testing.T) {
	line := `{"count":42,"active":true,"msg":"short"}`
	out, err := truncateLine(line, map[string]bool{}, 3, "...")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["count"] != float64(42) {
		t.Errorf("numeric field should be unchanged")
	}
	if obj["active"] != true {
		t.Errorf("bool field should be unchanged")
	}
}

func TestTruncateLine_InvalidJSON(t *testing.T) {
	_, err := truncateLine("not json", map[string]bool{}, 10, "...")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestTruncateLine_ExactLength_NoTruncation(t *testing.T) {
	line := `{"msg":"hello"}`
	out, err := truncateLine(line, map[string]bool{}, 5, "...")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	json.Unmarshal([]byte(out), &obj)
	if obj["msg"] != "hello" {
		t.Errorf("value at exact max length should not be truncated, got %q", obj["msg"])
	}
}

func TestTruncateLine_EmptySuffix(t *testing.T) {
	line := `{"msg":"abcdefgh"}`
	out, err := truncateLine(line, map[string]bool{}, 4, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"abcd"`) {
		t.Errorf("expected truncated value 'abcd', got %s", out)
	}
}
