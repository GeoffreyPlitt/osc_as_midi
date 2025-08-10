package main

import (
	"testing"
)

func TestExtractChannel(t *testing.T) {
	tests := []struct {
		address  string
		expected uint8
	}{
		{"/midi/0/note_on", 0},
		{"/midi/5/note_off", 5},
		{"/midi/15/cc", 15},
		{"/midi/16/note_on", 0}, // Should wrap or handle invalid channel
		{"/invalid/path", 0},    // Should handle gracefully
	}

	bridge := &Bridge{}
	for _, tt := range tests {
		result := bridge.extractChannel(tt.address)
		if result != tt.expected {
			t.Errorf("extractChannel(%s) = %d, expected %d", tt.address, result, tt.expected)
		}
	}
}

func TestBridgeCleanup(t *testing.T) {
	// Create a bridge with nil fields
	bridge := &Bridge{}

	// Should not panic even with nil fields
	bridge.Cleanup()

	// Create a bridge with initialized channel
	bridge = &Bridge{
		eventQueue: make(chan *MidiEvent),
	}

	// Should handle cleanup gracefully
	bridge.Cleanup()
}

func TestMidiEvent(t *testing.T) {
	// Test MidiEvent structure
	event := &MidiEvent{
		midiData: nil,
	}

	// Ensure the structure exists and can be created
	if event.midiData != nil {
		t.Error("Expected nil midiData")
	}
}
