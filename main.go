package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"time"
)

type ClusterConfig struct {
	Name         string `yaml:"name"`
	Namespace    string `yaml:"namespace"`
	UsePooler    bool   `yaml:"use_pooler"`
	UsePrimary   bool   `yaml:"use_primary"`
	UseSecondary bool   `yaml:"use_secondary"`
	UseInstance  string `yaml:"use_instance"`
	Checks       []struct {
		Name      string        `yaml:"name"`
		Statement string        `yaml:"statement"`
		Enabled   bool          `yaml:"enabled"`
		Every     time.Duration `yaml:"every"`
		Rise      int           `yaml:"rise"`
		Fall      int           `yaml:"fall"`
	} `yaml:"checks"`
}

type DefaultHealthChecks struct {
	CheckClusterEndpointHealth bool `yaml:"check_cluster_endpoint_health"`
	CheckWalReplicationHealth  bool `yaml:"check_wal_replication_health"`
	CheckLastBackupHealth      bool `yaml:"check_last_backup_health"`
	CheckDbInstanceHealth      bool `yaml:"check_db_instance_health"`
	CheckRPOHealth             bool `yaml:"check_rpo_health"`
}

type Config struct {
	Happy struct {
		Default struct {
			HealthChecks DefaultHealthChecks `yaml:"healthchecks"`
			Threshold    Threshold           `yaml:"threshold"`
		} `yaml:"defaults"`
		Clusters []ClusterConfig `yaml:"clusters"`
	} `yaml:"happy"`
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

	slog.Info("main.go", "msg", "validated config")
	// check required feilds here
	return happyConfig
}

func (c *Cluster) StartDefaultChecks() {
	c.WalReplicationHealth.Toggle("wal_replication_check_enabled",
		c.DefaultChecks.CheckWalReplicationHealth, true)
	c.LastBackupHealth.Toggle("backups_check_enabled",
		c.DefaultChecks.CheckLastBackupHealth, true)
	c.RPOHealth.Toggle("rpo_check_enabled",
		c.DefaultChecks.CheckRPOHealth, true)
	c.DbInstancesHealth.Toggle("db_instance_check_enabled",
		c.DefaultChecks.CheckDbInstanceHealth, true)

	// go startDefaultChecks?
}

func startHappy() {
	cfg := validateConfig()

	ch := make(chan *Cluster)

	// should scale for n clusters
	for _, clusterConfig := range cfg.Happy.Clusters {
		cluster := Cluster{
			ClusterConfig:    clusterConfig,
			DefaultThreshold: cfg.Happy.Default.Threshold,
			DefaultChecks:    cfg.Happy.Default.HealthChecks,
		}
		// is run for every cluster not required but will do for now
		cluster.StartDefaultChecks()

		go cluster.StartClusterStatusChecks(ch)
		go cluster.StartCustomChecks()

		slog.Info("main.go", cluster.Name, "happy started")
	}

	happy := make(map[string]*Cluster, len(cfg.Happy.Clusters))

	slog.Info("main.go", "msg", "happy started")

	for cluster := range ch {
		cluster.IsHappy()

		// swap cluster state by name
		happy[cluster.Name] = cluster
		currentCluster := happy[cluster.Name]

		slog.Info("main.go", currentCluster.Name, *happy[cluster.Name].Status)

		slo := 100 * (float64(currentCluster.Succeeded) / float64(currentCluster.Polled))

		slog.Info("main.go", currentCluster.Name,
			fmt.Sprintf("succeeded/polled: %v/%v, happy %v, slo: %.3f%%",
				currentCluster.Succeeded, currentCluster.Polled,
				currentCluster.Happy, slo))
	}
}

func main() {
	startHappy()
}
