package main

import (
	"errors"

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
	if b.oscServer == nil && b.midiDriver == nil {
		return errors.New("bridge not initialized")
	}
	return nil
}

func (b *Bridge) Cleanup() {
	if b.oscServer != nil {
		// b.oscServer.Close() // Uncomment when method exists
	}
	if b.midiOut != nil {
		b.midiOut.Close()
	}
	if b.midiDriver != nil {
		b.midiDriver.Close()
	}
}
