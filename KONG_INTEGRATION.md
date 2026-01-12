# Kong Gateway Integration - Enterprise API Management

## 🌟 Overview

GoFlow now includes **Kong Gateway** integration, transforming your iPaaS into an **enterprise-grade API management platform**. Kong acts as the "front door" to your workflows, providing security, rate limiting, caching, and protocol bridging.

---

## 🏗️ Architecture

```
External API Calls
        ↓
   Kong Gateway (8000)
   ├─ Authentication (API Keys, OAuth2)
   ├─ Rate Limiting (100 req/sec)
   ├─ Caching (5 min TTL)
   └─ Protocol Bridge (SOAP → REST)
        ↓
   GoFlow Backend (8080)
   ├─ Workflow Execution
   ├─ Database Operations
   └─ Third-party API Calls
        ↓
   External Services (Slack, Discord, etc.)
```

---

## 🚀 5 Enterprise Use Cases

### 1. **Protocol Bridge** (Modernizing Legacy Systems)

**Problem**: Your mobile app needs data from a 20-year-old SOAP service.

**Solution**: 
- Kong exposes a clean REST endpoint at `http://localhost:8000/api/legacy-data`
- GoFlow's SOAP connector converts REST → SOAP
- Legacy system returns XML
- GoFlow converts XML → JSON
- Kong returns clean JSON to your app

**Benefits**:
- ✅ Hide legacy system complexity
- ✅ Modern developers use REST/JSON only
- ✅ No changes to legacy system required

**Setup**:
```bash
# 1. Create a SOAP workflow in GoFlow
POST /api/workflows
{
  "name": "Legacy SOAP Bridge",
  "action_type": "soap_call",
  "config_json": {
    "soap_endpoint": "http://legacy-system.com/service.asmx",
    "soap_method": "GetCustomerData",
    "soap_namespace": "http://tempuri.org/",
    "soap_parameters": {
      "customerId": "{{customer_id}}"
    }
  }
}

# 2. Create Kong service with protocol bridge template
POST /api/kong/templates
{
  "workflow_id": "wf_123",
  "use_case": "protocol_bridge"
}

# 3. Your app calls Kong
GET http://localhost:8000/api/legacy-data?customer_id=12345
```

---

### 2. **Webhooks to Workflow** (Asynchronous Processing)

**Problem**: Stripe sends 1,000 webhooks per second during Black Friday. Your server crashes.

**Solution**:
- Kong enforces rate limiting: 100 req/sec
- Excess requests get `429 Too Many Requests`
- GoFlow processes webhooks asynchronously via worker pool
- Multi-step workflows execute in background

**Benefits**:
- ✅ Prevent webhook flooding
- ✅ Graceful degradation under load
- ✅ Automatic retry for failed webhooks

**Setup**:
```bash
# 1. Create a webhook workflow in GoFlow
POST /api/workflows
{
  "name": "Stripe Payment Handler",
  "trigger_type": "webhook",
  "action_type": "slack_message",
  "config_json": {
    "slack_message": "Payment received: {{amount}} from {{customer.email}}"
  }
}

# 2. Apply webhook handler template (adds rate limiting)
POST /api/kong/templates
{
  "workflow_id": "wf_456",
  "use_case": "webhook_handler"
}

# 3. Stripe sends webhooks to Kong
POST http://localhost:8000/webhooks/stripe
```

**Rate Limiting Configuration**:
- **Free Tier**: 100 requests/minute
- **Pro Tier**: 1,000 requests/minute
- **Enterprise**: Custom limits per tenant

---

### 3. **Smart API Aggregator** (Data Orchestration)

**Problem**: Your mobile app needs data from 3 APIs to render a dashboard. This causes 3 network round-trips.

**Solution**:
- Create a GoFlow workflow that calls 3 APIs in parallel
- Kong caches the response for 5 minutes
- Subsequent users get instant cached results

**Benefits**:
- ✅ 1 network call instead of 3
- ✅ 99% faster for cached responses
- ✅ Reduced load on backend APIs

**Setup**:
```bash
# 1. Create an aggregator workflow
POST /api/workflows
{
  "name": "Dashboard Aggregator",
  "action_type": "aggregator",
  "config_json": {
    "endpoints": [
      {"url": "https://api.salesforce.com/accounts", "key": "accounts"},
      {"url": "https://api.weather.com/current", "key": "weather"},
      {"url": "https://internal-db.com/metrics", "key": "metrics"}
    ]
  }
}

# 2. Apply aggregator template (adds caching)
POST /api/kong/templates
{
  "workflow_id": "wf_789",
  "use_case": "aggregator"
}

# 3. App calls once, gets all data
GET http://localhost:8000/api/dashboard
```

