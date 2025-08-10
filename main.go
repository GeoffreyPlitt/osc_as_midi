package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/GeoffreyPlitt/debuggo"
)

var debugMain = debuggo.Debug("osc-midi-bridge:main")

func main() {
	// Command-line flags
	var (
		oscPort    = flag.Int("osc-port", 9000, "UDP port for OSC messages")
		clientName = flag.String("client-name", "osc-midi-bridge", "JACK client name")
		portName   = flag.String("port-name", "midi_out", "JACK MIDI output port name")
		bufferSize = flag.Int("buffer-size", 64, "JACK buffer size in frames (64=1.3ms @ 48kHz)")
		listPorts  = flag.Bool("list-ports", false, "List available MIDI ports and exit")
	)

	flag.Parse()

	// Handle list-ports flag
	if *listPorts {
		if err := ListJackPorts(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	// Create bridge instance
	bridge, err := NewBridge(*oscPort, *clientName, *portName)
	if err != nil {
		log.Fatal(err)
	}

	// Setup signal handling
	setupSignalHandler(bridge)

	// Start the bridge
	debugMain("Starting OSC-MIDI bridge on port %d", *oscPort)
	fmt.Printf("OSC-MIDI Bridge started\n")
	fmt.Printf("  OSC Port: %d\n", *oscPort)
	fmt.Printf("  JACK Client: %s\n", *clientName)
	fmt.Printf("  MIDI Port: %s\n", *portName)
	fmt.Printf("  Buffer Size: %d frames\n", *bufferSize)

	if err := bridge.Start(); err != nil {
		log.Fatal(err)
	}
}

func setupSignalHandler(bridge *Bridge) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		debugMain("Received shutdown signal")
		bridge.Cleanup()
		os.Exit(0)
	}()
}
