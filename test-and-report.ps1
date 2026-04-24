# PowerShell script for testing and coverage report

Write-Host "🧪 Running tests..." -ForegroundColor Cyan
go test ./... -coverprofile=coverage.out -timeout 120s

Write-Host "`n📊 Generating HTML report..." -ForegroundColor Cyan
go tool cover -html=coverage.out -o coverage_report.html

Write-Host "`n✅ Done!" -ForegroundColor Green
Write-Host "📄 Files created:"
Write-Host "   - coverage_report.html (HTML interactive report)"
Write-Host "   - coverage.out (Coverage data)"

Write-Host "`n📈 Overall Coverage:"
go tool cover -func=coverage.out | Select-String "^total"

Write-Host "`n🌐 Opening report..." -ForegroundColor Cyan
Start-Process coverage_report.html
