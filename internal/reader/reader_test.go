package reader

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestReader_ReadDirectory(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create test directory
	tmpDir, cleanup := testDir(t)
	defer cleanup()

	// Create test files
	files := map[string]string{
		".gitignore":      "*.tmp\n*.log",
		"main.go":         "package main\n\nfunc main() {}\n",
		"test/file.txt":   "test content",
		"test/ignore.tmp": "should be ignored",
		"test/binary.exe": string([]byte{0x00, 0x01, 0x02}),
	}
	createTestFiles(t, tmpDir, files)

	// Create reader

	reader, err := New(false, tmpDir)
	if err != nil {
		t.Fatalf("failed to create reader: %v", err)
	}

	// Read directory
	contents := reader.ReadDirectory(tmpDir, false)

	// Verify results
	expectedFiles := []string{"main.go", "test/file.txt"}
	if len(contents) != len(expectedFiles) {
		fmt.Println(contents)
		t.Errorf("got %d files, want %d", len(contents), len(expectedFiles))
	}

	// Check that ignored/binary files are not included
	for _, content := range contents {
		if filepath.Ext(content) == ".tmp" {
			t.Error("found .tmp file that should be ignored")
		}
		if filepath.Ext(content) == ".exe" {
			t.Error("found .exe file that should be excluded")
		}
	}
}
