package main

import (
	"log/slog"
	"os"
)

// want shared default config and threshold pointer

func main() {
	happy, err := newHappy()
	if err != nil {
		slog.Error("failed to create happy instance", "error", err)
		os.Exit(1)
	}

	if err := happy.Start(); err != nil {
		slog.Error("failed to start happy instance", "error", err)
		os.Exit(1)
	}
}
