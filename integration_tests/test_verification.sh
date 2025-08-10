#!/bin/bash
set -e

output_file=$(mktemp)
timeout 20 honcho start 2>&1 | tee "$output_file" || {
    exit_code=$?
    [[ $exit_code -ne 124 && $exit_code -ne 0 ]] && exit $exit_code
}

# Count expected events
bridge_events=$(grep -c "NOTE-O[NF]" "$output_file" || echo "0")
midi_events=$(grep -E "(90.*3c.*7f|80.*3c.*00|95.*48.*64|85.*48.*00)" "$output_file" | wc -l)
test_complete=$(grep -c "TEST_COMPLETE" "$output_file" || echo "0")

total_found=$((bridge_events + midi_events + test_complete))

if [[ $total_found -eq 9 ]]; then
    echo "✓ Integration tests passed ($total_found/9 events)"
    rm -f "$output_file"
    exit 0
else
    echo "✗ Integration tests failed ($total_found/9 events)"
    rm -f "$output_file"
    exit 1
fi