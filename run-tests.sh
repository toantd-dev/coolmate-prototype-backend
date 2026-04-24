#!/bin/bash

# Generate test report and HTML coverage
echo "🧪 Running tests with coverage..."
go test ./... -coverprofile=coverage.out -timeout 120s

echo "📊 Generating HTML report..."
go tool cover -html=coverage.out -o coverage_report.html

echo ""
echo "✅ Done! Report files created:"
echo "   📄 coverage_report.html  - HTML interactive report"
echo "   📊 coverage.out          - Coverage data"
echo ""
echo "📈 Overall Coverage:"
go tool cover -func=coverage.out | tail -1
echo ""
echo "🌐 Open report: start coverage_report.html"
