package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var headLines int

var headCmd = &cobra.Command{
	Use:   "head [file]",
	Short: "Output the first N lines of a log file",
	Args:  cobra.MaximumNArgs(1),
	RunE:  RunHead,
}

func init() {
	headCmd.Flags().IntVarP(&headLines, "lines", "n", 10, "Number of lines to show from the beginning")
	rootCmd.AddCommand(headCmd)
}

func RunHead(cmd *cobra.Command, args []string) error {
	r, err := openInput(args)
	if err != nil {
		return err
	}
	if rc, ok := r.(io.Closer); ok {
		defer rc.Close()
	}

	scanner := bufio.NewScanner(r)
	count := 0
	for scanner.Scan() && count < headLines {
		fmt.Fprintln(os.Stdout, scanner.Text())
		count++
	}
	return scanner.Err()
}
