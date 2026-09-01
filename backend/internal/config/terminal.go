package config

import (
	"os"
	"strings"
)

// GetAllowedCommands returns a slice of command prefixes that are permitted for remote terminal execution.
// The list can be configured via the environment variable TERMINAL_ALLOWLIST as a comma‑separated list.
// If the variable is not set, a sensible default set is returned.
func GetAllowedCommands() []string {
	env := os.Getenv("TERMINAL_ALLOWLIST")
	if env == "" {
		// Default safe commands commonly needed in enterprise environments.
		return []string{"docker", "kubectl", "systemctl", "cat", "ls", "ps", "top", "whoami", "uptime"}
	}
	parts := strings.Split(env, ",")
	var cleaned []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}
