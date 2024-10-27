package template

import (
	"strings"
	"testing"
)

func TestTemplate_Format(t *testing.T) {
	tests := []struct {
		name     string
		input    FileData
		expected string
	}{
		{
			name: "formats basic file data",
			input: FileData{
				Path:    "/test/file.txt",
				Content: "hello world",
			},
			expected: "=== /test/file.txt ===\n\nhello world\n",
		},
		{
			name: "handles empty content",
			input: FileData{
				Path:    "/empty.txt",
				Content: "",
			},
			expected: "=== /empty.txt ===\n\n\n",
		},
		{
			name: "handles multiline content",
			input: FileData{
				Path:    "/multi.txt",
				Content: "line1\nline2\nline3",
			},
			expected: "=== /multi.txt ===\n\nline1\nline2\nline3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := New()
			if err != nil {
				t.Fatalf("failed to create template: %v", err)
			}

			got, err := tmpl.Format(tt.input)
			if err != nil {
				t.Fatalf("Format() error = %v", err)
			}

			// Normalize line endings for comparison
			got = strings.ReplaceAll(got, "\r\n", "\n")
			expected := strings.ReplaceAll(tt.expected, "\r\n", "\n")

			if got != expected {
				t.Errorf("Format() = %q, want %q", got, expected)
			}
		})
	}
}
