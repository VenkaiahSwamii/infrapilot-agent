package cli

import (
	"infrapilot/agent/internal/app"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start InfraPilot Agent",
	Run: func(cmd *cobra.Command, args []string) {
		app.RunAgent()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
