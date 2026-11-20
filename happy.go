package main

import (
	"fmt"
	"log/slog"
)

type Happy struct {
	config      Config
	clusters    map[string]*Cluster
	pollPending chan Poller
}

func NewHappy() (*Happy, error) {
	happy := &Happy{}
	if err := happy.validateConfig(); err != nil {
		return nil, fmt.Errorf("failed to validate config", err)
	}

	slog.Info("main.go", "msg", "validated config")

	// RegsiterCheck to add a check dynamically?
	for _, clusterConfig := range happy.config.Clusters {
		cluster := newCluster(clusterConfig)

		happy.clusters[cluster.Name] = &cluster
		slog.Info("cluster.go", cluster.Name, "happy started")
	}

	return happy, nil
}

func (h *Happy) initializeAndStartChecks() error {
	// init

	h.ready()

	return nil
}

func (h *Happy) ready() {
	for check := range h.pollPending {
		check.Poll()
	}
	//		slog.Info("main.go", currentCluster.Name, *happyClusters[cluster.Name].Status)
	//
	//		slo := 100 * (float64(currentCluster.Succeeded) / float64(currentCluster.Polled))
	//
	//		slog.Info("main.go", currentCluster.Name,
	//			fmt.Sprintf("succeeded/polled: %v/%v, happy %v, slo: %.3f%%",
	//				currentCluster.Succeeded, currentCluster.Polled,
	//				currentCluster.Happy, slo))
	//	}
	//
	// return nil
}

func (h *Happy) Start() error {

	h.initializeAndStartChecks()

	// should block here?
	return nil
}
