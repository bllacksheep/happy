package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type Config struct {
	Happy struct {
		Cluster struct {
			Name     string `yaml:"name"`
			Backups  string `yaml:"backups"`
			Uri      string `yaml:"uri"`
			Strategy string `yaml:"strategy"`
		} `yaml:"cluster"`
		HealthChecks struct {
			CustomChecks []struct {
				Name      string `yaml:"name"`
				Statement string `yaml:"statement"`
				Return    string `yaml:"return"`
			} `yaml:"custom_checks"`
		} `yaml:"healthchecks"`
	} `yaml:"happy"`
}

type Happy struct {
	Name   string `json:"name"`
	Result any    `json:"result"`
}

func validate() Config {
	config := Config{}
	configPath := "happy.yaml"

	configData, err := os.ReadFile(configPath)

	if err != nil {
		slog.Error("file", "main.go", "couldn't read file", configPath, err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		slog.Error("file", "main.go", "error unmarshalling yaml", configPath, err)
		os.Exit(1)
	}

	// check required feilds here


	return config
}

func main() {
	config := validate()

	customChecks := config.Happy.HealthChecks.CustomChecks
	numCustomChecks := len(customChecks)
	happy := make([]Happy, numCustomChecks) 

	db_connection_string := os.Getenv("DB_CONNECTION_URI")
	if db_connection_string == "" {
		db_connection_string = config.Happy.Cluster.Uri
	}

	conn, err := pgx.Connect(context.Background(), db_connection_string)
	if err != nil {
		slog.Error("file", "main.go", "error connecting to db", db_connection_string, err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	for i := 0; i < numCustomChecks; i++ {
		hc := customChecks[i]
		happy[i].Name = hc.Name
		err = conn.QueryRow(context.Background(), hc.Statement).Scan(&happy[i].Result)
		if err != nil {
			slog.Error("file", "main.go", "error performing query", hc.Statement, err)
		}
	}

	jsonData, err := json.MarshalIndent(happy, "", " ")
	if err != nil {
		slog.Error("file", "main.go", "error marshing json", "error", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
	slog.Info("file", "main.go", "health", "results", jsonData)
	CheckClusterHealth()
}
