package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var pivotCmd = &cobra.Command{
	Use:   "pivot",
	Short: "Pivot a field's values into columns, counting or aggregating occurrences",
	RunE: func(cmd *cobra.Command, args []string) error {
		rowField, _ := cmd.Flags().GetString("row")
		colField, _ := cmd.Flags().GetString("col")
		valField, _ := cmd.Flags().GetString("val")
		input, _ := cmd.Flags().GetString("input")

		if rowField == "" || colField == "" {
			return fmt.Errorf("--row and --col are required")
		}

		r, err := openInput(input)
		if err != nil {
			return err
		}
		if closer, ok := r.(io.Closer); ok {
			defer closer.Close()
		}

		return RunPivot(r, os.Stdout, rowField, colField, valField)
	},
}

func init() {
	pivotCmd.Flags().String("row", "", "Field to use as row key")
	pivotCmd.Flags().String("col", "", "Field whose values become columns")
	pivotCmd.Flags().String("val", "", "Field to sum per cell (omit to count)")
	pivotCmd.Flags().StringP("input", "i", "", "Input file (default: stdin)")
	rootCmd.AddCommand(pivotCmd)
}

// RunPivot reads JSON lines and produces a pivot table written to w.
func RunPivot(r io.Reader, w io.Writer, rowField, colField, valField string) error {
	table, rowOrder, colSet, err := collectPivot(r, rowField, colField, valField)
	if err != nil {
		return err
	}

	cols := sortedKeys(colSet)

	// header
	header := append([]string{rowField}, cols...)
	fmt.Fprintln(w, strings.Join(header, "\t"))

	for _, row := range rowOrder {
		parts := []string{row}
		for _, col := range cols {
			v := table[row][col]
			parts = append(parts, fmt.Sprintf("%g", v))
		}
		fmt.Fprintln(w, strings.Join(parts, "\t"))
	}
	return nil
}

func sortedKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
