package cmd

import (
	"strings"
	"testing"
)

func TestCollectSummary_BasicFields(t *testing.T) {
	input := strings.NewReader(
		`{"duration":10,"status":200}` + "\n" +
			`{"duration":20,"status":404}` + "\n" +
			`{"duration":30,"status":200}` + "\n",
	)
	result := collectSummary(input, []string{"duration"}, "")
	s, ok := result["_all"]["duration"]
	if !ok {
		t.Fatal("expected stats for 'duration'")
	}
	if s.Count != 3 {
		t.Errorf("expected count 3, got %d", s.Count)
	}
	if s.Sum != 60 {
		t.Errorf("expected sum 60, got %f", s.Sum)
	}
	if s.Min != 10 {
		t.Errorf("expected min 10, got %f", s.Min)
	}
	if s.Max != 30 {
		t.Errorf("expected max 30, got %f", s.Max)
	}
}

func TestCollectSummary_GroupBy(t *testing.T) {
	input := strings.NewReader(
		`{"svc":"api","latency":5}` + "\n" +
			`{"svc":"db","latency":15}` + "\n" +
			`{"svc":"api","latency":10}` + "\n",
	)
	result := collectSummary(input, []string{"latency"}, "svc")
	if len(result) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(result))
	}
	api := result["api"]["latency"]
	if api.Count != 2 || api.Sum != 15 {
		t.Errorf("api group unexpected: count=%d sum=%f", api.Count, api.Sum)
	}
	db := result["db"]["latency"]
	if db.Count != 1 || db.Sum != 15 {
		t.Errorf("db group unexpected: count=%d sum=%f", db.Count, db.Sum)
	}
}

func TestCollectSummary_SkipsNonJSON(t *testing.T) {
	input := strings.NewReader("not json\n" + `{"val":5}` + "\n")
	result := collectSummary(input, []string{"val"}, "")
	s := result["_all"]["val"]
	if s.Count != 1 {
		t.Errorf("expected count 1, got %d", s.Count)
	}
}

func TestCollectSummary_MissingField(t *testing.T) {
	input := strings.NewReader(`{"other":1}` + "\n")
	result := collectSummary(input, []string{"val"}, "")
	if _, ok := result["_all"]["val"]; ok {
		t.Error("expected no stats for missing field")
	}
}

func TestFieldStats_Avg(t *testing.T) {
	s := newFieldStats()
	if s.avg() != 0 {
		t.Error("expected avg 0 for empty stats")
	}
	s.record(10)
	s.record(20)
	if s.avg() != 15 {
		t.Errorf("expected avg 15, got %f", s.avg())
	}
}
