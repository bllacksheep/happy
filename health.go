package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"os"
	"time"
)

type Result struct {
	Name   string `json:"name"`
	Result any    `json:"result"`
}

const (
	Fail    int = 0
	Success int = 3

	Acceptable   float64 = 1.0
	Unacceptable float64 = 4.0
)

type ClusterState struct {
	Status int
	Weight float64
}

// quicker loop writing some quick tests for this part
func (c *Cluster) IsHappy() {
	// customChecks()

	if c.Status == nil {
		slog.Info("health.go", c.Name, "pending response from client")
		return
	}

	c.Polled++
	slog.Info("health.go", c.Name, fmt.Sprintf("poll count: %v", c.Polled))

	clusterStatus := c.Status.Status

	// conditions first
	// some weight potentially
	// can you weight a check so that if it's failing it skews overall health?

	// something like... success contributes less, skew if fail

	// failed or failing? needs something else
	clusterStates := map[string]ClusterState{
		"ClusterIsReady":               {Status: Success, Weight: Acceptable},
		"ClusterIsNotReady":            {Status: Fail, Weight: Unacceptable},
		"ContinuousArchivingSucceeded": {Status: Success, Weight: Acceptable},
		"ContinuousArchivingFailing":   {Status: Fail, Weight: Acceptable},
		"LastBackupSucceeded":          {Status: Success, Weight: Acceptable},
		"LastBackupFailed":             {Status: Fail, Weight: Acceptable},
	}

	score := 0.0
	weight := 0.0

	// total score but need individual scores numerator
	for _, current_condition := range clusterStatus.Conditions {
		state := clusterStates[current_condition.Reason]
		// weighted average
		score += float64(state.Status) * state.Weight
		weight += state.Weight
		switch current_condition.Reason {

		// set individual properties
		case "ContinuousArchivingSucceeded":
			slog.Info("health.go", c.Name, "wal replication is healthy")
			c.WalReplicationHealth.Healthy = true
		}
	}

	c.OverallHealthScore = score / weight
	slog.Info("health.go", c.Name, fmt.Sprintf("health score: %v", c.OverallHealthScore))

	// increment overall happy score if cluster in good health
	if c.OverallHealthScore == float64(Success) {
		c.Happy = true
		c.Succeeded++
		slog.Info("health.go", c.Name, fmt.Sprintf("happy count: %v", c.Succeeded))
		return
	}

	slog.Info("health.go", c.Name, fmt.Sprintf("happy count: %v", c.Succeeded))
	c.Happy = false
}

/*
allow each check to take advantage of % scoring via happy() interface
custom checks will need mutex on cluster struct

need to wrap each check in a Check type - will do later
give check an arbitrary function that executes it's statement?
*/

func (c *Cluster) StartCustomChecks() {

	customChecks := c.ClusterConfig.Checks
	numCustomChecks := len(customChecks)
	happy := make([]Result, numCustomChecks)

	// log connected to
	db_connection_string := os.Getenv("DB_CONNECTION_URI")
	if db_connection_string == "" {
		// needs to come from status probe
	}

	for i, check := range customChecks {
		happy[i].Name = check.Name
		go func() {

			conn, err := pgx.Connect(context.Background(), db_connection_string)
			if err != nil {
				slog.Error("health.go", check.Name, err)
				os.Exit(1)
			}

			defer conn.Close(context.Background())

			ticker := time.NewTicker(check.Every)

			// for loop and start in go routine for each check bounded by its interval
			for range ticker.C {
				// r/w on happy[] map at this point completely unsafe!
				err = conn.QueryRow(context.Background(), check.Statement).Scan(&happy[i].Result)
				if err != nil {
					slog.Error("health.go", check.Name, err)
				}

				jsonData, err := json.MarshalIndent(happy, "", " ")
				if err != nil {
					slog.Error("health.go", check.Name, err)
					os.Exit(1)
				}

				fmt.Println(string(jsonData))
				slog.Info("health.go", c.Name, jsonData)
				//c.getClusterStatus(client, req)
				//ch <- c
				slog.Info("health.go", check.Name, "wait for next poll")
			}
		}()
	}

}
