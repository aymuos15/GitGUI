package io

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// ReadDiff reads diff content from stdin or runs git diff
func ReadDiff() ([]string, error) {
	var input io.Reader

	// Check if stdin has data
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Data is being piped to stdin
		input = os.Stdin
	} else {
		// No piped data, run git diff
		cmd := exec.Command("git", "diff")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, fmt.Errorf("failed to create pipe: %w", err)
		}

		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("failed to run git diff: %w", err)
		}

		input = stdout
		defer cmd.Wait()
	}

	var lines []string
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading diff: %w", err)
	}

	return lines, nil
}
