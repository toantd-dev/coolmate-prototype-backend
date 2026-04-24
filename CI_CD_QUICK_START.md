# 🚀 CI/CD Quick Start

## 1️⃣ Setup (Chỉ làm 1 lần)

```bash
# 1. Push project to GitHub
git remote add origin https://github.com/YOUR_USERNAME/coolmate-backend.git
git push -u origin main

# 2. Create branches
git checkout -b staging
git checkout -b develop

# 3. Add GitHub Secrets
# Go to: GitHub Repo → Settings → Secrets
# Add: DOCKER_USERNAME, DOCKER_PASSWORD, etc.
```

## 2️⃣ Start Using CI/CD

```bash
# Feature branch (chỉ test, không deploy)
git checkout develop
git commit -m "your changes"
git push origin develop
# ✅ Tests run automatically

# Staging deploy (test + deploy to staging)
git checkout staging
git merge develop
git push origin staging
# ✅ Tests run + auto deploy to staging

# Production deploy (test + deploy to prod)
git checkout main
git merge staging
git push origin main
# ✅ Tests run + auto deploy to production
```

## 3️⃣ Monitor

**GitHub Actions:**
- Repo → Actions tab → View logs

**Coverage:**
- codecov.io → Select repo

**Docker Hub:**
- hub.docker.com → Your images

---

## 📊 Coverage Requirements

- **Minimum:** 60%
- **Current:** 71.8% ✅
- **Target:** 80%+

If coverage < 60% → Deployment blocked ❌

---

## 🧪 Test Files

```
internal/services/*_test.go     (75 tests, 71.8% coverage)
internal/handlers/*_test.go     (some tests)
```

Run locally:
```bash
go test ./... -cover
```

---

## 🐳 Docker

Images pushed automatically:
- `coolmate-backend:latest` (main/staging)
- `coolmate-backend:abc123def` (commit hash)

---

## 🆘 If Tests Fail

1. **Local:** Run `go test ./...`
2. **GitHub:** Check Actions logs
3. **Coverage:** Use `go tool cover -html=coverage.out`

---

**Questions?** See `CI_CD_SETUP.md` for detailed guide
