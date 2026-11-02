package main

import (
	"crypto/tls"
	"crypto/x509"
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

type Condition struct {
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	Message            string    `json:"message"`
	Reason             string    `json:"reason"`
	Status             string    `json:"status"`
	Type               string    `json:"type"`
}

// api response
type ClusterStatus struct {
	Status struct {
		Conditions []Condition `json:"conditions"`
	} `json:"status"`
}

func getClusterStatus(ns, name string, interval time.Duration) {
	//ch := make(chan ClusterStatus)

	tokenData, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		slog.Error("status.go", "msg", err)
		os.Exit(1)
	}

	caCertData, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	if err != nil {
		slog.Error("status.go", "msg", err)
		os.Exit(1)
	}
	endpoint_uri := fmt.Sprintf("apis/postgresql.cnpg.io/v1/namespaces/%s/clusters/%s", ns, name)
	endpoint := fmt.Sprintf("https://%s:%s/%s", os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT"), endpoint_uri)

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertData)

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		slog.Error("status.go", "msg", err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenData))

	ticker := time.NewTicker(interval)

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 30,
	}

	for range ticker.C {
		slog.Info("status.go", "msg", fmt.Sprintf("get endpoint %s", endpoint))
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("status.go", "msg", err)
			os.Exit(1)
		}
		var cs ClusterStatus
		if err := json.NewDecoder(resp.Body).Decode(&cs); err == nil {
			slog.Info("status.go", "msg", fmt.Sprintf("status: %s", cs))
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
