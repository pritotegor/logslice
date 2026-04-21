package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func BenchmarkDedupeLines_NoField(b *testing.B) {
	lines := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		lines[i] = fmt.Sprintf(`{"msg":"event","id":%d}`, i%100)
	}
	input := strings.Join(lines, "\n") + "\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out bytes.Buffer
		_ = dedupeLines(strings.NewReader(input), &out, "", 0)
	}
}

func BenchmarkDedupeLines_ByField(b *testing.B) {
	lines := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		lines[i] = fmt.Sprintf(`{"msg":"event %d","id":%d}`, i, i%100)
	}
	input := strings.Join(lines, "\n") + "\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out bytes.Buffer
		_ = dedupeLines(strings.NewReader(input), &out, "id", 0)
	}
}

func BenchmarkDedupeLines_WithWindow(b *testing.B) {
	lines := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		lines[i] = fmt.Sprintf(`{"msg":"event","id":%d}`, i%50)
	}
	input := strings.Join(lines, "\n") + "\n"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out bytes.Buffer
		_ = dedupeLines(strings.NewReader(input), &out, "", 100)
	}
}

func BenchmarkDedupeKey_FullLine(b *testing.B) {
	line := `{"level":"info","msg":"something happened","ts":"2024-01-01T00:00:00Z"}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dedupeKey(line, "")
	}
}

func BenchmarkDedupeKey_ByField(b *testing.B) {
	line := `{"level":"info","msg":"something happened","ts":"2024-01-01T00:00:00Z"}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dedupeKey(line, "level")
	}
}
