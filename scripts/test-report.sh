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

# Run all tests
echo -e "${GREEN}Running Unit Tests...${NC}"
echo -e "${GRAY}-------------------------------------------${NC}"

# Run tests and capture output
TEST_OUTPUT=$(go test -v -cover ./tests/unit/... 2>&1)
EXIT_CODE=$?

# Count results - use awk to ensure we get a number
PASSED=$(echo "$TEST_OUTPUT" | grep -c "^--- PASS:" 2>/dev/null)
FAILED=$(echo "$TEST_OUTPUT" | grep -c "^--- FAIL:" 2>/dev/null)

# Ensure we have numbers
PASSED=${PASSED:-0}
FAILED=${FAILED:-0}

# Remove leading zeros to avoid octal interpretation
PASSED=$((10#$PASSED))
FAILED=$((10#$FAILED))

TOTAL=$((PASSED + FAILED))

# Extract coverage
COVERAGE=$(echo "$TEST_OUTPUT" | grep -o 'coverage: [0-9.]*%' | tail -1 | grep -o '[0-9.]*')
COVERAGE=${COVERAGE:-N/A}

# Print summary
echo ""
echo -e "${CYAN}==========================================${NC}"
echo -e "${CYAN}           Test Summary Report${NC}"
echo -e "${CYAN}==========================================${NC}"
echo ""

if [ $TOTAL -gt 0 ]; then
    STATUS_COLOR=$GREEN
    STATUS="PASSED"
    if [ $FAILED -gt 0 ]; then
        STATUS_COLOR=$RED
        STATUS="FAILED"
    fi
    
    echo -e "  ${NC}Unit Tests:${NC}"
    echo -e "    Status:   ${STATUS_COLOR}$STATUS${NC}"
    echo -e "    Total:    ${GRAY}$TOTAL tests${NC}"
    echo -e "    Passed:   ${GREEN}$PASSED${NC}"
    echo -e "    Failed:   ${RED}$FAILED${NC}"
    echo -e "    Coverage: ${CYAN}$COVERAGE%${NC}"
    echo ""
fi

# Check for integration tests
if [ -d "./tests/integration" ]; then
    echo -e "${YELLOW}Running Integration Tests...${NC}"
    echo -e "${GRAY}-------------------------------------------${NC}"
    
    INT_OUTPUT=$(go test -v -cover ./tests/integration/... 2>&1)
    INT_EXIT=$?
    
    INT_PASSED=$(echo "$INT_OUTPUT" | grep -c "^--- PASS:" 2>/dev/null)
    INT_FAILED=$(echo "$INT_OUTPUT" | grep -c "^--- FAIL:" 2>/dev/null)
    
    INT_PASSED=${INT_PASSED:-0}
    INT_FAILED=${INT_FAILED:-0}
    INT_PASSED=$((10#$INT_PASSED))
    INT_FAILED=$((10#$INT_FAILED))
    
    INT_TOTAL=$((INT_PASSED + INT_FAILED))
    
    if [ $INT_TOTAL -gt 0 ]; then
        INT_STATUS_COLOR=$GREEN
        INT_STATUS="PASSED"
        if [ $INT_FAILED -gt 0 ]; then
            INT_STATUS_COLOR=$RED
            INT_STATUS="FAILED"
        fi
        
        echo -e "  ${NC}Integration Tests:${NC}"
        echo -e "    Status:   ${INT_STATUS_COLOR}$INT_STATUS${NC}"
        echo -e "    Total:    ${GRAY}$INT_TOTAL tests${NC}"
        echo -e "    Passed:   ${GREEN}$INT_PASSED${NC}"
        echo -e "    Failed:   ${RED}$INT_FAILED${NC}"
        echo ""
        
        TOTAL=$((TOTAL + INT_TOTAL))
        PASSED=$((PASSED + INT_PASSED))
        FAILED=$((FAILED + INT_FAILED))
    fi
fi

# Check for API tests
if [ -d "./tests/api" ]; then
    echo -e "\033[0;35mRunning API Tests...${NC}"
    echo -e "${GRAY}-------------------------------------------${NC}"
    
    API_OUTPUT=$(go test -v -cover ./tests/api/... 2>&1)
    API_EXIT=$?
    
    API_PASSED=$(echo "$API_OUTPUT" | grep -c "^--- PASS:" 2>/dev/null)
    API_FAILED=$(echo "$API_OUTPUT" | grep -c "^--- FAIL:" 2>/dev/null)
    
    API_PASSED=${API_PASSED:-0}
    API_FAILED=${API_FAILED:-0}
    API_PASSED=$((10#$API_PASSED))
    API_FAILED=$((10#$API_FAILED))
    
    API_TOTAL=$((API_PASSED + API_FAILED))
    
    if [ $API_TOTAL -gt 0 ]; then
        API_STATUS_COLOR=$GREEN
        API_STATUS="PASSED"
        if [ $API_FAILED -gt 0 ]; then
            API_STATUS_COLOR=$RED
            API_STATUS="FAILED"
        fi
        
        echo -e "  ${NC}API Tests:${NC}"
        echo -e "    Status:   ${API_STATUS_COLOR}$API_STATUS${NC}"
        echo -e "    Total:    ${GRAY}$API_TOTAL tests${NC}"
        echo -e "    Passed:   ${GREEN}$API_PASSED${NC}"
        echo -e "    Failed:   ${RED}$API_FAILED${NC}"
        echo ""
        
        TOTAL=$((TOTAL + API_TOTAL))
        PASSED=$((PASSED + API_PASSED))
        FAILED=$((FAILED + API_FAILED))
    fi
fi

# Overall summary
echo -e "${GRAY}-------------------------------------------${NC}"
echo -e "${NC}Overall Results:${NC}"
echo ""
echo -e "  Total Tests:  ${NC}$TOTAL${NC}"
echo -e "  Passed:       ${GREEN}$PASSED${NC}"
echo -e "  Failed:       ${RED}$FAILED${NC}"
echo ""

# Pass rate
if [ $TOTAL -gt 0 ]; then
    PASS_RATE=$((PASSED * 100 / TOTAL))
    PASS_RATE_COLOR=$GREEN
    if [ $PASS_RATE -lt 90 ]; then
        PASS_RATE_COLOR=$YELLOW
    fi
    if [ $PASS_RATE -lt 70 ]; then
        PASS_RATE_COLOR=$RED
    fi
    echo -e "  Pass Rate:    ${PASS_RATE_COLOR}${PASS_RATE}%${NC}"
fi

echo ""
echo -e "${CYAN}==========================================${NC}"

# Exit with error if any tests failed
if [ $FAILED -gt 0 ]; then
    echo ""
    echo -e "${RED}Tests FAILED!${NC}"
    exit 1
else
    echo ""
    echo -e "${GREEN}All tests PASSED!${NC}"
    exit 0
fi
