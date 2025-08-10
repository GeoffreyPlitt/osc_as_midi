#!/bin/bash

echo "Starting jackd server..."
ulimit -l unlimited 2>/dev/null || true

# Trap signals and kill jackd forcefully
cleanup() {
    echo "Force killing jackd..."
    if [[ -n "$jackd_pid" ]] && kill -0 "$jackd_pid" 2>/dev/null; then
        kill -9 "$jackd_pid" 2>/dev/null || true
    fi
    # Kill any remaining jackd processes
    pkill -9 jackd 2>/dev/null || true
    exit 0
}
trap cleanup SIGTERM SIGINT

JACK_NO_AUDIO_RESERVATION=1 jackd -r -d dummy --rate 48000 --period 1024 &
jackd_pid=$!

# Wait for jackd and handle signals
wait $jackd_pid 2>/dev/null || true