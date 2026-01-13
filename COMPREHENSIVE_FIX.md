# 🔴 COMPREHENSIVE FIX - "Failed to fetch" Issue

## ✅ Backend is Working!

I just tested the backend API and it's working perfectly:

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# ✅ Response: {"token":"...","user":{...}}
```

**The problem is the frontend is not connecting to the backend correctly.**

---

## 🔧 Solution: Force Rebuild Frontend

The frontend Docker container is cached with the wrong API URL. Let's fix this:

### **Step 1: Stop and Remove Frontend Container**

```bash
cd /Users/alex.macdonald/simple-ipass

# Stop and remove the frontend container completely
docker compose rm -sf frontend
```

### **Step 2: Rebuild with No Cache**

```bash
# Rebuild frontend from scratch (no cache)
docker compose build --no-cache --pull frontend
```

This will take 2-3 minutes.

### **Step 3: Start Frontend**

```bash
# Start the frontend
docker compose up -d frontend
```

### **Step 4: Wait and Test**

```bash
# Wait 30 seconds for frontend to start
sleep 30

# Check if it's running
docker compose ps frontend

# Should show "Up" and "(healthy)" or "Up (health: starting)"
```

### **Step 5: Test in Browser**

1. **Clear your browser cache** (important!):
   - **Chrome/Edge**: Cmd+Shift+Delete → Clear cache
   - **Safari**: Cmd+Option+E → Empty Caches
   - **Or**: Open in Incognito/Private window

2. **Open**: http://localhost:3000

3. **Try to register**:
   - Email: `newuser@example.com`
   - Password: `password123`

**Should work now!** ✅

---

## 🚀 Alternative: Run Frontend Locally (Faster for Testing)

If Docker is being problematic, run the frontend locally:

```bash
# Stop Docker frontend
docker compose stop frontend

# Navigate to frontend directory
cd /Users/alex.macdonald/simple-ipass/frontend

# Create environment file
cat > .env.local << 'EOF'
NEXT_PUBLIC_API_URL=http://localhost:8080/api
EOF

# Install dependencies (if needed)
npm install

# Run development server
npm run dev
```

**Open**: http://localhost:3000

**Advantages**:
- Changes reflect instantly
- No Docker caching issues
- Faster development iteration

---

## 🔍 Debugging: Check What URL Frontend is Using

### **Option 1: Check Browser Network Tab**

1. Open http://localhost:3000
2. Open DevTools (F12 or Cmd+Option+I)
3. Go to **Network** tab
4. Try to register
5. Look for the failed request
6. Check the **Request URL**

**Should be**: `http://localhost:8080/api/auth/register`
**If different**: Frontend has wrong API URL

### **Option 2: Check Frontend Build Logs**

```bash
# Check if NEXT_PUBLIC_API_URL was set during build
docker compose logs frontend 2>&1 | grep -i "API"
```

---

## 🐛 Still Not Working? Try These

### **1. Nuclear Option: Rebuild Everything**

```bash
cd /Users/alex.macdonald/simple-ipass

# Stop everything
docker compose down

# Remove all containers and volumes
docker compose down -v

# Rebuild everything from scratch
docker compose build --no-cache --pull

# Start everything
docker compose up -d

# Wait 2 minutes
sleep 120

# Check status
docker compose ps
```

### **2. Check Browser Console for Exact Error**

1. Open http://localhost:3000
2. Open DevTools → **Console** tab
3. Try to register
4. Copy the exact error message

Common errors:
- `net::ERR_CONNECTION_REFUSED` → Backend not running
- `CORS error` → Backend not allowing frontend origin
- `404 Not Found` → Wrong API URL
- `Failed to fetch` → Network issue or wrong URL

### **3. Verify Backend from Browser**

Open this URL in your browser:
```
http://localhost:8080/health
```

**Should see**: `{"status":"healthy","version":"0.2.0"}`

If this doesn't work, backend is not accessible.

### **4. Check if Frontend Can Reach Backend**

```bash
# Exec into frontend container
docker exec -it ipaas-frontend sh

# Try to reach backend from inside container
wget -O- http://host.docker.internal:8080/health

# Should return: {"status":"healthy"...}

# Exit container
exit
```

If this fails, there's a Docker networking issue.

---

## 📋 Complete Rebuild Script

Copy and paste this entire script:

```bash
#!/bin/bash
cd /Users/alex.macdonald/simple-ipass

echo "🔧 Stopping frontend..."
docker compose stop frontend

echo "🗑️  Removing frontend container..."
docker compose rm -f frontend

echo "🏗️  Rebuilding frontend (this takes 2-3 minutes)..."
docker compose build --no-cache --pull frontend

echo "🚀 Starting frontend..."
docker compose up -d frontend

echo "⏳ Waiting 30 seconds for frontend to start..."
sleep 30

echo "✅ Checking status..."
docker compose ps frontend

echo ""
echo "🌐 Frontend should be ready at: http://localhost:3000"
echo ""
echo "📝 Next steps:"
echo "1. Clear your browser cache (Cmd+Shift+Delete)"
echo "2. Open http://localhost:3000"
echo "3. Try to register with a new email"
echo ""
echo "If it still doesn't work, check browser console (F12) for errors."
```

Save as `fix_frontend.sh`, then run:

```bash
chmod +x fix_frontend.sh
./fix_frontend.sh
```

---

## 🎯 What We Know

✅ **Backend is working** - Tested and responding correctly  
✅ **CORS is configured** - Allows http://localhost:3000  
✅ **API endpoint exists** - `/api/auth/register` works  
❌ **Frontend can't connect** - Something wrong with frontend config

**Most likely cause**: Frontend was built with wrong `NEXT_PUBLIC_API_URL` or browser cache.

---

## 📞 Need More Help?

If none of this works, please provide:

1. **Browser console error** (F12 → Console tab → screenshot or text)
2. **Network tab details** (F12 → Network tab → click failed request → screenshot)
3. **Frontend logs**: `docker compose logs frontend --tail=50`
4. **Docker status**: `docker compose ps`

This will help me diagnose the exact issue!

---

## ✅ Expected Working Flow

Once fixed, this should happen:

1. Open http://localhost:3000
2. Click "Register"
3. Enter email + password
4. Click "Register" button
5. ✅ **Success!** → Redirected to login
6. Login with same credentials
7. ✅ **Success!** → Redirected to dashboard

---

**Try the comprehensive rebuild steps above and let me know if the issue persists!** 🚀

