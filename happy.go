package main

import (
	"fmt"
	"log/slog"
)

type Happy struct {
	Config Config
}


func newHappy() (*Happy, error) {
	//TODO: check here
	return &Happy{}, nil
}

func (h *Happy) Start() error {
	if err := h.validateConfig(); err != nil {
		return fmt.Errorf("failed to validate config", err)
	}

	defaults := &SharedDefaultConfig{
		h.Config.Happy.DefaultConfig.Checks,
		h.Config.Happy.DefaultConfig.Threshold,
	}

	for _, clusterConfig := range h.Config.Happy.Clusters {
		cluster, err := newCluster(clusterConfig, defaults)

		if err != nil {
			return fmt.Errorf("failed to create cluster", err)
		}

		// RegsiterCheck to add a check dynamically?
		cluster.StartChecks()

		slog.Info("main.go", cluster.Name, "happy started")
	}
	return nil
}



/*
	ch := make(chan *Cluster)

	happy := make(map[string]*Cluster, len(h.Config.Happy.Clusters))

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
*/
