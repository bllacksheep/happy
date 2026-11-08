package main

import (
	"flag"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Check struct {
	Name        string `json:"name" yaml:"name"`
	ClusterName string

	Slo float64 `json:"slo"`

	Polled   int `json:"polled"`
	Suceeded int `json:"succeeded"`

	Healthy   bool `json:"healthy"`
	Enabled   bool `json:"enabled" yaml:"enabled"`
	IsDefault bool `json:"is_default"`

	// check type might be defined by http: or sql: yaml parent
	UsePooler    bool   `json:"use_pooler" yaml:"use_pooler"`
	UsePrimary   bool   `json:"use_primary" yaml:"use_primary"`
	UseSecondary bool   `json:"use_secondary" yaml:"use_secondary"`
	UseInstance  string `json:"use_instance" yaml:"use_instance"`

	Threshold Threshold `json:"threshold" yaml:"threshold"`

	LastFailedAt  time.Time `json:"last_failed_at"`
	LastCheckedAt time.Time `json:"last_checked_at"`
}

type HttpCheck struct {
	Check
	Url string
}
type SqlCheck struct {
	Check
	Statement string
}

func (c *Check) Happy() {
	if c.Enabled {
		c.Healthy = c.Refresh()
		c.LastCheckedAt = time.Now()

		if !c.Healthy {
			c.LastFailedAt = time.Now()
		}
	}
}

func (c *HttpCheck) Refresh() bool {
	return true
}

func (c *SqlCheck) Refresh() bool {
	return true
}

func (c *Check) Toggle() {
	if c.IsDefault {
		c.Name = "default_" + c.Name
	}
	slog.Info("config.go", c.Name, strconv.FormatBool(c.Enabled))
}

type Health interface {
	Happy()
	Toggle()
	Refresh() bool
}

type ClusterConfig struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	// if this was type check, wouldn't need to set so much
	Checks []Check `yaml:"checks"`
}

type DefaultHealthCheckToggles struct {
	CheckClusterEndpointHealth bool `yaml:"check_cluster_endpoint_health"`
	CheckWalReplicationHealth  bool `yaml:"check_wal_replication_health"`
	CheckLastBackupHealth      bool `yaml:"check_last_backup_health"`
	CheckDbInstanceHealth      bool `yaml:"check_db_instance_health"`
	CheckRPOHealth             bool `yaml:"check_rpo_health"`
}

type Config struct {
	Happy struct {
		Defaults struct {
			HealthChecks DefaultHealthCheckToggles `yaml:"healthchecks"`
			Threshold    Threshold                 `yaml:"threshold"`
		} `yaml:"defaults"`
		Clusters []ClusterConfig `yaml:"clusters"`
	} `yaml:"happy"`
}

// return valid config
func validateConfig() Config {
	happyConfig := Config{}

	configPath := flag.String("c", "happy.yaml", "config file")
	flag.Parse()
	configData, err := os.ReadFile(*configPath)

	if err != nil {
		slog.Error("main.go", "msg", "couldn't read file", "err", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(configData, &happyConfig)
	if err != nil {
		slog.Error("main.go", "msg", "error unmarshalling yaml", "err", err)
		os.Exit(1)
	}

	slog.Info("main.go", "msg", "validated config")
	// check required feilds here
	return happyConfig
}
