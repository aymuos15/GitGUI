package diff

import (
	"strings"

	"gg/src/models"
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

			// Start new file with default "Modified" status
			currentFile = &models.FileDiff{
				Name:    fileName,
				Content: []string{line},
				Status:  "Modified",
			}
		} else if currentFile != nil {
			currentFile.Content = append(currentFile.Content, line)
		}
	}

	// Don't forget the last file
	if currentFile != nil {
		files = append(files, *currentFile)
	}

	// Initialize syntax highlighting, calculate stats, and detect status for all files
	for i := range files {
		files[i].InitSyntaxHighlighting()
		files[i].CalculateStats()
		detectFileStatus(&files[i])
	}

	return files
}

// detectFileStatus parses the file content to determine its status
func detectFileStatus(file *models.FileDiff) {
	for _, line := range file.Content {
		if strings.HasPrefix(line, "new file mode") {
			file.Status = "New"
			return
		} else if strings.HasPrefix(line, "deleted file mode") {
			file.Status = "Deleted"
			return
		} else if strings.HasPrefix(line, "rename from") {
			file.Status = "Renamed"
			return
		}
	}
	// Keep default "Modified" status if no special status detected
}

// CreateUntrackedFileDiffs converts a list of untracked file paths to FileDiff objects
func CreateUntrackedFileDiffs(untrackedPaths []string) []models.FileDiff {
	var files []models.FileDiff

	for _, filePath := range untrackedPaths {
		file := models.FileDiff{
			Name:    filePath,
			Content: []string{}, // No diff content for untracked files
			Status:  "Untracked",
		}
		file.InitSyntaxHighlighting()
		files = append(files, file)
	}

	return files
}
