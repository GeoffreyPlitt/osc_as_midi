#!/bin/bash
set -e

output_file=$(mktemp)
timeout 35 honcho start 2>&1 | tee "$output_file" || {
    exit_code=$?
    [[ $exit_code -ne 124 && $exit_code -ne 0 ]] && exit $exit_code
}

# Count events - avoid bash trace lines that start with '+'
bridge_events=$(grep -c "NOTE-O[NF]" "$output_file" || echo "0")
midi_events=$(grep -E "(90.*3c.*7f|80.*3c.*00|95.*48.*64|85.*48.*00)" "$output_file" | wc -l)
osc_rx_events=$(grep -c "OSC_RX:.*\/midi\/[0-9]*\/note_o[nf]" "$output_file" || echo "0")
test_complete=$(grep -v "^[^|]*| +" "$output_file" | grep -c -E "(TEST_COMPLETE|MIDI_TEST_COMPLETE)" || echo "0")

# Strip any whitespace and newlines
bridge_events=$(echo $bridge_events | tr -d '\n\r ')
midi_events=$(echo $midi_events | tr -d '\n\r ')
osc_rx_events=$(echo $osc_rx_events | tr -d '\n\r ')
test_complete=$(echo $test_complete | tr -d '\n\r ')

total_found=$((bridge_events + midi_events + osc_rx_events + test_complete))


if [[ $total_found -ge 14 ]]; then
    echo "=================================================="
    echo "✓ Integration tests PASSED ($total_found events, expected ≥14)"
    echo "  OSC→MIDI: $bridge_events bridge + $midi_events MIDI = $((bridge_events + midi_events)) events"
    echo "  MIDI→OSC: $osc_rx_events OSC messages received" 
    echo "  Completions: $test_complete"
    echo ""
    rm -f "$output_file"
    exit 0
else
    printf "=================================================="
    echo "✗ Integration tests FAILED ($total_found events, expected ≥14)"
    echo ""
    rm -f "$output_file"
    exit 1
fi