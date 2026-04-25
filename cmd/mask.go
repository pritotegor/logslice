package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var maskFields []string
var maskPattern string
var maskReplacement string

func init() {
	maskCmd := &cobra.Command{
		Use:   "mask [file]",
		Short: "Mask or redact sensitive field values in structured logs",
		Long:  `Redact sensitive fields by replacing their values with a placeholder string. Supports JSON field targeting and regex pattern matching.`,
		RunE:  RunMask,
	}

	maskCmd.Flags().StringSliceVar(&maskFields, "field", nil, "JSON field names to redact (e.g. --field password,token)")
	maskCmd.Flags().StringVar(&maskPattern, "pattern", "", "Regex pattern to match and redact within values")
	maskCmd.Flags().StringVar(&maskReplacement, "replace", "***", "Replacement string for masked values")

	rootCmd.AddCommand(maskCmd)
}

func RunMask(cmd *cobra.Command, args []string) error {
	var reader io.Reader
	var err error

	if len(args) > 0 {
		reader, err = openInput(args[0])
		if err != nil {
			return err
		}
	} else {
		reader = os.Stdin
	}

	var re *regexp.Regexp
	if maskPattern != "" {
		re, err = regexp.Compile(maskPattern)
		if err != nil {
			return fmt.Errorf("invalid pattern: %w", err)
		}
	}

	scanner := bufio.NewScanner(reader)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		masked := maskLine(line, maskFields, re, maskReplacement)
		fmt.Fprintln(writer, masked)
	}

	return scanner.Err()
}

func maskLine(line string, fields []string, re *regexp.Regexp, replacement string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		if re != nil {
			return re.ReplaceAllString(line, replacement)
		}
		return line
	}

	for _, f := range fields {
		if _, ok := obj[f]; ok {
			obj[f] = replacement
		}
	}

	if re != nil {
		for k, v := range obj {
			if s, ok := v.(string); ok {
				obj[k] = re.ReplaceAllString(s, replacement)
			}
		}
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}
