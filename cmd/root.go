package cmd

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:   "alien-invasion",
		Short: "Alien invasion siulation util",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(analyzeCmd)
}
