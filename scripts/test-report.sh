#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
GRAY='\033[0;37m'
NC='\033[0m' # No Color

echo ""
echo -e "${CYAN}==========================================${NC}"
echo -e "${CYAN}  Running All Tests with Coverage...${NC}"
echo -e "${CYAN}==========================================${NC}"
echo ""

# Create temp directory
TEMP_DIR="/tmp/test-results-$$"
mkdir -p "$TEMP_DIR"

# Function to run tests for a category
run_test_category() {
    local category=$1
    local path=$2
    local color=$3
    
    echo ""
    echo -e "${color}Running $category Tests...${NC}"
    echo -e "${GRAY}-------------------------------------------${NC}"
    
    local json_file="$TEMP_DIR/$category.json"
    local start_time=$(date +%s)
    
    # Check if path exists
    local test_path="${path//\/\.\.\./}"
    if [ ! -d "$test_path" ]; then
        echo -e "${GRAY}  No tests found (path doesn't exist)${NC}"
        echo "0|0|0|0|0|N/A|0"
        return
    fi
    
    # Run tests
    gotestsum --format dots-v2 --jsonfile "$json_file" -- -cover "$path" > /dev/null 2>&1
    local exit_code=$?
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    # Parse results
    if [ -f "$json_file" ]; then
        local passed=$(grep -c '"Action":"pass".*"Test":' "$json_file" 2>/dev/null || echo 0)
        local failed=$(grep -c '"Action":"fail".*"Test":' "$json_file" 2>/dev/null || echo 0)
        local skipped=$(grep -c '"Action":"skip".*"Test":' "$json_file" 2>/dev/null || echo 0)
        local total=$((passed + failed + skipped))
        
        # Extract coverage
        local coverage=$(grep -o 'coverage: [0-9.]*%' "$json_file" | tail -1 | grep -o '[0-9.]*' || echo "N/A")
        
        echo "$passed|$failed|$skipped|$total|$duration|$coverage|$exit_code"
    else
        echo "0|0|0|0|$duration|N/A|$exit_code"
    fi
}

# Run different test categories
unit_result=$(run_test_category "Unit" "./tests/unit/..." "$GREEN")
integration_result=$(run_test_category "Integration" "./tests/integration/..." "$YELLOW")
api_result=$(run_test_category "API" "./tests/api/..." "\033[0;35m")

# Parse results
IFS='|' read -r unit_passed unit_failed unit_skipped unit_total unit_duration unit_coverage unit_exit <<< "$unit_result"
IFS='|' read -r int_passed int_failed int_skipped int_total int_duration int_coverage int_exit <<< "$integration_result"
IFS='|' read -r api_passed api_failed api_skipped api_total api_duration api_coverage api_exit <<< "$api_result"

# Calculate totals
total_passed=$((unit_passed + int_passed + api_passed))
total_failed=$((unit_failed + int_failed + api_failed))
total_skipped=$((unit_skipped + int_skipped + api_skipped))
total_tests=$((total_passed + total_failed + total_skipped))
total_duration=$((unit_duration + int_duration + api_duration))

# Print detailed summary
echo ""
echo -e "${CYAN}==========================================${NC}"
echo -e "${CYAN}           Test Summary Report${NC}"
echo -e "${CYAN}==========================================${NC}"
echo ""

echo -e "${NC}By Category:${NC}"
echo ""

# Print category results
print_category() {
    local category=$1
    local passed=$2
    local failed=$3
    local skipped=$4
    local total=$5
    local duration=$6
    local coverage=$7
    
    if [ "$total" -gt 0 ]; then
        local status_color=$GREEN
        local status="PASSED"
        if [ "$failed" -gt 0 ]; then
            status_color=$RED
            status="FAILED"
        fi
        
        echo -e "  ${NC}$category Tests:${NC}"
        echo -e "    Status:   ${status_color}$status${NC}"
        echo -e "    Total:    ${GRAY}$total tests${NC}"
        echo -e "    Passed:   ${GREEN}$passed${NC}"
        echo -e "    Failed:   ${RED}$failed${NC}"
        echo -e "    Skipped:  ${YELLOW}$skipped${NC}"
        echo -e "    Coverage: ${CYAN}$coverage%${NC}"
        echo -e "    Duration: ${GRAY}${duration}s${NC}"
        echo ""
    fi
}

print_category "Unit" "$unit_passed" "$unit_failed" "$unit_skipped" "$unit_total" "$unit_duration" "$unit_coverage"
print_category "Integration" "$int_passed" "$int_failed" "$int_skipped" "$int_total" "$int_duration" "$int_coverage"
print_category "API" "$api_passed" "$api_failed" "$api_skipped" "$api_total" "$api_duration" "$api_coverage"

# Overall summary
echo -e "${GRAY}-------------------------------------------${NC}"
echo -e "${NC}Overall Results:${NC}"
echo ""
echo -e "  Total Tests:  ${NC}$total_tests${NC}"
echo -e "  Passed:       ${GREEN}$total_passed${NC}"
echo -e "  Failed:       ${RED}$total_failed${NC}"
echo -e "  Skipped:      ${YELLOW}$total_skipped${NC}"
echo -e "  Duration:     ${GRAY}${total_duration}s${NC}"
echo ""

# Pass rate
if [ "$total_tests" -gt 0 ]; then
    pass_rate=$((total_passed * 100 / total_tests))
    pass_rate_color=$GREEN
    if [ "$pass_rate" -lt 90 ]; then
        pass_rate_color=$YELLOW
    fi
    if [ "$pass_rate" -lt 70 ]; then
        pass_rate_color=$RED
    fi
    echo -e "  Pass Rate:    ${pass_rate_color}${pass_rate}%${NC}"
fi

echo ""
echo -e "${CYAN}==========================================${NC}"

# Cleanup
rm -rf "$TEMP_DIR"

# Exit with error if any tests failed
if [ "$total_failed" -gt 0 ]; then
    echo ""
    echo -e "${RED}Tests FAILED!${NC}"
    exit 1
else
    echo ""
    echo -e "${GREEN}All tests PASSED!${NC}"
    exit 0
fi
