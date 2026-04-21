package cmd

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var dedupeCmd = &cobra.Command{
	Use:   "dedupe [file]",
	Short: "Remove duplicate log lines based on a field or full content",
	RunE:  RunDedupe,
}

var dedupeField string
var dedupeWindow int

func init() {
	dedupeCmd.Flags().StringVarP(&dedupeField, "field", "f", "", "JSON field to deduplicate on (default: full line)")
	dedupeCmd.Flags().IntVarP(&dedupeWindow, "window", "w", 0, "Only deduplicate within a sliding window of N lines (0 = unlimited)")
	rootCmd.AddCommand(dedupeCmd)
}

func RunDedupe(cmd *cobra.Command, args []string) error {
	var r io.Reader
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
		defer f.Close()
		r = f
	} else {
		r = os.Stdin
	}

	return dedupeLines(r, os.Stdout, dedupeField, dedupeWindow)
}

func dedupeLines(r io.Reader, w io.Writer, field string, window int) error {
	scanner := bufio.NewScanner(r)
	seen := make(map[string]struct{})
	var order []string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		key := dedupeKey(line, field)

		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		order = append(order, key)

		if window > 0 && len(order) > window {
			oldKey := order[0]
			order = order[1:]
			delete(seen, oldKey)
		}

		fmt.Fprintln(w, line)
	}

	return scanner.Err()
}

func dedupeKey(line, field string) string {
	if field != "" {
		fields, err := ParseLine(line)
		if err == nil && fields != nil {
			if val, ok := fields[field]; ok {
				return fmt.Sprintf("%v", val)
			}
		}
	}
	h := md5.Sum([]byte(line))
	return fmt.Sprintf("%x", h)
}
