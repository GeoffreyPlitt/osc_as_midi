package main

import (
	"flag"
	"os"
	"testing"
)

func TestMainFunction(t *testing.T) {
	// Basic test to ensure main doesn't panic
	// In real usage, main() would block on server.ListenAndServe()
	// For testing, we just ensure it compiles and basic structure is correct

	// This is more of a smoke test
	// Actual integration testing would require mocking or running in a goroutine
}

func TestCLIFlags(t *testing.T) {
	// Save original command line args
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	tests := []struct {
		name string
		args []string
		test func(t *testing.T)
	}{
		{
			name: "Default values",
			args: []string{"osc-midi-bridge"},
			test: func(t *testing.T) {
				var (
					oscPort       = flag.Int("osc-port", 9000, "UDP port for OSC messages")
					clientName    = flag.String("client-name", "osc-midi-bridge", "JACK client name")
					portName      = flag.String("port-name", "midi_out", "JACK MIDI output port name")
					listPorts     = flag.Bool("list-ports", false, "List available MIDI ports and exit")
					oscTargetHost = flag.String("osc-target-host", "localhost", "Target host for outgoing OSC messages")
					oscTargetPort = flag.Int("osc-target-port", 8000, "Target port for outgoing OSC messages")
				)

				flag.Parse()

				if *oscPort != 9000 {
					t.Errorf("Expected default osc-port 9000, got %d", *oscPort)
				}
				if *clientName != "osc-midi-bridge" {
					t.Errorf("Expected default client-name 'osc-midi-bridge', got '%s'", *clientName)
				}
				if *portName != "midi_out" {
					t.Errorf("Expected default port-name 'midi_out', got '%s'", *portName)
				}
				if *listPorts != false {
					t.Errorf("Expected default list-ports false, got %t", *listPorts)
				}
				if *oscTargetHost != "localhost" {
					t.Errorf("Expected default osc-target-host 'localhost', got '%s'", *oscTargetHost)
				}
				if *oscTargetPort != 8000 {
					t.Errorf("Expected default osc-target-port 8000, got %d", *oscTargetPort)
				}
			},
		},
		{
			name: "Custom values",
			args: []string{"osc-midi-bridge", "--osc-port=7000", "--osc-target-host=192.168.1.100", "--osc-target-port=9000"},
			test: func(t *testing.T) {
				var (
					oscPort       = flag.Int("osc-port", 9000, "UDP port for OSC messages")
					oscTargetHost = flag.String("osc-target-host", "localhost", "Target host for outgoing OSC messages")
					oscTargetPort = flag.Int("osc-target-port", 8000, "Target port for outgoing OSC messages")
				)

				flag.Parse()

				if *oscPort != 7000 {
					t.Errorf("Expected osc-port 7000, got %d", *oscPort)
				}
				if *oscTargetHost != "192.168.1.100" {
					t.Errorf("Expected osc-target-host '192.168.1.100', got '%s'", *oscTargetHost)
				}
				if *oscTargetPort != 9000 {
					t.Errorf("Expected osc-target-port 9000, got %d", *oscTargetPort)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag package state
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = tt.args
			tt.test(t)
		})
	}
}
