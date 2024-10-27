package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/lkarasinski/cbtxt/internal/reader"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cbtxt [directory]",
	Short: "cbtxt is a directory text content tool",
	Long: `cbtxt transforms your codebase into a single string - useful for
copying to LLM's.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runRoot,
}

var (
	noGitignore bool
)

func init() {
	rootCmd.Flags().BoolVar(&noGitignore, "no-gitignore", false, "Ignore .gitignore rules when processing directory")
}

func runRoot(cmd *cobra.Command, args []string) {
	r, err := reader.New(noGitignore, ".")
	if err != nil {
		os.Exit(1)
	}

	dir := r.ProjectRoot
	if len(args) > 0 {
		dir = args[0]
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("Error: Directory '%s' does not exist\n", dir)
		os.Exit(1)
	}

	fileContents := r.ReadDirectory(dir, noGitignore)

	err = clipboard.WriteAll(strings.Join(fileContents, "\n"))

	if err != nil {
		fmt.Println(fmt.Errorf("could not copy to clipboard: %v", err))
		os.Exit(1)
	}
}

func Execute() error {
	return rootCmd.Execute()
}
