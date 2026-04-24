@echo off
setlocal enabledelayedexpansion

echo 🧪 Running tests with coverage...
go test ./... -coverprofile=coverage.out -timeout 120s

echo.
echo 📊 Generating HTML report...
go tool cover -html=coverage.out -o coverage_report.html

echo.
echo ✅ Done! Report files created:
echo    📄 coverage_report.html  - HTML interactive report
echo    📊 coverage.out          - Coverage data
echo.
echo 📈 Overall Coverage:
for /f "tokens=*" %%a in ('go tool cover -func^=coverage.out ^| findstr /R "total"') do (
    echo %%a
)
echo.
echo 🌐 Opening HTML report...
start coverage_report.html
