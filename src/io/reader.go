package io

import (
	"bufio"
	"fmt"
	"os/exec"
)

// ReadDiff reads diff content by running git diff command
func ReadDiff() ([]string, error) {
	// Run git diff command
	cmd := exec.Command("git", "diff")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to run git diff: %w", err)
	}

	var lines []string
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading diff: %w", err)
	}

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("git diff command failed: %w", err)
	}

	return lines, nil
}
