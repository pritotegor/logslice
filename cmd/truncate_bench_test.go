package cmd

import (
	"fmt"
	"testing"
)

func buildTruncateLines(n int, msgLen int) []string {
	lines := make([]string, n)
	msg := ""
	for i := 0; i < msgLen; i++ {
		msg += "a"
	}
	for i := 0; i < n; i++ {
		lines[i] = fmt.Sprintf(`{"msg":%q,"level":"info","index":%d}`, msg, i)
	}
	return lines
}

func BenchmarkTruncateLine_AllFields_1000(b *testing.B) {
	lines := buildTruncateLines(1000, 200)
	fields := map[string]bool{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, line := range lines {
			_, _ = truncateLine(line, fields, 50, "...")
		}
	}
}

func BenchmarkTruncateLine_SpecificField_1000(b *testing.B) {
	lines := buildTruncateLines(1000, 200)
	fields := map[string]bool{"msg": true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, line := range lines {
			_, _ = truncateLine(line, fields, 50, "...")
		}
	}
}

func BenchmarkTruncateLine_NoTruncation_1000(b *testing.B) {
	lines := buildTruncateLines(1000, 10)
	fields := map[string]bool{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, line := range lines {
			_, _ = truncateLine(line, fields, 100, "...")
		}
	}
}

func BenchmarkTruncateLine_AllFields_10000(b *testing.B) {
	lines := buildTruncateLines(10000, 200)
	fields := map[string]bool{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, line := range lines {
			_, _ = truncateLine(line, fields, 50, "...")
		}
	}
}
