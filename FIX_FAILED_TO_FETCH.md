# 🔧 Quick Fix: "Failed to fetch" Error on Registration

## Problem
When trying to create an account at http://localhost:3000/register, you get:
```
Error
Failed to fetch
```

## Root Cause
The frontend Docker container was built without the correct API URL environment variable. Next.js requires `NEXT_PUBLIC_*` variables to be set at **build time**, not runtime.

## ✅ Solution: Rebuild Frontend

Run these commands:

```bash
cd /Users/alex.macdonald/simple-ipass

# Stop and rebuild the frontend
docker compose stop frontend
docker compose build --no-cache frontend
docker compose up -d frontend

# Wait 30 seconds for frontend to start, then test
```

## ✅ Test the Fix

1. **Open the frontend**: http://localhost:3000

2. **Try to register**:
   - Email: `test@example.com`
   - Password: `password123`

3. **Should see success** and redirect to login

## 🔍 Verify Backend is Accessible

Before registering, verify the backend is working:

```bash
# Test backend health
curl http://localhost:8080/health

# Should return:
# {"status":"healthy","version":"0.2.0"}
```

If the backend is not accessible, restart it:

```bash
docker compose restart backend
```

## 🌐 What Was Fixed

### **Before** ❌
```dockerfile
# Frontend Dockerfile was missing API URL during build
RUN npm run build  # No environment variable!
```

```yaml
# docker-compose.yml was setting it at runtime (too late!)
environment:
  - NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

### **After** ✅
```dockerfile
# Frontend Dockerfile now sets API URL during build
ARG NEXT_PUBLIC_API_URL=http://localhost:8080/api
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
RUN npm run build  # Now it knows the API URL!
```

```yaml
# docker-compose.yml passes build arg
build:
  args:
    NEXT_PUBLIC_API_URL: http://localhost:8080/api
```

## 🔄 Alternative: Run Frontend Locally (For Development)

If you want faster iteration without Docker:

```bash
cd /Users/alex.macdonald/simple-ipass/frontend

# Create .env.local file
echo "NEXT_PUBLIC_API_URL=http://localhost:8080/api" > .env.local

# Install dependencies (first time only)
npm install

# Run development server
npm run dev
```

Then open: http://localhost:3000

**Advantage**: Changes are reflected instantly without rebuilding Docker image.

## ✅ Verification Steps

After rebuilding, verify everything works:

### 1. **Check Services Status**
```bash
docker compose ps
```

All services should show "Up" and "healthy" (or "(healthy)" in STATUS column).

### 2. **Test Backend API Directly**
```bash
# Register via API
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Should return:
# {"message":"User registered successfully"}
```

### 3. **Test Frontend**
Open http://localhost:3000 and try to register again.

## 🐛 Still Having Issues?

### **Check Browser Console**
1. Open DevTools (F12 or Cmd+Option+I)
2. Go to Console tab
3. Try to register again
4. Look for error messages

Common errors:
- **CORS error**: Backend not allowing frontend origin
- **Network error**: Backend not running
- **404 error**: Wrong API URL

### **Check Backend Logs**
```bash
docker compose logs backend -f
```

Look for:
- ✅ `Server starting on :8080`
- ✅ `Database connection successful`
- ❌ Any error messages

### **Check Frontend Logs**
```bash
docker compose logs frontend -f
```

Look for:
- ✅ `ready - started server on 0.0.0.0:3000`
- ❌ Build errors
- ❌ API connection errors

### **Nuclear Option: Full Restart**
```bash
cd /Users/alex.macdonald/simple-ipass

# Stop everything
docker compose down

# Rebuild everything (takes 5-10 minutes)
docker compose build --no-cache

# Start everything
docker compose up -d

# Wait 1-2 minutes, then test
```

## 📚 Related Documentation

- `COMPONENT_RUNNING_GUIDE.md` - Complete guide to running each component
- `TESTING_QUICK_START.md` - Testing guide
- `README.md` - Platform overview

## ✅ Expected Result

After the fix, registration should work:

1. Go to http://localhost:3000
2. Click "Register" (or go to /register)
3. Enter email and password
4. Click "Register"
5. **Success!** Redirected to login page
6. Login with same credentials
7. **Success!** Redirected to dashboard

## 🎉 Success!

Once you see the dashboard, your platform is working correctly! 

Next steps:
1. Connect to Slack/Discord (add webhook URLs)
2. Create your first workflow
3. View execution logs

---

**Quick Command Reference:**

```bash
# Fix the issue
docker compose build --no-cache frontend && docker compose up -d frontend

# Check status
docker compose ps

# View logs
docker compose logs frontend -f

# Test backend
curl http://localhost:8080/health

# Open frontend
open http://localhost:3000
```

---

**Need more help?** Check `COMPONENT_RUNNING_GUIDE.md` for detailed troubleshooting steps!

