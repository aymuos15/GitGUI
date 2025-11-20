package diff

import (
	"strings"

	"diffview/src/models"
)

// ParseDiffIntoFiles parses git diff output into separate file diffs
func ParseDiffIntoFiles(lines []string) []models.FileDiff {
	var files []models.FileDiff
	var currentFile *models.FileDiff

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			// Extract file name
			parts := strings.Fields(line)
			fileName := "unknown"
			if len(parts) >= 4 {
				fileName = parts[3]
				if strings.HasPrefix(fileName, "b/") {
					fileName = strings.TrimPrefix(fileName, "b/")
				}
			}

			// Save previous file if exists
			if currentFile != nil {
				files = append(files, *currentFile)
			}

			// Start new file
			currentFile = &models.FileDiff{
				Name:    fileName,
				Content: []string{line},
			}
		} else if currentFile != nil {
			currentFile.Content = append(currentFile.Content, line)
		}
	}

	// Don't forget the last file
	if currentFile != nil {
		files = append(files, *currentFile)
	}

	// Initialize syntax highlighting and calculate stats for all files
	for i := range files {
		files[i].InitSyntaxHighlighting()
		files[i].CalculateStats()
	}

	return files
}
