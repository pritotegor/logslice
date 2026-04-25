package cmd

import (
	"strings"
	"testing"
)

func TestDiffLines_IdenticalJSON(t *testing.T) {
	l := `{"level":"info","msg":"ok"}`
	diffs := diffLines(l, l, nil)
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %v", diffs)
	}
}

func TestDiffLines_DifferentFieldValue(t *testing.T) {
	l1 := `{"level":"info","msg":"ok"}`
	l2 := `{"level":"error","msg":"ok"}`
	diffs := diffLines(l1, l2, nil)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d: %v", len(diffs), diffs)
	}
	if !strings.Contains(diffs[0], "level") {
		t.Errorf("expected diff to mention 'level', got %q", diffs[0])
	}
}

func TestDiffLines_MissingFieldInSecond(t *testing.T) {
	l1 := `{"level":"info","trace":"abc"}`
	l2 := `{"level":"info"}`
	diffs := diffLines(l1, l2, nil)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if !strings.Contains(diffs[0], "<missing>") {
		t.Errorf("expected <missing> marker, got %q", diffs[0])
	}
}

func TestDiffLines_FilteredFields(t *testing.T) {
	l1 := `{"level":"info","msg":"a","code":200}`
	l2 := `{"level":"error","msg":"b","code":200}`
	// only compare "code" — should see no diffs
	diffs := diffLines(l1, l2, []string{"code"})
	if len(diffs) != 0 {
		t.Errorf("expected 0 diffs when filtering on equal field, got %v", diffs)
	}
}

func TestDiffLines_NonJSONFallback(t *testing.T) {
	diffs := diffLines("hello world", "hello earth", nil)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff for raw lines, got %d", len(diffs))
	}
}

func TestDiffLines_NonJSONIdentical(t *testing.T) {
	diffs := diffLines("same line", "same line", nil)
	if len(diffs) != 0 {
		t.Errorf("expected 0 diffs for identical raw lines, got %v", diffs)
	}
}

func TestDiffReaders_OnlyChanged(t *testing.T) {
	r1 := strings.NewReader("{\"a\":1}\n{\"a\":2}\n")
	r2 := strings.NewReader("{\"a\":1}\n{\"a\":9}\n")
	var out strings.Builder
	err := diffReaders(&out, r1, r2, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	result := out.String()
	if !strings.Contains(result, "line 2") {
		t.Errorf("expected diff on line 2, got: %s", result)
	}
	if strings.Contains(result, "line 1") {
		t.Errorf("line 1 should be identical, not reported")
	}
}

func TestDiffReaders_UnequalLength(t *testing.T) {
	r1 := strings.NewReader("{\"a\":1}\n")
	r2 := strings.NewReader("{\"a\":1}\n{\"a\":2}\n")
	var out strings.Builder
	err := diffReaders(&out, r1, r2, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	// line 2 exists only in r2; should produce a diff
	if !strings.Contains(out.String(), "line 2") {
		t.Errorf("expected diff for extra line in r2, got: %s", out.String())
	}
}
