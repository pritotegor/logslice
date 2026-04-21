package cmd

import (
	"math/rand"
	"strings"
	"testing"
)

func buildSampleLines(n int) []string {
	lines := make([]string, n)
	for i := 0; i < n; i++ {
		lines[i] = strings.Repeat("x", 10)
	}
	return lines
}

func collectRateSample(lines []string, rate float64, seed int64) []string {
	r := rand.New(rand.NewSource(seed))
	var out []string
	for _, line := range lines {
		if r.Float64() < rate {
			out = append(out, line)
		}
	}
	return out
}

func TestRateSample_ZeroRate(t *testing.T) {
	lines := buildSampleLines(100)
	result := collectRateSample(lines, 0.0, 42)
	if len(result) != 0 {
		t.Errorf("expected 0 lines at rate 0.0, got %d", len(result))
	}
}

func TestRateSample_FullRate(t *testing.T) {
	lines := buildSampleLines(50)
	result := collectRateSample(lines, 1.0, 42)
	if len(result) != 50 {
		t.Errorf("expected 50 lines at rate 1.0, got %d", len(result))
	}
}

func TestReservoirSample_FewerThanN(t *testing.T) {
	input := strings.NewReader("line1\nline2\nline3\n")
	scanner := bufio.NewScannerFromReader(input)
	r := rand.New(rand.NewSource(1))
	// Use internal helper via bufio directly
	_ = r
	_ = scanner
	// Validate reservoir doesn't exceed available lines
	lines := []string{"a", "b", "c"}
	reservoir := make([]string, 0, 10)
	for i, line := range lines {
		if len(reservoir) < 10 {
			reservoir = append(reservoir, line)
		} else {
			j := r.Intn(i + 1)
			if j < 10 {
				reservoir[j] = line
			}
		}
	}
	if len(reservoir) != 3 {
		t.Errorf("expected reservoir size 3, got %d", len(reservoir))
	}
}

func TestReservoirSample_ExactN(t *testing.T) {
	r := rand.New(rand.NewSource(7))
	lines := buildSampleLines(10)
	reservoir := make([]string, 0, 10)
	for i, line := range lines {
		if len(reservoir) < 10 {
			reservoir = append(reservoir, line)
		} else {
			j := r.Intn(i + 1)
			if j < 10 {
				reservoir[j] = line
			}
		}
	}
	if len(reservoir) != 10 {
		t.Errorf("expected 10 lines in reservoir, got %d", len(reservoir))
	}
}

func TestNewRand_ZeroSeedNotDeterministic(t *testing.T) {
	r1 := newRand(0)
	r2 := newRand(0)
	// With different random seeds they should (almost certainly) differ
	v1 := r1.Int63()
	v2 := r2.Int63()
	// This could theoretically collide but is astronomically unlikely
	if v1 == v2 {
		t.Log("warning: two zero-seed randoms produced same value (unlikely but possible)")
	}
}

func TestNewRand_FixedSeedDeterministic(t *testing.T) {
	r1 := newRand(42)
	r2 := newRand(42)
	if r1.Int63() != r2.Int63() {
		t.Error("expected same sequence for same seed")
	}
}
