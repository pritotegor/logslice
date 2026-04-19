package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// OutputFormat defines how matched log lines are written.
type OutputFormat string

const (
	FormatRaw  OutputFormat = "raw"
	FormatJSON OutputFormat = "json"
	FormatCSV  OutputFormat = "csv"
)

// Writer wraps an io.Writer and formats output lines.
type Writer struct {
	w      io.Writer
	format OutputFormat
	fields []string
}

// NewWriter creates a new output Writer.
func NewWriter(w io.Writer, format OutputFormat, fields []string) *Writer {
	return &Writer{w: w, format: format, fields: fields}
}

// WriteLine writes a single log line in the configured format.
func (wr *Writer) WriteLine(line string) error {
	switch wr.format {
	case FormatJSON:
		return wr.writeJSON(line)
	case FormatCSV:
		return wr.writeCSV(line)
	default:
		_, err := fmt.Fprintln(wr.w, line)
		return err
	}
}

func (wr *Writer) writeJSON(line string) error {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		// Not valid JSON — emit as-is wrapped in a message field.
		wrapped, _ := json.Marshal(map[string]string{"message": line})
		_, err2 := fmt.Fprintln(wr.w, string(wrapped))
		return err2
	}
	if len(wr.fields) > 0 {
		filtered := make(map[string]interface{}, len(wr.fields))
		for _, f := range wr.fields {
			if v, ok := obj[f]; ok {
				filtered[f] = v
			}
		}
		obj = filtered
	}
	out, _ := json.Marshal(obj)
	_, err := fmt.Fprintln(wr.w, string(out))
	return err
}

func (wr *Writer) writeCSV(line string) error {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		_, err2 := fmt.Fprintln(wr.w, line)
		return err2
	}
	keys := wr.fields
	if len(keys) == 0 {
		for k := range obj {
			keys = append(keys, k)
		}
	}
	vals := make([]string, len(keys))
	for i, k := range keys {
		if v, ok := obj[k]; ok {
			vals[i] = fmt.Sprintf("%v", v)
		}
	}
	_, err := fmt.Fprintln(wr.w, strings.Join(vals, ","))
	return err
}
