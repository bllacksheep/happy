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

// cluster status for n clusters
type Cluster struct {
	Name          string
	ClusterStatus chan ClusterStatus

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

// generic kube-api-server response
type ClusterStatus struct {
	Status struct {
		Conditions []Condition `json:"conditions"`
	} `json:"status"`
}

func getClusterStatus(client *http.Client, endpoint string, ch chan ClusterStatus) {
	slog.Info("status.go", "msg", fmt.Sprintf("get endpoint %s", endpoint))

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("status.go", "msg", err)
		os.Exit(1)
	}

	var cs ClusterStatus
	if err := json.NewDecoder(resp.Body).Decode(&cs); err == nil {
		slog.Info("status.go", "msg", fmt.Sprintf("status: %s", cs))
	}
	resp.Body.Close()
	ch <- cs
}

func cnpgClusterRequest(ns, name string) *http.Request {
	tokenData, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		slog.Error("status.go", "msg", err)
		os.Exit(1)
	}
	endpoint_uri := fmt.Sprintf("apis/postgresql.cnpg.io/v1/namespaces/%s/clusters/%s", ns, name)
	endpoint := fmt.Sprintf("https://%s:%s/%s", os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT"), endpoint_uri)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		slog.Error("status.go", "msg", err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenData))

	return req
}

func newTlsClient() *http.Client {
	caCertData, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	if err != nil {
		slog.Error("status.go", "msg", err)
		os.Exit(1)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertData)

	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{
		Transport: transport,
		Timeout:   time.Second * 30,
	}
}

func (c Cluster) StartCheckClusterStatus() {
	client := newTlsClient()
	req := cnpgClusterRequest(c.ns, c.name)

	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		getClusterStatus(client, req, c.ClusterStatus)
		/*
			data err := json.Marshal(status)
			if err != nil {
				slog.Error("file", "status.go", "failed to unmarshal status", status, err)
				os.Exit(1)
			}
		*/
	}
}
