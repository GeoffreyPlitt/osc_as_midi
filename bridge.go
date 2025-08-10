package main

import (
	_ "github.com/GeoffreyPlitt/debuggo"
	"github.com/hypebeast/go-osc/osc"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type Bridge struct {
	oscServer   *osc.Server
	midiDriver  drivers.Driver
	midiOut     drivers.Out
	debug       func(format string, args ...interface{})
	monitorMode bool
}

func (b *Bridge) Start() error {
	// TODO: Implementation
	return nil
}

func (b *Bridge) Cleanup() {
	// TODO: Implementation
}
