package main

import (
	"log/slog"
	"strconv"
	"time"
)

// trying a happy interface
type Check struct {
	Name        string    `json:"name"`
	Healthy     bool      `json:"healthy"`
	Enabled     bool      `json:"enabled"`
	LastChecked time.Time `json:"last_checked"`
	// other
}

type Health interface {
	Happy(string, func() bool) bool
	Toggle(bool)
}

// custom health check router?
func (c Check) Happy(name string, fn func() bool) bool {
	// stuff
	return fn()
}

func (c *Check) Toggle(name string, toggle bool) {
	c.Name = name
	c.Enabled = toggle
	slog.Info("happy.go", c.Name, strconv.FormatBool(c.Enabled))
}
