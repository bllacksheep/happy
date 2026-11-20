package main

import (
	"net/http"
	"time"
)

type Cluster struct {
	Name        string
	Primary     string
	Namespace   string
	ClusterName string

	Endpoint    *http.Request

	// http or sql pollers
	// defaults will either be in here or not at start up
	Checks []Poller

	// timer that will track recovery time
	RTO time.Duration

	OverallHealthScore float64
	Polled             int
	Succeeded          int

	Happy bool

	Status *CnpgClusterStatus
}

func newCluster(cfg ClusterConfig) (*Cluster, error) {

	cluster := &Cluster{
		Name:      cfg.Name,
		Namespace: cfg.Namespace,
	}

	endpoint, err := newCnpgClusterRequest(cluster.Name, cluster.Namespace)

	if err != nil {
		return nil, err
	}

	cluster.Endpoint = endpoint

	defaults := defaultChecks(cfg, endpoint)

	cluster.Checks = append(cluster.Checks, defaults...)

	return cluster, nil
}
