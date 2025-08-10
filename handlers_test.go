package main

import (
	"testing"

	"github.com/hypebeast/go-osc/osc"
)

func TestHandleNoteOn(t *testing.T) {
	bridge := &Bridge{
		monitorMode: true, // Don't send actual MIDI
	}

	tests := []struct {
		name      string
		message   *osc.Message
		shouldLog bool
	}{
		{
			name: "valid note on",
			message: &osc.Message{
				Address:   "/midi/0/note_on",
				Arguments: []interface{}{60, 127},
			},
			shouldLog: true,
		},
		{
			name: "missing velocity",
			message: &osc.Message{
				Address:   "/midi/0/note_on",
				Arguments: []interface{}{60},
			},
			shouldLog: false,
		},
		{
			name: "no arguments",
			message: &osc.Message{
				Address:   "/midi/0/note_on",
				Arguments: []interface{}{},
			},
			shouldLog: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			bridge.handleNoteOn(tt.message)
		})
	}
}

func TestHandleNoteOff(t *testing.T) {
	bridge := &Bridge{
		monitorMode: true,
	}

	tests := []struct {
		name    string
		message *osc.Message
	}{
		{
			name: "valid note off",
			message: &osc.Message{
				Address:   "/midi/0/note_off",
				Arguments: []interface{}{60, 0},
			},
		},
		{
			name: "with velocity",
			message: &osc.Message{
				Address:   "/midi/5/note_off",
				Arguments: []interface{}{72, 64},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			bridge.handleNoteOff(tt.message)
		})
	}
}

func TestSetupOSCHandlers(t *testing.T) {
	bridge := &Bridge{}

	// Should not panic even without oscServer initialized
	bridge.setupOSCHandlers()
}
