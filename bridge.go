package main

import (
	"errors"
	"fmt"

	"github.com/GeoffreyPlitt/debuggo"
	"github.com/hypebeast/go-osc/osc"
	"github.com/xthexder/go-jack"
)

var debugBridge = debuggo.Debug("bridge")

type MidiEvent struct {
	midiData *jack.MidiData
}

type Bridge struct {
	oscServer   *osc.Server
	jackClient  *jack.Client
	midiOutPort *jack.Port
	eventQueue  chan *MidiEvent
}

func NewBridge(oscPort int, clientName string, portName string) (*Bridge, error) {
	// Create OSC server with dispatcher
	dispatcher := osc.NewStandardDispatcher()
	server := &osc.Server{
		Addr:       fmt.Sprintf(":%d", oscPort),
		Dispatcher: dispatcher,
	}

	// Connect to JACK
	client, status := jack.ClientOpen(clientName, jack.NoStartServer)
	if status != 0 {
		return nil, fmt.Errorf("cannot connect to JACK server (status %d): %s\n\nPlease ensure JACK is running. Start it with:\n  jackd -d dummy -r 48000 -p 64\n\nFor even lower latency, try:\n  jackd -d dummy -r 48000 -p 32  # 0.67ms latency", status, jack.StrError(status))
	}

	// Create MIDI output port
	midiPort := client.PortRegister(portName, jack.DEFAULT_MIDI_TYPE, jack.PortIsOutput, 0)
	if midiPort == nil {
		client.Close()
		return nil, errors.New("failed to create MIDI output port")
	}

	b := &Bridge{
		oscServer:   server,
		jackClient:  client,
		midiOutPort: midiPort,
		eventQueue:  make(chan *MidiEvent, 1024), // Pre-allocated queue
	}

	// Set up process callback
	if code := client.SetProcessCallback(b.process); code != 0 {
		client.Close()
		return nil, fmt.Errorf("failed to set process callback: %s", jack.StrError(code))
	}

	// Set up OSC handlers
	b.setupOSCHandlers()

	return b, nil
}

func (b *Bridge) Start() error {
	if b.oscServer == nil || b.jackClient == nil {
		return errors.New("bridge not initialized")
	}

	// Activate JACK client
	if code := b.jackClient.Activate(); code != 0 {
		return fmt.Errorf("failed to activate JACK client: %s", jack.StrError(code))
	}

	// Start OSC server
	debugBridge("Starting OSC server on %s", b.oscServer.Addr)
	return b.oscServer.ListenAndServe()
}

func (b *Bridge) Cleanup() {
	debugBridge("Cleaning up bridge resources")

	if b.oscServer != nil {
		// OSC server cleanup happens when ListenAndServe returns
	}

	if b.jackClient != nil {
		b.jackClient.Close()
	}

	// Don't close the channel if it's already nil or closed
	// The channel will be garbage collected when the Bridge is freed
}

// JACK process callback - called by JACK in real-time thread
func (b *Bridge) process(nframes uint32) int {
	buffer := b.midiOutPort.MidiClearBuffer(nframes)

	processed := 0
	for processed < 32 { // Process max 32 events per cycle
		select {
		case event := <-b.eventQueue:
			event.midiData.Time = 0 // Immediate dispatch
			if err := b.midiOutPort.MidiEventWrite(event.midiData, buffer); err != 0 {
				debugBridge("Failed to write MIDI event: %v", err)
			}
			processed++
		default:
			return 0
		}
	}

	if processed == 32 {
		debugBridge("MIDI queue overflow, processed 32 events")
	}

	return 0
}

// List available JACK MIDI ports
func ListJackPorts() error {
	client, status := jack.ClientOpen("osc-midi-bridge-list", jack.NoStartServer)
	if status != 0 {
		return fmt.Errorf("cannot connect to JACK server: %s", jack.StrError(status))
	}
	defer client.Close()

	// Get all MIDI ports
	ports := client.GetPorts("", jack.DEFAULT_MIDI_TYPE, 0)

	fmt.Println("Available JACK MIDI ports:")
	fmt.Println("==========================")

	if len(ports) == 0 {
		fmt.Println("No MIDI ports found")
	} else {
		for i, port := range ports {
			// For now, just print the port name
			// TODO: Add direction detection when available in go-jack
			fmt.Printf("%d. %s\n", i+1, port)
		}
	}

	return nil
}
