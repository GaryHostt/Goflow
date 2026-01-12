# 🧪 How to Run the GoFlow Tests

Since Docker is not available in the Cursor sandbox, you'll need to run the tests manually in your terminal. Here's exactly how to do it:

---

## 🚀 Complete Test Execution (Recommended)

### Option 1: Automated Test Runner (Easiest)

Open a terminal and run:

```bash
cd /Users/alex.macdonald/simple-ipass
./scripts/run_all_tests.sh
```

This script will:
1. ✅ Check if Docker is running
2. ✅ Start all platform services
3. ✅ Wait for services to be healthy
4. ✅ Configure Kong ELK integration
5. ✅ Run connector tests (18 connectors)
6. ✅ Run Kong Gateway tests (5 patterns)
7. ✅ Display summary

**Duration**: ~3-4 minutes total

---

### Option 2: Step-by-Step Manual Testing

#### Step 1: Start the Platform

```bash
cd /Users/alex.macdonald/simple-ipass

# Start all services
docker compose up -d

# Wait for services to be healthy (60 seconds)
sleep 60

# Check service status
docker compose ps
```

**Expected Output**:
```
NAME                    STATUS
ipaas-backend           Up (healthy)
ipaas-elasticsearch     Up (healthy)
ipaas-frontend          Up (healthy)
ipaas-kibana           Up
ipaas-logstash         Up
ipaas-postgres         Up (healthy)
kong-database          Up (healthy)
kong-gateway           Up (healthy)
```

---

#### Step 2: Configure Kong ELK Integration

```bash
./scripts/configure_kong_elk.sh
```

**Expected Output**:
```
🔧 Configuring Kong to ship logs to ELK...
⏳ Waiting for Kong Admin API...
✅ Kong is ready
📤 Installing http-log plugin for ELK integration...
✅ Kong http-log plugin installed
📋 Installing request-transformer for enhanced tracking...
✅ Kong log shipping configured!

📊 Kong logs will now appear in:
   - Elasticsearch: http://localhost:9200/kong-logs-*
   - Kibana: http://localhost:5601
```

---

#### Step 3: Run Connector Tests

```bash
go run scripts/connector_test.go
```

**Expected Output**:
```
🚀 GoFlow Connector Test Suite
================================

📋 Running Connector Tests...
  ✅ Slack (Webhook): 45ms
  ✅ Discord (Webhook): 38ms
  ⚠️  Twilio (SMS): SKIPPED (API key required)
  ✅ OpenWeather API: 156ms
  ⚠️  NewsAPI: SKIPPED (API key required)
  ✅ The Cat API: 234ms
  ✅ Fake Store API: 189ms
  ✅ SOAP Connector: 12ms
  ✅ SWAPI (Star Wars API): 267ms
  ✅ Salesforce: 8ms
  ✅ PokeAPI: 178ms
  ✅ Bored API: 134ms
  ✅ Numbers API: 98ms
  ✅ NASA API: 456ms
  ✅ REST Countries: 245ms
  ✅ Dog CEO API: 167ms

📊 Test Summary
================================
Total Tests: 16
✅ Passed: 12
❌ Failed: 0
⚠️  Skipped: 4 (require API keys)

🎉 All critical tests passed!

📤 Shipping test results to ELK...
✅ Test results shipped to ELK!
   View in Kibana: http://localhost:5601
```

**Duration**: ~30 seconds

---

#### Step 4: Run Kong Gateway Tests

```bash
go run scripts/kong_test.go
```

**Expected Output**:
```
🚀 Kong Gateway Integration Test Suite
========================================

⏳ Waiting for Kong to be ready...
✅ Kong Gateway is ready

📋 Test 1: Protocol Bridge (SOAP to REST)
  ✅ Created Kong service
  ✅ Created Kong route
  ✅ Added request-transformer plugin
  ✅ Protocol Bridge validated (Status: 200)

📋 Test 2: Webhook Rate Limiting
  ✅ Rate limiting configured (10 req/min)
  ✅ Rate limiting validated (3/3 requests passed)

📋 Test 3: Smart API Aggregator
  ✅ API aggregator with caching configured

📋 Test 4: Federated Security (Auth Overlay)
  ✅ Key-based authentication configured
  ✅ Auth protection validated (401 without key)

📋 Test 5: Usage-Based Tracking
  ✅ Usage tracking headers configured
  ℹ️  View logs in ELK for full tracking data

🧹 Cleaning up test resources...
✅ Cleanup complete

🎉 All Kong Gateway tests passed!
```

