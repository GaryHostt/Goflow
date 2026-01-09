# Quick Start Guide

Welcome to your iPaaS platform! Follow these steps to get started.

## Prerequisites

Before starting, ensure you have:
- **Go 1.21+** installed ([Download](https://go.dev/dl/))
- **Node.js 18+** installed ([Download](https://nodejs.org/))

## Option 1: Automated Start (Recommended)

Run the startup script:

```bash
chmod +x start.sh
./start.sh
```

This will:
1. Check dependencies
2. Install Go modules
3. Build the backend
4. Install npm packages
5. Start both backend and frontend
6. Open your browser to http://localhost:3000

## Option 2: Manual Start

### Terminal 1 - Backend

```bash
# Install dependencies
go mod download

# Run the backend
go run cmd/api/main.go
```

Backend will be available at http://localhost:8080

### Terminal 2 - Frontend

```bash
# Navigate to frontend
cd frontend

# Install dependencies (first time only)
npm install

# Start development server
npm run dev
```

Frontend will be available at http://localhost:3000

## First Steps

1. **Register** a new account at http://localhost:3000/register
2. **Connect services** (Slack, Discord, or OpenWeather)
3. **Create a workflow**
4. **Test it** and view logs!

## Test with Demo Data

Generate sample workflows and logs:

```bash
go run scripts/generate_test_data.go
```

Then login with:
- Email: `demo@ipaas.com`
- Password: `password123`

## Example Workflow: Webhook to Slack

1. Get a Slack incoming webhook URL from https://api.slack.com/messaging/webhooks
2. Go to **Connections** and add your Slack webhook
3. Go to **Workflows** > **Create Workflow**
4. Configure:
   - Name: "Test Webhook"
   - Trigger: Webhook
   - Action: Send Slack Message
   - Message: "Hello from iPaaS!"
5. Copy the generated webhook URL
6. Test it:
   ```bash
   curl -X POST http://localhost:8080/api/webhooks/{workflow_id}
   ```
7. Check Slack and the **Logs** page!

## Project Structure

```
simple-ipass/
├── cmd/api/main.go         # Backend entry point
├── internal/               # Backend code
├── frontend/               # Next.js frontend
├── schema.sql             # Database schema
├── README.md              # Full documentation
├── MIGRATION.md           # Multi-tenant guide
└── start.sh               # Startup script
```

## Need Help?

- Read the full [README.md](README.md)
- Check [MIGRATION.md](MIGRATION.md) for multi-tenant architecture
- Review code comments marked with `// TODO: MULTI-TENANT`

## Stop the Platform

Press `Ctrl+C` in the terminal running the services.

Or manually:
```bash
pkill -f "cmd/api/main.go"
pkill -f "next dev"
```

Happy integrating! 🚀

