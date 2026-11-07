package main

import (
	"log/slog"
	"strconv"
	"time"
)

// trying a happy interface
// all 'checks' will satisfy this
type Check struct {
	Name string `json:"name"`

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

type Health interface {
	Happy(func(string) bool) bool
	Toggle(string, bool)
}

// custom health check router?
func (c Check) Happy(name string, fn func(string) bool) bool {
	return fn(c.Statement)
}

func (c *Check) Toggle(name string, toggle bool, isdef bool) {
	c.Name = name
	c.Enabled = toggle
	c.IsDefault = isdef
	slog.Info("happy.go", c.Name, strconv.FormatBool(c.Enabled))
}
