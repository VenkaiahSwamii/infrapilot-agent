package config

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// Spring 8.8: TLS & mTLS support
	TLSEnabled  bool   `yaml:"tls_enabled" json:"tls_enabled"`
	MTLSEnabled bool   `yaml:"mtls_enabled" json:"mtls_enabled"`
	CertPath    string `yaml:"cert_path" json:"cert_path"`
	KeyPath     string `yaml:"key_path" json:"key_path"`
	CAPath      string `yaml:"ca_path" json:"ca_path"`

	MachineID          string `yaml:"machine_id" json:"machine_id"`
	APIKey             string `yaml:"api_key" json:"api_key"`
	KeyVersion         int    `yaml:"key_version" json:"key_version"`
	EnrollmentToken    string `yaml:"enrollment_token" json:"enrollment_token"`
	BackendURL         string `yaml:"backend_url" json:"backend_url"`
	Interval           int    `yaml:"interval" json:"interval"`
	RegisteredHostname string `yaml:"registered_hostname" json:"registered_hostname"`
}

var AppConfig *Config

func Get() *Config {
	if AppConfig == nil {
		// Config should be loaded before use
		return &Config{}
	}
	return AppConfig
}

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	AppConfig = &cfg
	return AppConfig, nil
}

func LoadConfigJSON(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	AppConfig = &cfg
	return AppConfig, nil
}

func SaveConfigJSON(filePath string, cfg *Config) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}

func GetAllowedCommands() []string {
	return []string{
		"echo", "ls", "dir", "cd", "pwd",
		"whoami", "hostname", "date", "time",
		"systeminfo", "ipconfig", "netstat", "ping",
		"tracert", "nslookup", "tasklist", "ps",
		"df", "free", "uptime", "uname",
		"cat", "head", "tail", "grep", "find",
		"wc", "sort", "uniq", "awk", "sed",
		"docker", "kubectl", "kubens", "kubectx", "helm",
		"git", "npm", "node", "python", "python3",
		"curl", "wget", "ssh", "scp", "rsync",
	}
}
