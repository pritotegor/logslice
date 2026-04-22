package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	grepPattern   string
	grepField     string
	grepInvert    bool
	grepIgnoreCase bool
)

func init() {
	grepCmd := &cobra.Command{
		Use:   "grep",
		Short: "Search log lines by regex pattern",
		Long:  `Filter log lines matching a regular expression, optionally scoped to a specific JSON field.`,
		RunE:  RunGrep,
	}

	grepCmd.Flags().StringVarP(&grepPattern, "pattern", "p", "", "Regex pattern to match (required)")
	grepCmd.Flags().StringVarP(&grepField, "field", "f", "", "JSON field to search within (searches full line if omitted)")
	grepCmd.Flags().BoolVarP(&grepInvert, "invert", "v", false, "Invert match: output lines that do NOT match")
	grepCmd.Flags().BoolVarP(&grepIgnoreCase, "ignore-case", "i", false, "Case-insensitive matching")
	_ = grepCmd.MarkFlagRequired("pattern")

	rootCmd.AddCommand(grepCmd)
}

func RunGrep(cmd *cobra.Command, args []string) error {
	pattern := grepPattern
	if grepIgnoreCase {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	input, err := openInput(args)
	if err != nil {
		return err
	}
	if closer, ok := input.(interface{ Close() error }); ok {
		defer closer.Close()
	}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		matched := grepMatch(re, line, grepField)
		if grepInvert {
			matched = !matched
		}
		if matched {
			fmt.Fprintln(os.Stdout, line)
		}
	}
	return scanner.Err()
}

func grepMatch(re *regexp.Regexp, line, field string) bool {
	if field == "" {
		return re.MatchString(line)
	}
	fields, err := ParseLine(line)
	if err != nil || fields == nil {
		return false
	}
	val, ok := fields[field]
	if !ok {
		return false
	}
	return re.MatchString(fmt.Sprintf("%v", val))
}
