package main

import "time"

type Cluster struct {
	Name        string
	Primary     string
	Namespace   string
	ClusterName string

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
