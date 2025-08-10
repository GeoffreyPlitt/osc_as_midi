package main

import (
	"strconv"
	"strings"

	"github.com/hypebeast/go-osc/osc"
)

func (b *Bridge) handleNoteOn(msg *osc.Message) {
	// TODO: Implementation
}

func (b *Bridge) handleNoteOff(msg *osc.Message) {
	// TODO: Implementation
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
	// TODO: Implementation
}