**Duration**: ~45 seconds

---

#### Step 5: View Results in Kibana

Open your browser:
```
http://localhost:5601
```

1. Navigate to **Stack Management > Data Views**
2. Create data view: `connector-tests-*`
3. Create data view: `kong-logs-*`
4. Navigate to **Analytics > Discover**
5. Select your data view and explore!

---

## 🎯 Using Makefile Commands

You can also use the Makefile commands:

```bash
# Test all connectors
make test-connectors

# Test Kong Gateway
make test-kong

# Run full test suite
make test-full

# Configure Kong ELK
make configure-kong-elk
```

---

## 🐛 Troubleshooting

### Issue: "Docker is not running"
**Solution**: 
```bash
# Start Docker Desktop
open -a Docker

# Wait 30 seconds for Docker to start
sleep 30

# Try again
docker compose up -d
```

---

### Issue: "Connection refused" in tests
**Solution**:
```bash
# Check if services are healthy
docker compose ps

# View logs for any failing service
docker compose logs backend
docker compose logs kong

# Restart services
docker compose restart
```

---

### Issue: "Kong Admin API not accessible"
**Solution**:
```bash
# Check Kong logs
docker compose logs kong

# Verify Kong is healthy
curl http://localhost:8001/status

# Restart Kong if needed
docker compose restart kong
sleep 10
```

---

### Issue: "Elasticsearch not available"
**Solution**:
```bash
# Check Elasticsearch
curl http://localhost:9200/_cluster/health

# View logs
docker compose logs elasticsearch

# Restart if needed
docker compose restart elasticsearch
sleep 30
```

---

### Issue: "go not found"
**Solution**:
```bash
# Check Go installation
go version

# If not installed, install Go 1.21+
brew install go  # macOS

# Verify installation
go version
```

---

## 📊 Expected Test Results

### Connector Tests
- **Total**: 16 connectors
- **Expected Pass**: 12 (75%)
- **Expected Skip**: 4 (Twilio, NewsAPI, OpenWeather, Salesforce - require API keys)
- **Duration**: 30 seconds

### Kong Gateway Tests
- **Total**: 5 patterns
- **Expected Pass**: 5 (100%)
- **Duration**: 45 seconds

### Overall
- **Total Tests**: 22
- **Expected Pass**: 18 (82%)
- **Total Duration**: ~90 seconds

---

## ✅ Success Indicators

You'll know the tests passed successfully if you see:

1. **Connector Tests**:
   - ✅ "All critical tests passed!"
   - ✅ "Test results shipped to ELK!"
   - ✅ Pass rate ≥ 75%

2. **Kong Gateway Tests**:
   - ✅ "All Kong Gateway tests passed!"
   - ✅ "Cleanup complete"
   - ✅ All 5 patterns validated

3. **Kibana**:
   - ✅ Data appears in `connector-tests-*` index
   - ✅ Data appears in `kong-logs-*` index
   - ✅ Dashboards can be created

---

## 🚀 Quick Start (Copy-Paste)

Open a terminal and run these commands one by one:

```bash
# Navigate to project
cd /Users/alex.macdonald/simple-ipass

# Start platform
docker compose up -d

# Wait for services
sleep 60

# Check status
docker compose ps

# Configure Kong
./scripts/configure_kong_elk.sh

# Run connector tests
go run scripts/connector_test.go

# Run Kong tests
go run scripts/kong_test.go

# Open Kibana
open http://localhost:5601
```

---

## 📚 Next Steps After Tests Pass

1. **Create Kibana Dashboards**:
   - Connector performance metrics
   - Kong traffic patterns
   - Error analysis

2. **Set Up Alerts**:
   - Failed connector tests
   - High error rates
   - Slow response times

3. **Integrate with CI/CD**:
   - Add to GitHub Actions
   - Run on every PR
   - Block merge on failures

---

**You're all set!** Run the tests and watch your platform validate itself! 🚀

