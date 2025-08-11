#!/bin/bash
set -e

echo "Waiting for bridge to initialize..."
sleep 5

echo "Starting OSC dump on port 8000..."
oscdump -L 8000 | while read line; do
    echo "OSC_RX: $line"
done