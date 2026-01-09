# Simple iPaaS - Integration Platform as a Service

A full-stack iPaaS (Integration Platform as a Service) built with **Go** backend and **Next.js** frontend. This is a POC demonstrating core iPaaS concepts including webhook triggers, scheduled tasks, third-party connectors, and multi-user architecture with a clear path to multi-tenant migration.

## Features

### Core Functionality
- ✅ **User Authentication** - JWT-based auth with register/login
- ✅ **Workflow Management** - Create, enable/disable, delete workflows
- ✅ **Multiple Triggers** - Webhook and scheduled (polling) triggers
- ✅ **Third-Party Connectors** - Slack, Discord, OpenWeather API
- ✅ **Execution Logs** - Track all workflow executions with filtering
- ✅ **Encrypted Credentials** - AES-256 encryption for API keys
- ✅ **Background Scheduler** - Goroutine-based polling for scheduled tasks

### Architecture
- **Backend**: Go with gorilla/mux router, SQLite database
- **Frontend**: Next.js 14 with App Router, Tailwind CSS, Shadcn/UI
- **Database**: SQLite with multi-user design (ready for multi-tenant)
- **Concurrency**: Go routines for async workflow execution

## Tech Stack

### Backend
- Go 1.21+
- gorilla/mux - HTTP routing
- mattn/go-sqlite3 - SQLite driver
- golang-jwt/jwt - JWT authentication
- golang.org/x/crypto - Password hashing & encryption

### Frontend
- Next.js 14 (React 18)
- TypeScript
- Tailwind CSS
- Shadcn/UI components

## Project Structure

```
simple-ipass/
├── cmd/api/main.go              # Backend entry point
├── internal/
│   ├── db/database.go           # Database layer with repositories
│   ├── models/models.go         # Data models
│   ├── middleware/auth.go       # JWT authentication middleware
│   ├── handlers/                # HTTP request handlers
│   │   ├── auth.go
│   │   ├── credentials.go
│   │   ├── workflows.go
│   │   ├── webhooks.go
│   │   └── logs.go
│   ├── engine/                  # Execution engine
│   │   ├── executor.go          # Workflow execution logic
│   │   ├── scheduler.go         # Background scheduler
│   │   └── connectors/          # Third-party integrations
│   │       ├── slack.go
│   │       ├── discord.go
│   │       └── openweather.go
│   └── crypto/encrypt.go        # Encryption utilities
├── frontend/                    # Next.js frontend
│   ├── app/                     # App router pages
│   │   ├── login/
│   │   ├── register/
│   │   └── dashboard/           # Protected dashboard area
│   ├── components/              # Reusable components
│   └── lib/                     # Utilities and API client
├── scripts/
│   └── generate_test_data.go   # Test data generator
├── schema.sql                   # Database schema
├── MIGRATION.md                 # Multi-tenant migration guide
└── README.md                    # This file
```

## Installation & Setup

### Prerequisites
- Go 1.21 or higher
- Node.js 18+ and npm
- Git

### Backend Setup

1. **Install Go** (if not already installed):
   ```bash
   # macOS
   brew install go
   
   # Or download from https://go.dev/dl/
   ```

2. **Initialize and build the backend**:
   ```bash
   cd simple-ipass
   go mod download
   go build -o bin/api cmd/api/main.go
   ```

3. **Run the backend**:
   ```bash
   ./bin/api
   # Server will start on http://localhost:8080
   ```

   Or run directly:
   ```bash
   go run cmd/api/main.go
   ```

### Frontend Setup

1. **Install dependencies**:
   ```bash
   cd frontend
   npm install
   ```

2. **Run the development server**:
   ```bash
   npm run dev
   # Frontend will start on http://localhost:3000
   ```

### Generate Test Data (Optional)

```bash
go run scripts/generate_test_data.go
```

This creates a demo user and sample workflows:
- **Email**: demo@ipaas.com
- **Password**: password123

## Usage Guide

### 1. Register/Login
- Navigate to `http://localhost:3000`
- Create an account or use the demo credentials

### 2. Connect Services
- Go to **Connections** page
- Add your Slack webhook URL, Discord webhook, or OpenWeather API key

#### Getting Credentials:
- **Slack**: Create an incoming webhook at https://api.slack.com/messaging/webhooks
- **Discord**: Create a webhook in Server Settings > Integrations
- **OpenWeather**: Get a free API key at https://openweathermap.org/api

### 3. Create a Workflow

Example workflows:

