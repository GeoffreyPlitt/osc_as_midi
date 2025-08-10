#!/bin/bash
set -e

echo "Waiting for bridge to create MIDI port..."
sleep 5

echo "Starting jack_midi_dump..."
jack_midi_dump osc-midi-bridge:midi_out