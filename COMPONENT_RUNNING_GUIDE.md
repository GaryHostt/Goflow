# 🚀 GoFlow iPaaS - Component Running Guide

**Complete guide to running each component of the GoFlow platform**

---

## 🎯 Quick Start (Recommended)

### **Option 1: Run Everything with Docker** (Easiest!)

```bash
# Navigate to project directory
cd /Users/alex.macdonald/simple-ipass

# Start all services
docker compose up -d

# Wait 30 seconds for services to initialize, then open:
# Frontend: http://localhost:3000
# Backend: http://localhost:8080
# Kong: http://localhost:8002
# Kibana: http://localhost:5601
```

**That's it!** All 8 services will start automatically.

---

### **Option 2: Run Components Individually** (For Development)

If you want to run services outside Docker for faster development:

```bash
# 1. Start dependencies (Postgres, Elasticsearch, Kong) in Docker
docker compose up -d postgres elasticsearch kibana kong kong-database

# 2. Run backend locally
cd /Users/alex.macdonald/simple-ipass
go run cmd/api/main.go

# 3. Run frontend locally (in a new terminal)
cd /Users/alex.macdonald/simple-ipass/frontend
npm run dev
```

---

## 📋 Component-by-Component Guide

### **1. PostgreSQL Database**

#### **Start with Docker** (Recommended)
```bash
docker compose up -d postgres
```

#### **Verify**
```bash
docker compose ps postgres
docker compose logs postgres
```

#### **Connect**
```bash
# Connection details:
Host: localhost
Port: 5432
Database: ipaas
User: ipaas_user
Password: ipaas_password

# Test connection:
docker exec -it ipaas-postgres psql -U ipaas_user -d ipaas
```

---

### **2. Go Backend API**

#### **Start with Docker** (Recommended)
```bash
docker compose up -d backend
```

#### **Run Locally** (For Development)
```bash
cd /Users/alex.macdonald/simple-ipass

# Install dependencies
go mod download

# Run the server
go run cmd/api/main.go
```

#### **Verify**
```bash
# Health check
curl http://localhost:8080/health

# Expected response:
# {"status":"healthy","version":"0.2.0"}
```

#### **View Logs**
```bash
# Docker logs
docker compose logs backend -f

# Or if running locally, logs appear in terminal
```

#### **Environment Variables** (if running locally)
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=ipaas
export DB_USER=ipaas_user
export DB_PASSWORD=ipaas_password
export JWT_SECRET=your-secret-key-change-in-production
```

---

### **3. Next.js Frontend**

#### **Start with Docker** (Recommended)
```bash
docker compose up -d frontend
```

#### **Run Locally** (For Development)
```bash
cd /Users/alex.macdonald/simple-ipass/frontend

# First time only: Install dependencies
npm install

# Run development server
npm run dev
```

#### **Verify**
```bash
# Open in browser:
open http://localhost:3000

# Or curl:
curl http://localhost:3000
```

#### **Environment Variables**

Create `frontend/.env.local` (for local development):
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

**Note**: If running frontend in Docker, the API URL is already configured in `docker-compose.yml`.

---

### **4. Elasticsearch**

#### **Start with Docker** (Recommended)
```bash
docker compose up -d elasticsearch
```

#### **Verify**
```bash
# Health check
curl http://localhost:9200/_cluster/health

# Expected response:
# {"cluster_name":"docker-cluster","status":"green",...}
```

#### **View Indexes**
```bash
curl http://localhost:9200/_cat/indices?v
```

---

### **5. Kibana**

#### **Start with Docker** (Recommended)
```bash
docker compose up -d kibana
```

#### **Access**
```bash
# Open in browser:
open http://localhost:5601

# Wait 1-2 minutes for Kibana to fully start
```

#### **Verify**
```bash
# Check status
curl http://localhost:5601/api/status
```

---

### **6. Kong Gateway**

#### **Start with Docker** (Recommended)
```bash
# Kong requires its database and migration first:
docker compose up -d kong-database
docker compose up -d kong-migration
docker compose up -d kong
```

#### **Access Points**
- **Admin API**: http://localhost:8001
- **Proxy API**: http://localhost:8000  
- **Kong Manager**: http://localhost:8002

#### **Verify**
```bash
# Check Kong status
curl http://localhost:8001/status

# List services
curl http://localhost:8001/services

# List routes
curl http://localhost:8001/routes
```

---

### **7. Logstash** (Optional - for Kong logs)

#### **Start with Docker**
```bash
docker compose up -d logstash
```

#### **Verify**
```bash
docker compose logs logstash

# Should see: "Successfully started Logstash API endpoint"
```

---

## 🔍 Troubleshooting Guide

### **Problem: "Failed to fetch" error on registration**

**Cause**: Frontend can't connect to backend.

**Solution 1: Check services are running**
```bash
docker compose ps

# All services should show "Up" and "healthy"
```

**Solution 2: Rebuild frontend with correct API URL**
```bash
docker compose down frontend
docker compose build --no-cache frontend
docker compose up -d frontend
```

**Solution 3: Check backend is accessible from browser**
```bash
# Open in browser (should return JSON):
open http://localhost:8080/health

# If this fails, backend is not running correctly
docker compose logs backend
```

**Solution 4: Clear browser cache and try again**
```bash
# In browser: Clear cache and hard reload (Cmd+Shift+R on Mac)
```

---

### **Problem: Backend won't start**

**Symptoms**: `docker compose logs backend` shows errors

**Solution 1: Check database is ready**
```bash
docker compose ps postgres

# Should show "healthy"
```

**Solution 2: Restart with fresh database**
```bash
docker compose down -v  # ⚠️ This deletes all data!
docker compose up -d
```

**Solution 3: Check for port conflicts**
```bash
# Check if something else is using port 8080
lsof -i :8080

