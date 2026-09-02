package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var (
	remoteHostFlag     string
	remotePortFlag     int
	remoteUserFlag     string
	remotePasswordFlag string
	remoteKeyFlag      string
	remoteServerFlag   string
	remoteTokenFlag    string
	testOnlyFlag       bool
)

var deployRemoteCmd = &cobra.Command{
	Use:   "deploy-remote",
	Short: "Remotely install and enroll InfraPilot Agent onto a target machine using SSH credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		if remoteHostFlag == "" {
			return fmt.Errorf("--host is required (e.g. --host 192.168.1.50)")
		}
		if remoteUserFlag == "" {
			return fmt.Errorf("--user is required (e.g. --user root or --user ubuntu)")
		}
		if remotePasswordFlag == "" && remoteKeyFlag == "" {
			return fmt.Errorf("either --password or --ssh-key must be provided")
		}
		if remoteServerFlag == "" {
			remoteServerFlag = "http://localhost:8080"
		}

		payload := map[string]interface{}{
			"host":         remoteHostFlag,
			"port":         remotePortFlag,
			"username":     remoteUserFlag,
			"password":     remotePasswordFlag,
			"ssh_key":      remoteKeyFlag,
			"server_url":   remoteServerFlag,
			"enroll_token": remoteTokenFlag,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}

		endpoint := fmt.Sprintf("%s/api/v1/agent/remote-deploy/execute", remoteServerFlag)
		if testOnlyFlag {
			endpoint = fmt.Sprintf("%s/api/v1/agent/remote-deploy/test", remoteServerFlag)
			fmt.Printf("🔍 Testing SSH connectivity to %s:%d as '%s'...\n", remoteHostFlag, remotePortFlag, remoteUserFlag)
		} else {
			fmt.Printf("🚀 Deploying InfraPilot Agent to %s:%d as '%s'...\n", remoteHostFlag, remotePortFlag, remoteUserFlag)
		}

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to communicate with InfraPilot backend: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		var parsed map[string]interface{}
		_ = json.Unmarshal(body, &parsed)

		if testOnlyFlag {
			if success, ok := parsed["success"].(bool); ok && success {
				fmt.Printf("✅ Connection Successful! Hostname: %v (OS: %v %v)\n", parsed["hostname"], parsed["os"], parsed["arch"])
			} else {
				fmt.Printf("❌ Connection Failed: %v\n", parsed["error"])
			}
			return nil
		}

		if success, ok := parsed["success"].(bool); ok && success {
			fmt.Println("==================================================")
			fmt.Println("🎉 Remote Deployment Completed Successfully!")
			fmt.Printf("Target Host:  %v\n", parsed["host"])
			fmt.Printf("Hostname:     %v\n", parsed["hostname"])
			fmt.Printf("Machine ID:   %v\n", parsed["machine_id"])
			fmt.Printf("Total Time:   %v ms\n", parsed["total_time_ms"])
			fmt.Println("==================================================")

			if logs, ok := parsed["logs"].([]interface{}); ok {
				fmt.Println("\nExecution Logs:")
				for _, log := range logs {
					fmt.Println(log)
				}
			}
		} else {
			fmt.Printf("❌ Remote Deployment Failed: %v\n", parsed["error"])
			if logs, ok := parsed["logs"].([]interface{}); ok {
				fmt.Println("\nLogs:")
				for _, log := range logs {
					fmt.Println(log)
				}
			}
		}

		return nil
	},
}

func init() {
	deployRemoteCmd.Flags().StringVar(&remoteHostFlag, "host", "", "Target machine IP or hostname (e.g. 192.168.1.50)")
	deployRemoteCmd.Flags().IntVar(&remotePortFlag, "port", 22, "Target SSH port (default: 22)")
	deployRemoteCmd.Flags().StringVar(&remoteUserFlag, "user", "root", "Target machine username (e.g. root, ubuntu)")
	deployRemoteCmd.Flags().StringVar(&remotePasswordFlag, "password", "", "Target machine SSH password")
	deployRemoteCmd.Flags().StringVar(&remoteKeyFlag, "ssh-key", "", "Target machine SSH private key")
	deployRemoteCmd.Flags().StringVar(&remoteServerFlag, "server", "http://localhost:8080", "InfraPilot backend server URL")
	deployRemoteCmd.Flags().StringVar(&remoteTokenFlag, "token", "", "Optional custom enrollment token")
	deployRemoteCmd.Flags().BoolVar(&testOnlyFlag, "test-only", false, "Test connection and credentials only without installing")

	rootCmd.AddCommand(deployRemoteCmd)
}
