package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Render log lines using a Go template",
	Long:  `Apply a Go text/template to each JSON log line, outputting the rendered result.`,
	RunE:  RunTemplate,
}

func init() {
	templateCmd.Flags().StringP("tmpl", "t", "", "Go template string to apply to each line (required)")
	templateCmd.Flags().StringP("file", "f", "", "Input file (default: stdin)")
	templateCmd.Flags().BoolP("skip-invalid", "s", false, "Skip lines that fail to parse or render")
	_ = templateCmd.MarkFlagRequired("tmpl")
	rootCmd.AddCommand(templateCmd)
}

func RunTemplate(cmd *cobra.Command, args []string) error {
	tmplStr, _ := cmd.Flags().GetString("tmpl")
	filePath, _ := cmd.Flags().GetString("file")
	skipInvalid, _ := cmd.Flags().GetBool("skip-invalid")

	tmpl, err := template.New("log").Option("missingkey=zero").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	reader, err := openInput(filePath)
	if err != nil {
		return err
	}
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	return applyTemplate(reader, os.Stdout, tmpl, skipInvalid)
}

func applyTemplate(r io.Reader, w io.Writer, tmpl *template.Template, skipInvalid bool) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var fields map[string]interface{}
		if err := json.Unmarshal([]byte(line), &fields); err != nil {
			if skipInvalid {
				continue
			}
			return fmt.Errorf("failed to parse line as JSON: %w", err)
		}

		var sb strings.Builder
		if err := tmpl.Execute(&sb, fields); err != nil {
			if skipInvalid {
				continue
			}
			return fmt.Errorf("template execution failed: %w", err)
		}

		fmt.Fprintln(w, sb.String())
	}
	return scanner.Err()
}
