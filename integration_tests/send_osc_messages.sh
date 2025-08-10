#!/bin/bash
set -e

echo "Waiting for all services to be ready..."
sleep 8

echo "Connecting JACK MIDI ports..."
if jack_lsp -c | grep -q "osc-midi-bridge:midi_out:input"; then
    jack_connect osc-midi-bridge:midi_out osc-midi-bridge:midi_out:input || true
    echo "JACK ports connected"
fi

echo "Sending test messages..."

set -x
oscsend localhost 9000 /midi/0/note_on ii 60 127
sleep 0.5
oscsend localhost 9000 /midi/0/note_off ii 60 0  
sleep 0.5
oscsend localhost 9000 /midi/5/note_on ii 72 100
sleep 0.5
oscsend localhost 9000 /midi/5/note_off ii 72 0
sleep 2

echo "TEST_COMPLETE: OSC message sequence finished"

# Exit cleanly so honcho can shut down all processes
exit 0