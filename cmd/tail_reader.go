package cmd

import (
	"bufio"
	"io"
	"os"
	"time"
)

// readLastN reads all lines from r and returns the last n lines.
func readLastN(r io.Reader, n int) ([]string, error) {
	scanner := bufio.NewScanner(r)
	buf := make([]string, 0, n+1)

	for scanner.Scan() {
		line := scanner.Text()
		buf = append(buf, line)
		if len(buf) > n {
			buf = buf[1:]
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return buf, nil
}

// followFile continuously reads new lines appended to path and writes them to w.
// It polls the file every 250ms, similar to tail -f behaviour.
func followFile(path string, w io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Seek to end before following.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		for {
			line, err := reader.ReadString('\n')
			if len(line) > 0 {
				if _, werr := io.WriteString(w, line); werr != nil {
					return werr
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}
