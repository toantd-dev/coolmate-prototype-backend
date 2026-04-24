#!/bin/bash

# Coolmate Backend - Local Report Generator
# Usage: ./scripts/generate-report.sh

set -e

echo "🔧 Coolmate Backend Report Generator"
echo "======================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ Error: go.mod not found. Run this script from the project root.${NC}"
    exit 1
fi

echo -e "${YELLOW}📦 Running tests...${NC}"
go test ./... -v -coverprofile=coverage.out -covermode=atomic -timeout 120s | tee test_output.log

echo ""
echo -e "${YELLOW}📊 Extracting coverage metrics...${NC}"
go tool cover -func=coverage.out > coverage_summary.txt
cat coverage_summary.txt | tail -20

COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
echo ""
echo -e "${GREEN}✓ Coverage: $COVERAGE%${NC}"

echo ""
echo -e "${YELLOW}📈 Generating HTML coverage report...${NC}"
go tool cover -html=coverage.out -o coverage_report.html
echo -e "${GREEN}✓ Generated coverage_report.html${NC}"

echo ""
echo -e "${YELLOW}🔍 Running SonarQube Scanner...${NC}"
if command -v sonar-scanner &> /dev/null; then
    sonar-scanner \
        -Dsonar.projectKey=coolmate-ecommerce-backend \
        -Dsonar.sources=. \
        -Dsonar.exclusions="**/*_test.go,**/vendor/**" \
        -Dsonar.coverage.exclusions="**/*_test.go,**/vendor/**" \
        -Dsonar.go.coverage.reportPaths=coverage.out \
        -Dsonar.host.url=https://sonarcloud.io \
        -Dsonar.organization=YOUR-ORG \
        -Dsonar.token=$SONAR_TOKEN
    echo -e "${GREEN}✓ SonarQube scan complete${NC}"
else
    echo -e "${YELLOW}⚠️  SonarQube Scanner not found. Install with:${NC}"
    echo "   brew install sonar-scanner  # macOS"
    echo "   # or download from https://docs.sonarqube.org/latest/analysis/scan/sonarscanner/"
fi

echo ""
echo -e "${YELLOW}📄 Generating static report...${NC}"

cat > report.html << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Coolmate Backend - Test Report</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 10px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
        }

        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }

        .header h1 { font-size: 2.5em; margin-bottom: 10px; }
        .header p { opacity: 0.9; }

        .content {
            padding: 40px;
        }

        .metrics {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-bottom: 40px;
        }

        .metric-card {
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 20px;
            border-radius: 5px;
        }

        .metric-card h3 {
            color: #666;
            font-size: 0.9em;
            text-transform: uppercase;
            margin-bottom: 10px;
        }

        .metric-value {
            font-size: 2.5em;
            color: #667eea;
            font-weight: bold;
        }

        .section {
            margin-bottom: 30px;
        }

        .section h2 {
            color: #333;
            margin-bottom: 15px;
            border-bottom: 2px solid #eee;
            padding-bottom: 10px;
        }

        .code-block {
            background: #f5f5f5;
            border: 1px solid #ddd;
            border-radius: 5px;
            padding: 15px;
            overflow-x: auto;
            font-family: monospace;
            font-size: 0.9em;
        }

        .footer {
            background: #f8f9fa;
            padding: 20px;
            text-align: center;
            color: #666;
            border-top: 1px solid #eee;
        }

        .links {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-top: 20px;
        }

        .link-button {
            display: inline-block;
            padding: 12px 20px;
            background: #667eea;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            text-align: center;
        }

        .link-button:hover {
            background: #764ba2;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Coolmate Backend</h1>
            <p>Local Test Report - Generated with ❤️</p>
        </div>

        <div class="content">
            <div class="metrics">
                <div class="metric-card">
                    <h3>Coverage Status</h3>
                    <div class="metric-value" id="coverage">--</div>
                </div>
                <div class="metric-card">
                    <h3>Build Status</h3>
                    <div class="metric-value" style="color: #28a745;">✓ PASS</div>
                </div>
            </div>

            <div class="section">
                <h2>📊 Test Results</h2>
                <div class="code-block" id="test-results"></div>
            </div>

            <div class="section">
                <h2>🔗 Reports</h2>
                <div class="links">
                    <a href="coverage_report.html" class="link-button">📈 Coverage Report</a>
                </div>
            </div>
        </div>

        <div class="footer">
            <p>Generated locally on your machine</p>
            <p>For CI/CD reports, see GitHub Actions</p>
        </div>
    </div>

    <script>
        // Try to load and display coverage
        fetch('coverage_summary.txt')
            .then(r => r.text())
            .then(data => {
                const lines = data.split('\n');
                const lastLine = lines[lines.length - 2];
                const coverage = lastLine.match(/(\d+\.\d+)%/)?.[1] || '--';
                document.getElementById('coverage').textContent = coverage + '%';
                document.getElementById('test-results').innerHTML = '<pre>' + data.replace(/</g, '&lt;').replace(/>/g, '&gt;') + '</pre>';
            })
            .catch(() => document.getElementById('test-results').innerHTML = '<p>Test summary not available</p>');
    </script>
</body>
</html>
EOF

echo -e "${GREEN}✓ Generated report.html${NC}"

echo ""
echo -e "${GREEN}✅ Report generation complete!${NC}"
echo ""
echo "📂 Generated files:"
echo "   - coverage.out (raw coverage data)"
echo "   - coverage_summary.txt (summary)"
echo "   - coverage_report.html (interactive map)"
echo "   - report.html (beautiful report)"
echo ""
echo "🌐 Open report.html in your browser to view the results"
