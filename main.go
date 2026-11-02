package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
)

type ClusterConfig struct {
	Name         string `yaml:"name"`
	Namespace    string `yaml:"namespace"`
	UsePooler    bool   `yaml:"use_pooler"`
	UsePrimary   bool   `yaml:"use_primary"`
	UseSecondary bool   `yaml:"use_secondary"`
	UseInstance  string `yaml:"use_instance"`
	Checks       []struct {
		Name      string `yaml:"name"`
		Statement string `yaml:"statement"`
		Enabled   bool   `yaml:"enabled"`
		Every     string `yaml:"every"`
		Rise      string `yaml:"rise"`
		Fall      string `yaml:"fall"`
	} `yaml:"checks"`
}

type Config struct {
	Happy struct {
		Default struct {
			HealthChecks struct {
				CheckClusterEndpointHealth bool `yaml:"check_cluster_endpoint_health"`
				CheckWalReplicationHealth  bool `yaml:"check_wal_replication_health"`
				CheckLastBackupHealth      bool `yaml:"check_last_backup_health"`
				CheckDbInstanceHealth      bool `yaml:"check_db_instance_health"`
				CheckRPOHealth             bool `yaml:"check_rpo_health"`
			} `yaml:"healthchecks"`
			Threshold struct {
				Every string `yaml:"every"`
				Rise  string `yaml:"rise"`
				Fall  string `yaml:"fall"`
			}
		} `yaml:"default"`
		Clusters []ClusterConfig `yaml:"clusters"`
	} `yaml:"happy"`
}

type Happy struct {
	Name   string `json:"name"`
	Result any    `json:"result"`
}

// return valid config
func validateConfig() Config {
	happyConfig := Config{}

	configPath := flag.String("c", "happy.yaml", "config file")
	flag.Parse()
	configData, err := os.ReadFile(*configPath)

	if err != nil {
		slog.Error("main.go", "msg", "couldn't read file", "err", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(configData, &happyConfig)
	if err != nil {
		slog.Error("main.go", "msg", "error unmarshalling yaml", "err", err)
		os.Exit(1)
	}

	// check required feilds here

	return happyConfig
}

func main() {
	config := validateConfig()

	// will scale for n clusters
	for cluster := range config.Happy.Clusters {
		c := Cluster{
			Name:      cluster.Name,
			Namespace: cluster.Namespace,
			// use pooler
		}
		go c.StartClusterStatusChecks()
	}

	// custom checks below this line

	customChecks := config.Happy.Clusters[0].Checks
	numCustomChecks := len(customChecks)
	happy := make([]Happy, numCustomChecks)

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

	// block
	select {}
}
