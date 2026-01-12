# Multi-Step Workflows (Action Chaining)

## 🔗 Overview

GoFlow now supports **multi-step workflows** where you can chain multiple actions together. Instead of just executing one action, you can create sequences like:

**Schedule → Check Weather → Send to Discord → Send SMS**

This solves the problem of wanting to use data from one action (like weather data) in subsequent actions (like sending that data to Discord and Twilio).

---

## 🌟 Key Features

- ✅ **Sequential Execution**: Actions execute one after another
- ✅ **Data Passing**: Use data from previous action in next action
- ✅ **Context-Aware**: Respects cancellation throughout the chain
- ✅ **Flexible**: Chain any combination of messaging actions
- ✅ **Performance Tracking**: Each chain step is logged separately
- ✅ **Failure Handling**: Chain continues even if one step fails

---

## 📋 Supported Chain Actions

Currently, these action types can be chained:
- ✅ `slack_message` - Send to Slack
- ✅ `discord_post` - Send to Discord  
- ✅ `twilio_sms` - Send SMS

**Primary Actions** (can be used as first action):
- All 12 connectors: weather_check, swapi_fetch, salesforce, news_fetch, etc.

---

## 🎯 Use Case: Weather to Multiple Channels

### Problem
You want to check the weather every hour and send that data to both Discord AND send an SMS.

### Solution
Create a workflow with action chaining:

```json
{
  "name": "Hourly Weather Multi-Channel",
  "trigger_type": "schedule",
  "action_type": "weather_check",
  "config_json": {
    "interval": 60,
    "city": "London"
  },
  "action_chain": [
    {
      "action_type": "discord_post",
      "config": {
        "discord_message": "🌤️ Weather Update: {{weather.main}} in {{name}}, Temp: {{main.temp}}°C"
      },
      "use_data_from": "previous"
    },
    {
      "action_type": "twilio_sms",
      "config": {
        "twilio_to": "+1-555-1234",
        "twilio_message": "Weather: {{weather.main}}, {{main.temp}}°C in {{name}}"
      },
      "use_data_from": "previous"
    }
  ]
}
```

### What Happens
1. **Step 1 (Primary)**: Check weather in London
2. **Step 2 (Chain)**: Send weather data to Discord
3. **Step 3 (Chain)**: Send weather data via SMS
4. **Result**: Weather data delivered to 2 channels automatically!

---

## 🔧 How to Create Multi-Step Workflows

### Method 1: API Request

```bash
curl -X POST http://localhost:8080/api/workflows \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Weather to Discord and SMS",
    "trigger_type": "schedule",
    "action_type": "weather_check",
    "config_json": "{\"interval\":60,\"city\":\"London\"}",
    "action_chain": [
      {
        "action_type": "discord_post",
        "config": {
          "discord_message": "Weather: {{weather.main}} in {{name}}"
        },
        "use_data_from": "previous"
      },
      {
        "action_type": "twilio_sms",
        "config": {
          "twilio_to": "+1-555-1234",
          "twilio_message": "Weather update: {{weather.main}}"
        },
        "use_data_from": "previous"
      }
    ]
  }'
```

### Method 2: Frontend (Coming Soon)
The workflow creation UI will have a visual chain builder where you can:
1. Select primary action (e.g., "Check Weather")
2. Click "+ Add Action" to add chain steps
3. Configure each step
4. See the flow diagram update in real-time

---

## 📊 Action Chain Structure

### ChainedAction Schema
```json
{
  "action_type": "slack_message|discord_post|twilio_sms",
  "config": {
    // Action-specific configuration
    "slack_message": "Your message here",
    "discord_message": "Your message here",
    "twilio_to": "+1-555-1234",
    "twilio_message": "Your SMS here"
  },
  "use_data_from": "previous" // Optional: use data from previous action
}
```

### Key Fields
- **`action_type`**: The type of action to execute
- **`config`**: Configuration specific to that action type
- **`use_data_from`**: If set to `"previous"`, the action will receive data from the previous action's result

