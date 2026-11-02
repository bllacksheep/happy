package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
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

// api response
type ClusterStatus struct {
	Conditions []Condition
}

type Condition struct {
	LastTransitionTime time.Time
	Message            string
	Reason             string
	Status             string
	Type               string
}

// state swap
type Cluster struct {
	ClusterStatus
}

func getClusterStatus(ns, name string, interval time.Duration) {
	//ch := make(chan ClusterStatus)

	tokenData, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	caCertPath := "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"

	endpoint_uri := fmt.Sprintf("/apis/postgresql.cnpg.io/v1/namespaces/%s/clusters/%s", ns, name)
	endpoint := fmt.Sprintf("https://%s:%s/%s", os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT"), endpoint_uri)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		slog.Error("status.go", "msg", "error building new request")
		os.Exit(1)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenData))
	req.Header.Set("CACERT", caCertPath)

	ticker := time.NewTicker(interval)

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	for range ticker.C {
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("status.go", "msg", "error reaching endpoint")
			os.Exit(1)
		}
		var cs ClusterStatus
		if err := json.NewDecoder(resp.Body).Decode(&cs); err == nil {
			fmt.Println("got: ")
			fmt.Println(cs)
			//c.Status = cs // full state swap
		}
		resp.Body.Close()
		//ch <- cs
	}

}

func CheckClusterHealth() {

	getClusterStatus("db", "cluster-0001", time.Second*10)

	/*
		status, err := os.ReadFile("status.json")
		if err != nil {
			slog.Error("file", "status.go", "failed to read", "status.json", err)
			os.Exit(1)
		}

		err = json.Unmarshal(status, &cs)
		if err != nil {
			slog.Error("file", "status.go", "failed to unmarshal status", status, err)
			os.Exit(1)
		}
	*/
}
