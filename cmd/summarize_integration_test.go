package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func writeTempSummarizeFile(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp("", "summarize-*.log")
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func captureSummarizeOutput(t *testing.T, args []string) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := summarizeCmd
	cmd.ResetFlags()
	init_summarize(cmd)
	cmd.SetArgs(args)
	cmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func init_summarize(cmd interface{}) {
	// flags already registered via init()
}

func TestRunSummarize_BasicOutput(t *testing.T) {
	lines := []string{
		`{"duration":100,"size":512}`,
		`{"duration":200,"size":1024}`,
		`{"duration":300,"size":256}`,
	}
	file := writeTempSummarizeFile(t, lines)

	r := strings.NewReader(strings.Join(lines, "\n") + "\n")
	result := collectSummary(r, []string{"duration", "size"}, "")

	d := result["_all"]["duration"]
	if d == nil || d.Count != 3 {
		t.Fatalf("expected 3 duration entries, got %v", d)
	}
	if d.Sum != 600 {
		t.Errorf("expected sum 600, got %f", d.Sum)
	}
	_ = file
}

func TestRunSummarize_GroupBy_Output(t *testing.T) {
	lines := []string{
		`{"env":"prod","cpu":80}`,
		`{"env":"dev","cpu":20}`,
		`{"env":"prod","cpu":60}`,
	}
	r := strings.NewReader(strings.Join(lines, "\n") + "\n")
	result := collectSummary(r, []string{"cpu"}, "env")

	prod := result["prod"]["cpu"]
	if prod.Count != 2 || prod.Sum != 140 {
		t.Errorf("prod: count=%d sum=%f", prod.Count, prod.Sum)
	}
	dev := result["dev"]["cpu"]
	if dev.Count != 1 || dev.Sum != 20 {
		t.Errorf("dev: count=%d sum=%f", dev.Count, dev.Sum)
	}
}
