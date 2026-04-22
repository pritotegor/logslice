package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// openInputFn is a package-level variable to allow test overrides.
// In production it delegates to openInput.
var openInputFn = openInput

// runGrepOnReader runs the grep logic directly against an io.Reader,
// capturing stdout and returning the result as a string.
func runGrepOnReader(input io.Reader) string {
	pattern := grepPattern
	if grepIgnoreCase {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

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

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
