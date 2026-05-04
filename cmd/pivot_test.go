package cmd

import (
	"strings"
	"testing"
)

func TestCollectPivot_CountMode(t *testing.T) {
	input := strings.NewReader(
		`{"env":"prod","level":"error"}` + "\n" +
			`{"env":"prod","level":"info"}` + "\n" +
			`{"env":"dev","level":"error"}` + "\n" +
			`{"env":"prod","level":"error"}` + "\n",
	)
	table, rowOrder, colSet, err := collectPivot(input, "env", "level", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rowOrder) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rowOrder))
	}
	if _, ok := colSet["error"]; !ok {
		t.Error("expected 'error' column")
	}
	if table["prod"]["error"] != 2 {
		t.Errorf("expected prod/error=2, got %g", table["prod"]["error"])
	}
	if table["prod"]["info"] != 1 {
		t.Errorf("expected prod/info=1, got %g", table["prod"]["info"])
	}
	if table["dev"]["error"] != 1 {
		t.Errorf("expected dev/error=1, got %g", table["dev"]["error"])
	}
}

func TestCollectPivot_SumMode(t *testing.T) {
	input := strings.NewReader(
		`{"svc":"api","region":"us","latency":10}` + "\n" +
			`{"svc":"api","region":"eu","latency":20}` + "\n" +
			`{"svc":"api","region":"us","latency":5}` + "\n",
	)
	table, _, _, err := collectPivot(input, "svc", "region", "latency")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if table["api"]["us"] != 15 {
		t.Errorf("expected api/us=15, got %g", table["api"]["us"])
	}
	if table["api"]["eu"] != 20 {
		t.Errorf("expected api/eu=20, got %g", table["api"]["eu"])
	}
}

func TestCollectPivot_SkipsNonJSON(t *testing.T) {
	input := strings.NewReader(
		"not json\n" +
			`{"env":"prod","level":"info"}` + "\n",
	)
	_, rowOrder, _, err := collectPivot(input, "env", "level", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rowOrder) != 1 {
		t.Errorf("expected 1 row, got %d", len(rowOrder))
	}
}

func TestCollectPivot_MissingRowField(t *testing.T) {
	input := strings.NewReader(
		`{"level":"info"}` + "\n",
	)
	_, rowOrder, _, err := collectPivot(input, "env", "level", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rowOrder) != 0 {
		t.Errorf("expected 0 rows, got %d", len(rowOrder))
	}
}

func TestSortedKeys_Order(t *testing.T) {
	m := map[string]struct{}{"banana": {}, "apple": {}, "cherry": {}}
	keys := sortedKeys(m)
	expected := []string{"apple", "banana", "cherry"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %s got %s", i, expected[i], k)
		}
	}
}
