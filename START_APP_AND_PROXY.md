# 🚀 Start GoFlow with Kong Gateway Proxy

## Quick Start - Full Stack

### Prerequisites
✅ Docker Desktop must be running

---

## Method 1: Docker Compose (Recommended - All Services)

Starts everything together:
- ✅ Backend API (Go)
- ✅ Frontend (Next.js)
- ✅ Kong Gateway (API Proxy)
- ✅ PostgreSQL (Database)
- ✅ Elasticsearch + Logstash + Kibana (Monitoring)

```bash
cd /Users/alex.macdonald/simple-ipass

# Start all services
docker compose up --build

# Or in detached mode (background):
docker compose up -d --build
```

**Wait 2-3 minutes for all services to start.**

### Check Status
```bash
# View running containers
docker compose ps

# View logs
docker compose logs -f

# View specific service logs
docker compose logs -f backend
docker compose logs -f kong
docker compose logs -f frontend
```

---

## Method 2: Individual Services (For Development)

If you want more control, start services individually:

### Step 1: Start Supporting Services
```bash
cd /Users/alex.macdonald/simple-ipass

# Start PostgreSQL, Elasticsearch, Logstash, Kibana only
docker compose up -d postgres elasticsearch logstash kibana
```

### Step 2: Start Backend (Local)
```bash
cd /Users/alex.macdonald/simple-ipass

# Option A: Docker
docker compose up -d backend

# Option B: Local Go (faster for development)
go run cmd/api/main.go
```

### Step 3: Start Kong Gateway
```bash
cd /Users/alex.macdonald/simple-ipass

docker compose up -d kong
```

### Step 4: Start Frontend
```bash
cd /Users/alex.macdonald/simple-ipass

# Option A: Docker
docker compose up -d frontend

# Option B: Local Next.js (with Dev Mode enabled)
./scripts/run_frontend_locally.sh
```

---

## Service URLs

After starting all services:

| Service | URL | Description |
|---------|-----|-------------|
| **Frontend** | http://localhost:3000 | GoFlow Dashboard |
| **Backend API** | http://localhost:8080 | Direct API access |
| **Kong Gateway** | http://localhost:8000 | Proxied API (recommended) |
| **Kong Admin** | http://localhost:8001 | Kong configuration |
| **Kibana** | http://localhost:5601 | Logs & Monitoring |
| **Elasticsearch** | http://localhost:9200 | Search & Analytics |

---

## Testing the Setup

### 1. Check Backend Health
```bash
# Direct API
curl http://localhost:8080/health

# Via Kong Gateway
curl http://localhost:8000/health
```

### 2. Check Kong Gateway
```bash
# Kong health
curl http://localhost:8001/status

# List Kong services
curl http://localhost:8001/services

# List Kong routes
curl http://localhost:8001/routes
```

### 3. Access Frontend
Open browser: http://localhost:3000

Click "Skip Login - Dev Mode" to quickly test!

---

## Configure Kong (First Time Setup)

After services start, configure Kong Gateway patterns:

```bash
cd /Users/alex.macdonald/simple-ipass

# Configure Kong with integration patterns
./scripts/configure_kong_elk.sh
```

This sets up:
- ✅ Rate limiting
- ✅ Request logging to ELK
- ✅ CORS headers
- ✅ API key authentication (optional)

---

## Stopping Services

### Stop All Services
```bash
cd /Users/alex.macdonald/simple-ipass

# Stop and keep data
docker compose stop

# Stop and remove containers (keeps volumes)
docker compose down

# Stop and remove everything (including data)
docker compose down -v
```

### Stop Individual Service
```bash
docker compose stop backend
docker compose stop kong
docker compose stop frontend
```

---

## Troubleshooting

### Docker Not Running
**Error:** `Cannot connect to the Docker daemon`

**Solution:**
1. Open Docker Desktop application
2. Wait for it to fully start
3. Retry the command

