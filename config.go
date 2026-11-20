package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type ClusterConfig struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`

	CheckClusterEndpointHealth bool `yaml:"cluster_endpoint_health"`
	CheckWalReplicationHealth  bool `yaml:"wal_replication_health"`
	CheckLastBackupHealth      bool `yaml:"last_backup_health"`
	CheckDbInstanceHealth      bool `yaml:"db_instance_health"`
	CheckRPOHealth             bool `yaml:"rpo_health"`

	Checks []Check `yaml:"checks"`
}

type Config struct {
	Clusters []ClusterConfig `yaml:"clusters"`
}

// // custom health check router?
//
//	func (c Check) Happy(name string, fn func(string) bool) bool {
//		return fn(c.Statement)
//	}
//
//	func (c *Check) Toggle(name string, toggle bool, isdef bool) {
//		c.Name = name
//		c.Enabled = toggle
//		c.IsDefault = isdef
//		slog.Info("happy.go", c.Name, strconv.FormatBool(c.Enabled))
//	}

func (h *Happy) validateConfig() error {
	configPath := flag.String("c", "happy.yaml", "config file")
	flag.Parse()
	configData, err := os.ReadFile(*configPath)

	if err != nil {
		return fmt.Errorf("couldn't read file", err)
	}

	err = yaml.Unmarshal(configData, &h.config)
	if err != nil {
		return fmt.Errorf("error unmarhsalling yaml", err)
	}

	if len(h.config.Clusters) == 0 {
		return fmt.Errorf("no clusters have been configured")
	}

	// check required feilds here
	return nil
}
