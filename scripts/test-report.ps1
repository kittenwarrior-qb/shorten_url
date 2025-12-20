# Test Report Script for Windows

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  Running All Tests with Coverage..." -ForegroundColor Cyan  
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Create temp directory for results
$tempDir = "$env:TEMP\test-results"
New-Item -ItemType Directory -Force -Path $tempDir | Out-Null

# Function to run tests for a specific category
function Run-TestCategory {
    param(
        [string]$Category,
        [string]$Path,
        [string]$Color
    )
    
    Write-Host ""
    Write-Host "Running $Category Tests..." -ForegroundColor $Color
    Write-Host "-------------------------------------------" -ForegroundColor Gray
    
    $jsonFile = "$tempDir\$Category.json"
    $startTime = Get-Date
    
    # Check if path exists
    if (-not (Test-Path $Path.Replace("/...", ""))) {
        Write-Host "  No tests found (path doesn't exist)" -ForegroundColor Gray
        return @{
            Category = $Category
            Passed = 0
            Failed = 0
            Skipped = 0
            Total = 0
            Duration = 0
            Coverage = "N/A"
            ExitCode = 0
        }
    }
    
    # Run tests
    $output = gotestsum --format dots-v2 --jsonfile $jsonFile -- -cover $Path 2>&1
    $exitCode = $LASTEXITCODE
    
    $duration = (Get-Date) - $startTime
    
    # Parse results if file exists
    if (Test-Path $jsonFile) {
        $results = Get-Content $jsonFile | ConvertFrom-Json
        $passed = ($results | Where-Object { $_.Action -eq "pass" -and $_.Test }).Count
        $failed = ($results | Where-Object { $_.Action -eq "fail" -and $_.Test }).Count
        $skipped = ($results | Where-Object { $_.Action -eq "skip" -and $_.Test }).Count
        
        # Extract coverage if available
        $coverageLine = $output | Select-String "coverage:" | Select-Object -Last 1
        $coverage = if ($coverageLine) { 
            ($coverageLine -replace '.*coverage:\s*(\d+\.\d+)%.*', '$1') 
        } else { 
            "N/A" 
        }
        
        return @{
            Category = $Category
            Passed = $passed
            Failed = $failed
            Skipped = $skipped
            Total = $passed + $failed + $skipped
            Duration = $duration.TotalSeconds
            Coverage = $coverage
            ExitCode = $exitCode
        }
    } else {
        return @{
            Category = $Category
            Passed = 0
            Failed = 0
            Skipped = 0
            Total = 0
            Duration = $duration.TotalSeconds
            Coverage = "N/A"
            ExitCode = $exitCode
        }
    }
}

# Run different test categories
$unitResults = Run-TestCategory -Category "Unit" -Path "./tests/unit/..." -Color "Green"
$integrationResults = Run-TestCategory -Category "Integration" -Path "./tests/integration/..." -Color "Yellow"
$apiResults = Run-TestCategory -Category "API" -Path "./tests/api/..." -Color "Magenta"

# Calculate totals
$allResults = @($unitResults, $integrationResults, $apiResults)
$totalPassed = 0
$totalFailed = 0
$totalSkipped = 0
$totalTests = 0
$totalDuration = 0

foreach ($result in $allResults) {
    $totalPassed += $result.Passed
    $totalFailed += $result.Failed
    $totalSkipped += $result.Skipped
    $totalTests += $result.Total
    $totalDuration += $result.Duration
}

# Print detailed summary
Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "           Test Summary Report" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Category breakdown
Write-Host "By Category:" -ForegroundColor White
Write-Host ""
foreach ($result in $allResults) {
    if ($result.Total -gt 0) {
        $statusColor = if ($result.Failed -gt 0) { "Red" } else { "Green" }
        $status = if ($result.Failed -gt 0) { "FAILED" } else { "PASSED" }
        
        Write-Host "  $($result.Category) Tests:" -ForegroundColor White
        Write-Host "    Status:   $status" -ForegroundColor $statusColor
        Write-Host "    Total:    $($result.Total) tests" -ForegroundColor Gray
        Write-Host "    Passed:   $($result.Passed)" -ForegroundColor Green
        Write-Host "    Failed:   $($result.Failed)" -ForegroundColor Red
        Write-Host "    Skipped:  $($result.Skipped)" -ForegroundColor Yellow
        Write-Host "    Coverage: $($result.Coverage)%" -ForegroundColor Cyan
        Write-Host "    Duration: $([math]::Round($result.Duration, 2))s" -ForegroundColor Gray
        Write-Host ""
    }
}

# Overall summary
Write-Host "-------------------------------------------" -ForegroundColor Gray
Write-Host "Overall Results:" -ForegroundColor White
Write-Host ""
Write-Host "  Total Tests:  $totalTests" -ForegroundColor White
Write-Host "  Passed:       $totalPassed" -ForegroundColor Green
Write-Host "  Failed:       $totalFailed" -ForegroundColor Red
Write-Host "  Skipped:      $totalSkipped" -ForegroundColor Yellow
Write-Host "  Duration:     $([math]::Round($totalDuration, 2))s" -ForegroundColor Gray
Write-Host ""

# Pass rate
if ($totalTests -gt 0) {
    $passRate = [math]::Round(($totalPassed / $totalTests) * 100, 2)
    $passRateColor = if ($passRate -ge 90) { "Green" } elseif ($passRate -ge 70) { "Yellow" } else { "Red" }
    Write-Host "  Pass Rate:    $passRate%" -ForegroundColor $passRateColor
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan

# Cleanup
Remove-Item -Recurse -Force $tempDir -ErrorAction SilentlyContinue

# Exit with error if any tests failed
$hasFailures = ($allResults | Where-Object { $_.Failed -gt 0 }).Count -gt 0
if ($hasFailures) {
    Write-Host ""
    Write-Host "Tests FAILED!" -ForegroundColor Red
    exit 1
} else {
    Write-Host ""
    Write-Host "All tests PASSED!" -ForegroundColor Green
    exit 0
}
