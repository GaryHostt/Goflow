# 🚀 Dev Mode - Skip Login Feature

## Overview

**Dev Mode** is a development-only feature that allows you to bypass the login/registration process and instantly access the platform with a pre-configured test account.

**Perfect for:**
- 🏃 Quick integration testing
- 🔧 Rapid workflow development
- 🐛 Debugging without authentication hassles
- 🎨 UI/UX iteration

---

## ✅ How to Enable

### **Frontend (Local Development)**

The dev mode button is automatically enabled when running locally:

```bash
cd /Users/alex.macdonald/simple-ipass
./scripts/run_frontend_locally.sh
```

This sets `NEXT_PUBLIC_DEV_MODE=true` in your `.env.local` file.

### **Manual Setup**

If you prefer manual setup:

```bash
cd /Users/alex.macdonald/simple-ipass/frontend

# Create or edit .env.local
cat > .env.local << 'EOF'
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_DEV_MODE=true
EOF

npm run dev
```

---

## 🎯 How to Use

1. **Start backend** (in Docker or locally)
   ```bash
   docker compose up -d backend postgres
   # OR
   go run cmd/api/main.go
   ```

2. **Start frontend** with dev mode
   ```bash
   ./scripts/run_frontend_locally.sh
   ```

3. **Open**: http://localhost:3000

4. **Click**: "Skip Login - Dev Mode" button (orange button with ⚡ icon)

5. **Done!** You're instantly logged in and redirected to the dashboard

---

## 🔐 What Happens Behind the Scenes

1. **Backend creates/fetches dev user**:
   - Email: `dev@goflow.local`
   - Password: `dev123`
   - Auto-creates if doesn't exist

2. **Generates JWT token**

3. **Frontend stores token** in localStorage

4. **Redirects to dashboard**

---

## 📊 Dev User Credentials

If you need to manually login as the dev user:

```
Email: dev@goflow.local
Password: dev123
```

---

## 🔒 Security

**Dev Mode is ONLY available in development:**

- ✅ **Backend**: Only enabled when `ENVIRONMENT != production`
- ✅ **Frontend**: Only shows button when `NEXT_PUBLIC_DEV_MODE=true`
- ✅ **Production**: `/api/auth/dev-login` endpoint is **disabled**

**Never set `NEXT_PUBLIC_DEV_MODE=true` in production!**

---

## 🎨 UI Design

The dev mode button features:
- ⚡ **Lightning bolt icon** for quick recognition
- 🟧 **Orange color** to indicate "developer tool"
- 📏 **Separator line** to visually distinguish from normal login
- 💬 **Clear label**: "Skip Login - Dev Mode"

---

## 🧪 Testing Workflow

### **Quick Integration Testing Loop:**

```bash
# 1. Start everything
./scripts/run_frontend_locally.sh

# 2. Open browser: http://localhost:3000

# 3. Click "Skip Login - Dev Mode"

# 4. You're in! Now:
   - Create workflows
   - Test connectors
   - Debug issues
   - Iterate quickly

# 5. Make code changes → See results immediately (hot reload!)
```

---

## 🔄 Toggling Dev Mode On/Off

### **Disable Dev Mode**

Edit `.env.local`:
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_DEV_MODE=false  # or remove this line
```

Then restart: `npm run dev`

### **Re-enable Dev Mode**

Edit `.env.local`:
```bash
NEXT_PUBLIC_DEV_MODE=true
```

Then restart: `npm run dev`

---

## 🐛 Troubleshooting

### **Dev Mode button not showing**

**Check 1**: Verify environment variable
```bash
cd frontend
cat .env.local | grep DEV_MODE

# Should show:
# NEXT_PUBLIC_DEV_MODE=true
```

**Check 2**: Restart frontend
```bash
# Stop frontend (Ctrl+C)
npm run dev
```

**Check 3**: Clear browser cache
- Hard reload: Cmd+Shift+R (Mac) or Ctrl+Shift+R (Windows)

### **"Dev login failed" error**

**Check 1**: Backend is running
```bash
curl http://localhost:8080/health

