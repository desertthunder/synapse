CREATE TABLE tasks (
    id INTEGER PRIMARY KEY,
    -- The actual command or script to be executed, corresponds to a map key
    name VARCHAR(255) NOT NULL,
    -- Lower number => higher priority
    -- Default of 5 represents medium priority
    priority INTEGER NOT NULL DEFAULT 5,
    scheduled_time DATETIME,
    retry_delay_seconds INTEGER NOT NULL DEFAULT 45,
    -- Options: 'pending', 'running', 'completed', 'failed', 'cancelled'
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at DATETIME,
    completed_at DATETIME,
);

CREATE INDEX idx_tasks_status_priority ON tasks (status, priority);

CREATE INDEX idx_tasks_created_at ON tasks (created_at);

CREATE INDEX idx_tasks_name ON tasks (task_name);
