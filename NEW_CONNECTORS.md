# 🎉 New Connectors & Dynamic Field Mapping

## Overview

GoFlow now supports **4 new connectors** and **dynamic field mapping** for Slack and Discord messages! This allows you to create truly dynamic workflows that map data from webhook payloads into your actions.

---

## ✨ New Feature: Dynamic Field Mapping

### **Template Syntax**

Use `{{path.to.field}}` in your Slack/Discord/Twilio messages to dynamically insert data from webhook payloads.

### **Example**

**Webhook Payload:**
```json
{
  "user": {
    "name": "Alex",
    "email": "alex@example.com"
  },
  "order": {
    "id": "12345",
    "total": 99.99
  }
}
```

**Slack Message Template:**
```
Hello {{user.name}}! Your order #{{order.id}} for ${{order.total}} has been placed.
```

**Result:**
```
Hello Alex! Your order #12345 for $99.99 has been placed.
```

### **Supported in:**
- ✅ Slack messages (`slack_message`)
- ✅ Discord messages (`discord_post`)
- ✅ Twilio SMS (`twilio_sms`)

---

## 📱 1. Twilio SMS Connector

Send SMS messages using [Twilio](https://www.twilio.com/).

### **Action Type:** `twilio_sms`

### **Credentials Format:**

```json
{
  "account_sid": "ACxxxxxxxxxxxxxxxxxxxxx",
  "auth_token": "your_auth_token",
  "from_number": "+15551234567"
}
```

### **Configuration:**

```json
{
  "twilio_to": "+15559876543",
  "twilio_message": "Hello {{user.name}}, your code is {{verification.code}}"
}
```

### **Features:**
- ✅ Dynamic template mapping for phone number and message
- ✅ Context-aware execution with timeouts
- ✅ Full error handling

### **Use Cases:**
- Send verification codes
- Order notifications
- Alert messages
- Two-factor authentication

---

## 📰 2. News API Connector

Fetch news articles from [NewsAPI](https://newsapi.org/).

### **Action Type:** `news_fetch`

### **Credentials:**
- Service Name: `newsapi`
- API Key: Your News API key from https://newsapi.org/

### **Configuration:**

```json
{
  "news_query": "bitcoin",
  "news_country": "us",
  "news_category": "technology",
  "news_page_size": 10
}
```

### **Parameters:**
- `news_query` - Search term (e.g., "bitcoin", "ai")
- `news_country` - Country code (e.g., "us", "gb", "ca")
- `news_category` - Category: "business", "entertainment", "health", "science", "sports", "technology"
- `news_page_size` - Number of articles (1-100, default: 10)

### **Response Data:**
```json
{
  "total_results": 156,
  "articles": [
    {
      "source": {"name": "TechCrunch"},
      "author": "John Doe",
      "title": "Breaking Tech News",
      "description": "...",
      "url": "https://...",
      "publishedAt": "2026-01-09T00:00:00Z"
    }
  ]
}
```

### **Use Cases:**
- Daily news digests
- Market monitoring
- Topic alerts
- Competitor tracking

---

## 🐱 3. The Cat API Connector

Fetch adorable cat images from [The Cat API](https://thecatapi.com/).

### **Action Type:** `cat_fetch`

### **Credentials (Optional):**
- Service Name: `catapi`
- API Key: Optional (get from https://thecatapi.com/)

### **Configuration:**

```json
{
  "cat_limit": 3,
  "cat_has_breeds": true,
  "cat_breed_id": "beng",
  "cat_category": "boxes"
}
```

### **Parameters:**
- `cat_limit` - Number of cat images (1-10, default: 1)
- `cat_has_breeds` - Filter to cats with breed information
- `cat_breed_id` - Specific breed (e.g., "beng" for Bengal, "pers" for Persian)
- `cat_category` - Category (e.g., "boxes", "hats", "sinks")

### **Response Data:**
```json
{
  "cats": [
    {
      "id": "MTY3ODIyMQ",
      "url": "https://cdn2.thecatapi.com/images/...",
      "width": 1080,
      "height": 1080,
      "breeds": [
        {
          "name": "Bengal",
          "temperament": "Alert, Agile, Energetic",
          "origin": "United States",
          "description": "..."
        }
      ]
    }
  ]
}
```

### **Use Cases:**
- Daily cat picture bot
- Fun Slack/Discord integrations
- Website content generation
- Social media automation

---

## 🛒 4. Fake Store API Connector

Fetch mock e-commerce data from [Fake Store API](https://fakestoreapi.com/).

### **Action Type:** `fakestore_fetch`

### **No Credentials Required** - Public API

### **Configuration:**

```json
{
  "fakestore_endpoint": "products",
  "fakestore_limit": 5,
  "fakestore_category": "electronics"
}
```

### **Parameters:**
- `fakestore_endpoint` - Endpoint: "products", "users", "carts", "categories"
- `fakestore_limit` - Number of items (1-20)
- `fakestore_category` - Product category: "electronics", "jewelery", "men's clothing", "women's clothing"

### **Response Data (Products):**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Fjallraven Backpack",
      "price": 109.95,
      "description": "...",
      "category": "men's clothing",
      "image": "https://fakestoreapi.com/img/...",
      "rating": {
        "rate": 3.9,
        "count": 120
      }
    }
  ]
}
```

### **Use Cases:**
- Testing e-commerce integrations
- Demo applications
- UI prototyping
- Training and education

---

## 🎯 Complete Workflow Examples

### Example 1: Order Notification via Twilio

**Workflow:**
- **Trigger:** Webhook (order placed)
- **Action:** Twilio SMS

**Payload:**
```json
{
  "customer": {
    "name": "Sarah",
    "phone": "+15559876543"
  },
  "order": {
    "id": "ORD-12345",
    "total": 149.99,
    "items": 3
  }
}
```

**Configuration:**
```json
{
  "twilio_to": "{{customer.phone}}",
  "twilio_message": "Hi {{customer.name}}! Your order {{order.id}} (${{order.total}}) is confirmed. Thank you!"
}
```

**Result SMS:**
```
Hi Sarah! Your order ORD-12345 ($149.99) is confirmed. Thank you!
```

---

### Example 2: Slack Alert for Tech News

**Workflow:**
- **Trigger:** Schedule (every 6 hours)
- **Action 1:** News API (fetch tech news)
- **Action 2:** Slack (send summary)

**News Config:**
```json
{
  "news_query": "artificial intelligence",
  "news_category": "technology",
  "news_page_size": 3
}
```

---

### Example 3: Daily Cat Picture Bot

**Workflow:**
- **Trigger:** Schedule (daily at 9 AM)
- **Action 1:** Cat API (fetch cat)
- **Action 2:** Discord (post image)

**Cat Config:**
```json
{
  "cat_limit": 1,
  "cat_has_breeds": true
}
```

**Discord Message:**
```json
{
  "discord_message": "Good morning! 🐱 Here's your daily cat: {{cats.0.url}}"
}
```

---

## 📊 Connector Summary

| Connector | Action Type | Auth Required | Template Support | Use Case |
|-----------|-------------|---------------|------------------|----------|
| **Slack** | `slack_message` | Webhook URL | ✅ Yes | Team notifications |
| **Discord** | `discord_post` | Webhook URL | ✅ Yes | Community updates |
| **Twilio** | `twilio_sms` | Account SID + Token | ✅ Yes | SMS alerts |
| **News API** | `news_fetch` | API Key | ❌ No | News aggregation |
| **Cat API** | `cat_fetch` | Optional | ❌ No | Fun content |
| **Fake Store** | `fakestore_fetch` | None | ❌ No | Testing/demos |
| **OpenWeather** | `weather_check` | API Key | ❌ No | Weather data |

---

## 🔧 API Endpoints

### Create Twilio Credentials

```bash
curl -X POST http://localhost:8080/api/credentials \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "twilio",
    "api_key": "{\"account_sid\":\"AC...\",\"auth_token\":\"...\",\"from_number\":\"+15551234567\"}"
  }'
```

### Create News API Credentials

```bash
curl -X POST http://localhost:8080/api/credentials \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "newsapi",
    "api_key": "your_news_api_key"
  }'
