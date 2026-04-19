package cmd

import (
	"github.com/spf13/cobra"
)

var (
	filePath  string
	fromTime  string
	toTime    string
	fieldFilter string
)

var rootCmd = &cobra.Command{
	Use:   "logslice",
	Short: "A fast CLI tool for filtering and slicing structured log files",
	Long: `logslice filters and slices structured (JSON) log files
by time range or field patterns.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runFilter(filePath, fromTime, toTime, fieldFilter)
	},
}

func init() {
	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "path to log file (required)")
	rootCmd.Flags().StringVar(&fromTime, "from", "", "start time filter (RFC3339)")
	rootCmd.Flags().StringVar(&toTime, "to", "", "end time filter (RFC3339)")
	rootCmd.Flags().StringVarP(&fieldFilter, "field", "q", "", "field filter in key=value format")
	_ = rootCmd.MarkFlagRequired("file")
}

func Execute() error {
	return rootCmd.Execute()
}
