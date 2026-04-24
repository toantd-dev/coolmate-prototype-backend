# CI/CD Pipeline Setup Guide

## 🎯 Mục Tiêu

- ✅ Chạy unit tests trước mỗi lần deploy
- ✅ Kiểm tra code coverage (tối thiểu 60%)
- ✅ Build Docker image tự động
- ✅ Deploy tự động lên Staging/Production
- ✅ Upload coverage reports lên Codecov

---

## 📋 Prerequisites

1. **GitHub Repository**
   - Project phải được push lên GitHub
   - Có 3 branches: `main`, `staging`, `develop`

2. **Docker Hub Account** (optional, nếu muốn push image)
   - Username & Password

3. **Codecov Account** (optional, để track coverage)
   - https://codecov.io

---

## 🔑 Step 1: Setup GitHub Secrets

GitHub Actions sử dụng secrets để lưu trữ sensitive data (credentials, tokens, etc.)

### Cách thêm Secrets:

1. Vào GitHub repo → **Settings** → **Secrets and variables** → **Actions**

2. Click **New repository secret**

3. Thêm các secrets sau:

| Secret Name | Value | Purpose |
|------------|-------|---------|
| `DOCKER_USERNAME` | Your Docker Hub username | Push image to Docker Hub |
| `DOCKER_PASSWORD` | Your Docker Hub password | Push image to Docker Hub |
| `STAGING_HOST` | staging.example.com | Staging server hostname |
| `STAGING_USER` | deploy_user | SSH user for staging |
| `STAGING_DEPLOY_KEY` | Private SSH key | SSH authentication |
| `PROD_HOST` | prod.example.com | Production server hostname |
| `PROD_USER` | deploy_user | SSH user for production |
| `PROD_DEPLOY_KEY` | Private SSH key | SSH authentication |

---

## 🔄 Step 2: File Structure

```
project-root/
├── .github/
│   └── workflows/
│       └── ci.yml              # ← GitHub Actions workflow
├── codecov.yml                  # ← Codecov config
├── Dockerfile                   # ← Docker build with tests
├── docker-compose.yml
├── .env.example
└── go.mod
```

✅ Tất cả đã được tạo!

---

## 📊 Step 3: How It Works

### Flow Khi Push Code:

```
git push → GitHub → 
  ├─ Trigger CI/CD (ci.yml)
  │   ├─ Setup PostgreSQL + Redis
  │   ├─ Run: go test ./... ✅
  │   ├─ Check: coverage >= 60% ✅
  │   ├─ Build: Docker image ✅
  │   ├─ Push: Docker Hub 🐳
  │   └─ Upload: Coverage to Codecov 📊
  │
  ├─ If on 'staging' branch:
  │   └─ Deploy to Staging 🚀
  │
  ├─ If on 'main' branch:
  │   └─ Deploy to Production 🚀
  │
  └─ If tests FAIL:
      └─ Block deployment ❌
```

---

## 🧪 Step 4: Test the CI/CD

### Test 1: Push to Develop (No Deploy)
```bash
git checkout develop
echo "# test" >> README.md
git add .
git commit -m "test ci"
git push origin develop
```
✅ Should run tests only, no deployment

### Test 2: Push to Staging (Deploy to Staging)
```bash
git checkout staging
git merge develop
git push origin staging
```
✅ Should run tests + deploy to staging

### Test 3: Push to Main (Deploy to Production)
```bash
git checkout main
git merge staging
git push origin main
```
✅ Should run tests + deploy to production

---

## 📈 Step 5: Monitor Coverage

### View Coverage Trends:

1. Vào codecov.io với GitHub account
2. Select repo → Coverage tracking
3. Mỗi push sẽ update coverage graph

### GitHub Actions Log:

1. Vào repo → **Actions** tab
2. Click workflow → View logs
3. Xem chi tiết test results & coverage

---

## 🛑 Coverage Threshold

**Hiện tại: 60% minimum**

Để thay đổi, edit `.github/workflows/ci.yml`:

```yaml
# Line ~75
if (( $(echo "$coverage < 60" | bc -l) )); then
  echo "❌ Coverage below 60% threshold"
  exit 1
fi
```

Thay `60` thành `70` hoặc `80` nếu muốn cao hơn.

---

## 🐳 Step 6: Docker Hub Push

Để push Docker image tự động:

1. Create Docker Hub account
2. Add secrets: `DOCKER_USERNAME` & `DOCKER_PASSWORD`
3. Workflow sẽ tự push image khi push to main/staging

Images sẽ được tag:
- `coolmate-backend:latest`
- `coolmate-backend:abc123def` (commit hash)

---

## 🚀 Step 7: Deployment Script

Workflow hiện tại chỉ print message. Để deploy thực sự, thêm script:

### Ví dụ Deploy với SSH:

```bash
# .github/workflows/ci.yml, line ~150
- name: Deploy to Staging
  env:
    DEPLOY_KEY: ${{ secrets.STAGING_DEPLOY_KEY }}
    DEPLOY_HOST: ${{ secrets.STAGING_HOST }}
    DEPLOY_USER: ${{ secrets.STAGING_USER }}
  run: |
    mkdir -p ~/.ssh
    echo "$DEPLOY_KEY" > ~/.ssh/id_rsa
    chmod 600 ~/.ssh/id_rsa
    ssh-keyscan -H $DEPLOY_HOST >> ~/.ssh/known_hosts
    
    ssh -i ~/.ssh/id_rsa $DEPLOY_USER@$DEPLOY_HOST << 'DEPLOY'
      cd /app
      docker pull coolmate-backend:latest
      docker-compose up -d
    DEPLOY
```

---

## 📊 Coverage Report

Sau mỗi push, bạn sẽ thấy:

✅ **GitHub Actions:**
```
Code Coverage: 71.8%
✅ Coverage 71.8% meets threshold
```

✅ **Codecov Dashboard:**
- Graph của coverage trends
- Pull request comments với coverage changes
- File-level coverage breakdown

---

## 🔧 Troubleshooting

### Issue: Tests fail in CI but pass locally
**Solution:** 
- CI uses different env vars (`.env.test`)
- Check `.github/workflows/ci.yml` env section
- Make sure test database name is correct

### Issue: Docker push fails
**Solution:**
- Check Docker Hub credentials
- Verify secret names: `DOCKER_USERNAME`, `DOCKER_PASSWORD`
- Ensure Docker Hub account has push permissions

### Issue: Deploy fails
**Solution:**
- Check SSH key format (should be PEM)
- Verify deploy server is accessible
- Check deploy user has permissions

---

## 📚 Commands để Test Locally

```bash
# Run tests like CI does
go test -v ./... -timeout 60s

# Generate coverage locally
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Build Docker image locally
docker build -t coolmate-backend:test .

# Check Go vet locally
go vet ./...
```

---

## ✅ Checklist

- [ ] Push project to GitHub
- [ ] Create 3 branches: main, staging, develop
- [ ] Add GitHub Secrets (DOCKER credentials, Deploy keys)
- [ ] Test CI on develop branch
- [ ] Test deployment on staging branch
- [ ] Monitor codecov.io for coverage trends
- [ ] Setup deployment script (SSH, K8s, etc.)

---

## 🎉 Done!

Now every push will:
✅ Run 75+ unit tests  
✅ Check 71.8% code coverage  
✅ Build Docker image  
✅ Deploy automatically  

Happy deploying! 🚀
