# ✅ GoFlow Testing Suite - Quick Start

## **Run All Tests (Automated)**

```bash
cd /Users/alex.macdonald/simple-ipass
./scripts/run_all_tests.sh
```

This script:
1. ✅ Checks if Docker is running
2. ✅ Starts all services (`docker compose up -d`)
3. ✅ Waits for services to be healthy
4. ✅ Configures Kong ELK integration
5. ✅ Runs connector validation tests
6. ✅ Runs Kong Gateway validation tests

---

## **Run Individual Tests**

### **1. Connector Validation**
```bash
cd /Users/alex.macdonald/simple-ipass
go run scripts/validate_connectors.go
```

**Tests:**
- ✅ PokeAPI
- ✅ Bored API
- ✅ Numbers API
- ✅ Dog CEO API
- ✅ REST Countries
- ✅ SWAPI (Star Wars)
- ✅ Cat API
- ✅ Fake Store API
- ✅ NASA API (with DEMO_KEY)
- ⚠️  Slack, Discord, Twilio (require webhook URLs)
- ⚠️  OpenWeather, NewsAPI (require API keys)
- ⚠️  Salesforce (requires OAuth)
- ✅ SOAP Connector (structure validated)

---

### **2. Kong Gateway Validation**
```bash
cd /Users/alex.macdonald/simple-ipass
go run scripts/validate_kong.go
```

**Tests:**
- 🔄 Protocol Bridge (SOAP to REST)
- 🚦 Webhook Rate Limiting
- 🔀 Smart API Aggregator
- 🔐 Federated Security (Auth Overlay)
- 📊 Usage Tracking (ELK logging)

---

## **Manual Testing via UI**

### **Kong Manager**
http://localhost:8002

**Create integration patterns:**
1. Click "Services" → "New Service"
2. Configure routes and plugins
3. Test via Proxy: http://localhost:8000

---

### **API Management UI**
http://localhost:3000/dashboard/api-management

**Use templates:**
- Protocol Bridge (SOAP → REST)
- Webhook Handler (Rate limiting)
- Smart Aggregator (Multi-source)
- Auth Overlay (Security)
- Usage Tracker (Billing)

---

## **View Results**

### **Console Output**
```
✅ Passed: 9
❌ Failed: 0
⚠️  Skipped: 9 (require API keys)
```

### **Kibana Dashboards**
http://localhost:5601

**Create Index Pattern:**
1. Stack Management → Data Views
2. Create: `kong-logs-*`
3. View logs in Discover

---

## **Troubleshooting**

### **Issue: "Docker is not running"**
```bash
# Start Docker Desktop manually
open -a Docker
# Wait 30 seconds, then retry
./scripts/run_all_tests.sh
```

### **Issue: "go: cannot run *_test.go files"**
✅ **FIXED!** Files renamed:
- `connector_test.go` → `validate_connectors.go`
- `kong_test.go` → `validate_kong.go`

### **Issue: "Kong is not available"**
```bash
# Check Kong health
docker compose ps kong

# View Kong logs
docker compose logs kong

# Restart Kong
docker compose restart kong
```

### **Issue: "Service unhealthy"**
```bash
# Check all services
docker compose ps

# View specific logs
docker compose logs frontend
docker compose logs backend
docker compose logs logstash
```

---

## **What's Tested**

| Category | Count | Status |
|----------|-------|--------|
| **Public APIs** | 9 | ✅ Fully tested |
| **Auth APIs** | 6 | ⚠️  Require keys |
| **SOAP Connector** | 1 | ✅ Structure validated |
| **Kong Patterns** | 5 | ✅ Configuration checked |
| **ELK Integration** | 1 | ✅ Logs shipping |

---

## **Next Steps**

1. ✅ Add API keys to test authenticated connectors
2. ✅ Create Kong services via API Management UI
3. ✅ View real-time logs in Kibana
4. ✅ Build your first multi-step workflow!

---

**All tests passing?** 🎉 Your GoFlow iPaaS is production-ready!

