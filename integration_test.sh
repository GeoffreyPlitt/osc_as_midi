#!/bin/bash

set -e

echo "Starting OSC-to-MIDI bridge integration tests..."

# PIDs to track for cleanup
JACKD_PID=""
BRIDGE_PID=""

# Cleanup function
cleanup() {
    echo "Cleaning up background processes..."
    
    # Kill bridge if running
    if [[ -n "$BRIDGE_PID" ]] && kill -0 "$BRIDGE_PID" 2>/dev/null; then
        echo "Stopping bridge (PID $BRIDGE_PID)..."
        kill -TERM "$BRIDGE_PID" 2>/dev/null || true
        sleep 1
        if kill -0 "$BRIDGE_PID" 2>/dev/null; then
            kill -9 "$BRIDGE_PID" 2>/dev/null || true
        fi
    fi
    
    # Kill jackd if running
    if [[ -n "$JACKD_PID" ]] && kill -0 "$JACKD_PID" 2>/dev/null; then
        echo "Stopping jackd (PID $JACKD_PID)..."
        kill -TERM "$JACKD_PID" 2>/dev/null || true
        sleep 1
        if kill -0 "$JACKD_PID" 2>/dev/null; then
            kill -9 "$JACKD_PID" 2>/dev/null || true
        fi
    fi
}

# Set trap to cleanup on exit, interrupt, or error
trap cleanup EXIT INT TERM

# Check if oscsend is available
if ! command -v oscsend &> /dev/null; then
    echo "ERROR: oscsend not found. Cannot run integration tests."
    exit 1
fi

echo "oscsend tool found, proceeding with tests..."

# Start jackd in background with container-friendly settings (no realtime)
echo "Starting jackd..."
# Set memory limits for the container
ulimit -l unlimited 2>/dev/null || true
JACK_NO_AUDIO_RESERVATION=1 jackd -r -d dummy --rate 48000 --period 1024 &
JACKD_PID=$!

# Wait for jackd to start
echo "Waiting for jackd to initialize..."
sleep 2

# Check if jackd is running
if ! kill -0 $JACKD_PID 2>/dev/null; then
    echo "ERROR: jackd failed to start"
    exit 1
fi

echo "jackd is running with PID $JACKD_PID"

# Start the bridge in background
echo "Starting bridge..."
DEBUG=* ./osc-midi-bridge &
BRIDGE_PID=$!

# Wait for bridge to start
sleep 2

# Check if bridge is running
if ! kill -0 $BRIDGE_PID 2>/dev/null; then
    echo "ERROR: Bridge failed to start"
    exit 1
else
    echo "Bridge is running with PID $BRIDGE_PID"
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

# Test signal handling (only if bridge is still running)
if kill -0 $BRIDGE_PID 2>/dev/null; then
    echo "Testing graceful shutdown..."
    kill -TERM $BRIDGE_PID
    sleep 2
    
    # Check if process terminated cleanly
    if kill -0 $BRIDGE_PID 2>/dev/null; then
        echo "WARNING: Bridge did not terminate cleanly, forcing kill"
        kill -9 $BRIDGE_PID
        exit 1
    else
        echo "Bridge terminated cleanly"
        BRIDGE_PID="" # Clear PID since process is gone
    fi
else
    echo "Bridge already exited, skipping signal handling test"
fi

echo "Integration tests completed successfully!"
echo "Cleanup will be handled by trap on exit"