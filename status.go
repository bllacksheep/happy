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

type Threshold struct {
	Every time.Duration `json:"every" yaml:"every"`
	Rise  int           `json:"rise" yaml:"rise"`
	Fall  int           `json:"fall" yaml:"fall"`
}

// cluster status for n clusters
type Cluster struct {
	ClusterConfig
	DefaultThreshold Threshold
	DefaultChecks    DefaultHealthChecks

	Status *ClusterStatus

	WalReplicationHealthy bool `json:"wal_replication_healthy"`
	DbInstancesHealthy    bool `json:"db_instances_healthy"`
	LastBackupHealthy     bool `json:"last_backup_healthy"`
	RPOHealthy            bool `json:"rpo_healthy"`

	// timer that will track recovery time
	RTO time.Duration `json:"rto"`

	Happy         int     `json:"happy"`
	OverallHealth float64 `json:"overall_health"`
	Polled        int     `json:"polled"`
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

func (c *Cluster) getClusterStatus(client *http.Client, r *http.Request) {
	resp, err := client.Do(r)
	if err != nil {
		slog.Error("status.go", "msg", err)
		os.Exit(1)
	}

	var cs ClusterStatus
	err = json.NewDecoder(resp.Body).Decode(&cs)
	if err != nil {
		slog.Error("status.go", "msg", err)
		os.Exit(1)
	}

	resp.Body.Close()
	c.Status = &cs
}

func newCnpgClusterRequest(name, ns string) *http.Request {
	tokenData, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		slog.Error("status.go", "msg", err)
		os.Exit(1)
	}
	endpoint_uri := fmt.Sprintf("apis/postgresql.cnpg.io/v1/namespaces/%s/clusters/%s", ns, name)
	endpoint := fmt.Sprintf("https://%s:%s/%s", os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT"), endpoint_uri)
	slog.Info("status.go", "msg", fmt.Sprintf("build endpoint %s", endpoint))

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

	slog.Info("status.go", "msg", "created new tls client %s")

	return &http.Client{
		Transport: transport,
		Timeout:   time.Second * 30,
	}
}

func (c Cluster) StartClusterStatusChecks(ch chan Cluster) {
	client := newTlsClient()
	req := newCnpgClusterRequest(c.Name, c.Namespace)

	ticker := time.NewTicker(c.DefaultThreshold.Every)

	slog.Info("status.go", "msg", fmt.Sprintf("watcher starting for %s", c.Name))
	for range ticker.C {
		c.getClusterStatus(client, req)
		ch <- c
		slog.Info("status.go", "msg", "wait for next poll")
	}
}
