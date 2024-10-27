# cbtxt

cbtxt is a simple CLI tool that transforms your codebase into a single text file, making it easy to share code with Large Language Models (LLMs). I created this, while learning golang.

## Features

- Converts an entire codebase into a single, formatted text file
- Respects `.gitignore` rules (can be disabled)
- Automatically detects and skips binary files
- Copies output directly to clipboard
- Project root detection based on `.gitignore` location

## Installation

```bash
go install github.com/lkarasinski/cbtxt@latest
```

## Usage

Basic usage:
```bash
cbtxt [directory]
```

By default, CBTXT will use the current directory if no directory is specified.

Options:
- `--no-gitignore`: Ignore `.gitignore` rules when processing directory

Example:
```bash
# Process current project
cbtxt

# Process specific directory
cbtxt ./my-project

# Process directory ignoring .gitignore rules
cbtxt --no-gitignore ./my-project
```

## How It Works

1. Locates the project root by finding the `.gitignore` file
2. Walks through the directory tree
3. Skips files that are:
   - Matched by `.gitignore` patterns
   - Binary files
   - Common binary extensions (e.g., .pdf, .jpg, .exe)
4. Formats each file's content with its path
5. Copies the combined output to clipboard

## Development Status

This is a learning project I created while learning Go programming. While functional, it might not follow all best practices and could have room for improvements.
## License

MIT License - see the [LICENSE](LICENSE) file for details
