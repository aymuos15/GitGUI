package io

import (
	"bufio"
	"fmt"
	"os/exec"
)

// ReadDiff reads diff content by running git diff command
// Falls back to staged changes if working tree is empty
func ReadDiff() ([]string, string, error) {
	// Try working tree changes first
	lines, err := runGitDiff("git", "diff")
	if err != nil {
		return nil, "", err
	}

	// If working tree is empty, fall back to staged changes
	if len(lines) == 0 {
		lines, err = runGitDiff("git", "diff", "--cached")
		if err != nil {
			return nil, "", err
		}
		return lines, "staged", nil
	}

	return lines, "working", nil
}

// runGitDiff executes a git diff command and returns the output lines
func runGitDiff(args ...string) ([]string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to run git command: %w", err)
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
		return nil, fmt.Errorf("git command failed: %w", err)
	}

	return lines, nil
}

// ReadUntrackedFiles reads untracked files using git ls-files
func ReadUntrackedFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to run git command: %w", err)
	}

	var files []string
	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		files = append(files, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading untracked files: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("git command failed: %w", err)
	}

	return files, nil
}
