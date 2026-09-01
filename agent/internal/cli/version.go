package cli

import (
	"fmt"

	"infrapilot/agent/internal/system"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show InfraPilot version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("InfraPilot Agent")
		fmt.Println("Version:", system.GetVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
