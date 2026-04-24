# SonarQube + GitHub Actions Setup Guide

## 📋 Prerequisites

1. **GitHub Account** (already have)
2. **SonarCloud Account** (free tier available)
3. **Organization on SonarCloud** (or personal)

---

## 🚀 Step 1: Create SonarCloud Account

1. Go to [https://sonarcloud.io](https://sonarcloud.io)
2. Sign up with GitHub account
3. Choose your GitHub organization or personal account

---

## 🔑 Step 2: Generate SONAR_TOKEN

1. Go to [SonarCloud Account Security](https://sonarcloud.io/account/security/)
2. Generate a new token:
   - Token name: `GitHub Actions`
   - Click **Generate**
   - Copy the token

---

## 🔐 Step 3: Add Secret to GitHub Repository

1. Go to your GitHub repository
2. **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Name: `SONAR_TOKEN`
5. Value: Paste the token from SonarCloud
6. Click **Add secret**

---

## 📝 Step 4: Update sonar-project.properties

Edit `sonar-project.properties`:

```properties
sonar.projectKey=coolmate-ecommerce-backend
sonar.organization=YOUR-ORG-KEY  # Change this!
sonar.projectName=Coolmate E-commerce Backend
```

Replace `YOUR-ORG-KEY` with your SonarCloud organization key.

---

## 📄 Step 5: Enable GitHub Pages

1. Go to Repository **Settings** → **Pages**
2. Under "Build and deployment":
   - Source: **Deploy from a branch**
   - Branch: **gh-pages**
   - Folder: **/ (root)**
3. Click **Save**

---

## ✅ Step 6: Push and Test

```bash
git add .
git commit -m "Setup SonarQube and GitHub Pages integration"
git push origin main
```

The GitHub Actions workflow will:
- ✅ Run all tests
- 📊 Generate coverage reports
- 🔍 Scan code with SonarQube
- 📄 Create static HTML report
- 📈 Deploy to GitHub Pages

---

## 📊 Access Your Reports

After the workflow completes:

1. **SonarQube Analysis**
   - URL: `https://sonarcloud.io/dashboard?id=coolmate-ecommerce-backend`

2. **GitHub Pages Report**
   - URL: `https://USERNAME.github.io/REPO-NAME/report.html`

3. **GitHub Actions Artifacts**
   - Settings → Actions → Artifacts

---

## 🎯 What Gets Analyzed

### Code Coverage
- Overall coverage percentage
- Coverage by file and function
- Coverage trend over time

### Code Quality
- Code smells
- Security vulnerabilities
- Bugs and hotspots
- Complexity metrics

### Test Results
- Unit test count
- Pass/fail status
- Test execution time

---

## 📱 PR Integration

Every Pull Request will show:

```
📊 Test & Coverage Report

total: (statements) 70.4%

✅ All Tests Passed

Links:
- 📈 View Full Report
- 🔍 SonarQube Analysis
```

---

## 🔧 Workflow Features

### Automated Testing
- Runs on every push to `main` or `develop`
- Runs on all Pull Requests
- Uses PostgreSQL + Redis for integration tests
- 120-second timeout

### Coverage Metrics
- Atomic coverage mode
- HTML report generation
- Codecov integration
- Coverage summary extraction

### SonarQube Analysis
- Code quality metrics
- Security scanning
- Complexity analysis
- Technical debt assessment

### Reports Generated
1. `report.html` - Beautiful static report page
2. `coverage_report.html` - Interactive coverage map
3. `coverage_summary.txt` - Summary statistics
4. `test_report.json` - Raw test data

---

## 📈 Monitoring Coverage

View coverage trends on SonarQube:
- Dashboard shows history
- Track improvements over time
- Set quality gates
- Enforce standards

---

## ⚡ Troubleshooting

### SonarQube scan fails
- ✓ Check SONAR_TOKEN secret exists
- ✓ Verify organization key in sonar-project.properties
- ✓ Check network connectivity

### GitHub Pages not updating
- ✓ Verify gh-pages branch exists
- ✓ Check Pages settings point to gh-pages
- ✓ Wait 1-2 minutes for build

### Coverage not showing
- ✓ Ensure coverage.out is generated
- ✓ Check test execution isn't skipped
- ✓ Verify Go version is 1.22+

---

## 📚 Resources

- [SonarCloud Docs](https://docs.sonarcloud.io/)
- [GitHub Pages Docs](https://docs.github.com/en/pages)
- [Go Coverage Reports](https://golang.org/cmd/go/)
- [GitHub Actions Docs](https://docs.github.com/en/actions)

---

## 🎉 You're All Set!

Your CI/CD pipeline now includes:
- ✅ Automated testing
- ✅ Code coverage reporting
- ✅ Code quality analysis
- ✅ Beautiful static reports
- ✅ GitHub Pages integration

Happy coding! 🚀
