package main

import (
	_ "flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create bridge instance
	bridge := &Bridge{}

	// Setup signal handling
	setupSignalHandler(bridge)

	// TODO: Parse flags and start bridge
}

func setupSignalHandler(bridge *Bridge) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		// TODO: Add debug logging when implemented
		bridge.Cleanup()
		os.Exit(0)
	}()
}
