package cmd

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func buildPivotLines(n int) string {
	levels := []string{"info", "warn", "error", "debug"}
	envs := []string{"prod", "staging", "dev"}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		env := envs[i%len(envs)]
		level := levels[i%len(levels)]
		fmt.Fprintf(&sb, `{"env":%q,"level":%q,"latency":%d}`+"\n", env, level, i%200)
	}
	return sb.String()
}

func BenchmarkCollectPivot_Count_1000(b *testing.B) {
	data := buildPivotLines(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = collectPivot(strings.NewReader(data), "env", "level", "")
	}
}

func BenchmarkCollectPivot_Sum_1000(b *testing.B) {
	data := buildPivotLines(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = collectPivot(strings.NewReader(data), "env", "level", "latency")
	}
}

func BenchmarkCollectPivot_Count_10000(b *testing.B) {
	data := buildPivotLines(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = collectPivot(strings.NewReader(data), "env", "level", "")
	}
}

func BenchmarkRunPivot_Count_1000(b *testing.B) {
	data := buildPivotLines(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RunPivot(strings.NewReader(data), io.Discard, "env", "level", "")
	}
}
