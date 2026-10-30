package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// mine
type ClusterHealth struct {
	WalReplicationHealthy bool `json:"wal_replication_healthy"`
	DbInstancesHealthy    bool `json:"db_instances_healthy"`
	LastBackupHealthy     bool `json:"last_backup_healthy"`
	RPOHealthy            bool `json:"rpo_healthy"`

	Happy bool `json:"happy"`
}

// theirs
type ClusterStatus struct {
	Conditions []struct {
		LastTransitionTime time.Time
		Message            string
		Reason             string
		Status             string
		Type               string
	}
}

func CheckWalReplicationHealth() {
	fmt.Println("checking wal replication health")
}

func (cs *ClusterStatus) CheckEndpointHealth()   {}
func (cs *ClusterStatus) CheckDbInstanceHealth() {}
func (cs *ClusterStatus) CheckLastBackupHealth() {}
func (cs *ClusterStatus) CheckRPOHealth()        {}
func (cs *ClusterStatus) CheckAreWeHappy()       {}

type HealthChecks []func()

func addChecks() HealthChecks {
	var hc HealthChecks
	hc = append(hc, CheckWalReplicationHealth)
	return hc

}

func initHealthChecks() {
	healthChecks := addChecks()
	for _, hc := range healthChecks {
		go hc()
	}
}

func CheckClusterHealth() {
	initHealthChecks()

	cs := ClusterStatus{}
	//ch := ClusterHealth{}

	status, err := os.ReadFile("status.json")
	if err != nil {
		slog.Error("file", "cluster.go", "failed to read", "status.json", err)
		os.Exit(1)
	}

	err = json.Unmarshal(status, &cs)
	if err != nil {
		slog.Error("file", "cluster.go", "failed to unmarshal status", status, err)
		os.Exit(1)
	}

	/*
		cs.CheckEndpointHealth(status)
		cs.CheckWalReplicationHealth()
		cs.CheckDBInstanaceHealth(status)
		cs.CheckLastBackupHealth(status)
		cs.CheckRPOHealth(status)

		cs.CheckAreWeHappy()
	*/
	fmt.Println(cs)
}
