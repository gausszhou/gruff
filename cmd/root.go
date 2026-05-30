package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gruff",
	Short: "Render Markdown to ANSI terminal output",
	Long: `gruff is a fast Markdown-to-ANSI renderer for the terminal.

It reads a Markdown file and prints colorized, word-wrapped output 
to stdout, with automatic terminal width detection.`,
	Args:              cobra.MinimumNArgs(1),
	TraverseChildren:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return renderFile(cmd, args)
	},
}

var (
	width int
	light bool
)

func init() {
	rootCmd.PersistentFlags().IntVarP(&width, "width", "w", 0, "word wrap width (0 = auto-detect)")
	rootCmd.PersistentFlags().BoolVarP(&light, "light", "l", false, "use light theme")
}

func Execute() error {
	return rootCmd.Execute()
}
