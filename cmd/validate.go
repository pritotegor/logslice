package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate that each line is well-formed JSON",
	Long:  `Reads log lines and reports any that are not valid JSON. Exits with code 1 if any invalid lines are found.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunValidate(cmd, args)
	},
}

func init() {
	validateCmd.Flags().BoolP("strict", "s", false, "Exit immediately on first invalid line")
	validateCmd.Flags().BoolP("quiet", "q", false, "Suppress per-line error output")
	rootCmd.AddCommand(validateCmd)
}

func RunValidate(cmd *cobra.Command, args []string) error {
	strict, _ := cmd.Flags().GetBool("strict")
	quiet, _ := cmd.Flags().GetBool("quiet")

	reader, err := openInput(args)
	if err != nil {
		return err
	}
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	invalidCount, totalCount, err := validateLines(reader, os.Stderr, strict, quiet)
	if err != nil {
		return err
	}

	if !quiet {
		fmt.Fprintf(os.Stderr, "Validated %d lines: %d invalid\n", totalCount, invalidCount)
	}

	if invalidCount > 0 {
		return fmt.Errorf("%d invalid line(s) found", invalidCount)
	}
	return nil
}

func validateLines(r io.Reader, errOut io.Writer, strict, quiet bool) (invalid, total int, err error) {
	scanner := bufio.NewScanner(r)
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		lineNum++
		total++
		if !json.Valid([]byte(line)) {
			invalid++
			if !quiet {
				fmt.Fprintf(errOut, "line %d: invalid JSON: %s\n", lineNum, line)
			}
			if strict {
				return invalid, total, nil
			}
		}
	}
	return invalid, total, scanner.Err()
}
