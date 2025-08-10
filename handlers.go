package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/GeoffreyPlitt/debuggo"
	"github.com/hypebeast/go-osc/osc"
	"github.com/xthexder/go-jack"
)

var debugHandlers = debuggo.Debug("osc-midi-bridge:handlers")

func (b *Bridge) handleNoteOn(msg *osc.Message) error {
	if len(msg.Arguments) < 2 {
		return errors.New("note_on requires at least 2 arguments: note and velocity")
	}

	channel := b.extractChannel(msg.Address)
	note := toUint8(msg.Arguments[0])
	velocity := toUint8(msg.Arguments[1])

	// Create MIDI note on message: 0x90 | channel, note, velocity
	midiBytes := []byte{
		byte(0x90 | (channel & 0x0F)),
		byte(note & 0x7F),
		byte(velocity & 0x7F),
	}

	select {
	case b.eventQueue <- &MidiEvent{
		midiData: &jack.MidiData{
			Time:   0,
			Buffer: midiBytes,
		},
	}:
		debugHandlers("note-on ch:%d note:%d vel:%d", channel, note, velocity)
	default:
		return errors.New("MIDI queue full")
	}

	return nil
}

func (b *Bridge) handleNoteOff(msg *osc.Message) error {
	if len(msg.Arguments) < 2 {
		return errors.New("note_off requires at least 2 arguments: note and velocity")
	}

	channel := b.extractChannel(msg.Address)
	note := toUint8(msg.Arguments[0])
	velocity := toUint8(msg.Arguments[1])

	// Create MIDI note off message: 0x80 | channel, note, velocity
	midiBytes := []byte{
		byte(0x80 | (channel & 0x0F)),
		byte(note & 0x7F),
		byte(velocity & 0x7F),
	}

	select {
	case b.eventQueue <- &MidiEvent{
		midiData: &jack.MidiData{
			Time:   0,
			Buffer: midiBytes,
		},
	}:
		debugHandlers("note-off ch:%d note:%d vel:%d", channel, note, velocity)
	default:
		return errors.New("MIDI queue full")
	}

	return nil
}

func (b *Bridge) extractChannel(address string) uint8 {
	parts := strings.Split(address, "/")
	if len(parts) >= 3 && parts[1] == "midi" {
		ch, err := strconv.Atoi(parts[2])
		if err == nil && ch >= 0 && ch <= 15 {
			return uint8(ch)
		}
	}
	return 0
}

func (b *Bridge) setupOSCHandlers() {
	// Check if oscServer exists
	if b.oscServer == nil || b.oscServer.Dispatcher == nil {
		debugHandlers("OSC server not initialized")
		return
	}

	// Get the dispatcher from the server
	dispatcher, ok := b.oscServer.Dispatcher.(*osc.StandardDispatcher)
	if !ok {
		debugHandlers("Failed to get standard dispatcher")
		return
	}

	// Handle note on messages: /midi/{channel}/note_on
	// Using wildcard pattern for channels 0-15
	for i := 0; i < 16; i++ {
		path := fmt.Sprintf("/midi/%d/note_on", i)
		dispatcher.AddMsgHandler(path, func(msg *osc.Message) {
			if err := b.handleNoteOn(msg); err != nil {
				debugHandlers("Error handling note_on: %v", err)
			}
		})
	}

	// Handle note off messages: /midi/{channel}/note_off
	for i := 0; i < 16; i++ {
		path := fmt.Sprintf("/midi/%d/note_off", i)
		dispatcher.AddMsgHandler(path, func(msg *osc.Message) {
			if err := b.handleNoteOff(msg); err != nil {
				debugHandlers("Error handling note_off: %v", err)
			}
		})
	}

	debugHandlers("OSC handlers configured for /midi/{0-15}/note_on and /midi/{0-15}/note_off")
}
