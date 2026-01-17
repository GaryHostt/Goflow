# 🎯 Ready to Start GoFlow Platform!

## ⚠️ Important: Docker Must Be Running

Before starting, ensure **Docker Desktop is running** on your Mac.

---

## 🚀 Three Ways to Start

### Option 1: Quick Start Script (Recommended)
```bash
cd /Users/alex.macdonald/simple-ipass
./scripts/start_platform.sh
```

✅ Checks Docker is running  
✅ Starts all services  
✅ Shows status and URLs  
✅ Most user-friendly!

---

### Option 2: Make Command
```bash
cd /Users/alex.macdonald/simple-ipass

# Start all services
make docker-up-build

# Or without rebuild (faster)
make docker-up

# Or use the convenience command
make start-platform
```

---

### Option 3: Direct Docker Compose
```bash
cd /Users/alex.macdonald/simple-ipass

# Start all services
docker compose up -d --build

# View logs
docker compose logs -f
```

---

## 📦 What Gets Started

When you run the platform, these services start:

| Service | Port | Description |
|---------|------|-------------|
| **Frontend** | 3000 | Next.js web interface |
| **Backend** | 8080 | Go API server |
| **Kong Gateway** | 8000 | API proxy with rate limiting |
| **Kong Admin** | 8001 | Kong configuration API |
| **PostgreSQL** | 5432 | Production database |
| **Elasticsearch** | 9200 | Log storage & search |
| **Logstash** | 5000 | Log processing |
| **Kibana** | 5601 | Log visualization |

---

## ✅ After Starting

### 1. Wait 2-3 Minutes
Services need time to initialize (especially on first run).

### 2. Check Status
```bash
docker compose ps
```

All services should show "running" status.

### 3. Access the Platform
Open browser: **http://localhost:3000**

### 4. Quick Login
Click **"Skip Login - Dev Mode"** for instant access!

---

## 🧪 Verify Everything Works

### Test 1: Health Check
```bash
# Backend health
curl http://localhost:8080/health

# Kong health  
curl http://localhost:8000/health
```

### Test 2: Access Frontend
Visit: http://localhost:3000

Should see GoFlow landing page.

### Test 3: Check Kong Admin
```bash
curl http://localhost:8001/status
```

Should return Kong status.

### Test 4: View Logs in Kibana
Visit: http://localhost:5601

(May take 3-4 minutes for Elasticsearch to be ready)

---

## 🔧 Configure Kong (Optional)

After services are running, configure Kong Gateway patterns:

```bash
cd /Users/alex.macdonald/simple-ipass
./scripts/configure_kong_elk.sh
```

Or use the convenience command:
```bash
make start-platform-kong
```

This sets up:
- ✅ Rate limiting (100 req/min)
- ✅ Request logging to ELK
- ✅ CORS headers
- ✅ Request ID tracking

---

## 🛑 Stopping Services

### Stop All Services
```bash
docker compose stop
# or
make docker-down
```

### Stop and Remove Everything
```bash
docker compose down -v
```

### Stop One Service
```bash
docker compose stop frontend
```

---

## 📊 Viewing Logs

### All Services
```bash
docker compose logs -f
```

### Specific Service
```bash
docker compose logs -f backend
docker compose logs -f kong
docker compose logs -f frontend
```

### Just Errors
```bash
docker compose logs -f | grep -i error
```

---

## 🐛 Troubleshooting

### Problem: Docker Not Running
**Error:** `Cannot connect to the Docker daemon`

**Solution:**
1. Open **Docker Desktop** application
2. Wait for it to show "Docker Desktop is running"
3. Try again

---

### Problem: Port Already in Use
**Error:** `port is already allocated`

**Solution:**
```bash
# Find what's using the port
lsof -i :8080  # or :3000, :8000, etc.

# Kill the process
kill -9 <PID>
```

---

### Problem: Service Won't Start
**Error:** Container keeps restarting

**Solution:**
```bash
# Check logs for the service
docker compose logs <service-name>

# Try rebuilding
docker compose down
docker compose up --build
```

---

### Problem: Frontend Can't Connect
**Error:** "Failed to fetch" in browser console

**Solution:**
```bash
# 1. Verify backend is running
curl http://localhost:8080/health

# 2. Rebuild frontend
docker compose up -d --build frontend

# 3. Or run frontend locally
./scripts/run_frontend_locally.sh
```

---

### Problem: Kong Not Working
**Error:** Kong Admin API unreachable

**Solution:**
```bash
# Restart Kong
docker compose restart kong

# Check Kong logs
docker compose logs kong

# Verify Kong is healthy
curl http://localhost:8001/status
```

---

### Problem: Database Issues
**Error:** "failed to ping database"

**Solution:**
```bash
# Restart PostgreSQL
docker compose restart postgres

# Check if it's running
docker compose ps postgres

# View logs
docker compose logs postgres
```

---

## 🧹 Clean Slate (If Everything Breaks)

Start completely fresh:

```bash
cd /Users/alex.macdonald/simple-ipass

# Stop and remove everything
docker compose down -v --rmi all

# Rebuild from scratch
docker compose up --build

# Wait 2-3 minutes...
```

---

## 🎯 Quick Reference

```bash
# START
./scripts/start_platform.sh

# STOP
docker compose down

# RESTART ONE SERVICE
docker compose restart backend

# VIEW LOGS
docker compose logs -f

# STATUS
docker compose ps

# HEALTH CHECKS
curl http://localhost:8080/health  # Backend
curl http://localhost:8000/health  # Kong
curl http://localhost:8001/status  # Kong Admin

# CLEAN START
docker compose down -v && docker compose up --build
```

---

## 📚 More Documentation

- **START_APP_AND_PROXY.md** - Complete startup guide
- **COMPONENT_RUNNING_GUIDE.md** - Run individual components
- **TESTING_VALIDATION.md** - Testing procedures
- **README.md** - Project overview
- **QUICKSTART.md** - Quick start guide

---

## 🎉 You're Ready!

Once Docker Desktop is running:

1. Run: `./scripts/start_platform.sh`
2. Wait 2-3 minutes
3. Open: http://localhost:3000
4. Click: "Skip Login - Dev Mode"
5. Start building integrations!

---

**Need Help?** Check the docs above or run:
```bash
make help
```
