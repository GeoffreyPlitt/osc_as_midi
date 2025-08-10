package main

import (
	"testing"

	"github.com/hypebeast/go-osc/osc"
)

func TestHandleNoteOn(t *testing.T) {
	bridge := &Bridge{
		eventQueue: make(chan *MidiEvent, 10), // Buffer for test
	}

	tests := []struct {
		name        string
		message     *osc.Message
		shouldError bool
	}{
		{
			name: "valid note on",
			message: &osc.Message{
				Address:   "/midi/0/note_on",
				Arguments: []interface{}{60, 127},
			},
			shouldError: false,
		},
		{
			name: "missing velocity",
			message: &osc.Message{
				Address:   "/midi/0/note_on",
				Arguments: []interface{}{60},
			},
			shouldError: true,
		},
		{
			name: "no arguments",
			message: &osc.Message{
				Address:   "/midi/0/note_on",
				Arguments: []interface{}{},
			},
			shouldError: true,
		},
		{
			name: "channel 15",
			message: &osc.Message{
				Address:   "/midi/15/note_on",
				Arguments: []interface{}{64, 100},
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear the channel before test
			for len(bridge.eventQueue) > 0 {
				<-bridge.eventQueue
			}

			err := bridge.handleNoteOn(tt.message)
			if (err != nil) != tt.shouldError {
				t.Errorf("handleNoteOn() error = %v, shouldError %v", err, tt.shouldError)
			}

			// If no error, check that event was queued
			if err == nil {
				select {
				case event := <-bridge.eventQueue:
					if event == nil || event.midiData == nil {
						t.Error("Expected valid MIDI event in queue")
					}
					// Verify MIDI message structure
					if len(event.midiData.Buffer) != 3 {
						t.Errorf("Expected 3-byte MIDI message, got %d bytes", len(event.midiData.Buffer))
					}
				default:
					t.Error("Expected event in queue but found none")
				}
			}
		})
	}
}

func TestHandleNoteOff(t *testing.T) {
	bridge := &Bridge{
		eventQueue: make(chan *MidiEvent, 10),
	}

	tests := []struct {
		name        string
		message     *osc.Message
		shouldError bool
	}{
		{
			name: "valid note off",
			message: &osc.Message{
				Address:   "/midi/0/note_off",
				Arguments: []interface{}{60, 0},
			},
			shouldError: false,
		},
		{
			name: "with velocity",
			message: &osc.Message{
				Address:   "/midi/5/note_off",
				Arguments: []interface{}{72, 64},
			},
			shouldError: false,
		},
		{
			name: "missing arguments",
			message: &osc.Message{
				Address:   "/midi/0/note_off",
				Arguments: []interface{}{60},
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear the channel before test
			for len(bridge.eventQueue) > 0 {
				<-bridge.eventQueue
			}

			err := bridge.handleNoteOff(tt.message)
			if (err != nil) != tt.shouldError {
				t.Errorf("handleNoteOff() error = %v, shouldError %v", err, tt.shouldError)
			}

			// If no error, check that event was queued
			if err == nil {
				select {
				case event := <-bridge.eventQueue:
					if event == nil || event.midiData == nil {
						t.Error("Expected valid MIDI event in queue")
					}
					// Verify it's a note off message (0x80)
					if (event.midiData.Buffer[0] & 0xF0) != 0x80 {
						t.Error("Expected note off status byte")
					}
				default:
					t.Error("Expected event in queue but found none")
				}
			}
		})
	}
}

func TestSetupOSCHandlers(t *testing.T) {
	// Create a bridge without OSC server
	bridge := &Bridge{}

	// Should not panic even without oscServer initialized
	bridge.setupOSCHandlers()

	// Create a bridge with OSC server but wrong dispatcher type
	bridge = &Bridge{
		oscServer: &osc.Server{},
	}

	// Should handle gracefully
	bridge.setupOSCHandlers()
}

func TestQueueOverflow(t *testing.T) {
	// Create bridge with small queue
	bridge := &Bridge{
		eventQueue: make(chan *MidiEvent, 1),
	}

	msg := &osc.Message{
		Address:   "/midi/0/note_on",
		Arguments: []interface{}{60, 127},
	}

	// First message should succeed
	err := bridge.handleNoteOn(msg)
	if err != nil {
		t.Errorf("First message failed: %v", err)
	}

	// Second message should fail (queue full)
	err = bridge.handleNoteOn(msg)
	if err == nil {
		t.Error("Expected error for full queue")
	}
}

func TestMidiMessageFormat(t *testing.T) {
	bridge := &Bridge{
		eventQueue: make(chan *MidiEvent, 10),
	}

	// Test note on
	noteOnMsg := &osc.Message{
		Address:   "/midi/3/note_on",
		Arguments: []interface{}{int32(64), int32(100)},
	}

	bridge.handleNoteOn(noteOnMsg)

	event := <-bridge.eventQueue
	if event.midiData.Buffer[0] != 0x93 { // 0x90 | 0x03
		t.Errorf("Expected note on status 0x93, got 0x%02X", event.midiData.Buffer[0])
	}
	if event.midiData.Buffer[1] != 64 {
		t.Errorf("Expected note 64, got %d", event.midiData.Buffer[1])
	}
	if event.midiData.Buffer[2] != 100 {
		t.Errorf("Expected velocity 100, got %d", event.midiData.Buffer[2])
	}

	// Test note off
	noteOffMsg := &osc.Message{
		Address:   "/midi/7/note_off",
		Arguments: []interface{}{uint(72), float32(50.5)},
	}

	bridge.handleNoteOff(noteOffMsg)

	event = <-bridge.eventQueue
	if event.midiData.Buffer[0] != 0x87 { // 0x80 | 0x07
		t.Errorf("Expected note off status 0x87, got 0x%02X", event.midiData.Buffer[0])
	}
	if event.midiData.Buffer[1] != 72 {
		t.Errorf("Expected note 72, got %d", event.midiData.Buffer[1])
	}
	if event.midiData.Buffer[2] != 50 {
		t.Errorf("Expected velocity 50, got %d", event.midiData.Buffer[2])
	}
}
