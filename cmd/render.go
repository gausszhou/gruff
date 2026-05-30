package cmd

import (
	"fmt"
	"os"

	"github.com/gausszhou/gruff/gruff"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var renderCmd = &cobra.Command{
	Use:   "render <file.md>",
	Short: "Render a Markdown file to ANSI",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return renderFile(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
}

func renderFile(_ *cobra.Command, args []string) error {
	b, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	w := width
	if w <= 0 {
		if w2, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w2 > 0 {
			w = w2
		} else {
			w = 80
		}
	}

	opts := []gruff.Option{gruff.WithWordWrap(w)}
	if light {
		opts = append(opts, gruff.WithLight())
	} else {
		opts = append(opts, gruff.WithDark())
	}

	out, err := gruff.Render(string(b), opts...)
	if err != nil {
		return err
	}
	fmt.Print(out)
	return nil
}
