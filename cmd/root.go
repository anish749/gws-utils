package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gws_utils",
	Short: "Utilities and wrappers for Google Workspace CLI (gws)",
}

func Execute() error {
	return rootCmd.Execute()
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