---

## 🔄 Data Flow

### Without `use_data_from`
Each action is independent:
```
Weather Check → {temp: 15, city: "London"}
  ↓
Discord Post → "Generic message" (no weather data)
  ↓
Twilio SMS → "Generic message" (no weather data)
```

### With `use_data_from: "previous"`
Data flows through the chain:
```
Weather Check → {temp: 15, city: "London", weather: {main: "Clouds"}}
  ↓ (passes data)
Discord Post → "Weather: Clouds in London" ✅
  ↓ (passes data)
Twilio SMS → "Weather update: Clouds" ✅
```

---

## 🎨 Template Mapping in Chains

When `use_data_from: "previous"` is set, you can use template syntax to access data from the previous action:

### Weather API Response
```json
{
  "name": "London",
  "weather": [
    {"main": "Clouds", "description": "overcast clouds"}
  ],
  "main": {
    "temp": 15.2,
    "feels_like": 14.5,
    "humidity": 82
  }
}
```

### Template Usage in Chain
```json
{
  "action_type": "slack_message",
  "config": {
    "slack_message": "🌤️ {{name}}: {{weather.main}}, {{main.temp}}°C, Humidity: {{main.humidity}}%"
  },
  "use_data_from": "previous"
}
```

### Result
`"🌤️ London: Clouds, 15.2°C, Humidity: 82%"`

---

## 📋 Complete Examples

### Example 1: SWAPI to Multiple Channels
Fetch Star Wars data and share to Discord and Slack:

```json
{
  "name": "Star Wars Character to Social",
  "trigger_type": "schedule",
  "action_type": "swapi_fetch",
  "config_json": {
    "interval": 60,
    "swapi_resource": "people",
    "swapi_id": "1"
  },
  "action_chain": [
    {
      "action_type": "discord_post",
      "config": {
        "discord_message": "⭐ Character: {{name}}, Height: {{height}}cm, Homeworld: {{homeworld}}"
      },
      "use_data_from": "previous"
    },
    {
      "action_type": "slack_message",
      "config": {
        "slack_message": "Star Wars Trivia: {{name}} was born in {{birth_year}}"
      },
      "use_data_from": "previous"
    }
  ]
}
```

---

### Example 2: Salesforce Query to Notifications
Query Salesforce and notify multiple channels:

```json
{
  "name": "High-Value Accounts Alert",
  "trigger_type": "schedule",
  "action_type": "salesforce",
  "config_json": {
    "interval": 1440,
    "salesforce_operation": "query",
    "salesforce_query": "SELECT Id, Name, AnnualRevenue FROM Account WHERE AnnualRevenue > 1000000 LIMIT 5"
  },
  "action_chain": [
    {
      "action_type": "slack_message",
      "config": {
        "slack_message": "💰 High-value accounts found! Check your dashboard."
      }
    },
    {
      "action_type": "twilio_sms",
      "config": {
        "twilio_to": "+1-555-SALES",
        "twilio_message": "New high-value accounts detected. Review now!"
      }
    }
  ]
}
```

---

### Example 3: News to Multi-Channel
Fetch news and distribute:

```json
{
  "name": "Breaking News Multi-Channel",
  "trigger_type": "schedule",
  "action_type": "news_fetch",
  "config_json": {
    "interval": 30,
    "news_query": "technology",
    "news_country": "us",
    "news_page_size": 5
  },
  "action_chain": [
    {
      "action_type": "discord_post",
      "config": {
        "discord_message": "📰 Latest Tech News: Check your feed!"
      }
    },
    {
      "action_type": "slack_message",
      "config": {
        "slack_message": "🔔 New tech news articles available"
      }
    }
  ]
}
```

---

## 🏗️ Database Schema

