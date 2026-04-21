package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var tailLines int
var tailFollow bool

var tailCmd = &cobra.Command{
	Use:   "tail [file]",
	Short: "Output the last N lines of a log file, optionally following",
	Args:  cobra.MaximumNArgs(1),
	RunE:  RunTail,
}

func init() {
	tailCmd.Flags().IntVarP(&tailLines, "lines", "n", 10, "Number of lines to show from the end")
	tailCmd.Flags().BoolVarP(&tailFollow, "follow", "f", false, "Follow the file for new lines (like tail -f)")
	rootCmd.AddCommand(tailCmd)
}

func RunTail(cmd *cobra.Command, args []string) error {
	r, err := openInput(args)
	if err != nil {
		return err
	}
	if rc, ok := r.(io.Closer); ok {
		defer rc.Close()
	}

	lines, err := readLastN(r, tailLines)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	for _, line := range lines {
		fmt.Fprintln(os.Stdout, line)
	}

	if tailFollow {
		if len(args) == 0 {
			return fmt.Errorf("--follow requires a file argument")
		}
		return followFile(args[0], os.Stdout)
	}
	return nil
}
