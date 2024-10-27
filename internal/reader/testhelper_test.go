package reader

import (
	"os"
	"path/filepath"
	"testing"
)

// testDir creates a temporary directory with test files
func testDir(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "reader-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// createTestFiles creates test files with given content
func createTestFiles(t *testing.T, dir string, files map[string]string) {
	t.Helper()

	for path, content := range files {
		fullPath := filepath.Join(dir, path)

		// Create parent directories if needed
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("failed to create directories: %v", err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write file %s: %v", path, err)
		}
	}
}
