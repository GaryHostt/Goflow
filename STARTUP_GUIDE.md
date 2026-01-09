# 🚀 GoFlow Startup Guide

## Quick Fix for Your Issues

### ✅ Issue 1: Backend - FIXED!
The `go.sum` file was corrupted. I've regenerated it with the correct checksums.

### ✅ Issue 2: Frontend - Easy Fix!
You ran the comment `#` as part of the command. See correct commands below.

---

## 🎯 How to Start GoFlow

### Terminal 1: Backend (Go API)

```bash
cd /Users/alex.macdonald/simple-ipass
go run cmd/api/main.go
```

**Expected output:**
```
2026/01/09 00:00:00 🚀 GoFlow API starting on :8080
2026/01/09 00:00:00 📊 Scheduler started
```

**Test it:**
```bash
curl http://localhost:8080/health
```

---

### Terminal 2: Frontend (Next.js)

**First time only - Install dependencies:**
```bash
cd /Users/alex.macdonald/simple-ipass/frontend
npm install
```

**Then run the dev server:**
```bash
npm run dev
```

**Expected output:**
```
  ▲ Next.js 14.0.4
  - Local:        http://localhost:3000
  - Network:      http://192.168.1.x:3000

 ✓ Ready in 2.3s
```

**Open in browser:**
```
http://localhost:3000
```

---

## ⚠️ Common Mistakes to Avoid

### ❌ DON'T DO THIS:
```bash
npm install  # First time only
```
The `#` is a comment and will cause an error!

### ✅ DO THIS INSTEAD:
```bash
# Run these commands separately:
npm install
npm run dev
```

---

## 🔍 Troubleshooting

### Backend Issues

**"go: command not found"**
```bash
# Install Go (if not installed)
brew install go

# Verify installation
go version
```

**"cannot find package"**
```bash
cd /Users/alex.macdonald/simple-ipass
go mod download
go run cmd/api/main.go
```

**"database is locked"**
```bash
# Remove the old database and restart
rm ipaas.db
go run cmd/api/main.go
```

---

### Frontend Issues

**"npm: command not found"**
```bash
# Install Node.js (if not installed)
brew install node

# Verify installation
node -v
npm -v
```

**"next: command not found"**
```bash
# You need to run npm install first!
cd frontend
npm install
npm run dev
```

**"Module not found: Can't resolve 'clsx'"**
```bash
# Install missing dependencies
cd frontend
npm install clsx tailwind-merge class-variance-authority
npm run dev
```

---

## 🎉 Success Checklist

Once both servers are running, verify:

### ✅ Backend Health Check
```bash
curl http://localhost:8080/health
```

**Expected:**
```json
{
  "status": "healthy",
  "version": "0.3.0",
  "uptime": "5s",
  "timestamp": "2026-01-09T00:00:00Z",
  "checks": {
    "database": "ok",
    "runtime": "ok"
  }
}
```

### ✅ Frontend Running
Open browser: **http://localhost:3000**

You should see:
- 🎨 GoFlow logo on login page
- 🔐 Login/Register forms
- 💅 Beautiful Tailwind CSS styling

---

## 🐳 Alternative: Docker Compose (Easiest!)

If you want to run everything with one command:

```bash
cd /Users/alex.macdonald/simple-ipass

# Start everything
docker-compose up

# Or run in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop everything
docker-compose down
```

**Services available:**
- Frontend: http://localhost:3000
- Backend: http://localhost:8080
- Kibana: http://localhost:5601
- Elasticsearch: http://localhost:9200

---

## 📝 Quick Test Workflow

Once both servers are running:

### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@goflow.dev","password":"password123"}'
```

Save the `token` from the response.

### 2. Create a Credential
```bash
TOKEN="your-token-here"

curl -X POST http://localhost:8080/api/credentials \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "slack",
    "api_key": "https://hooks.slack.com/services/TEST/TEST/TEST"
  }'
```

### 3. Create a Workflow
```bash
curl -X POST http://localhost:8080/api/workflows \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Workflow",
    "trigger_type": "webhook",
    "action_type": "slack_message",
    "config_json": "{\"credential_id\":\"cred_abc\",\"slack_message\":\"Hello!\"}"
  }'
```

### 4. Test in UI
1. Open http://localhost:3000
2. Login with `demo@goflow.dev` / `password123`
3. Navigate to Dashboard
4. See your workflow!

---

## 🎯 Summary

**To start GoFlow:**

```bash
# Terminal 1: Backend
cd /Users/alex.macdonald/simple-ipass
go run cmd/api/main.go

# Terminal 2: Frontend
cd /Users/alex.macdonald/simple-ipass/frontend
npm install  # First time only
npm run dev  # Run this after install completes
```

**Then visit:** http://localhost:3000

---

## 🆘 Still Having Issues?

1. **Check Go version:** `go version` (need 1.21+)
2. **Check Node version:** `node -v` (need 18+)
3. **Check ports:** Make sure 3000 and 8080 are free
4. **Try Docker:** `docker-compose up` for zero-config startup

---

**Your S-Tier enterprise platform is ready to run!** 🚀

**Files Fixed:**
- ✅ `go.sum` regenerated with correct checksums
- ✅ All dependencies properly specified

**Next Step:** Run the commands above in two separate terminals!

