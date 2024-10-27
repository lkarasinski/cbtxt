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

var (
	gitIgnore *gitignore.GitIgnore
)

type Reader struct {
	tmpl        *template.Template
	projectRoot string
}

func New() (*Reader, error) {
	tmpl, err := template.New()
	if err != nil {
		return nil, err
	}

	return &Reader{tmpl: tmpl}, nil
}

func (r *Reader) ReadFile(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("could not get absolute path %s: %v", path, err)
	}

	content, err := os.ReadFile(path)

	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	data := template.FileData{
		Path:    strings.TrimPrefix(absPath, filepath.Dir(r.projectRoot)),
		Content: string(content),
	}

	return r.tmpl.Format(data)
}

func (r *Reader) FindProjectRoot(dir string) (string, error) {
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
			r.projectRoot = absPath
			return dir, nil
		}
	}

	// If we're at the root directory, stop searching
	parent := filepath.Dir(absPath)
	if parent == absPath {
		return "", nil
	}

	// Recursively check parent directory
	return r.FindProjectRoot(parent)
}

func (r *Reader) ReadDirectory(dir string, noGitignore bool) {
	if !noGitignore {
		var err error
		gitIgnore, err = gitignore.New(filepath.Join(dir, ".gitignore"))

		if err != nil {
			fmt.Printf("Error loading .gitignore: %v", err)
			return
		}
	}

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if gitIgnore.ShouldIgnore(path) {
			return nil
		}

		if !info.IsDir() {
			file, err := r.ReadFile(path)

			if err != nil {
				fmt.Println(fmt.Errorf("could read file: %v", err))
				os.Exit(1)
			}

			fmt.Println(file)
		}
		return nil
	})

	if err != nil {
		fmt.Println(fmt.Errorf("could not walk directory: %v", err))
		os.Exit(1)
	}
}
