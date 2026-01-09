# GoFlow API Documentation

This directory contains the OpenAPI Specification (OAS 3.0) for the GoFlow API.

---

## 📄 Files

- **`openapi.yaml`** - Complete OpenAPI 3.0 specification

---

## 🌐 View the API Documentation

### Option 1: Swagger UI (Recommended)

Visit the online Swagger Editor:
```
https://editor.swagger.io/
```

Then paste the contents of `openapi.yaml` or upload the file.

### Option 2: Redoc (Beautiful Documentation)

```bash
# Install redoc-cli
npm install -g redoc-cli

# Generate HTML documentation
redoc-cli bundle openapi.yaml -o api-docs.html

# Open in browser
open api-docs.html
```

### Option 3: Local Swagger UI with Docker

```bash
docker run -p 8081:8080 -e SWAGGER_JSON=/openapi.yaml -v $(pwd)/openapi.yaml:/openapi.yaml swaggerapi/swagger-ui
```

Then visit: http://localhost:8081

### Option 4: VS Code Extension

Install the "OpenAPI (Swagger) Editor" extension in VS Code and open `openapi.yaml`.

---

## 📚 API Overview

### Authentication

All protected endpoints require a JWT token obtained from:
- `POST /api/auth/register` - Create account
- `POST /api/auth/login` - Get token

Include the token in requests:
```bash
Authorization: Bearer <your-jwt-token>
```

---

## 🚀 Quick Examples

### 1. Register and Login

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### 2. Create Workflow

```bash
TOKEN="your-jwt-token-here"

curl -X POST http://localhost:8080/api/workflows \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Slack Alert",
    "trigger_type": "webhook",
    "action_type": "slack_message",
    "config_json": "{\"credential_id\":\"cred_123\",\"slack_message\":\"Alert!\"}"
  }'
```

### 3. Test Workflow (Dry Run)

```bash
curl -X POST http://localhost:8080/api/workflows/dry-run \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action_type": "slack_message",
    "config_json": "{\"credential_id\":\"cred_123\",\"slack_message\":\"Test\"}"
  }'
```

### 4. Trigger Webhook

```bash
# No authentication required for webhook triggers
curl -X POST http://localhost:8080/api/webhooks/wf_abc123 \
  -H "Content-Type: application/json" \
  -H "X-Idempotency-Key: unique-uuid-12345" \
  -d '{"event":"test","data":"hello"}'
```

### 5. Health Check

```bash
# Comprehensive health check
curl http://localhost:8080/health

# Kubernetes liveness
curl http://localhost:8080/health/live

# Kubernetes readiness
curl http://localhost:8080/health/ready
```

---

## 📊 API Endpoints Summary

### Public Endpoints (No Auth Required)
- `POST /api/auth/register` - Create account
- `POST /api/auth/login` - Get JWT token
- `POST /api/webhooks/{id}` - Trigger webhook workflow
- `GET /health` - Health check
- `GET /health/live` - Liveness probe
- `GET /health/ready` - Readiness probe

### Protected Endpoints (JWT Required)
- `POST /api/credentials` - Save credentials
- `GET /api/credentials` - List credentials
- `POST /api/workflows` - Create workflow
- `GET /api/workflows` - List workflows
- `POST /api/workflows/dry-run` - Test workflow
- `PUT /api/workflows/{id}/toggle` - Enable/disable workflow
- `DELETE /api/workflows/{id}` - Delete workflow
- `GET /api/logs` - Get execution logs

---

## 🔐 Security Features

- **JWT Authentication** - Secure token-based auth
- **Rate Limiting** - 5 req/sec (free), 50 req/sec (paid)
- **Idempotency Keys** - Prevent duplicate operations
- **CORS** - Configurable origins
- **Secret Masking** - Logs never contain credentials

---

## 🌟 Advanced Features

### Idempotency

Prevent duplicate webhook executions:
```bash
curl -X POST http://localhost:8080/api/webhooks/wf_123 \
  -H "X-Idempotency-Key: uuid-12345" \
  -d '{"data":"test"}'
```

Same key within 24 hours returns cached result.

### Rate Limiting

Response headers:
```
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 3
Retry-After: 1
```

### Template Variables

Use dynamic data in workflows:
```json
{
  "slack_message": "Hello {{user.name}}, order {{order.id}} shipped!"
}
```

---

## 📖 Response Format

All responses follow a consistent envelope:

**Success:**
```json
{
  "success": true,
  "data": { ... }
}
```

**Error:**
```json
{
  "success": false,
  "error": "Error message"
}
```

---

## 🧪 Testing

### Postman Collection

Import `openapi.yaml` into Postman:
1. Open Postman
2. File → Import
3. Select `openapi.yaml`
4. All endpoints automatically configured!

### Insomnia

Import into Insomnia:
1. Create → Import From → URL or File
2. Select `openapi.yaml`

---

## 🔄 API Versioning

Current version: **v0.3.0**

All endpoints are prefixed with `/api/` to allow future versioning:
- Current: `/api/workflows`
- Future: `/api/v2/workflows`

---

## 📝 Changelog

### v0.3.0 (2026-01-08)
- ✅ Added idempotency support
- ✅ Added rate limiting
- ✅ Added template engine
- ✅ Enhanced health checks
- ✅ Circuit breaker pattern
- ✅ Secret masking

### v0.2.0 (2026-01-07)
- ✅ Production hardening
- ✅ Worker pool
- ✅ Context-aware execution
- ✅ Repository pattern

### v0.1.0 (2026-01-06)
- ✅ Initial release
- ✅ Basic authentication
- ✅ Workflow management
- ✅ Three connectors

---

## 🤝 Contributing

Found an issue with the API? Please open an issue on GitHub.

---

## 📄 License

MIT License - See LICENSE file for details

---

**GoFlow** - S-Tier Enterprise Integration Platform  
Documentation generated from OpenAPI 3.0 specification

