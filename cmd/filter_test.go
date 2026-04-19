package cmd

import (
	"testing"
	"time"
)

func TestMatchesField(t *testing.T) {
	entry := map[string]interface{}{"level": "error", "service": "api"}

	if !matchesField(entry, "level", "error") {
		t.Error("expected match for level=error")
	}
	if matchesField(entry, "level", "info") {
		t.Error("expected no match for level=info")
	}
	if matchesField(entry, "missing", "value") {
		t.Error("expected no match for missing key")
	}
}

func TestMatchesTime_NoTimestamp(t *testing.T) {
	entry := map[string]interface{}{"msg": "hello"}
	if !matchesTime(entry, time.Time{}, time.Time{}) {
		t.Error("entry with no timestamp should pass through")
	}
}

func TestMatchesTime_WithinRange(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	entry := map[string]interface{}{"time": now.Format(time.RFC3339)}

	from := now.Add(-time.Minute)
	to := now.Add(time.Minute)

	if !matchesTime(entry, from, to) {
		t.Error("expected entry within range to match")
	}
}

func TestMatchesTime_OutOfRange(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	old := now.Add(-2 * time.Hour)
	entry := map[string]interface{}{"time": old.Format(time.RFC3339)}

	from := now.Add(-time.Minute)
	to := now.Add(time.Minute)

	if matchesTime(entry, from, to) {
		t.Error("expected old entry to be filtered out")
	}
}
