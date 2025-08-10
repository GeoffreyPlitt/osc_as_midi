#!/bin/bash

set -e

echo "Starting OSC-to-MIDI bridge integration tests..."

# Start the bridge in background
./osc-midi-bridge --debug &
BRIDGE_PID=$!

# Wait for bridge to start
sleep 2

# Check if bridge is running
if ! kill -0 $BRIDGE_PID 2>/dev/null; then
    echo "ERROR: Bridge failed to start"
    exit 1
fi

echo "Testing OSC messages..."

# Test note on
echo "Sending note_on to channel 0..."
oscsend localhost 9000 /midi/0/note_on ii 60 127
sleep 0.5

# Test note off
echo "Sending note_off to channel 0..."
oscsend localhost 9000 /midi/0/note_off ii 60 0
sleep 0.5

# Test different channel
echo "Sending note_on to channel 5..."
oscsend localhost 9000 /midi/5/note_on ii 72 100
sleep 0.5

echo "Sending note_off to channel 5..."
oscsend localhost 9000 /midi/5/note_off ii 72 0
sleep 0.5

# Test signal handling
echo "Testing graceful shutdown..."
kill -TERM $BRIDGE_PID
sleep 1

# Check if process terminated cleanly
if kill -0 $BRIDGE_PID 2>/dev/null; then
    echo "WARNING: Bridge did not terminate cleanly, forcing kill"
    kill -9 $BRIDGE_PID
    exit 1
fi

echo "Integration tests completed successfully!"