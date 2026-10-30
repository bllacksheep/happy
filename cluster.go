package main

import (
	"fmt"
	"json"
	"os"
)

type ClusterStatus struct {
	WalReplicationHealthy bool
	DbInstancesHealthy bool
	LastBackupHealthy bool
	RPOHealthy bool

	Happy bool
}

func (ClusterStatus *cs) CheckEndpointHealth() {}
func (ClusterStatus *cs) CheckWalReplicationHealth() {}
func (ClusterStatus *cs) CheckDbInstanceHealth() {}
func (ClusterStatus *cs) CheckLastBackupHealth() {}
func (ClusterStatus *cs) CheckRPOHealth() {}
func (ClusterStatus *cs) CheckAreWeHappy() {}


func CheckClusterHealth() {
	cs := &ClusterHealth{}

	status, err := os.ReadFile("status.json")

	cs.CheckEndpointHealth(status)
	cs.CheckWalReplicationHealth(status)
	cs.CheckDBInstanaceHealth(status)
	cs.CheckLastBackupHealth(status)
	cs.CheckRPOHealth(status)

	cs.CheckAreWeHappy()

	fmt.Println("watching")
}
