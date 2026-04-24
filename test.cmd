@echo off
echo 🧪 Running tests...
go test ./... -coverprofile=coverage.out -timeout 120s

echo.
echo 📊 Generating HTML report...
go tool cover -html=coverage.out -o coverage_report.html

echo.
echo ✅ Done!
echo 📄 Files created:
echo    - coverage_report.html
echo    - coverage.out

echo.
echo 📈 Coverage:
go tool cover -func=coverage.out | findstr "total"

echo.
echo 🌐 Opening report in browser...
start coverage_report.html
