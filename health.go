package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"os"
)

type Result struct {
	Name   string `json:"name"`
	Result any    `json:"result"`
}

const (
	Fail    int = 0
	Success int = 3

	Normal float64 = 1.0
	High   float64 = 4.0
)

type ClusterState struct {
	Status int
	Weight float64
}

// quicker loop writing some quick tests for this part
func (c *Cluster) IsHappy() {
	// customChecks()

	if c.Status == nil {
		slog.Info("health.go", "msg", fmt.Sprintf("pending response for %s", c.Name))
		return
	}

	c.Polled += 1

	clusterStatus := c.Status.Status

	// conditions first
	// some weight potentially
	// can you weight a check so that if it's failing it skews overall health?

	// something like... success contributes less, skew if fail
	clusterStates := map[string]ClusterState{
		"ClusterIsReady":               {Status: Success, Weight: Normal},
		"ClusterIsNotReady":            {Status: Fail, Weight: High},
		"ContinuousArchivingSucceeded": {Status: Success, Weight: Normal},
		"ContinuousArchivingFailing":   {Status: Fail, Weight: High},
		"LastBackupSucceeded":          {Status: Success, Weight: Normal},
		"LastBackupFailed":             {Status: Fail, Weight: High},
	}

	score := 0.0
	weight := 0.0

	// total score but need individual scores numerator
	for _, current_condition := range clusterStatus.Conditions {
		state := clusterStates[current_condition.Reason]
		score += float64(state.Status) * state.Weight
		weight += state.Weight
		switch current_condition.Reason {
		case "ContinuousArchivingSucceeded":
			c.WalReplicationHealthy = true
		}
	}

	totalScore := score / weight
	fmt.Println(totalScore)

	if totalScore == float64(Success) {
		c.Happy += 1
	}
}

// custom checks in here
func customChecks(cluster ClusterConfig) {
	customChecks := cluster.Checks
	numCustomChecks := len(customChecks)
	happy := make([]Result, numCustomChecks)

	db_connection_string := os.Getenv("DB_CONNECTION_URI")
	if db_connection_string == "" {
		// needs to come from status probe
	}

	fmt.Println(db_connection_string)
	conn, err := pgx.Connect(context.Background(), db_connection_string)
	if err != nil {
		slog.Error("health.go", "err", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	for i := 0; i < numCustomChecks; i++ {
		hc := customChecks[i]
		happy[i].Name = hc.Name
		err = conn.QueryRow(context.Background(), hc.Statement).Scan(&happy[i].Result)
		if err != nil {
			slog.Error("health.go", "err", err)
		}
	}

	jsonData, err := json.MarshalIndent(happy, "", " ")
	if err != nil {
		slog.Error("health.go", "err", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
	slog.Info("main.go", "msg", jsonData)
}
