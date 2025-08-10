#!/bin/bash
set -e

output_file=$(mktemp)
timeout 30 honcho start 2>&1 | tee "$output_file" || {
    exit_code=$?
    [[ $exit_code -ne 124 && $exit_code -ne 0 ]] && exit $exit_code
}

# Count expected events - avoid bash trace lines that start with '+'
bridge_events=$(grep -c "NOTE-O[NF]" "$output_file" || echo "0")
midi_events=$(grep -E "(90.*3c.*7f|80.*3c.*00|95.*48.*64|85.*48.*00)" "$output_file" | wc -l)
test_complete=$(grep -v "^[^|]*| +" "$output_file" | grep -c "TEST_COMPLETE: OSC message sequence finished" || echo "0")

total_found=$((bridge_events + midi_events + test_complete))

echo ""
echo ""
echo ""

if [[ $total_found -eq 9 ]]; then
    echo "=================================================="
    echo "✓ Integration tests PASSED ($total_found/9 events)"
    echo ""
    rm -f "$output_file"
    exit 0
else
    printf "=================================================="
    echo "✗ Integration tests FAILED ($total_found/9 events)"
    echo ""
    rm -f "$output_file"
    exit 1
fi