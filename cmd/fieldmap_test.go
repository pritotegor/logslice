package cmd

import (
	"testing"
)

func TestJsonTypeName_AllTypes(t *testing.T) {
	cases := []struct {
		input    interface{}
		expected string
	}{
		{"hello", "string"},
		{float64(42), "number"},
		{true, "bool"},
		{nil, "null"},
		{map[string]interface{}{"a": 1}, "object"},
		{[]interface{}{1, 2}, "array"},
	}
	for _, tc := range cases {
		got := jsonTypeName(tc.input)
		if got != tc.expected {
			t.Errorf("jsonTypeName(%v) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestJsonTypeName_Unknown(t *testing.T) {
	type custom struct{}
	got := jsonTypeName(custom{})
	if got != "unknown" {
		t.Errorf("expected unknown, got %q", got)
	}
}

func TestFieldInfo_CountAndTypes(t *testing.T) {
	info := &fieldInfo{Types: make(map[string]int)}
	info.Count++
	info.Types["string"]++
	info.Count++
	info.Types["number"]++

	if info.Count != 2 {
		t.Errorf("expected count 2, got %d", info.Count)
	}
	if info.Types["string"] != 1 {
		t.Errorf("expected string count 1, got %d", info.Types["string"])
	}
	if info.Types["number"] != 1 {
		t.Errorf("expected number count 1, got %d", info.Types["number"])
	}
}