```

### Create Cat API Credentials (Optional)

```bash
curl -X POST http://localhost:8080/api/credentials \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "catapi",
    "api_key": "your_cat_api_key"
  }'
```

---

## 🎓 Template Engine Details

### Supported Syntax

- **Simple paths:** `{{name}}`
- **Nested objects:** `{{user.email}}`
- **Array access:** `{{items.0.name}}`
- **Deep nesting:** `{{order.shipping.address.city}}`

### Invalid Paths

If a path doesn't exist in the payload, the template variable is left unchanged:

**Payload:**
```json
{"name": "Alex"}
```

**Template:**
```
Hello {{name}}, your order {{order.id}} is ready.
```

**Result:**
```
Hello Alex, your order {{order.id}} is ready.
```

---

## 📚 Getting API Keys

### Twilio
1. Sign up at https://www.twilio.com/
2. Get your Account SID and Auth Token from the console
3. Purchase a phone number

### News API
1. Sign up at https://newsapi.org/
2. Get your API key from the dashboard
3. Free tier: 100 requests/day

### The Cat API
1. Optional - works without key
2. Sign up at https://thecatapi.com/ for higher limits
3. Free tier: 10,000 requests/month

### Fake Store API
- No registration needed!
- Completely free and public

---

## 🚀 Next Steps

1. **Try the examples** - Create a workflow with dynamic templates
2. **Combine connectors** - Fetch cat pictures and post to Slack!
3. **Build complex workflows** - Chain multiple actions together
4. **Read the API docs** - See `openapi.yaml` for full API reference

---

## 🎉 Summary

**New Connectors: 4**
- ✅ Twilio SMS
- ✅ News API
- ✅ The Cat API
- ✅ Fake Store API

**Enhanced Connectors: 3**
- ✅ Slack (dynamic templates)
- ✅ Discord (dynamic templates)
- ✅ Twilio (dynamic templates)

**Total Connectors: 7**

**Your GoFlow platform is now even more powerful!** 🚀

---

**Date**: January 9, 2026  
**Version**: 0.4.0  
**Status**: Production-Ready with Enhanced Connectors ✅

