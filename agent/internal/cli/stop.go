package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop InfraPilot Agent",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping InfraPilot Agent...")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
