# Multi-Tenant Migration Strategy

This document outlines the migration path from the current **Multi-User** architecture to a production-ready **Multi-Tenant** architecture.

## Current Architecture (Multi-User)

In the current implementation, every resource (workflows, credentials, logs) is associated with a `user_id`. This is a simple approach suitable for individual users.

### Current Database Schema

```sql
users (id, email, password_hash)
credentials (id, user_id, service_name, encrypted_key)
workflows (id, user_id, name, trigger_type, action_type, config_json, is_active)
logs (id, workflow_id, status, message, executed_at)
```

### Current Query Pattern

```sql
SELECT * FROM workflows WHERE user_id = ?
```

## Target Architecture (Multi-Tenant)

In a multi-tenant architecture, resources belong to an **Organization** (tenant), and multiple users can be part of the same organization.

### Benefits of Multi-Tenant

1. **Team Collaboration**: Multiple users from the same company can manage shared workflows
2. **Enterprise Sales**: Sell to companies, not just individuals
3. **Centralized Billing**: Bill at the organization level
4. **Shared Resources**: One Slack connection can be used by all team members

## Migration Steps

### Phase 1: Add Tenants Table

Create a new `tenants` table to represent organizations:

```sql
CREATE TABLE tenants (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Add tenant_id to users table
ALTER TABLE users ADD COLUMN tenant_id TEXT;
ALTER TABLE users ADD FOREIGN KEY (tenant_id) REFERENCES tenants(id);

-- Create index for performance
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
```

### Phase 2: Migrate Existing Data

For existing users, create a tenant per user:

```sql
-- For each existing user, create a personal tenant
INSERT INTO tenants (id, name, created_at)
SELECT 
    'tenant_' || id,
    email || '''s Organization',
    created_at
FROM users;

-- Link users to their personal tenants
UPDATE users 
SET tenant_id = 'tenant_' || id;
```

### Phase 3: Update Foreign Keys

Change resource ownership from `user_id` to `tenant_id`:

```sql
-- Add tenant_id to workflows
ALTER TABLE workflows ADD COLUMN tenant_id TEXT;

-- Populate tenant_id from user_id
UPDATE workflows 
SET tenant_id = (SELECT tenant_id FROM users WHERE users.id = workflows.user_id);

-- Add foreign key constraint
ALTER TABLE workflows ADD FOREIGN KEY (tenant_id) REFERENCES tenants(id);

-- Repeat for credentials table
ALTER TABLE credentials ADD COLUMN tenant_id TEXT;
UPDATE credentials 
SET tenant_id = (SELECT tenant_id FROM users WHERE users.id = credentials.user_id);
ALTER TABLE credentials ADD FOREIGN KEY (tenant_id) REFERENCES tenants(id);

-- Create indexes
CREATE INDEX idx_workflows_tenant_id ON workflows(tenant_id);
CREATE INDEX idx_credentials_tenant_id ON credentials(tenant_id);
```

### Phase 4: Update Backend Code

#### 4.1 Update JWT Claims

```go
// TODO: MULTI-TENANT - Add tenant_id to JWT claims
claims := jwt.MapClaims{
    "user_id":  userID,
    "tenant_id": tenantID,  // Add this
    "exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
}
```

#### 4.2 Update Middleware

```go
// TODO: MULTI-TENANT - Extract tenant_id from JWT
func AuthMiddleware(next http.Handler) http.Handler {
    // ... existing code ...
    
    tenantID, ok := claims["tenant_id"].(string)
    if !ok {
        http.Error(w, "Invalid tenant_id in token", http.StatusUnauthorized)
        return
    }
    
    // Add both to context
    ctx := context.WithValue(r.Context(), UserIDKey, userID)
    ctx = context.WithValue(ctx, TenantIDKey, tenantID)
    next.ServeHTTP(w, r.WithContext(ctx))
}
```

#### 4.3 Update Repository Queries

Change all queries from:

```go
// OLD: Multi-User
SELECT * FROM workflows WHERE user_id = ?
```

To:

```go
// NEW: Multi-Tenant
SELECT * FROM workflows WHERE tenant_id = ?
```

Example:

```go
// TODO: MULTI-TENANT - Filter by tenant_id instead of user_id
func (db *Database) GetWorkflowsByTenantID(tenantID string) ([]models.Workflow, error) {
    query := `SELECT id, tenant_id, name, trigger_type, action_type, config_json, is_active, last_executed_at, created_at 
              FROM workflows WHERE tenant_id = ? ORDER BY created_at DESC`
    // ... rest of implementation
}
```

### Phase 5: Update Frontend

#### 5.1 Add Organization Switcher

Add a UI component to switch between organizations (for users who belong to multiple tenants):

```tsx
// TODO: MULTI-TENANT - Add organization switcher
<OrganizationSwitcher organizations={userOrganizations} />
```

#### 5.2 Update API Client

```typescript
// TODO: MULTI-TENANT - Include tenant_id in API headers
const headers: HeadersInit = {
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${token}`,
  'X-Tenant-ID': tenantID,  // Add this
};
```

### Phase 6: Add Tenant Management

Create new endpoints for tenant management:

- `POST /api/tenants` - Create new organization
- `POST /api/tenants/:id/invite` - Invite users to organization
- `GET /api/tenants/:id/members` - List organization members
- `PUT /api/tenants/:id/members/:userId/role` - Update member role

## Rollout Strategy

1. **Deploy Schema Changes**: Run migration scripts during maintenance window
2. **Deploy Backend Changes**: Update backend with tenant-aware code
3. **Deploy Frontend Changes**: Update frontend with tenant context
4. **Gradual Rollout**: Enable multi-tenant features for new signups first
5. **Migration Tool**: Provide a tool for existing users to create/join organizations

## Code Locations to Update

All locations marked with `// TODO: MULTI-TENANT` comments:

### Backend
- `internal/middleware/auth.go` - JWT extraction
- `internal/handlers/auth.go` - JWT generation
- `internal/db/database.go` - All repository queries
- `internal/handlers/*.go` - All handler user_id extraction

### Frontend
- `frontend/lib/api.ts` - API client headers
- All dashboard pages - Display tenant context

## Testing Multi-Tenant

1. Create two test organizations
2. Add users to each organization
3. Create workflows in each organization
4. Verify data isolation (org A cannot see org B's data)
5. Test cross-tenant scenarios (user belongs to multiple orgs)

## Backward Compatibility

The migration is designed to be backward compatible:
- Existing users get their own "personal" tenant
- Existing workflows continue to work
- No data loss during migration
- Users can later join/create new organizations

## Production Considerations

1. **Tenant Limits**: Add quotas per tenant (max workflows, max executions/month)
2. **Billing**: Implement tenant-level billing
3. **Analytics**: Track usage per tenant
4. **Security**: Ensure tenant isolation at all levels
5. **Performance**: Add proper indexes on tenant_id columns

