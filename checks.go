package main


type Check struct {
	Name string `json:"name"`
	ClusterName string `json:"cluster_name"`

	Slo  float64
	Rise int
	Fall int

	Polled    int
	Suceeded  int

	Healthy   bool   `json:"healthy"`
	Enabled   bool   `json:"enabled"`
	IsDefault bool   `json:"is_default"`
	Statement string `json:"statement"`

	Every     time.Duration
	Threshold time.Duration

	LastFailedAt  time.Time `json:"last_failed_at"`
	LastCheckedAt time.Time `json:"last_checked_at"`
}

type SharedDefaultConfig struct {
	Checks    []Poller
	Threshold Threshold
}

type Defaults interface {
	SharedDefaultChecks() []Poller
	SharedDefaultThreshold() Threshold
}

func (sc *SharedDefaultConfig) SharedDefaultChecks() []Poller {
	sc.Checks = defaultCheck(name, )
	return sc.Checks
}

func (sc *SharedDefaultConfig) SharedDefaultThreshold() Threshold {
	return sc.Threshold
}

func defaultCheck(cfg ClusterConfig) []Poller {
	return []Poller{
		&HttpCheck{
			Name:        "default_wal_replication_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   cfg.Threshold,
			Enabled:     cfg.CheckWalReplicationHealth,
		},
		&HttpCheck{
			Name:        "default_backups_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   cfg.Threshold,
			Enabled:     cfg.CheckLastBackupHealth,
		},
		&HttpCheck{
			Name:        "default_rpo_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   cfg.Threshold, 
			Enabled:     cfg.CheckRPOHealth,
		},
		&HttpCheck{
			Name:        "default_db_instance_check_enabled",
			ClusterName: c.Name,
			IsDefault:   true,
			Threshold:   cfg.Threshold, 
			Enabled:     ,
		},
	}
}


// load the checks into the app state
func (c *Cluster) registerChecks(dfg SharedDefaultConfig) {
	c.Checks = append(c.Checks, defaultCheck(c.ClusterName, dfg))
}

startchecks

customchecks

defaultchecks

registerchecks
