package reader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// maxBytesToCheck is the number of bytes to check for binary content
const maxBytesToCheck = 512

// isBinaryFile checks if a file is likely to be binary
func isBinaryFile(filepath string) (bool, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return false, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	// Read the first maxBytesToCheck bytes
	buf := make([]byte, maxBytesToCheck)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("read file: %w", err)
	}
	buf = buf[:n]

	// Check for null bytes and high concentration of non-text bytes
	numNonText := 0
	for _, b := range buf {
		if b == 0 {
			return true, nil // Null bytes strongly indicate binary
		}
		if b < 7 || b > 127 {
			numNonText++
		}
	}

	// If more than 30% non-text bytes, probably binary
	return float64(numNonText)/float64(len(buf)) > 0.3, nil
}

// isExcludedFile checks if a file should be excluded based on common binary extensions
func isExcludedFile(path string) bool {
	// Common binary file extensions
	binaryExtensions := map[string]bool{
		".pdf":   true,
		".png":   true,
		".jpg":   true,
		".jpeg":  true,
		".gif":   true,
		".zip":   true,
		".tar":   true,
		".gz":    true,
		".rar":   true,
		".exe":   true,
		".dll":   true,
		".so":    true,
		".pyc":   true,
		".o":     true,
		".class": true,
	}

	ext := strings.ToLower(filepath.Ext(path))
	return binaryExtensions[ext]
}
