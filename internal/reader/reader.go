package reader

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/lkarasinski/cbtxt/internal/gitignore"
	"github.com/lkarasinski/cbtxt/internal/template"
)

type Reader struct {
	tmpl        *template.Template
	ProjectRoot string
	gitIgnore   *gitignore.GitIgnore
}

func New(noGitignore bool, path string) (*Reader, error) {
	tmpl, err := template.New()
	if err != nil {
		return nil, err
	}

	projectRoot, err := findProjectRoot(path)

	if err != nil {
		return nil, err
	}

	gitIgnore, err := gitignore.New(filepath.Join(projectRoot, ".gitignore"), !noGitignore)

	if err != nil {
		return nil, err
	}

	return &Reader{tmpl: tmpl, ProjectRoot: projectRoot, gitIgnore: gitIgnore}, nil
}

func (r *Reader) ProcessFile(path string, content []byte) (string, error) {
	data := template.FileData{
		Path:    path,
		Content: string(content),
	}

	formatted, err := r.tmpl.Format(data)

	if err != nil {
		return "", err
	}

	return formatted, nil
}

func (r *Reader) FilterFile(path string, file *os.File) (bool, error) {
	if isExcludedFile(path) {
		return false, nil
	}

	isBinary, err := isBinaryFile(file)
	if err != nil {
		return false, fmt.Errorf("failed to check if binary: %w", err)
	}
	if isBinary {
		return false, nil
	}
	if r.gitIgnore.ShouldIgnore(path) {
		return false, nil
	}

	return true, nil
}

func (r *Reader) ReadFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	validFile, err := r.FilterFile(path, file)
	if err != nil {
		return "", fmt.Errorf("failed to filter file: %w", err)
	}

	if !validFile {
		return "", nil
	}

	content, err := os.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return r.ProcessFile(path, content)
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
			} else if file != "" {
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
