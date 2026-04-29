package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var limitCmd = &cobra.Command{
	Use:   "limit",
	Short: "Output at most N lines, with optional offset",
	Long:  `Reads lines from a log file (or stdin) and outputs at most --count lines, skipping the first --offset lines.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		count, _ := cmd.Flags().GetInt("count")
		offset, _ := cmd.Flags().GetInt("offset")
		input, _ := cmd.Flags().GetString("input")

		r, err := openInput(input)
		if err != nil {
			return fmt.Errorf("opening input: %w", err)
		}
		if closer, ok := r.(io.Closer); ok {
			defer closer.Close()
		}

		return RunLimit(r, os.Stdout, offset, count)
	},
}

func init() {
	limitCmd.Flags().IntP("count", "n", 100, "Maximum number of lines to output")
	limitCmd.Flags().IntP("offset", "o", 0, "Number of lines to skip before outputting")
	limitCmd.Flags().StringP("input", "i", "", "Input file (default: stdin)")
	rootCmd.AddCommand(limitCmd)
}

// RunLimit reads from r, skips the first offset lines, then writes at most count lines to w.
func RunLimit(r io.Reader, w io.Writer, offset, count int) error {
	scanner := bufio.NewScanner(r)
	skipped := 0
	written := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if skipped < offset {
			skipped++
			continue
		}
		if count >= 0 && written >= count {
			break
		}
		fmt.Fprintln(w, line)
		written++
	}
	return scanner.Err()
}
