package main

import (
	"fmt"
	"log/slog"
)

// load the checks into the app state
func (c *Cluster) RegisterChecks() {
	defaultChecks := []Poller{
		&HttpCheck{
			Name:              "default_wal_replication_check_enabled",
			ClusterName:       c.Name,
			IsDefault:         true,
			Threshold:         c.DefaultThreshold,
			Enabled:           true,
			UserDefinedToggle: c.DefaultCheckToggles.CheckWalReplicationHealth,
		},
		&HttpCheck{
			Name:              "default_backups_check_enabled",
			ClusterName:       c.Name,
			IsDefault:         true,
			Threshold:         c.DefaultThreshold,
			Enabled:           true,
			UserDefinedToggle: c.DefaultCheckToggles.CheckLastBackupHealth,
		},
		&HttpCheck{
			Name:              "default_rpo_check_enabled",
			ClusterName:       c.Name,
			IsDefault:         true,
			Threshold:         c.DefaultThreshold,
			Enabled:           true,
			UserDefinedToggle: c.DefaultCheckToggles.CheckRPOHealth,
		},
		&HttpCheck{
			Name:              "default_db_instance_check_enabled",
			ClusterName:       c.Name,
			IsDefault:         true,
			Threshold:         c.DefaultThreshold,
			Enabled:           true,
			UserDefinedToggle: c.DefaultCheckToggles.CheckDbInstanceHealth,
		},
	}

	c.Checks = defaultChecks


	// should be small constant
	for _, check := range c.Checks {
		check.Toggle()
	}
}

func startHappy() {
	cfg := validateConfig()

	ch := make(chan *Cluster)

	// should scale for n clusters
	for _, clusterConfig := range cfg.Happy.Clusters {
		cluster := Cluster{
			Name:      clusterConfig.Name,
			Namespace: clusterConfig.Namespace,

			DefaultCheckToggles: cfg.Happy.Defaults.HealthChecks,
			DefaultThreshold:    cfg.Happy.Defaults.Threshold,
		}

		// cluster.RegisterChecks should build up the []Poller
		// for all custom and default checks
		cluster.RegisterChecks()

		go cluster.StartStatusPoller(ch)
		go cluster.StartChecks()

		slog.Info("main.go", cluster.Name, "happy started")
	}

	happy := make(map[string]*Cluster, len(cfg.Happy.Clusters))

	slog.Info("main.go", "msg", "happy started")

	for cluster := range ch {
		cluster.IsHappy()

		// swap cluster state by name
		happy[cluster.Name] = cluster
		currentCluster := happy[cluster.Name]

		slog.Info("main.go", currentCluster.Name, *happy[cluster.Name].Status)

		slo := 100 * (float64(currentCluster.Succeeded) / float64(currentCluster.Polled))

		slog.Info("main.go", currentCluster.Name,
			fmt.Sprintf("succeeded/polled: %v/%v, happy %v, slo: %.3f%%",
				currentCluster.Succeeded, currentCluster.Polled,
				currentCluster.Happy, slo))
	}
}

func main() {
	startHappy()
}
