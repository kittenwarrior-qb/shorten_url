# Test Report Script for Windows

Write-Host ""
Write-Host "==============================" -ForegroundColor Cyan
Write-Host "  Running Unit Tests..." -ForegroundColor Cyan  
Write-Host "==============================" -ForegroundColor Cyan
Write-Host ""

# Run tests with JSON output
$jsonFile = "$env:TEMP\test-results.json"
gotestsum --format dots-v2 --jsonfile $jsonFile -- -cover ./tests/...

# Parse JSON results
$results = Get-Content $jsonFile | ConvertFrom-Json

$passed = ($results | Where-Object { $_.Action -eq "pass" -and $_.Test }).Count
$failed = ($results | Where-Object { $_.Action -eq "fail" -and $_.Test }).Count  
$skipped = ($results | Where-Object { $_.Action -eq "skip" -and $_.Test }).Count
$total = $passed + $failed + $skipped

Write-Host ""
Write-Host "==============================" -ForegroundColor Cyan
Write-Host "        Test Summary" -ForegroundColor Cyan
Write-Host "==============================" -ForegroundColor Cyan
Write-Host ""
Write-Host "  Total:    $total tests" -ForegroundColor White
Write-Host "  Passed:   $passed" -ForegroundColor Green
Write-Host "  Failed:   $failed" -ForegroundColor Red
Write-Host "  Skipped:  $skipped" -ForegroundColor Yellow
Write-Host ""
Write-Host "==============================" -ForegroundColor Cyan

# Exit with error code if tests failed
if ($failed -gt 0) {
    exit 1
}
