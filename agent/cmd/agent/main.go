package main

import (
	"log"

	"infrapilot/agent/internal/cli"
	"infrapilot/agent/internal/service"

	kservice "github.com/kardianos/service"
)

func main() {
	// If running as a Windows service (non-interactive, started by SCM)
	if !kservice.Interactive() {
		s, err := service.NewService()
		if err != nil {
			log.Fatal(err)
		}
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Otherwise, run the CLI
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
