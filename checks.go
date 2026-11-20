package main

import "time"

type Check struct {
	Name string `json:"name"`

	Slo  float64
	Rise int
	Fall int

	Polled   int
	Suceeded int

	Healthy   bool `json:"healthy"`
	Enabled   bool `json:"enabled"`
	IsDefault bool `json:"is_default"`

	// needs to only be in sql checks
	Statement string `json:"statement"`

	Threshold Threshold

	LastFailedAt  time.Time `json:"last_failed_at"`
	LastCheckedAt time.Time `json:"last_checked_at"`
}

type Poller interface {
	Poll() error
	MetadataResolver
}

type MetadataResolver interface {
	Owner() string
}

type Metadata struct {
	ClusterName string
}

func (cm Metadata) Owner() string {
	return cm.ClusterName
}

type HttpCheck struct {
	Url string

	MetadataResolver
	Check
}

type SqlCheck struct {
	Statement        string
	ConnectionString string

	MetadataResolver
	Check
}

func (c HttpCheck) Poll() error {
	return nil
}

func (c SqlCheck) Poll() error {
	return nil
}

// defaultChecks returns unique copied slice per cluster config
// of checks that satisfy the Poller interface with a default threshold
// call is free to murate this these checks at runtime
func defaultChecks(cfg ClusterConfig, endpoint string) []Poller {
	// endpoint for default checks is shared cluster endpoint
	// backup checks will need a differnt endpoint

	th := Threshold{
		Every: time.Second * 5,
		Rise:  3,
		Fall:  2,
	}

	// value in interface, can't mutate the config
	// can't mutate cluster but can be used to lookup the cluster by name
	return []Poller{
		HttpCheck{
			Name:      "default_wal_replication_check_enabled",
			Url:       endpoint,
			IsDefault: true,
			Threshold: th,
			Enabled:   cfg.CheckWalReplicationHealth,

			MetadataResolver: Metadata{ClusterName: cfg.Name},
		},
		HttpCheck{
			Name:      "default_backups_check_enabled",
			Url:       endpoint,
			IsDefault: true,
			Threshold: th,
			Enabled:   cfg.CheckLastBackupHealth,

			MetadataResolver: Metadata{ClusterName: cfg.Name},
		},
		HttpCheck{
			Name:      "default_rpo_check_enabled",
			Url:       endpoint,
			IsDefault: true,
			Threshold: th,
			Enabled:   cfg.CheckRPOHealth,

			MetadataResolver: Metadata{ClusterName: cfg.Name},
		},
		HttpCheck{
			Name:      "default_db_instance_check_enabled",
			Url:       endpoint,
			IsDefault: true,
			Threshold: th,
			Enabled:   cfg.CheckDbInstanceHealth,

			MetadataResolver: Metadata{ClusterName: cfg.Name},
		},
	}
}

// startchecks

// customchecks
