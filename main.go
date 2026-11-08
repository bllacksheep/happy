package main

import (
	"fmt"
	"log/slog"
)

func (c *Cluster) StartDefaultChecks() {
	c.WalReplicationHealth.Toggle("wal_replication_check_enabled",
		c.DefaultChecks.CheckWalReplicationHealth, true)
	c.LastBackupHealth.Toggle("backups_check_enabled",
		c.DefaultChecks.CheckLastBackupHealth, true)
	c.RPOHealth.Toggle("rpo_check_enabled",
		c.DefaultChecks.CheckRPOHealth, true)
	c.DbInstancesHealth.Toggle("db_instance_check_enabled",
		c.DefaultChecks.CheckDbInstanceHealth, true)

	// go startDefaultChecks?
}

func startHappy() {
	cfg := validateConfig()

	ch := make(chan *Cluster)

	// should scale for n clusters
	for _, clusterConfig := range cfg.Happy.Clusters {
		cluster := Cluster{
			ClusterConfig:    clusterConfig,
			DefaultThreshold: cfg.Happy.Default.Threshold,
			DefaultChecks:    cfg.Happy.Default.HealthChecks,
		}
		// is run for every cluster not required but will do for now
		cluster.StartDefaultChecks()

		go cluster.StartClusterStatusChecks(ch)
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
