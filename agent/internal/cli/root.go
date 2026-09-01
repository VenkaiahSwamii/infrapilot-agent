package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "infrapilot-agent",
	Short: "InfraPilot Enterprise Monitoring Agent",
}

func Execute() error {
	return rootCmd.Execute()
}
