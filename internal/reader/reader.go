package reader

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/lkarasinski/cbtxt/internal/gitignore"
	"github.com/lkarasinski/cbtxt/internal/template"
)

type Reader struct {
	tmpl        *template.Template
	ProjectRoot string
	gitIgnore   *gitignore.GitIgnore
}

func New(noGitignore bool) (*Reader, error) {
	tmpl, err := template.New()
	if err != nil {
		return nil, err
	}

	projectRoot, err := findProjectRoot(".")

	if err != nil {
		return nil, err
	}

	gitIgnore, err := gitignore.New(filepath.Join(projectRoot, ".gitignore"), noGitignore)

	if err != nil {
		return nil, err
	}

	return &Reader{tmpl: tmpl, ProjectRoot: projectRoot, gitIgnore: gitIgnore}, nil
}

func (r *Reader) ReadFile(path string) (string, error) {
	if isExcludedFile(path) {
		return "", fmt.Errorf("excluded file type: %s", path)
	}

	isBinary, err := isBinaryFile(path)
	if err != nil {
		return "", fmt.Errorf("check if binary: %w", err)
	}
	if isBinary {
		return "", fmt.Errorf("binary file detected: %s", path)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("could not get absolute path %s: %v", path, err)
	}

	content, err := os.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	data := template.FileData{
		Path:    strings.TrimPrefix(absPath, filepath.Dir(r.ProjectRoot)),
		Content: string(content),
	}

	return r.tmpl.Format(data)
}

func findProjectRoot(dir string) (string, error) {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("could not get absolute path for directory %s: %v", dir, err)
	}

	data, err := os.ReadDir(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to read directory %s: %v", absPath, err)
	}

	// Check current directory for .gitignore
	for _, file := range data {
		if file.Name() == ".gitignore" {
			return absPath, nil
		}
	}

	// If we're at the root directory, stop searching
	parent := filepath.Dir(absPath)
	if parent == absPath {
		return "", nil
	}

	// Recursively check parent directory
	return findProjectRoot(parent)
}

func (r *Reader) ReadDirectory(dir string, noGitignore bool) []string {
	fileContents := []string{}

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if r.gitIgnore.ShouldIgnore(path) {
			return nil
		}

		if !info.IsDir() {
			file, err := r.ReadFile(path)

			if err != nil {
				fmt.Println(fmt.Errorf("could read file: %v", err))
			} else {
				fileContents = append(fileContents, file)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println(fmt.Errorf("could not walk directory: %v", err))
		os.Exit(1)
	}

	return fileContents
}