**Webhook to Slack**
- Trigger: Webhook
- Action: Send Slack Message
- Result: Unique URL like `http://localhost:8080/api/webhooks/{workflow_id}`

**Scheduled Weather Check**
- Trigger: Schedule (every 10 minutes)
- Action: Check Weather
- Config: City name

### 4. Test Your Workflow

For webhook triggers:
```bash
curl -X POST http://localhost:8080/api/webhooks/{workflow_id}
```

For scheduled workflows, they run automatically based on the interval.

### 5. View Logs
- Go to **Logs** page
- Filter by success/failed status
- See execution history with timestamps

## API Endpoints

### Public Routes
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login and get JWT token
- `POST /api/webhooks/:id` - Trigger workflow via webhook

### Protected Routes (require JWT)
- `POST /api/credentials` - Save encrypted credentials
- `GET /api/credentials` - List user's credentials
- `POST /api/workflows` - Create workflow
- `GET /api/workflows` - List user's workflows
- `PUT /api/workflows/:id/toggle` - Enable/disable workflow
- `DELETE /api/workflows/:id` - Delete workflow
- `GET /api/logs` - Get execution logs

## Multi-Tenant Migration

This project is designed with a **multi-user** architecture that's ready to migrate to **multi-tenant**. See [MIGRATION.md](MIGRATION.md) for the complete migration strategy.

### Key Design Decisions for Migration Path:

1. All database queries filter by `user_id` (easily changed to `tenant_id`)
2. JWT middleware extracts user context (can add tenant context)
3. Code comments mark multi-tenant preparation points with `// TODO: MULTI-TENANT`

## Security Considerations

⚠️ **This is a POC/Learning Project**. For production use:

1. **Environment Variables**: Move secrets to environment variables
   - JWT secret
   - Encryption key
   - Database path

2. **HTTPS**: Use TLS/SSL in production

3. **Rate Limiting**: Add rate limiting to API endpoints

4. **Input Validation**: Enhance validation and sanitization

5. **Secrets Management**: Use proper secrets management (AWS Secrets Manager, HashiCorp Vault)

6. **Database**: Consider PostgreSQL for production

## Learning Outcomes

By building this iPaaS, you'll learn:

1. **Backend Patterns**
   - Repository pattern for data access
   - Middleware for authentication
   - Concurrent processing with goroutines
   - Background job scheduling

2. **Integration Concepts**
   - Webhook handling
   - Third-party API integration
   - Credential encryption
   - Error handling and logging

3. **Frontend Architecture**
   - Protected routes
   - JWT token management
   - API client abstraction
   - Modern UI with Tailwind + Shadcn

4. **Product Ownership**
   - Multi-user to multi-tenant evolution
   - Feature prioritization (simple connectors first)
   - Observability (logging and monitoring)
   - User experience design

## Extending the Platform

### Add New Connectors

1. Create connector in `internal/engine/connectors/`:
```go
type NewService struct {
    APIKey string
}

func (n *NewService) Execute() connectors.Result {
    // Your integration logic
}
```

2. Add case in `internal/engine/executor.go`
3. Add to frontend workflow builder

### Add More Trigger Types
- Cron expressions
- Event-based triggers (email, file upload)
- Conditional triggers (if/then logic)

### Add ELK Stack (Elasticsearch, Logstash, Kibana)

See the original plan for ELK integration with Docker Compose for advanced logging and analytics.

## Troubleshooting

### Backend won't start
- Ensure port 8080 is not in use
- Check that `schema.sql` exists in project root
- Verify Go dependencies are installed: `go mod download`

### Frontend won't connect to backend
- Check `NEXT_PUBLIC_API_URL` in `frontend/lib/api.ts`
- Ensure backend is running on port 8080
- Check browser console for CORS errors

### Database issues
- Delete `ipaas.db` and restart to reset
- Check file permissions on database file

## Contributing

This is a learning project, but contributions are welcome! Areas for improvement:
- Additional connectors (Google Sheets, GitHub, etc.)
- OAuth2 support for connectors
- Workflow templates
- Better error handling
- Unit and integration tests

## License

MIT License - Feel free to use this for learning and personal projects.

## Acknowledgments

Built as a comprehensive learning project to master:
- Backend systems (Go)
- Integration patterns (iPaaS)
- Modern frontend (Next.js)
- Product ownership and architecture

---

**Ready to build?** Start with `go run cmd/api/main.go` and `npm run dev` in the frontend folder!