**Caching Configuration**:
- **Cache TTL**: 300 seconds (5 minutes)
- **Cache Key**: URL + Query Params
- **Cache Bypass**: `Cache-Control: no-cache` header

---

### 4. **Federated Security** (Auth Overlay)

**Problem**: You have 10 internal tools with no authentication. Partners need secure access.

**Solution**:
- Kong handles OAuth2 or API Key authentication
- Once authenticated, Kong injects a "Trust Header" (`X-Authenticated-User: alice@company.com`)
- GoFlow workflows trust this header and execute business logic

**Benefits**:
- ✅ Centralized authentication
- ✅ No code changes to workflows
- ✅ Support OAuth2, API Keys, JWT, LDAP

**Setup**:
```bash
# 1. Create a workflow (no auth logic needed!)
POST /api/workflows
{
  "name": "Internal Report Generator",
  "action_type": "pdf_generator"
}

# 2. Apply auth overlay template
POST /api/kong/templates
{
  "workflow_id": "wf_321",
  "use_case": "auth_overlay"
}

# 3. Kong requires API key
GET http://localhost:8000/api/reports
X-API-Key: secret_key_123

# 4. Kong validates key, adds trust header
# GoFlow receives: X-Authenticated-User: alice@company.com
```

**Supported Auth Methods**:
- **API Key**: Simple key in header or query param
- **OAuth2**: Full OAuth2 flow (authorization code, client credentials)
- **JWT**: Validate JWT tokens with RS256/HS256
- **LDAP**: Enterprise directory integration

---

### 5. **Usage-Based Monetization**

**Problem**: You want to charge customers $0.01 per API call, but don't want to build billing infrastructure.

**Solution**:
- Kong tracks every API call with `key-auth` plugin
- Kong emits logs to Elasticsearch (ELK)
- GoFlow runs a nightly job to calculate bills
- Invoices generated in Stripe or QuickBooks

**Benefits**:
- ✅ Real-time usage tracking
- ✅ Automatic billing calculations
- ✅ Per-tenant or per-user limits

**Setup**:
```bash
# 1. Create a workflow
POST /api/workflows
{
  "name": "Data Sync API",
  "action_type": "database_sync"
}

# 2. Apply monetization template (adds rate limiting + usage tracking)
POST /api/kong/templates
{
  "workflow_id": "wf_555",
  "use_case": "monetization"
}

# 3. Customer calls API with their key
GET http://localhost:8000/api/data-sync
X-API-Key: customer_key_abc123

# 4. Kong logs to ELK:
# {
#   "api_key": "customer_key_abc123",
#   "tenant_id": "acme_corp",
#   "request_count": 1,
#   "timestamp": "2026-01-12T12:00:00Z"
# }

# 5. Nightly job calculates bill
# Acme Corp: 10,000 API calls × $0.01 = $100.00
```

**Billing Dashboard**:
```sql
-- ELK Query: Daily API usage per tenant
SELECT tenant_id, COUNT(*) as api_calls, DATE(timestamp) as date
FROM kong_logs
WHERE timestamp >= NOW() - INTERVAL 30 DAY
GROUP BY tenant_id, date
ORDER BY api_calls DESC
```

---

## 📦 Quick Start

### 1. Start All Services
```bash
cd /Users/alex.macdonald/simple-ipass
docker-compose up -d
```

**Services**:
- **GoFlow Backend**: http://localhost:8080
- **GoFlow Frontend**: http://localhost:3000
- **Kong Gateway**: http://localhost:8000
- **Kong Admin API**: http://localhost:8001
- **Kong Manager (GUI)**: http://localhost:8002
- **Kibana (Logs)**: http://localhost:5601

### 2. Create Your First Kong Service

**Option A: Via GoFlow UI**
1. Go to http://localhost:3000/dashboard/api-management
2. Select a workflow
3. Choose a use case template
4. Click "Create Kong Service"

**Option B: Via API**
```bash
# Get your JWT token
TOKEN=$(curl -X POST http://localhost:8080/api/auth/login \
  -d '{"email":"user@example.com","password":"password"}' | jq -r '.token')

# Create Kong service
curl -X POST http://localhost:8080/api/kong/templates \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "workflow_id": "wf_123",
    "use_case": "webhook_handler"
  }'
```

### 3. Test Your Kong Service
```bash
# Call via Kong Gateway
curl http://localhost:8000/api/webhooks/wf_123 \
  -d '{"test": "data"}'

# Check Kong logs
curl http://localhost:8001/status
```

---

## 🔧 Kong Plugins Reference

### Available Plugins

