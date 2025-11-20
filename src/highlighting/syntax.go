package highlighting

import "diffview/src/models"

// HighlightCode applies syntax highlighting to a line of code using file cache
func HighlightCode(code string, fileRef *models.FileDiff, lineIdx int) string {
	if fileRef == nil {
		// Fallback: no caching available
		return code
	}

	return fileRef.HighlightLine(lineIdx, code)
}
