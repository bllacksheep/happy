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

type Condition struct {
		LastTransitionTime time.Time
		Message            string
		Reason             string
		Status             string
		Type               string
}

// build up status
type Status struct {
	Conditions []Condition
}

// full state swap here
type Cluster struct {
	mu sync.RWMutex
	Status
}

type HealthCheck func()
type HealthChecks []HealthCheck

type kubeApiServerClient struct {
	Headers map[string]string
	*http.Client
}

func getToken() string {
	return os.Getenv("TOKEN")
}

func getCACert() string {
	return os.Getenv("CACERT")
}

func cnpgClusterStatusUri() string {
	return os.Getenv("CNPG_CLUSTER_OBJECT_EP")
}

func newClient() kubeApiServerClient {
	return kubeApiServerClient{
		map[string]string{"Authorization: Bearer ": getToken(), "CACERT": getCaCert},
		&http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (k kubeApiServerClient) newRequest(url string, method string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		slog.Error("file", "cluster.go", "error building new request")
		os.Exit(1)
	}

	for k, v := range k.Headers {
		req.Header.Add(k, v)
	}
	return req
}

func (c *Cluster) getClusterStatus(client kubeApiServerClient, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		req := client.newRequest(cnpgClusterStatusUri(), http.MethodGet)
		resp, err: = client.Do(req)
		if err != nil {
			slog.Error("file", "cluster.go", "error reachign endpoint")
			os.Exit(1)
		}
	}
	var cs ClusterStatus
	if err := json.NewDecoder(resp.Body).Decode(&cs); err == nil {
			s.mu.Lock()
			c.Status = cs
			s.mu.Unlock()
		}
	resp.Body.Close()	
}

func CheckWalReplicationHealth() {
	resp := getClusterStatus()

	slog.Info("file", "cluster.go", "checking wal replication health")

	cs := ClusterStatus{}
	err = json.Unmarshal(resp.Body, &cs)
	if err != nil {
		slog.Error("file", "cluster.go", "failed to unmarshal status", status, err)
		os.Exit(1)
	}

	fmt.Println(cs)

}

func (cs *ClusterStatus) CheckEndpointHealth()   {
	// validate certs
	resp := getClusterStatus()
}

func (cs *ClusterStatus) CheckDbInstanceHealth() {}
func (cs *ClusterStatus) CheckLastBackupHealth() {}
func (cs *ClusterStatus) CheckRPOHealth()        {}
func (cs *ClusterStatus) CheckAreWeHappy()       {}



func getBuiltinChecks() HealthChecks {
	var hc HealthChecks
	hc = append(hc, CheckWalReplicationHealth)
	return hc
}

func Poll(threadID int, hc HealthCheck) {
	for {
		msg := hc()
		/*
		Progress(threadID, msg)
		if msg.Status == "complete" {
			Progress(threadID, msg)
			return
		}
		*/

	}
}

func initHealthChecks() {
	builtinHealthChecks := getBuiltinChecks()

	//ui.Clear()
	
	for id, hc := range builtinHealthChecks {
		go func(tid int, check HealthCheck) {
			Poll(tid, check)
			wg.Done()
		}(id, hc)
	}

	//defer Return()
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
