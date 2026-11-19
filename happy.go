package main

import (
	"fmt"
	"log/slog"
)

type Happy struct {
	Config   Config
	Clusters []Cluster
}

func NewHappy() (*Happy, error) {
	happy := &Happy{}

	if err := happy.validateConfig(); err != nil {
		return nil, fmt.Errorf("failed to validate config", err)
	}

	// RegsiterCheck to add a check dynamically?
	for _, clusterConfig := range happy.Clusters {
		cluster := newCluster(clusterConfig)

		happy.Clusters = append(happy.Clusters, cluster)
		slog.Info("cluster.go", cluster.Name, "happy started")
	}

	return happy, nil
}


func (h *Happy) initializeAndStartChecks() {

}

func (h *Happy) Start() error {

	h.initializeAndStartChecks()

	// should block here?
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
