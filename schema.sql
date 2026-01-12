-- iPaaS SQLite Schema
-- Multi-User design with path to Multi-Tenant migration

-- 1. Users Table (The current top-level entity)
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Credentials Table (Encrypted API keys/Tokens)
-- We separate this so you can manage sensitive data independently
CREATE TABLE IF NOT EXISTS credentials (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    service_name TEXT NOT NULL, -- e.g., 'slack', 'discord', 'openweather'
    encrypted_key TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 3. Workflows/Integrations Table
CREATE TABLE IF NOT EXISTS workflows (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    trigger_type TEXT NOT NULL, -- 'webhook', 'schedule'
    action_type TEXT NOT NULL,  -- 'slack_message', 'discord_post', 'weather_check' (primary action)
    config_json TEXT NOT NULL,  -- Stores params like channel IDs or thresholds
    action_chain TEXT,          -- JSON array of additional actions to execute sequentially (new!)
    is_active BOOLEAN DEFAULT 1,
    last_executed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 4. Execution Logs (Crucial for iPaaS visibility)
CREATE TABLE IF NOT EXISTS logs (
    id TEXT PRIMARY KEY,
    workflow_id TEXT NOT NULL,
    status TEXT NOT NULL, -- 'success', 'failed'
    message TEXT,
    executed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_credentials_user_id ON credentials(user_id);
CREATE INDEX IF NOT EXISTS idx_workflows_user_id ON workflows(user_id);
CREATE INDEX IF NOT EXISTS idx_workflows_trigger_type ON workflows(trigger_type);
CREATE INDEX IF NOT EXISTS idx_logs_workflow_id ON logs(workflow_id);
CREATE INDEX IF NOT EXISTS idx_logs_executed_at ON logs(executed_at);

