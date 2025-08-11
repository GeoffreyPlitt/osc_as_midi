#!/bin/bash
set -e

echo "Waiting for jackd to be ready..."
sleep 3

echo "Starting bidirectional OSC-MIDI bridge..."
DEBUG=* ../osc-midi-bridge --osc-target-port 8000