# 🚀 Quick Start: New Connectors

## Twilio SMS

**Send an SMS with dynamic data:**

```bash
# 1. Add Twilio credentials
curl -X POST http://localhost:8080/api/credentials \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "service_name": "twilio",
    "api_key": "{\"account_sid\":\"ACxxx\",\"auth_token\":\"xxx\",\"from_number\":\"+15551234567\"}"
  }'

# 2. Create workflow
curl -X POST http://localhost:8080/api/workflows \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Order SMS Alert",
    "trigger_type": "webhook",
    "action_type": "twilio_sms",
    "config_json": "{\"twilio_to\":\"{{customer.phone}}\",\"twilio_message\":\"Hi {{customer.name}}! Order {{order.id}} confirmed.\"}"
  }'

# 3. Trigger with webhook
curl -X POST http://localhost:8080/api/webhooks/WORKFLOW_ID \
  -d '{"customer":{"name":"Alex","phone":"+15559876543"},"order":{"id":"12345"}}'
```

---

## News API

**Fetch latest tech news:**

```bash
# 1. Add News API key
curl -X POST http://localhost:8080/api/credentials \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "service_name": "newsapi",
    "api_key": "your_newsapi_key"
  }'

# 2. Create workflow
curl -X POST http://localhost:8080/api/workflows \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Tech News Digest",
    "trigger_type": "schedule",
    "action_type": "news_fetch",
    "config_json": "{\"news_query\":\"artificial intelligence\",\"news_category\":\"technology\",\"news_page_size\":5}"
  }'
```

---

## The Cat API

**Daily cat picture bot:**

```bash
# 1. Create workflow (no credentials needed!)
curl -X POST http://localhost:8080/api/workflows \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Daily Cat",
    "trigger_type": "schedule",
    "action_type": "cat_fetch",
    "config_json": "{\"cat_limit\":1,\"cat_has_breeds\":true}"
  }'
```

---

## Fake Store API

**Fetch product catalog:**

```bash
# 1. Create workflow (no credentials needed!)
curl -X POST http://localhost:8080/api/workflows \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Product Sync",
    "trigger_type": "schedule",
    "action_type": "fakestore_fetch",
    "config_json": "{\"fakestore_endpoint\":\"products\",\"fakestore_category\":\"electronics\",\"fakestore_limit\":10}"
  }'
```

---

## Dynamic Templates

**Use webhook data in messages:**

**Slack:**
```json
{
  "slack_message": "Hello {{user.name}}! Your order #{{order.id}} for ${{order.total}} is ready."
}
```

**Discord:**
```json
{
  "discord_message": "New signup: {{user.email}} from {{user.country}}"
}
```

**Twilio:**
```json
{
  "twilio_to": "{{customer.phone}}",
  "twilio_message": "Hi {{customer.name}}, your code is {{verification.code}}"
}
```

---

## Template Syntax

- Simple: `{{name}}`
- Nested: `{{user.email}}`
- Array: `{{items.0.title}}`
- Deep: `{{order.shipping.address.city}}`

---

**See NEW_CONNECTORS.md for detailed documentation!**