# Should return: {"status":"healthy","version":"0.2.0"}
```

**Check 2**: Backend is in development mode
```bash
docker compose logs backend | grep "Dev mode enabled"

# Should see: Dev mode enabled - /api/auth/dev-login endpoint available
```

**Check 3**: Check backend logs
```bash
docker compose logs backend --tail=50
```

### **Button shows but doesn't work**

**Check browser console** (F12 → Console tab):
```
Look for:
- Network errors
- API endpoint errors
- CORS errors
```

**Test endpoint directly**:
```bash
curl -X POST http://localhost:8080/api/auth/dev-login

# Should return: {"token":"...","user":{...}}
```

---

## 📝 Implementation Details

### **Backend Endpoint**

**File**: `internal/handlers/auth.go`

```go
func (h *AuthHandler) DevLogin(w http.ResponseWriter, r *http.Request) {
    const (
        devEmail    = "dev@goflow.local"
        devPassword = "dev123"
    )
    
    // Get or create dev user
    // Generate JWT token
    // Return auth response
}
```

**Registered in**: `cmd/api/main.go`

```go
if getEnv("ENVIRONMENT", "development") == "development" {
    router.HandleFunc("/api/auth/dev-login", authHandler.DevLogin).Methods("POST")
}
```

### **Frontend Component**

**File**: `frontend/app/login/page.tsx`

```typescript
const DEV_MODE = process.env.NEXT_PUBLIC_DEV_MODE === 'true'

const handleDevLogin = async () => {
    const response = await api.post('/auth/dev-login')
    setToken(response.token)
    router.push('/dashboard')
}
```

---

## 🎯 Use Cases

### **1. Quick Integration Development**
```
Click Dev Mode → Dashboard → Create Workflow → Test → Iterate
No login hassle every time!
```

### **2. Connector Testing**
```
Dev Mode → Connections → Add Slack → Test → Debug
Fast feedback loop!
```

### **3. UI/UX Development**
```
Dev Mode → Dashboard → Make UI changes → See results
No need to login after every refresh!
```

### **4. Demo Preparation**
```
Dev Mode → Pre-configure workflows → Ready to demo
Skip the setup!
```

---

## 🚫 What NOT to Do

❌ **Don't use in production**
- Security risk: Anyone can access the platform
- Dev endpoint is automatically disabled in production

❌ **Don't commit `.env.local` to git**
- Already in `.gitignore`
- Contains local development settings

❌ **Don't share dev user in multi-developer environment**
- Data will conflict
- Use normal accounts for shared environments

---

## ✅ Best Practices

✅ **Use for solo development**
- Perfect for quick testing and iteration

✅ **Use normal accounts for team work**
- Create separate accounts for each developer

✅ **Disable before production build**
- Set `NEXT_PUBLIC_DEV_MODE=false`
- Or remove the variable entirely

✅ **Keep dev user for testing**
- Useful for automated testing scripts
- Can be reset easily

---

## 🔄 Related Scripts

**Run frontend with dev mode**:
```bash
./scripts/run_frontend_locally.sh
```

**Run everything in Docker** (no dev mode):
```bash
docker compose up -d
```

**Fix frontend issues**:
```bash
./scripts/fix_frontend.sh
```

---

## 📖 Additional Documentation

- `RUN_FRONTEND_LOCALLY.md` - Running frontend locally
- `COMPONENT_RUNNING_GUIDE.md` - Running all components
- `COMPREHENSIVE_FIX.md` - Troubleshooting guide

---

## 🎉 Summary

**Dev Mode = Fastest Development Workflow!**

```
1 Click → Logged In → Building Integrations
```

No more:
- ❌ Creating accounts
- ❌ Remembering passwords
- ❌ Logging in every time
- ❌ Re-authenticating after errors

Just:
- ✅ Click
- ✅ Code
- ✅ Test
- ✅ Ship

---

**Happy developing!** 🚀⚡

