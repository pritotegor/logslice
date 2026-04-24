package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
)

var summarizeCmd = &cobra.Command{
	Use:   "summarize [file]",
	Short: "Summarize numeric field statistics across log lines",
	RunE:  RunSummarize,
}

func init() {
	summarizeCmd.Flags().StringSliceP("fields", "f", nil, "Numeric fields to summarize (required)")
	summarizeCmd.Flags().StringP("group-by", "g", "", "Group results by this field")
	rootCmd.AddCommand(summarizeCmd)
}

func RunSummarize(cmd *cobra.Command, args []string) error {
	fields, _ := cmd.Flags().GetStringSlice("fields")
	groupBy, _ := cmd.Flags().GetString("group-by")

	if len(fields) == 0 {
		return fmt.Errorf("at least one --fields value is required")
	}

	r, err := openInput(args)
	if err != nil {
		return err
	}
	if rc, ok := r.(io.Closer); ok {
		defer rc.Close()
	}

	results := collectSummary(r, fields, groupBy)
	printSummary(os.Stdout, results, fields, groupBy)
	return nil
}