| Plugin | Use Case | Configuration |
|--------|----------|---------------|
| **rate-limiting** | Prevent abuse | `second: 100, hour: 10000` |
| **key-auth** | API key authentication | `key_names: ["apikey", "X-API-Key"]` |
| **oauth2** | OAuth2 authentication | `scopes: ["read", "write"]` |
| **proxy-cache** | Response caching | `cache_ttl: 300` (5 min) |
| **request-size-limiting** | Prevent large payloads | `allowed_payload_size: 1` (1MB) |
| **cors** | Cross-origin requests | `origins: ["*"]` |
| **request-transformer** | Modify requests | `add.headers: ["X-Custom: value"]` |
| **response-transformer** | Modify responses | `add.json: ["status": "success"]` |

### Add Custom Plugins

```bash
# Example: Add IP restriction
curl -X POST http://localhost:8001/services/bridge-wf_123/plugins \
  -d "name=ip-restriction" \
  -d "config.whitelist=192.168.1.0/24"
```

---

## 📊 Monitoring & Analytics

### Kong Logs in Kibana

1. Open http://localhost:5601
2. Create data view: `kong-*`
3. Visualize:
   - **API call volume**: Line chart by timestamp
   - **Top consumers**: Pie chart by `api_key`
   - **Error rate**: Count where `status >= 400`

### Real-time Kong Status
```bash
curl http://localhost:8001/status
```

**Response**:
```json
{
  "database": {
    "reachable": true
  },
  "server": {
    "connections_accepted": 1234,
    "connections_handled": 1234,
    "total_requests": 5678
  }
}
```

---

## 🎯 Production Recommendations

### 1. Kong + PostgreSQL (Not SQLite)
```yaml
# docker-compose.yml (already configured!)
kong-database:
  image: postgres:16-alpine
  environment:
    POSTGRES_DB: kong
    POSTGRES_USER: kong
    POSTGRES_PASSWORD: strong_password_here
```

### 2. Enable Kong Plugins Globally
```bash
# Enable rate limiting for ALL services
curl -X POST http://localhost:8001/plugins \
  -d "name=rate-limiting" \
  -d "config.second=1000"
```

### 3. Use Kong Secrets Manager
```bash
# Store API keys securely
curl -X POST http://localhost:8001/consumers/alice/key-auth \
  -d "key=securely_generated_key"
```

### 4. Set Up Kong Clustering (Enterprise)
- **Control Plane**: Manages configuration
- **Data Plane**: Handles traffic
- **Hybrid Mode**: Decoupled for scale

---

## 🔐 Security Best Practices

### 1. Always Use HTTPS in Production
```yaml
# docker-compose.yml
kong:
  environment:
    KONG_SSL_CERT: /certs/server.crt
    KONG_SSL_CERT_KEY: /certs/server.key
```

### 2. Enable Request Validation
```bash
# Add request validator plugin
curl -X POST http://localhost:8001/services/my-service/plugins \
  -d "name=request-validator" \
  -d "config.body_schema=..."
```

### 3. Implement API Key Rotation
```bash
# Rotate customer keys monthly
curl -X DELETE http://localhost:8001/consumers/alice/key-auth/old_key_id
curl -X POST http://localhost:8001/consumers/alice/key-auth
```

---

## 🧪 Testing Kong Integration

### 1. Test Protocol Bridge
```bash
# Call Kong endpoint
curl http://localhost:8000/api/soap-bridge \
  -d '{"customer_id": "12345"}'

# Verify SOAP call in logs
docker logs kong-gateway | grep "SOAP"
```

### 2. Test Rate Limiting
```bash
# Spam Kong with requests
for i in {1..150}; do
  curl http://localhost:8000/api/test
done

# Should see 429 Too Many Requests after 100
```

### 3. Test Caching
```bash
# First call (miss)
time curl http://localhost:8000/api/cached-endpoint
# Takes 500ms

# Second call (hit)
time curl http://localhost:8000/api/cached-endpoint
# Takes 5ms (cached!)
```

---

## 📚 Additional Resources

- **Kong Documentation**: https://docs.konghq.com/
- **Kong Plugin Hub**: https://docs.konghq.com/hub/
- **GoFlow + Kong Examples**: See `/scripts/kong-examples.sh`

---

## 🎉 Summary

**What You've Built**:
- ✅ SOAP to REST protocol bridge
- ✅ Rate-limited webhook processing
- ✅ Smart API aggregation with caching
- ✅ Centralized authentication overlay
- ✅ Usage-based monetization tracking

**Production-Ready Features**:
- ✅ Kong Gateway with PostgreSQL backend
- ✅ 5 battle-tested use case templates
- ✅ Frontend UI for non-developers
- ✅ ELK integration for billing analytics
- ✅ Health checks & monitoring

**Your GoFlow platform is now an enterprise-grade API Gateway + iPaaS!** 🚀

---

**Next Steps**:
1. Deploy to Kubernetes (Kong has native k8s support)
2. Enable Kong Enterprise for RBAC & advanced analytics
3. Set up multi-region Kong clusters for global scale
4. Implement custom Kong plugins for your business logic

