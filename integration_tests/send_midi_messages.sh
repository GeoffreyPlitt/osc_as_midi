#!/bin/bash
set -e

echo "Waiting for MIDI setup..."
sleep 3

echo "Sending MIDI to bridge input..."
set -x

# Build our custom MIDI sender
go build -o midi_sender midi_sender.go

# Run our deterministic MIDI sender in background
./midi_sender &
midi_pid=$!

# Wait for it to activate and create ports, then connect
sleep 0.1
jack_connect midi_sender:out osc-midi-bridge:midi_in

# Wait for sequence to complete
wait $midi_pid

echo "MIDI_TEST_COMPLETE"