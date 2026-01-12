# Kong Gateway Quick Reference

## 🚀 Quick Commands

### Start Services
```bash
docker-compose up -d
```

### Access Points
- **Kong Gateway**: http://localhost:8000
- **Kong Admin**: http://localhost:8001
- **Kong Manager**: http://localhost:8002
- **GoFlow UI**: http://localhost:3000/dashboard/api-management

---

## 📋 5 Use Cases (TL;DR)

| Use Case | What It Does | When to Use |
|----------|--------------|-------------|
| **Protocol Bridge** | SOAP → REST | Legacy system modernization |
| **Webhook Handler** | Rate limiting | High-volume webhooks |
| **Aggregator** | Response caching | Reduce API chattiness |
| **Auth Overlay** | API key/OAuth2 | Secure internal tools |
| **Monetization** | Usage tracking | Pay-per-use billing |

---

## 🔧 Common API Calls

### Create Kong Service (Protocol Bridge)
```bash
curl -X POST http://localhost:8080/api/kong/templates \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "workflow_id": "wf_123",
    "use_case": "protocol_bridge"
  }'
```

### Create Kong Service (Rate Limiting)
```bash
curl -X POST http://localhost:8080/api/kong/templates \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "workflow_id": "wf_456",
    "use_case": "webhook_handler"
  }'
```

### List Kong Services
```bash
curl http://localhost:8001/services
```

### Delete Kong Service
```bash
curl -X DELETE http://localhost:8001/services/service_id
```

---

## 🎯 Example: SOAP Bridge in 3 Steps

### Step 1: Create SOAP Workflow
```bash
curl -X POST http://localhost:8080/api/workflows \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Legacy SOAP Bridge",
    "action_type": "soap_call",
    "trigger_type": "webhook",
    "config_json": {
      "soap_endpoint": "http://legacy.com/service.asmx",
      "soap_method": "GetData",
      "soap_namespace": "http://tempuri.org/",
      "soap_parameters": {"id": "{{customer_id}}"}
    }
  }'
```

### Step 2: Apply Kong Template
```bash
curl -X POST http://localhost:8080/api/kong/templates \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "workflow_id": "wf_xyz",
    "use_case": "protocol_bridge"
  }'
```

### Step 3: Call via Kong
```bash
curl http://localhost:8000/api/webhooks/wf_xyz \
  -d '{"customer_id": "12345"}'
```

---

## 📊 Monitoring

### Check Kong Health
```bash
curl http://localhost:8001/status
```

### View Kong Logs
```bash
docker logs kong-gateway
```

### ELK Dashboard
1. Open http://localhost:5601
2. Create index: `kong-*`
3. Visualize API usage

---

## 🔐 Security

### Add API Key Auth
```bash
curl -X POST http://localhost:8001/services/my-service/plugins \
  -d "name=key-auth"
```

### Create API Key
```bash
curl -X POST http://localhost:8001/consumers/alice/key-auth \
  -d "key=secret_key_123"
```

### Test with API Key
```bash
curl http://localhost:8000/my-endpoint \
  -H "X-API-Key: secret_key_123"
```

---

## ⚡ Performance

### Enable Caching (5 min)
```bash
curl -X POST http://localhost:8001/services/my-service/plugins \
  -d "name=proxy-cache" \
  -d "config.cache_ttl=300"
```

### Set Rate Limits
```bash
curl -X POST http://localhost:8001/services/my-service/plugins \
  -d "name=rate-limiting" \
  -d "config.second=100" \
  -d "config.hour=10000"
```

---

## 🐛 Troubleshooting

### Kong won't start
```bash
# Check database
docker logs kong-database

# Run migrations
docker-compose up kong-migration
```

### Can't connect to Kong Admin
```bash
# Verify port 8001 is exposed
docker ps | grep kong

# Check Kong health
curl http://localhost:8001/status
```

### 404 on Kong routes
```bash
# List all routes
curl http://localhost:8001/routes

# Check service configuration
curl http://localhost:8001/services/my-service
```

---

## 📚 Full Documentation
See [KONG_INTEGRATION.md](./KONG_INTEGRATION.md) for complete guide.

