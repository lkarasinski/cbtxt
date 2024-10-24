package gitignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
)

type GitIgnore struct {
	patterns []glob.Glob
}

func New(gitignorePath string) (*GitIgnore, error) {
	gi := &GitIgnore{
		patterns: make([]glob.Glob, 0),
	}

	// Read .gitignore file
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		// If .gitignore doesn't exist, just use default patterns
		return gi.addDefaultPatterns(), nil
	}

	// Process .gitignore content
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			gi.addPattern(line)
		}
	}

	return gi.addDefaultPatterns(), nil
}

func (gi *GitIgnore) addDefaultPatterns() *GitIgnore {
	defaultPatterns := []string{
		".git/",
		".gitignore",
		"*.lock",
	}

	for _, pattern := range defaultPatterns {
		gi.addPattern(pattern)
	}

	return gi
}

func (gi *GitIgnore) addPattern(pattern string) {
	if pattern == "" {
		return
	}

	// Handle different pattern types
	patterns := []string{}

	// Remove trailing slash if present
	pattern = strings.TrimSuffix(pattern, "/")

	// For each pattern, we'll create multiple variants to ensure comprehensive matching
	patterns = append(patterns,
		pattern,             // Original pattern
		"**/"+pattern,       // Match in any subdirectory
		pattern+"/**",       // Match all contents if it's a directory
		"**/"+pattern+"/**", // Match directory and contents anywhere
	)

	// Compile all pattern variants
	for _, p := range patterns {
		g, err := glob.Compile(p)
		if err != nil {
			continue // Skip invalid patterns
		}
		gi.patterns = append(gi.patterns, g)
	}
}

func (gi *GitIgnore) ShouldIgnore(path string) bool {
	// Convert Windows paths to forward slashes for consistency
	path = filepath.ToSlash(path)

	// Check both the path and its base name
	baseName := filepath.Base(path)

	for _, pattern := range gi.patterns {
		if pattern.Match(path) || pattern.Match(baseName) {
			return true
		}
	}

	return false
}
