package cli

import (
	"fmt"
	"os"

	"infrapilot/agent/internal/config"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show InfraPilot agent status",
	RunE: func(cmd *cobra.Command, args []string) error {

		if _, err := os.Stat("config.json"); err != nil {
			fmt.Println("Agent Status : NOT ENROLLED")
			return nil
		}

		cfg, err := config.LoadConfigJSON("config.json")
		if err != nil {
			return err
		}

		fmt.Println("=========== InfraPilot Agent ===========")
		fmt.Println("Status      : ENROLLED")
		fmt.Println("Machine ID  :", cfg.MachineID)
		fmt.Println("Backend     :", cfg.BackendURL)
		fmt.Println("Interval    :", cfg.Interval, "seconds")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
