package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func buildBenchLimitInput(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf(`{"ts":"2024-01-01T00:00:00Z","level":"info","msg":"event %d"}`, i))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func BenchmarkRunLimit_1000_NoOffset(b *testing.B) {
	data := buildBenchLimitInput(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out bytes.Buffer
		_ = RunLimit(strings.NewReader(data), &out, 0, 500)
	}
}

func BenchmarkRunLimit_10000_WithOffset(b *testing.B) {
	data := buildBenchLimitInput(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out bytes.Buffer
		_ = RunLimit(strings.NewReader(data), &out, 100, 200)
	}
}

func BenchmarkRunLimit_AllLines(b *testing.B) {
	data := buildBenchLimitInput(5000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var out bytes.Buffer
		_ = RunLimit(strings.NewReader(data), &out, 0, 5000)
	}
}
