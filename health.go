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

func Happy(cluster ClusterConfig) {
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
		slog.Error("main.go", "msg", "error connecting to db", "err", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	for i := 0; i < numCustomChecks; i++ {
		hc := customChecks[i]
		happy[i].Name = hc.Name
		err = conn.QueryRow(context.Background(), hc.Statement).Scan(&happy[i].Result)
		if err != nil {
			slog.Error("main.go", "msg", "error performing query", hc.Statement, err)
		}
	}

	jsonData, err := json.MarshalIndent(happy, "", " ")
	if err != nil {
		slog.Error("main.go", "msg", "error marshing json", "err", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
	slog.Info("main.go", "msg", "health", "results", jsonData)
}
