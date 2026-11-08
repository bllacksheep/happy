package main

import (
	"log/slog"
	"strconv"
	"time"
)

// trying a happy interface
// all 'checks' will satisfy this
type Check struct {
	Name    string `json:"name"`
	Cluster *Cluster

	Slo  float64
	Rise int
	Fall int

	Polled   int
	Suceeded int

	Healthy   bool `json:"healthy"`
	Enabled   bool `json:"enabled"`
	IsDefault bool `json:"is_default"`
	// check type might be defined by http: or sql: yaml parent
	Statement string `json:"statement"`

	Every     time.Duration
	Threshold time.Duration

	LastFailedAt  time.Time `json:"last_failed_at"`
	LastCheckedAt time.Time `json:"last_checked_at"`
}

type Health interface {
	Happy(func(string) bool)
	Toggle(string, bool)
}

// custom health check router?
func (c Check) Happy(check func(string) bool) {
	c.Healthy = check(c.Statement)
	c.LastCheckedAt = time.Now()

	if !c.Healthy {
		c.LastFailedAt = time.Now()
	}
}

func (c *Check) Toggle(name string, toggle bool, isdef bool) {
	c.Name = name
	c.Enabled = toggle
	c.IsDefault = isdef

	if c.IsDefault {
		c.Name = "default_" + c.Name
	}

	slog.Info("happy.go", c.Name, strconv.FormatBool(c.Enabled))
}
