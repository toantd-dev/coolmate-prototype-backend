@echo off
REM Coolmate Backend - Local Report Generator (Windows)
REM Usage: scripts\generate-report.bat

setlocal enabledelayedexpansion

echo.
echo ============================================================
echo   Coolmate Backend Report Generator
echo ============================================================
echo.

REM Check if go.mod exists
if not exist "go.mod" (
    echo ERROR: go.mod not found. Run this from the project root.
    exit /b 1
)

echo [1/4] Running tests...
go test ./... -v -coverprofile=coverage.out -covermode=atomic -timeout 120s > test_output.log 2>&1

echo [2/4] Extracting coverage metrics...
go tool cover -func=coverage.out > coverage_summary.txt
type coverage_summary.txt | findstr /R "^total"

echo.
echo [3/4] Generating HTML coverage report...
go tool cover -html=coverage.out -o coverage_report.html
echo OK - Generated coverage_report.html

echo.
echo [4/4] Generating static report...

REM Create report.html using PowerShell
powershell -NoProfile -ExecutionPolicy Bypass -File "scripts\generate-report.ps1"

echo.
echo ============================================================
echo   Report generation complete!
echo ============================================================
echo.
echo Generated files:
echo   - coverage.out (raw coverage data)
echo   - coverage_summary.txt (summary)
echo   - coverage_report.html (interactive map)
echo   - report.html (beautiful report)
echo.
echo Open report.html in your browser to view the results
echo.
