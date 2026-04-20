package cmd

import (
	"bufio"
	"fmt"
	"os"
)

// openInput returns a line scanner reading from the first positional argument
// (treated as a file path) or from stdin when no argument is provided.
func openInput(args []string) (*bufio.Scanner, error) {
	if len(args) == 0 {
		return bufio.NewScanner(os.Stdin), nil
	}

	f, err := os.Open(args[0])
	if err != nil {
		return nil, fmt.Errorf("cannot open file %q: %w", args[0], err)
	}

	// Wrap with a larger buffer to handle long log lines.
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	return scanner, nil
}
