package cmd

import (
	"testing"
)

func TestNewLogStats_InitialState(t *testing.T) {
	s := NewLogStats()
	if s.TotalLines != 0 {
		t.Errorf("expected TotalLines=0, got %d", s.TotalLines)
	}
	if s.ParsedLines != 0 {
		t.Errorf("expected ParsedLines=0, got %d", s.ParsedLines)
	}
	if len(s.FieldCounts) != 0 {
		t.Errorf("expected empty FieldCounts")
	}
}

func TestRecord_NilFields(t *testing.T) {
	s := NewLogStats()
	s.Record(nil)
	if s.TotalLines != 1 {
		t.Errorf("expected TotalLines=1, got %d", s.TotalLines)
	}
	if s.ParsedLines != 0 {
		t.Errorf("expected ParsedLines=0 for nil fields, got %d", s.ParsedLines)
	}
}

func TestRecord_WithFields(t *testing.T) {
	s := NewLogStats()
	s.Record(map[string]interface{}{"level": "info", "msg": "started"})
	s.Record(map[string]interface{}{"level": "error", "msg": "failed"})
	s.Record(map[string]interface{}{"level": "info", "msg": "retrying"})

	if s.TotalLines != 3 {
		t.Errorf("expected TotalLines=3, got %d", s.TotalLines)
	}
	if s.ParsedLines != 3 {
		t.Errorf("expected ParsedLines=3, got %d", s.ParsedLines)
	}

	lvl, ok := s.FieldCounts["level"]
	if !ok {
		t.Fatal("expected 'level' field in FieldCounts")
	}
	if lvl.Count != 3 {
		t.Errorf("expected level count=3, got %d", lvl.Count)
	}
	if lvl.Values["info"] != 2 {
		t.Errorf("expected info=2, got %d", lvl.Values["info"])
	}
	if lvl.Values["error"] != 1 {
		t.Errorf("expected error=1, got %d", lvl.Values["error"])
	}
}

func TestRecord_UniqueValues(t *testing.T) {
	s := NewLogStats()
	for i := 0; i < 5; i++ {
		s.Record(map[string]interface{}{"host": "server-1"})
	}
	s.Record(map[string]interface{}{"host": "server-2"})

	h := s.FieldCounts["host"]
	if h.Count != 6 {
		t.Errorf("expected count=6, got %d", h.Count)
	}
	if len(h.Values) != 2 {
		t.Errorf("expected 2 unique host values, got %d", len(h.Values))
	}
}
