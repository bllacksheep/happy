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

type Status int
type Weight int

const (
	Fail    Status = 0
	Success Status = 3

	Low     Weight = 1
	High    Weight = 4
)

type ClusterState struct {
	Ret    Status
	Weight Weight
}

// quicker loop writing some quick tests for this part
func (c *Cluster) IsHappy() {
	// customChecks()

	if c.Status == nil {
		slog.Info("health.go", "msg", fmt.Sprintf("pending response for %s", c.Name))
		return
	}

	//status := c.Status.Status

	// conditions first
	// some weight potentially
	// can you weight a check so that if it's failing it skews overall health?

	// something like... success contributes less, skew if fail
	states := map[string]ClusterState{
		"ClusterIsReady":               {Ret: Success, Weight: Low},
		"ClusterIsNotReady":            {Ret: Fail, Weight: High},
		"ContinuousArchivingSucceeded": {Ret: Success, Weight: Low},
		"ContinuousArchivingFailing":   {Ret: Fail, Weight: High},
	}

	_ = states

	/*
		score := 0.0
		for _, condition := range status.Conditions {
			score += weights[condition.Reason]
		}
	*/

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
