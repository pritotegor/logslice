package cmd

import (
	"encoding/json"
	"testing"
)

func TestCompactLine_NonJSON_Passthrough(t *testing.T) {
	result, err := compactLine("plain text line", nil, false)
	if err == nil {
		t.Fatal("expected error for non-JSON input")
	}
	if result != "plain text line" {
		t.Errorf("expected original line, got %q", result)
	}
}

func TestCompactLine_RemovesFields(t *testing.T) {
	input := `{"level":"info","msg":"hello","secret":"abc"}`
	result, err := compactLine(input, []string{"secret"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(result), &obj); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if _, ok := obj["secret"]; ok {
		t.Error("expected 'secret' field to be removed")
	}
	if obj["level"] != "info" {
		t.Errorf("expected level=info, got %v", obj["level"])
	}
}

func TestCompactLine_RemovesMultipleFields(t *testing.T) {
	input := `{"a":1,"b":2,"c":3}`
	result, err := compactLine(input, []string{"a", "c"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	json.Unmarshal([]byte(result), &obj)
	if _, ok := obj["a"]; ok {
		t.Error("expected 'a' to be removed")
	}
	if _, ok := obj["c"]; ok {
		t.Error("expected 'c' to be removed")
	}
	if obj["b"] == nil {
		t.Error("expected 'b' to remain")
	}
}

func TestCompactLine_DropEmpty_RemovesNullAndEmpty(t *testing.T) {
	input := `{"level":"info","msg":"","err":null,"code":200}`
	result, err := compactLine(input, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	json.Unmarshal([]byte(result), &obj)
	if _, ok := obj["msg"]; ok {
		t.Error("expected empty 'msg' to be dropped")
	}
	if _, ok := obj["err"]; ok {
		t.Error("expected null 'err' to be dropped")
	}
	if obj["level"] != "info" {
		t.Errorf("expected level=info, got %v", obj["level"])
	}
}

func TestCompactLine_DropEmpty_KeepsNonEmpty(t *testing.T) {
	input := `{"level":"warn","count":0}`
	result, err := compactLine(input, nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	json.Unmarshal([]byte(result), &obj)
	if obj["level"] != "warn" {
		t.Errorf("expected level=warn, got %v", obj["level"])
	}
	if obj["count"] == nil {
		t.Error("expected count=0 to be preserved")
	}
}

func TestCompactLine_NoFieldsNoDropEmpty_Unchanged(t *testing.T) {
	input := `{"level":"debug","msg":"ok"}`
	result, err := compactLine(input, nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var obj map[string]interface{}
	json.Unmarshal([]byte(result), &obj)
	if len(obj) != 2 {
		t.Errorf("expected 2 fields, got %d", len(obj))
	}
}
