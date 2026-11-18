package main

connecttok8s

newcluster


func newCluster(name, ns string) Cluster {
	return Cluster{
		Name: name,
		Namespace: ns,
	}

}

type Happy struct {
	Config Config
}

func newHappy() (*Happy, error) {
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

	ch := make(chan *Cluster)

	// should scale for n clusters
	for _, clusterConfig := range h.Config.Happy.Clusters {
		cluster := newCluster(clusterConfig.Name, clusterConfig.Namespace)
		cluster.Defaults = defaults
		// cluster := Cluster{
		// 	Name:      clusterConfig.Name,
		// 	Namespace: clusterConfig.Namespace,
		//
		// 	Defaults: defaults,
		// }

		// cluster.RegisterChecks should build up the []Poller
		// for all custom and default checks
		cluster.RegisterChecks()

		go cluster.StartStatusPoller(ch)
		go cluster.StartChecks()

		slog.Info("main.go", cluster.Name, "happy started")
	}

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
}

