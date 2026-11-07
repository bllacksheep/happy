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
	DefaultChecks    DefaultHealthChecks
	DefaultThreshold Threshold

	Status *ClusterStatus

	CurrentPrimary string `json:"current_primary"`

	// playing with idea of common check interface
	WalReplicationHealth Check `json:"wal_replication_health"`
	DbInstancesHealth    Check `json:"db_instances_health"`
	LastBackupHealth     Check `json:"last_backup_health"`
	RPOHealth            Check `json:"rpo_health"`

	// timer that will track recovery time
	RTO time.Duration `json:"rto"`

	OverallHealthScore float64 `json:"overall_health"`
	Polled             int     `json:"polled"`
	Succeeded          int     `json:"succeeded"`
	Happy              bool    `json:"happy"`
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
		slog.Error("status.go", c.Name, err)
		os.Exit(1)
	}

	var cs ClusterStatus
	err = json.NewDecoder(resp.Body).Decode(&cs)
	if err != nil {
		slog.Error("status.go", c.Name, err)
		os.Exit(1)
	}

	resp.Body.Close()
	c.Status = &cs
}

func newCnpgClusterRequest(name, ns string) *http.Request {
	tokenData, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		slog.Error("status.go", name, err)
		os.Exit(1)
	}
	endpoint_uri := "apis/postgresql.cnpg.io/v1/namespaces/" + ns + "/clusters/" + name
	endpoint := fmt.Sprintf("https://%s:%s/%s", os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT"), endpoint_uri)
	slog.Info("status.go", name, "endpoint "+endpoint)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		slog.Error("status.go", name, err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenData))

	return req
}

func newTlsClient(name string) *http.Client {
	caCertData, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	if err != nil {
		slog.Error("status.go", name, err)
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

func (c *Cluster) StartClusterStatusChecks(ch chan *Cluster) {
	client := newTlsClient(c.Name)
	slog.Info("status.go", c.Name, "created new tls client")
	req := newCnpgClusterRequest(c.Name, c.Namespace)

	ticker := time.NewTicker(c.DefaultThreshold.Every)

	slog.Info("status.go", c.Name, "started watcher")
	for range ticker.C {
		// view refresh
		c.getClusterStatus(client, req)
		ch <- c
		slog.Info("status.go", c.Name, "wait for next poll")
	}
}
