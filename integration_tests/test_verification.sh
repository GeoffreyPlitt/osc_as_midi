#!/bin/bash
set -e

printf '\n\nStarting integration tests ...'

# Run honcho and capture output with tee so we can see it AND process it
output_file=$(mktemp)
timeout 15 honcho start 2>&1 | tee "$output_file" || {
    exit_code=$?
    if [[ $exit_code -eq 124 ]]; then
        echo "Tests completed with timeout (expected - jackd cleanup issue)"
    elif [[ $exit_code -eq 0 ]]; then
        echo "Tests completed successfully"
    else
        echo "Tests failed with exit code $exit_code"
        exit $exit_code
    fi
}

echo ""
echo "Honcho execution completed, analyzing output..."

# Check for bridge output - NOTE-ON/NOTE-OFF messages
bridge_events=$(grep -c "NOTE-O[NF]" "$output_file" || echo "0")
echo "Bridge events detected: $bridge_events/4"

# Check for MIDI dump patterns - looking for hex patterns
midi_note_on_ch0=$(grep -c "90.*3c.*7f" "$output_file" || echo "0")
midi_note_off_ch0=$(grep -c "80.*3c.*00" "$output_file" || echo "0") 
midi_note_on_ch5=$(grep -c "95.*48.*64" "$output_file" || echo "0")
midi_note_off_ch5=$(grep -c "85.*48.*00" "$output_file" || echo "0")

total_midi_events=$((midi_note_on_ch0 + midi_note_off_ch0 + midi_note_on_ch5 + midi_note_off_ch5))
echo "MIDI events detected: $total_midi_events/4"

# Check for test sequence completion
test_complete=$(grep -c "TEST_COMPLETE" "$output_file" || echo "0")
echo "Test sequence completed: $test_complete/1"

echo ""
echo "=== Test Results ==="
echo "✓ Bridge NOTE-ON/OFF events: $bridge_events/4"
echo "✓ MIDI dump events: $total_midi_events/4" 
echo "✓ Test sequence completion: $test_complete/1"

# Determine overall success
total_expected=9  # 4 bridge + 4 midi + 1 completion
total_found=$((bridge_events + total_midi_events + test_complete))

if [[ $total_found -eq $total_expected ]]; then
    echo ""
    echo "🎉 ALL INTEGRATION TESTS PASSED! ($total_found/$total_expected events detected)"
    rm -f "$output_file"
    exit 0
else
    echo ""
    echo "❌ INTEGRATION TESTS FAILED: Only $total_found/$total_expected events detected"
    echo ""
    echo "=== Debug Information ==="
    echo "MIDI Note On Ch0 (90 3c 7f): $midi_note_on_ch0"
    echo "MIDI Note Off Ch0 (80 3c 00): $midi_note_off_ch0"
    echo "MIDI Note On Ch5 (95 48 64): $midi_note_on_ch5" 
    echo "MIDI Note Off Ch5 (85 48 00): $midi_note_off_ch5"
    echo ""
    rm -f "$output_file"
    exit 1
fi