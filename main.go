package main

import (
	"log/slog"
	"os"
)

func main() {
	happy, err := NewHappy()
	if err != nil {
		slog.Error("failed to create happy instance", "error", err)
		os.Exit(1)
	}

	if err := happy.Start(); err != nil {
		slog.Error("failed to start happy instance", "error", err)
		os.Exit(1)
	}
}
