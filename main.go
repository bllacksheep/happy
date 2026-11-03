package main

import (
	"flag"
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
		} `yaml:"default"`
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

	// check required feilds here

	return happyConfig
}

type HappyView struct {
	Name      string        `json:"name"`
	Namespace string        `json:"namespace"`
	Status    ClusterStatus `json:"cluster_status"`
}

func main() {
	config := validateConfig()
	ch := make(chan Cluster)

	// will scale for n clusters
	for _, cluster := range config.Happy.Clusters {
		clusterView := Cluster{
			ClusterConfig:    cluster,
			DefaultThreshold: config.Happy.Default.Threshold,
			DefaultChecks:    config.Happy.Default.HealthChecks,
		}
		go clusterView.StartClusterStatusChecks(ch)

		//Happy(cluster)
	}

	happy := make(map[string]*Cluster, len(config.Happy.Clusters))

	for {
		select {
		case cluster := <-ch:
			happy[cluster.Name] = &cluster
			slog.Info("main.go", "msg", happy)
		}
	}

}