### Workflows Table
```sql
CREATE TABLE workflows (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    trigger_type TEXT NOT NULL,
    action_type TEXT NOT NULL,
    config_json TEXT NOT NULL,
    action_chain TEXT,          -- JSON array of ChainedAction objects
    is_active BOOLEAN DEFAULT 1,
    last_executed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Action Chain JSON Format
```json
[
  {
    "action_type": "discord_post",
    "config": {"discord_message": "Message 1"},
    "use_data_from": "previous"
  },
  {
    "action_type": "twilio_sms",
    "config": {"twilio_to": "+1-555-1234", "twilio_message": "Message 2"},
    "use_data_from": "previous"
  }
]
```

---

## ⚡ Performance

### Execution Times
- **Primary Action**: Variable (50ms - 2000ms depending on connector)
- **Chain Action Overhead**: ~10ms per action
- **Total**: Primary + (Chain Count × Action Time) + Overhead

### Example
```
Weather Check: 500ms
Discord Post: 200ms
Twilio SMS: 300ms
Overhead: 20ms
---
Total: 1020ms (1 second)
```

---

## 🔍 Logging

### Primary Action Log
```json
{
  "workflow_id": "wf_123",
  "status": "success",
  "message": "Weather check completed | Chain: 2/2 actions succeeded",
  "data": {
    "weather": {...},
    "chain_results": [
      {"status": "success", "message": "Discord message sent"},
      {"status": "success", "message": "SMS sent to +1-555-1234"}
    ],
    "chain_count": 2
  }
}
```

### Individual Chain Logs
Each chain action is logged separately with:
- `chain_step`: 1, 2, 3, etc.
- `total_steps`: Total number of chain actions
- `action_type`: Type of chained action
- `status`: success/failed/cancelled

---

## 🚫 Limitations

### Current Limitations
1. **Chain actions are messaging only**: Only `slack_message`, `discord_post`, and `twilio_sms` supported in chains
2. **No branching**: Actions execute sequentially, no if/then logic
3. **No loops**: Chain executes once per trigger
4. **Maximum chain length**: Recommended max 10 actions for performance

### Future Enhancements
- [ ] Add conditional branching (if/then/else)
- [ ] Support all 12 connectors in chains
- [ ] Parallel execution option
- [ ] Error handling strategies (stop-on-error vs continue)
- [ ] Custom retry logic per chain step
- [ ] Visual chain builder in frontend

---

## 🎯 Best Practices

### 1. Keep Chains Short
✅ **Good**: 2-3 chain actions  
❌ **Bad**: 10+ chain actions (slow, hard to debug)

### 2. Use Meaningful Messages
✅ **Good**: `"Weather: {{weather.main}} in {{name}}"`  
❌ **Bad**: `"Update"` (no context)

### 3. Test with Dry Run
Before creating a chained workflow, test each action individually

### 4. Handle Missing Data
Templates gracefully handle missing fields by showing empty string

### 5. Monitor Chain Results
Check `chain_results` in logs to see which steps succeeded/failed

---

## 🧪 Testing

### Test Action Chain
```bash
# Create test workflow
POST /api/workflows
{
  "name": "Test Chain",
  "trigger_type": "webhook",
  "action_type": "weather_check",
  "config_json": "{\"city\":\"London\"}",
  "action_chain": [
    {
      "action_type": "slack_message",
      "config": {"slack_message": "Test: {{name}}"},
      "use_data_from": "previous"
    }
  ]
}

# Trigger workflow
POST /api/webhooks/wf_xyz
{}

# Check logs
GET /api/logs?workflow_id=wf_xyz
```

---

## 📚 Summary

**What You Can Do Now**:
- ✅ Chain up to 10 actions sequentially
- ✅ Pass data from primary action to chain actions
- ✅ Use template syntax to format messages
- ✅ Send same data to multiple channels
- ✅ Combine any primary action with messaging chains

**Use Cases Unlocked**:
- Weather → Discord + SMS
- SWAPI → Slack + Discord
- Salesforce Query → Slack + SMS
- News → Multiple channels
- Any data source → Multiple destinations

**Your GoFlow platform now supports sophisticated multi-step workflows!** 🔗🚀

