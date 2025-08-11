#!/bin/bash
set -e

echo "Waiting for MIDI setup..."
sleep 3

echo "Sending MIDI to bridge input..."
set -x
# Use jack_midiseq to send MIDI notes to bridge input
# Format: jack_midiseq name total_duration [startindex note duration] ...
# Sample rate is 48000, so 12000 samples = 0.25 seconds
jack_midiseq midi_sender 48000 \
    0 67 12000 \
    24000 79 12000 &
midiseq_pid=$!

# Connect the sequencer output to our bridge input  
sleep 1
jack_connect midi_sender:out osc-midi-bridge:midi_in

# Wait for sequence to complete
sleep 2
kill $midiseq_pid 2>/dev/null || true

echo "MIDI_TEST_COMPLETE"