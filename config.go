package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Poller interface {
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
	Url string

	Check
}

type SqlCheck struct {
	Statement        string
	ConnectionString string

	Check
}

func (c *HttpCheck) Poll() bool {
	return true
}

func (c *SqlCheck) Poll() bool {
	return true
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

// custom health check router?
func (c Check) Happy(name string, fn func(string) bool) bool {
	return fn(c.Statement)
}

func (c *Check) Toggle(name string, toggle bool, isdef bool) {
	c.Name = name
	c.Enabled = toggle
	c.IsDefault = isdef
	slog.Info("happy.go", c.Name, strconv.FormatBool(c.Enabled))
}

type ClusterConfig struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`

	// Defaults Defaults
	UseDefaults bool `yaml:"use_defaults"`

	Checks []Check `yaml:"checks"`
}

func (h *Happy) validateConfig() error {
	configPath := flag.String("c", "happy.yaml", "config file")
	flag.Parse()
	configData, err := os.ReadFile(*configPath)

	if err != nil {
		return fmt.Errorf("couldn't read file", err)
	}

	err = yaml.Unmarshal(configData, &h.Config)
	if err != nil {
		return fmt.Errorf("error unmarhsalling yaml", err)
	}

	slog.Info("main.go", "msg", "validated config")
	// check required feilds here
	return nil
}