# Kill the process if needed
kill -9 <PID>
```

---

### **Problem: Frontend won't load**

**Solution 1: Check logs**
```bash
docker compose logs frontend

# Look for build errors or startup errors
```

**Solution 2: Rebuild**
```bash
docker compose down frontend
docker compose build --no-cache frontend
docker compose up -d frontend
```

**Solution 3: Run locally**
```bash
cd frontend
rm -rf .next node_modules
npm install
npm run dev
```

---

### **Problem: Kong returns 404 or 503**

**Solution 1: Check Kong database migration**
```bash
docker compose logs kong-migration

# Should see: "Database is up-to-date"
```

**Solution 2: Restart Kong**
```bash
docker compose restart kong
```

**Solution 3: Recreate Kong services**
```bash
docker compose down kong kong-database kong-migration
docker compose up -d kong-database
# Wait 10 seconds
docker compose up -d kong-migration
# Wait for migration to complete
docker compose up -d kong
```

---

### **Problem: Elasticsearch not starting**

**Cause**: Not enough memory

**Solution 1: Increase Docker memory**
```bash
# Docker Desktop → Settings → Resources
# Set memory to at least 4GB
```

**Solution 2: Reduce Elasticsearch memory**
Edit `docker-compose.yml` line 76:
```yaml
- "ES_JAVA_OPTS=-Xms256m -Xmx256m"  # Reduced from 512m
```

Then restart:
```bash
docker compose restart elasticsearch
```

---

## 🧪 Testing Each Component

### **Backend API**
```bash
# Health check
curl http://localhost:8080/health

# Register a user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### **Frontend**
```bash
# Open in browser
open http://localhost:3000

# Should see GoFlow logo and login page
```

### **Elasticsearch**
```bash
# Cluster health
curl http://localhost:9200/_cluster/health?pretty

# List all indexes
curl http://localhost:9200/_cat/indices?v
```

### **Kong Gateway**
```bash
# Kong status
curl http://localhost:8001/status

# Create a test service
curl -X POST http://localhost:8001/services \
  -d "name=test-service" \
  -d "url=http://backend:8080"

# List services
curl http://localhost:8001/services
```

---

## 📊 Monitoring

### **View All Service Status**
```bash
docker compose ps
```

### **View Logs for All Services**
```bash
docker compose logs -f
```

### **View Logs for Specific Service**
```bash
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f kong
```

### **Check Resource Usage**
```bash
docker stats
```

---

## 🛑 Stopping Services

### **Stop All Services**
```bash
docker compose down
```

### **Stop and Remove All Data** (⚠️ Destructive!)
```bash
docker compose down -v
```

### **Stop Specific Service**
```bash
docker compose stop backend
docker compose stop frontend
```

### **Restart Specific Service**
```bash
docker compose restart backend
```

---

## 🔄 Updating Components

### **Rebuild Backend**
```bash
docker compose build backend
docker compose up -d backend
```

### **Rebuild Frontend**
```bash
docker compose build frontend
docker compose up -d frontend
```

### **Rebuild Everything**
```bash
docker compose build --no-cache
docker compose up -d
```

---

## 🌐 Access URLs

| Component | URL | Purpose |
|-----------|-----|---------|
| **Frontend** | http://localhost:3000 | Main UI |
| **Backend API** | http://localhost:8080 | REST API |
| **Backend Health** | http://localhost:8080/health | Health check |
| **Kong Proxy** | http://localhost:8000 | API Gateway |
| **Kong Admin** | http://localhost:8001 | Admin API |
| **Kong Manager** | http://localhost:8002 | Web UI |
| **Elasticsearch** | http://localhost:9200 | Search/Logs |
| **Kibana** | http://localhost:5601 | Log viewer |
| **PostgreSQL** | localhost:5432 | Database |

---

## 📝 Development Workflow

### **Typical Development Session**

1. **Start infrastructure**
   ```bash
   docker compose up -d postgres elasticsearch kibana kong kong-database
   ```

2. **Run backend locally** (for faster iteration)
   ```bash
   cd /Users/alex.macdonald/simple-ipass
   go run cmd/api/main.go
   ```

3. **Run frontend locally** (in new terminal)
   ```bash
   cd /Users/alex.macdonald/simple-ipass/frontend
   npm run dev
   ```

4. **Make changes and test**

5. **When done, stop everything**
   ```bash
   # Ctrl+C to stop go and npm
   docker compose down
   ```

---

## 🚀 Production Deployment

**For production, additional steps are required:**

1. **Environment Variables**
   - Set `JWT_SECRET` to a strong random value
   - Configure production database credentials
   - Set `NEXT_PUBLIC_API_URL` to your domain

2. **SSL/TLS**
   - Add certificates to Kong
   - Configure HTTPS for all services

3. **Scaling**
   - Increase worker pool size
   - Configure database connection pooling
   - Set up load balancing

4. **Monitoring**
   - Configure Kibana alerts
   - Set up uptime monitoring
   - Add error tracking (Sentry, etc.)

See `README.md` for complete production checklist.

---

## ✅ Quick Reference

**Start everything:**
```bash
docker compose up -d
```

**View status:**
```bash
docker compose ps
```

**View logs:**
```bash
docker compose logs -f [service]
```

**Restart service:**
```bash
docker compose restart [service]
```

**Stop everything:**
```bash
docker compose down
```

**Test backend:**
```bash
curl http://localhost:8080/health
```

**Open frontend:**
```bash
open http://localhost:3000
```

---

**Need help?** Check the troubleshooting section above or view detailed logs with `docker compose logs [service]`.

