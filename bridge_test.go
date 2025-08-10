package main

import (
	"testing"

	"github.com/hypebeast/go-osc/osc"
	"github.com/xthexder/go-jack"
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

	// Create a bridge with initialized channels
	bridge = &Bridge{
		eventQueue:  make(chan *MidiEvent),
		oscOutQueue: make(chan *osc.Message, 16),
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

func TestParseIncomingMIDI(t *testing.T) {
	bridge := &Bridge{}

	tests := []struct {
		name         string
		midiData     []byte
		expectedPath string
		expectedNote int32
		expectedVel  int32
		shouldBeNil  bool
	}{
		{
			name:         "Note On - Channel 0",
			midiData:     []byte{0x90, 0x3C, 0x7F}, // Note On, Middle C, full velocity
			expectedPath: "/midi/0/note_on",
			expectedNote: 60,
			expectedVel:  127,
			shouldBeNil:  false,
		},
		{
			name:         "Note Off - Channel 5",
			midiData:     []byte{0x85, 0x40, 0x40}, // Note Off, E4, velocity 64
			expectedPath: "/midi/5/note_off",
			expectedNote: 64,
			expectedVel:  64,
			shouldBeNil:  false,
		},
		{
			name:         "Note On with velocity 0 (becomes Note Off)",
			midiData:     []byte{0x90, 0x3C, 0x00}, // Note On, Middle C, velocity 0
			expectedPath: "/midi/0/note_off",
			expectedNote: 60,
			expectedVel:  0,
			shouldBeNil:  false,
		},
		{
			name:        "Invalid MIDI - too short",
			midiData:    []byte{0x90, 0x3C}, // Missing velocity byte
			shouldBeNil: true,
		},
		{
			name:        "Unsupported MIDI message - CC",
			midiData:    []byte{0xB0, 0x07, 0x7F}, // Control Change
			shouldBeNil: true,
		},
		{
			name:         "Note On - Channel 15",
			midiData:     []byte{0x9F, 0x24, 0x50}, // Note On, Channel 15 (0xF), C2, velocity 80
			expectedPath: "/midi/15/note_on",
			expectedNote: 36,
			expectedVel:  80,
			shouldBeNil:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &jack.MidiData{
				Buffer: tt.midiData,
			}

			result := bridge.parseIncomingMIDI(event)

			if tt.shouldBeNil {
				if result != nil {
					t.Errorf("Expected nil result for %s, got %v", tt.name, result)
				}
				return
			}

			if result == nil {
				t.Errorf("Expected non-nil result for %s", tt.name)
				return
			}

			if result.Address != tt.expectedPath {
				t.Errorf("Expected path %s, got %s", tt.expectedPath, result.Address)
			}

			if len(result.Arguments) != 2 {
				t.Errorf("Expected 2 arguments, got %d", len(result.Arguments))
				return
			}

			note, ok := result.Arguments[0].(int32)
			if !ok {
				t.Errorf("Expected note to be int32, got %T", result.Arguments[0])
				return
			}
			if note != tt.expectedNote {
				t.Errorf("Expected note %d, got %d", tt.expectedNote, note)
			}

			velocity, ok := result.Arguments[1].(int32)
			if !ok {
				t.Errorf("Expected velocity to be int32, got %T", result.Arguments[1])
				return
			}
			if velocity != tt.expectedVel {
				t.Errorf("Expected velocity %d, got %d", tt.expectedVel, velocity)
			}
		})
	}
}

func TestOSCQueueHandling(t *testing.T) {
	// Test that OSC queue can handle being closed
	bridge := &Bridge{
		oscOutQueue: make(chan *osc.Message, 16),
	}

	// Test sending message to queue
	msg := osc.NewMessage("/test", int32(1), int32(2))
	select {
	case bridge.oscOutQueue <- msg:
		// Message sent successfully
	default:
		t.Error("Failed to send message to OSC queue")
	}

	// Test cleanup closes the queue
	bridge.Cleanup()

	// After cleanup, the channel should be closed
	// We need to drain any existing message first
	select {
	case _, open := <-bridge.oscOutQueue:
		if open {
			// If we received a message, try again to check if channel is closed
			select {
			case _, open := <-bridge.oscOutQueue:
				if open {
					t.Error("Expected OSC queue to be closed after cleanup")
				}
			default:
				t.Error("Expected to receive close signal from channel")
			}
		}
	default:
		t.Error("Expected to receive from closed channel")
	}
}
