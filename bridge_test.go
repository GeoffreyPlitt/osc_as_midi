package main

import (
	"testing"
)

func TestBridgeStart(t *testing.T) {
	bridge := &Bridge{
		monitorMode: true, // Don't actually send MIDI in tests
	}

	// Test that bridge can start without error
	err := bridge.Start()
	if err == nil {
		t.Error("Expected error when starting bridge without proper initialization")
	}
}

func TestBridgeCleanup(t *testing.T) {
	bridge := &Bridge{}

	// Should not panic even with nil fields
	bridge.Cleanup()
}

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
