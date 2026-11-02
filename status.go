package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"log/slog"
	"os"
	"time"
	"net/http"
)

// mine
type ClusterHealth struct {
	WalReplicationHealthy bool `json:"wal_replication_healthy"`
	DbInstancesHealthy    bool `json:"db_instances_healthy"`
	LastBackupHealthy     bool `json:"last_backup_healthy"`
	RPOHealthy            bool `json:"rpo_healthy"`

	Happy bool `json:"happy"`
}

// api response
type Status struct {
	Conditions []Condition
}

type Condition struct {
		LastTransitionTime time.Time
		Message            string
		Reason             string
		Status             string
		Type               string
}

type kubeApiServerClient struct {
	Headers map[string]string
	*http.Client
}

// state swap
type Cluster struct {
	mu sync.RWMutex
	Status
}

func newClient() kubeApiServerClient {
	return kubeApiServerClient{
			map[string]string{"Authorization: Bearer ": getToken(), "CACERT": getCaCert()},
			&http.Client{
				Timeout: time.Second * 30,
			},
		}

	}
}

func getClusterStatus(endpoint string, interval time.Duration) {
	ch := make(chan ClusterStatus)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		slog.Error("file", "status.go", "error building new request")
		os.Exit(1)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", getKAPIToken()))
	req.Header.Set("CACERT", getCACert())

	ticker := time.NewTicker(interval)

	for range ticker.C {
		resp, err: = client.Do(req)
		if err != nil {
			slog.Error("file", "cluster.go", "error reachign endpoint")
			os.Exit(1)
		}
		var cs ClusterStatus
		if err := json.NewDecoder(resp.Body).Decode(&cs); err == nil {
			s.mu.Lock()
			c.Status = cs // full state swap
			s.mu.Unlock()
		}
		resp.Body.Close()	
		ch <- cs
	}

}


// need interface status fragement
func (c *Cluster) CheckWalReplicationHealth() {
	c.mu.RLock()
	defer s.mu.RUnlock()
	
	// return status fragment here?
	// update healthcheck to be func() StatusFragement?
	fmt.Println(c.Status.Conditions)
}

/*
// treat these more like the view
func (cs *ClusterStatus) CheckEndpointHealth()   {}
func (cs *ClusterStatus) CheckDbInstanceHealth() {}
func (cs *ClusterStatus) CheckLastBackupHealth() {}
func (cs *ClusterStatus) CheckRPOHealth()        {}
func (cs *ClusterStatus) CheckAreWeHappy()       {}
*/

/*
func getBuiltinChecks() HealthChecks {
	var hc HealthChecks
	hc = append(hc, CheckWalReplicationHealth)
	return hc
}

func initHealthChecks() {
	builtinHealthChecks := getBuiltinChecks()
	
	for id, hc := range builtinHealthChecks {
		go func(tid int, check HealthCheck) {
			// not right anymore
			// what now
			wg.Done()
		}(id, hc)
	}

}
*/

// whole piece
// assumes status populated
func CheckClusterHealth() {
	//initHealthChecks()

	//cs := ClusterStatus{}
	//ch := ClusterHealth{}

	getClusterStatus(url, time.Second * 10)

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
}

