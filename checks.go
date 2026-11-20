package main

import "time"

type Check struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`

	Slo  float64
	Rise int
	Fall int

	Polled   int
	Suceeded int

	Healthy   bool   `json:"healthy"`
	Enabled   bool   `json:"enabled"`
	IsDefault bool   `json:"is_default"`
	Statement string `json:"statement"`

	Threshold Threshold

	LastFailedAt  time.Time `json:"last_failed_at"`
	LastCheckedAt time.Time `json:"last_checked_at"`
}

// defaultChecks returns unique copied slice per cluster config
// of checks that satisfy the Poller interface with a default threshold
// call is free to murate this these checks at runtime
func defaultChecks(cfg ClusterConfig) []Poller {

	th := Threshold{
		Every: time.Second * 5,
		Rise:  3,
		Fall:  2,
	}

	// value in interface, can't mutate the config
	return []Poller{
		HttpCheck{
			Name:        "default_wal_replication_check_enabled",
			Url:         "",
			ClusterName: cfg.Name,
			IsDefault:   true,
			Threshold:   th,
			Enabled:     cfg.CheckWalReplicationHealth,
		},
		HttpCheck{
			Name:        "default_backups_check_enabled",
			ClusterName: cfg.Name,
			IsDefault:   true,
			Threshold:   th,
			Enabled:     cfg.CheckLastBackupHealth,
		},
		HttpCheck{
			Name:        "default_rpo_check_enabled",
			ClusterName: cfg.Name,
			IsDefault:   true,
			Threshold:   th,
			Enabled:     cfg.CheckRPOHealth,
		},
		HttpCheck{
			Name:        "default_db_instance_check_enabled",
			ClusterName: cfg.Name,
			IsDefault:   true,
			Threshold:   th,
			Enabled:     cfg.CheckDbInstanceHealth,
		},
	}
}

// startchecks

// customchecks
