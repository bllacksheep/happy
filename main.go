package main

import (
	"fmt"
	"log/slog"
)

func (c *Cluster) RegisterDefaultChecks() {
	// default checks are always enabled unless configured fasle
	defaultChecks := []Poller{
		&HttpCheck{
			Name:        "wal_replication_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   c.DefaultThreshold,
			Enabled:     true,
		},
		&HttpCheck{
			Name:        "backups_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   c.DefaultThreshold,
			Enabled:     true,
		},
		&HttpCheck{
			Name:        "rpo_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   c.DefaultThreshold,
			Enabled:     true,
		},
		&HttpCheck{
			Name:        "db_instance_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   c.DefaultThreshold,
			Enabled:     true,
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

			DefaultCheckToggles: cfg.Happy.Defaults.HealthChecks,
			DefaultThreshold:    cfg.Happy.Defaults.Threshold,
		}
		cluster.RegisterDefaultChecks()

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
