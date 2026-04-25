package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <file1> <file2>",
	Short: "Compare two structured log files and show differing fields",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunDiff(cmd, args[0], args[1])
	},
}

func init() {
	diffCmd.Flags().StringSliceP("fields", "f", nil, "fields to compare (default: all)")
	diffCmd.Flags().BoolP("only-changed", "c", false, "only output lines where fields differ")
	rootCmd.AddCommand(diffCmd)
}

func RunDiff(cmd *cobra.Command, path1, path2 string) error {
	fields, _ := cmd.Flags().GetStringSlice("fields")
	onlyChanged, _ := cmd.Flags().GetBool("only-changed")

	f1, err := os.Open(path1)
	if err != nil {
		return fmt.Errorf("opening %s: %w", path1, err)
	}
	defer f1.Close()

	f2, err := os.Open(path2)
	if err != nil {
		return fmt.Errorf("opening %s: %w", path2, err)
	}
	defer f2.Close()

	return diffReaders(cmd.OutOrStdout(), f1, f2, fields, onlyChanged)
}

func diffReaders(w io.Writer, r1, r2 io.Reader, fields []string, onlyChanged bool) error {
	sc1 := bufio.NewScanner(r1)
	sc2 := bufio.NewScanner(r2)

	lineNum := 0
	for {
		has1 := sc1.Scan()
		has2 := sc2.Scan()
		if !has1 && !has2 {
			break
		}
		lineNum++

		l1 := ""
		if has1 {
			l1 = sc1.Text()
		}
		l2 := ""
		if has2 {
			l2 = sc2.Text()
		}

		diffs := diffLines(l1, l2, fields)
		if onlyChanged && len(diffs) == 0 {
			continue
		}
		if len(diffs) > 0 {
			for _, d := range diffs {
				fmt.Fprintf(w, "line %d | %s\n", lineNum, d)
			}
		}
	}
	return nil
}

func diffLines(line1, line2 string, fields []string) []string {
	var m1, m2 map[string]interface{}
	if err := json.Unmarshal([]byte(line1), &m1); err != nil {
		if line1 != line2 {
			return []string{fmt.Sprintf("< %s | > %s", line1, line2)}
		}
		return nil
	}
	if err := json.Unmarshal([]byte(line2), &m2); err != nil {
		return []string{fmt.Sprintf("< %s | > %s", line1, line2)}
	}

	keys := fields
	if len(keys) == 0 {
		seen := map[string]bool{}
		for k := range m1 {
			seen[k] = true
		}
		for k := range m2 {
			seen[k] = true
		}
		for k := range seen {
			keys = append(keys, k)
		}
	}

	var diffs []string
	for _, k := range keys {
		v1, ok1 := m1[k]
		v2, ok2 := m2[k]
		switch {
		case !ok1 && !ok2:
			continue
		case !ok1:
			diffs = append(diffs, fmt.Sprintf("field %q: <missing> | %v", k, v2))
		case !ok2:
			diffs = append(diffs, fmt.Sprintf("field %q: %v | <missing>", k, v1))
		case fmt.Sprintf("%v", v1) != fmt.Sprintf("%v", v2):
			diffs = append(diffs, fmt.Sprintf("field %q: %v | %v", k, v1, v2))
		}
	}
	return diffs
}
