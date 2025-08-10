package main

import (
	"github.com/GeoffreyPlitt/debuggo"
	"github.com/hypebeast/go-osc/osc"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type Bridge struct {
	oscServer   *osc.Server
	midiDriver  midi.Driver
	midiOut     midi.Out
	debug       *debuggo.Debug
	monitorMode bool
}

func (b *Bridge) Start() error {
	// TODO: Implementation
	return nil
}

func (b *Bridge) Cleanup() {
	// TODO: Implementation
}