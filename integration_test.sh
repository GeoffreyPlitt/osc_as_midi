#!/bin/bash

set -e

echo "Starting OSC-to-MIDI bridge integration tests..."

# PIDs to track for cleanup
JACKD_PID=""
BRIDGE_PID=""
MIDI_DUMP_PID=""

# Cleanup function
cleanup() {
    echo "Cleaning up background processes..."
    
    # Kill MIDI dump if running
    if [[ -n "$MIDI_DUMP_PID" ]] && kill -0 "$MIDI_DUMP_PID" 2>/dev/null; then
        echo "Stopping jack_midi_dump (PID $MIDI_DUMP_PID)..."
        kill -TERM "$MIDI_DUMP_PID" 2>/dev/null || true
        sleep 1
        if kill -0 "$MIDI_DUMP_PID" 2>/dev/null; then
            kill -9 "$MIDI_DUMP_PID" 2>/dev/null || true
        fi
    fi
    
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

# Check available JACK ports first
echo "Checking available JACK ports..."
jack_lsp -A

# Start MIDI event monitoring
echo "Starting MIDI event monitoring..."
# Try alternative approach with jack_monitor first
if command -v jack_monitor &> /dev/null; then
    echo "Using jack_monitor for MIDI capture..."
    jack_monitor osc-midi-bridge:midi_out > midi_events.log 2>&1 &
    MIDI_DUMP_PID=$!
else
    echo "Using jack_midi_dump for MIDI capture..."
    jack_midi_dump osc-midi-bridge:midi_out > midi_events.log 2>&1 &
    MIDI_DUMP_PID=$!
fi

# Wait for MIDI dump to connect
sleep 3

# Check if MIDI dump is running
if ! kill -0 $MIDI_DUMP_PID 2>/dev/null; then
    echo "ERROR: MIDI monitoring failed to start"
    echo "MIDI dump log contents:"
    cat midi_events.log 2>/dev/null || echo "No log file found"
    exit 1
fi

echo "MIDI monitoring is running with PID $MIDI_DUMP_PID"

# Verify connection was established
echo "Checking JACK connections after monitoring start..."
jack_lsp -c

# Manually connect the bridge output to the midi_dump input if needed
echo "Attempting to connect JACK MIDI ports..."
if jack_lsp -c | grep -q "osc-midi-bridge:midi_out:input"; then
    echo "Connecting osc-midi-bridge:midi_out to osc-midi-bridge:midi_out:input"
    jack_connect osc-midi-bridge:midi_out osc-midi-bridge:midi_out:input || echo "Connection failed or already exists"
fi

echo "Final JACK connections:"
jack_lsp -c

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

# Give extra time for MIDI events to be processed and captured
echo "Waiting for MIDI events to be processed..."
sleep 2

# Stop MIDI monitoring and verify events
echo "Stopping MIDI event monitoring..."
if [[ -n "$MIDI_DUMP_PID" ]] && kill -0 "$MIDI_DUMP_PID" 2>/dev/null; then
    kill -TERM "$MIDI_DUMP_PID" 2>/dev/null || true
    sleep 1
    MIDI_DUMP_PID=""  # Clear PID since we stopped it
fi

# Verify MIDI events were captured
echo "Verifying captured MIDI events..."
if [[ ! -f midi_events.log ]]; then
    echo "ERROR: MIDI events log file not found"
    exit 1
fi

# Check for expected MIDI events
# Note On Ch 0, Note 60 (0x3C), Vel 127 (0x7F) = 90 3C 7F
# Note Off Ch 0, Note 60 (0x3C), Vel 0 (0x00) = 80 3C 00  
# Note On Ch 5, Note 72 (0x48), Vel 100 (0x64) = 95 48 64
# Note Off Ch 5, Note 72 (0x48), Vel 0 (0x00) = 85 48 00

events_found=0

if grep -qi "90.*3c.*7f" midi_events.log; then
    echo "✓ MIDI Note On Ch0 Note60 Vel127 detected"
    ((events_found++))
else
    echo "✗ MIDI Note On Ch0 Note60 Vel127 NOT detected"
fi

if grep -qi "80.*3c.*00" midi_events.log; then
    echo "✓ MIDI Note Off Ch0 Note60 Vel0 detected"  
    ((events_found++))
else
    echo "✗ MIDI Note Off Ch0 Note60 Vel0 NOT detected"
fi

if grep -qi "95.*48.*64" midi_events.log; then
    echo "✓ MIDI Note On Ch5 Note72 Vel100 detected"
    ((events_found++))
else
    echo "✗ MIDI Note On Ch5 Note72 Vel100 NOT detected"
fi

if grep -qi "85.*48.*00" midi_events.log; then
    echo "✓ MIDI Note Off Ch5 Note72 Vel0 detected"
    ((events_found++))
else
    echo "✗ MIDI Note Off Ch5 Note72 Vel0 NOT detected"
fi

echo "MIDI verification: $events_found/4 expected events detected"

if [[ $events_found -lt 4 ]]; then
    echo "ERROR: Not all expected MIDI events were detected"
    echo "MIDI events log contents:"
    cat midi_events.log
    echo "Raw hex dump of log file:"
    hexdump -C midi_events.log | head -20
    exit 1
fi

echo "✓ All MIDI events verified successfully"

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