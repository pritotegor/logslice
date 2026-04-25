package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func captureTruncateOutput(t *testing.T, args []string, input string) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := &cobra.Command{}
	truncateCmd.ResetFlags()
	truncateCmd.Flags().StringSliceP("fields", "f", nil, "")
	truncateCmd.Flags().IntP("max-len", "n", 100, "")
	truncateCmd.Flags().StringP("suffix", "s", "...", "")
	truncateCmd.Flags().StringP("input", "i", "", "")
	_ = cmd

	tf, _ := os.CreateTemp("", "truncate-*.log")
	tf.WriteString(input)
	tf.Close()
	defer os.Remove(tf.Name())

	fullArgs := append([]string{"--input", tf.Name()}, args...)
	truncateCmd.SetArgs(fullArgs)
	truncateCmd.Execute()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestRunTruncate_BasicOutput(t *testing.T) {
	input := fmt.Sprintf("%s\n%s\n",
		`{"msg":"this is a long message","level":"info"}`,
		`{"msg":"short","level":"debug"}`,
	)
	out := captureTruncateOutput(t, []string{"--max-len", "7", "--suffix", "~"}, input)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(lines))
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &obj); err != nil {
		t.Fatalf("line 1 not valid JSON: %v", err)
	}
	if obj["msg"] != "this is~" {
		t.Errorf("expected truncated msg, got %q", obj["msg"])
	}
}

func TestRunTruncate_SkipsEmptyLines(t *testing.T) {
	input := `{"msg":"hello world"}` + "\n\n" + `{"msg":"bye"}` + "\n"
	out := captureTruncateOutput(t, []string{"--max-len", "5"}, input)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 non-empty output lines, got %d", len(lines))
	}
}

func TestRunTruncate_SpecificFieldOnly(t *testing.T) {
	input := `{"msg":"this is long","tag":"this is also long"}` + "\n"
	out := captureTruncateOutput(t, []string{"--fields", "msg", "--max-len", "4", "--suffix", ""}, input)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &obj); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if obj["msg"] != "this" {
		t.Errorf("expected msg='this', got %q", obj["msg"])
	}
	if obj["tag"] != "this is also long" {
		t.Errorf("tag should be unchanged, got %q", obj["tag"])
	}
}
