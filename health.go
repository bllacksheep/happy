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

	Low  int = 1
	High int = 4
)

type ClusterState struct {
	Status int
	Weight int
}

// quicker loop writing some quick tests for this part
func (c *Cluster) IsHappy() {
	// customChecks()

	if c.Status == nil {
		slog.Info("health.go", "msg", fmt.Sprintf("pending response for %s", c.Name))
		return
	}

	clusterStatus := c.Status.Status

	// conditions first
	// some weight potentially
	// can you weight a check so that if it's failing it skews overall health?

	// something like... success contributes less, skew if fail
	clusterStates := map[string]ClusterState{
		"ClusterIsReady":               {Status: Success, Weight: Low},
		"ClusterIsNotReady":            {Status: Fail, Weight: High},
		"ContinuousArchivingSucceeded": {Status: Success, Weight: Low},
		"ContinuousArchivingFailing":   {Status: Fail, Weight: High},
	}

	score := 0
	weight := 0

	// total score but need individual scores numerator
	for _, current_condition := range clusterStatus.Conditions {
		state := clusterStates[current_condition.Reason]
		score += state.Status * state.Weight
		weight += state.Weight
	}

	totalScore := score / weight
	fmt.Println(totalScore)

	// test
	c.WalReplicationHealthy = true
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