### Port Already in Use
**Error:** `port is already allocated`

**Solution:**
```bash
# Find process using the port (example: port 8080)
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use a different port by editing docker-compose.yml
```

### Frontend Can't Connect to Backend
**Error:** `Failed to fetch` or `Network error`

**Solution:**
1. Check backend is running: `curl http://localhost:8080/health`
2. Check Kong is running: `curl http://localhost:8000/health`
3. Rebuild frontend: `docker compose up -d --build frontend`
4. Or use local frontend: `./scripts/run_frontend_locally.sh`

### Database Connection Failed
**Error:** `failed to ping database`

**Solution:**
```bash
# Restart PostgreSQL
docker compose restart postgres

# Check PostgreSQL logs
docker compose logs postgres

# Verify it's running
docker compose ps postgres
```

### Kong Gateway Not Working
**Error:** `Kong Admin API unreachable`

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

## Quick Commands Reference

```bash
# Start everything
docker compose up -d --build

# View logs
docker compose logs -f

# Check status
docker compose ps

# Restart a service
docker compose restart backend

# Stop everything
docker compose down

# Clean rebuild
docker compose down -v && docker compose up --build

# Configure Kong
./scripts/configure_kong_elk.sh

# Run tests
make test-full
```

---

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         Browser                              │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
         ┌─────────────────────────┐
         │   Frontend (Next.js)    │
         │   http://localhost:3000 │
         └─────────────┬───────────┘
                       │
                       ▼
         ┌─────────────────────────┐
         │  Kong Gateway (Proxy)   │◄──── Rate Limiting
         │   http://localhost:8000 │◄──── Auth Overlay
         └─────────────┬───────────┘◄──── Logging to ELK
                       │
                       ▼
         ┌─────────────────────────┐
         │    Backend API (Go)     │
         │   http://localhost:8080 │
         └─────────────┬───────────┘
                       │
         ┌─────────────┴───────────┐
         │                         │
         ▼                         ▼
┌──────────────┐         ┌──────────────────┐
│  PostgreSQL  │         │  ELK Stack       │
│  (Database)  │         │  (Monitoring)    │
└──────────────┘         └──────────────────┘
```

---

## Development Workflow

### For Backend Development
```bash
# Run backend locally (faster iteration)
go run cmd/api/main.go

# Keep other services in Docker
docker compose up -d postgres kong elasticsearch logstash kibana
```

### For Frontend Development
```bash
# Run frontend locally (with hot reload)
./scripts/run_frontend_locally.sh

# Keep other services in Docker
docker compose up -d backend postgres kong
```

### For Kong Configuration
```bash
# Edit Kong config
./scripts/configure_kong_elk.sh

# Test Kong patterns
make test-kong
```

---

## What's Next?

After starting:

1. ✅ **Access Frontend**: http://localhost:3000
2. ✅ **Use Dev Mode**: Click "Skip Login - Dev Mode"
3. ✅ **Create Connections**: Configure your API keys
4. ✅ **Build Workflows**: Create your first integration
5. ✅ **Monitor**: View logs in Kibana (http://localhost:5601)

---

## Advanced Options

### Clean Slate (Reset Everything)
```bash
# Stop and remove all containers, volumes, and images
docker compose down -v --rmi all

# Rebuild from scratch
docker compose up --build
```

### View Real-Time Logs
```bash
# All services
docker compose logs -f

# Just backend and Kong
docker compose logs -f backend kong

# Just errors
docker compose logs -f | grep -i error
```

### Scale Services
```bash
# Run 3 backend instances
docker compose up -d --scale backend=3

# Let Kong load balance across them
```

---

**🎉 Your integration platform with enterprise-grade API Gateway is ready to go!**

Need help? Check:
- `README.md` - Project overview
- `COMPONENT_RUNNING_GUIDE.md` - Detailed component guide
- `TESTING_VALIDATION.md` - Testing procedures
- `QUICKSTART.md` - Quick start guide
