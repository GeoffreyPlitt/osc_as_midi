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

var debugMain = debuggo.Debug("main")

func main() {
	// Command-line flags
	var (
		oscPort       = flag.Int("osc-port", 9000, "UDP port for OSC messages")
		clientName    = flag.String("client-name", "osc-midi-bridge", "JACK client name")
		portName      = flag.String("port-name", "midi_out", "JACK MIDI output port name")
		listPorts     = flag.Bool("list-ports", false, "List available MIDI ports and exit")
		oscTargetHost = flag.String("osc-target-host", "localhost", "Target host for outgoing OSC messages")
		oscTargetPort = flag.Int("osc-target-port", 8000, "Target port for outgoing OSC messages")
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
	bridge, err := NewBridge(*oscPort, *clientName, *portName, *oscTargetHost, *oscTargetPort)
	if err != nil {
		log.Fatal(err)
	}

	// Setup signal handling
	setupSignalHandler()

	// Start the bridge
	debugMain("Starting OSC-MIDI bridge on port %d", *oscPort)
	fmt.Printf("OSC-MIDI Bridge started\n")
	fmt.Printf("  OSC Port: %d\n", *oscPort)
	fmt.Printf("  JACK Client: %s\n", *clientName)
	fmt.Printf("  MIDI Port: %s\n", *portName)

	if err := bridge.Start(); err != nil {
		log.Fatal(err)
	}
}

func setupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("Received SIGTERM, exiting.")
		os.Exit(0)
	}()
}
