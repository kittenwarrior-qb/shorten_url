#!/bin/bash

echo "Running Unit Tests..."
echo "========================"
echo ""

# Run tests and capture JSON output
gotestsum --format dots-v2 --jsonfile /tmp/test-results.json -- -cover ./tests/... 2>/dev/null

# Parse results
PASSED=$(grep -c '"Action":"pass"' /tmp/test-results.json 2>/dev/null | head -1)
FAILED=$(grep -c '"Action":"fail"' /tmp/test-results.json 2>/dev/null | head -1)
SKIPPED=$(grep -c '"Action":"skip"' /tmp/test-results.json 2>/dev/null | head -1)

# Handle empty values
PASSED=${PASSED:-0}
FAILED=${FAILED:-0}
SKIPPED=${SKIPPED:-0}

TOTAL=$((PASSED + FAILED + SKIPPED))

echo ""
echo "========================"
echo "📊 Test Summary"
echo "========================"
echo "Total:   $TOTAL tests"
echo "Pass:   $PASSED"
echo "Fail:   $FAILED"
echo "Skip:   $SKIPPED"
echo "========================"

# Exit with error if any test failed
if [ "$FAILED" -gt 0 ]; then
    exit 1
fi
