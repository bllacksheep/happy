package main

import (
	"fmt"
	"log/slog"
)

func (c *Cluster) StartDefaultChecks(toggles DefaultHealthCheckToggles) {
	defaultChecks := []Check{
		{
			Name:        "wal_replication_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   c.DefaultThreshold,
			Enabled:     toggles.CheckWalReplicationHealth,
		},
		{
			Name:        "backups_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   c.DefaultThreshold,
			Enabled:     toggles.CheckWalReplicationHealth,
		},
		{
			Name:        "rpo_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   c.DefaultThreshold,
			Enabled:     toggles.CheckWalReplicationHealth,
		},
		{
			Name:        "db_instance_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   c.DefaultThreshold,
			Enabled:     toggles.CheckWalReplicationHealth,
		},
	}
	c.Checks = defaultChecks

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

			DefaultThreshold: cfg.Happy.Defaults.Threshold,
		}
		cluster.StartDefaultChecks(cfg.Happy.Defaults.HealthChecks)

		go cluster.StartStatusPoller(ch)
		go cluster.StartCustomChecks()

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
