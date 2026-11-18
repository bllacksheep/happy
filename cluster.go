package main

import "time"

type Cluster struct {
	Name        string
	Primary     string
	Namespace   string
	ClusterName string

	// http or sql pollers
	Checks []Poller

	// timer that will track recovery time
	RTO time.Duration

	OverallHealthScore float64
	Polled             int
	Succeeded          int

	Happy bool

	Status *CnpgClusterStatus
}

type Condition struct {
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	Message            string    `json:"message"`
	Reason             string    `json:"reason"`
	Status             string    `json:"status"`
	Type               string    `json:"type"`
}

// generic kube-api-server response
type CnpgClusterStatus struct {
	Status struct {
		Conditions []Condition `json:"conditions"`
	} `json:"status"`
}
