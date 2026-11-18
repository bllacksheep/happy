package main

import (
	"log/slog"
	"os"
)

// load the checks into the app state
func (c *Cluster) RegisterChecks() {
	defaultChecks := []Poller{
		&HttpCheck{
			Name:        "default_wal_replication_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   c.Defaults.SharedDefaultThreshold(),
			Enabled:     c.Defaults.Enabled(),
		},
		&HttpCheck{
			Name:        "default_backups_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   c.Defaults.SharedDefaultThreshold(),
			Enabled:     c.Defaults.Enabled(),
		},
		&HttpCheck{
			Name:        "default_rpo_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   c.Defaults.SharedDefaultThreshold(),
			Enabled:     c.Defaults.Enabled(),
		},
		&HttpCheck{
			Name:        "default_db_instance_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   c.Defaults.SharedDefaultThreshold(),
			Enabled:     c.Defaults.Enabled(),
		},
	}
	c.Checks = defaultChecks

	// should be small constant
	for _, check := range c.Checks {
		check.Toggle()
	}
}

// want shared default config and threshold pointer

type SharedDefaultConfig struct {
	Checks    []Check
	Threshold Threshold
}

type Defaults interface {
	Enabled() bool
	SharedDefaultChecks() []Check
	SharedDefaultThreshold() Threshold
}

func (sc *SharedDefaultConfig) SharedDefaultChecks() []Check {
	return sc.Checks
}

func (sc *SharedDefaultConfig) SharedDefaultThreshold() Threshold {
	return sc.Threshold
}

func (sc *SharedDefaultConfig) Enabled() bool {
	return true
}

func main() {
	happy, err := newHappy()
	if err != nil {
		slog.Error("failed to create happy instance", "error", err)
		os.Exit(1)
	}

	if err := happy.Start(); err != nil {
		slog.Error("failed to start happy instance", "error", err)
		os.Exit(1)
	}
}
