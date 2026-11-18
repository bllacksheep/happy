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

func (h *Happy) validateConfig() error {
	configPath := flag.String("c", "happy.yaml", "config file")
	flag.Parse()
	configData, err := os.ReadFile(*configPath)

	if err != nil {
		slog.Error("main.go", "msg", "couldn't read file", "err", err)
		return err
	}

	err = yaml.Unmarshal(configData, &h.Config)
	if err != nil {
		slog.Error("main.go", "msg", "error unmarshalling yaml", "err", err)
		return err
	}

	slog.Info("main.go", "msg", "validated config")
	// check required feilds here
	return nil
}
