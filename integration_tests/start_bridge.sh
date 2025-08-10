#!/bin/bash
set -e

echo "Waiting for jackd to be ready..."
sleep 3

echo "Starting OSC-MIDI bridge..."
DEBUG=* ../osc-midi-bridge