package main

import (
	"flag"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Checker interface {
	Happy(func() bool)
	Toggle()
}

type Poller interface {
	Checker
	Poll() bool
}

type Check struct {
	Name        string `json:"name" yaml:"name"`
	ClusterName string

	Slo float64 `json:"slo"`

	Polled   int `json:"polled"`
	Suceeded int `json:"succeeded"`

	Healthy           bool `json:"healthy"`
	Enabled           bool `json:"enabled" yaml:"enabled"`
	UserDefinedToggle *bool
	IsDefault         bool `json:"is_default"`

	// check type might be defined by http: or sql: yaml parent
	UsePooler    bool   `json:"use_pooler" yaml:"use_pooler"`
	UsePrimary   bool   `json:"use_primary" yaml:"use_primary"`
	UseSecondary bool   `json:"use_secondary" yaml:"use_secondary"`
	UseInstance  string `json:"use_instance" yaml:"use_instance"`

	Threshold Threshold `json:"threshold" yaml:"threshold"`

	LastFailedAt  time.Time `json:"last_failed_at"`
	LastCheckedAt time.Time `json:"last_checked_at"`
}

func (cm *Check) Happy(fn func() bool) {
	if cm.Enabled {
		cm.Healthy = fn()
		cm.LastCheckedAt = time.Now()

		if !cm.Healthy {
			cm.LastFailedAt = time.Now()
		}
	}
}

func (c *Check) Toggle() {
	if c.UserDefinedToggle != nil {
		c.Enabled = *c.UserDefinedToggle
	}
	slog.Info("config.go", c.Name, strconv.FormatBool(c.Enabled))
}

type HttpCheck struct {
	Check
	Url string
}

type SqlCheck struct {
	Check
	Statement        string
	ConnectionString string
}

func (c *HttpCheck) Poll() bool {
	return true
}

func (c *SqlCheck) Poll() bool {
	return true
}

type ClusterConfig struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`

	Defaults Defaults

	Checks []Check `yaml:"checks"`
}

// optioanl default global config
type DefaultConfig struct {
	CheckClusterEndpointHealth *bool `yaml:"cluster_endpoint_health"`
	CheckWalReplicationHealth  *bool `yaml:"wal_replication_health"`
	CheckLastBackupHealth      *bool `yaml:"last_backup_health"`
	CheckDbInstanceHealth      *bool `yaml:"db_instance_health"`
	CheckRPOHealth             *bool `yaml:"rpo_health"`

	Checks    []Check   `yaml:"checks"`
	Threshold Threshold `yaml:"threshold"`
}

type Config struct {
	Happy struct {
		DefaultConfig `yaml:"defaults"`
		Clusters      []ClusterConfig `yaml:"clusters"`
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

