-- Create river_jobs table
CREATE TABLE river_jobs (
    id BIGSERIAL PRIMARY KEY,
    args JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    kind TEXT NOT NULL,
    max_attempts INTEGER NOT NULL,
    priority INTEGER NOT NULL,
    queue TEXT NOT NULL,
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    state TEXT NOT NULL,
    tags TEXT[] NOT NULL DEFAULT '{}',
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    attempt INTEGER NOT NULL DEFAULT 0,
    attempted_at TIMESTAMP WITH TIME ZONE,
    attempted_by TEXT,
    discarded_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    error_backtrace TEXT,
    error_count INTEGER NOT NULL DEFAULT 0,
    finished_at TIMESTAMP WITH TIME ZONE
);

-- Create river_jobs indices
CREATE INDEX river_jobs_scheduled_at_idx ON river_jobs (scheduled_at) WHERE state = 'available';
CREATE INDEX river_jobs_state_idx ON river_jobs (state);
CREATE INDEX river_jobs_kind_idx ON river_jobs (kind);
CREATE INDEX river_jobs_queue_idx ON river_jobs (queue);
CREATE INDEX river_jobs_priority_scheduled_at_id_idx ON river_jobs (priority DESC, scheduled_at, id) WHERE state = 'available';

-- Create river_job_archival table for job archival
CREATE TABLE river_job_archival_jobs (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    job_id BIGINT NOT NULL,
    unique_job_key TEXT NOT NULL
);

-- Create river_leader_elections table for leader election
CREATE TABLE river_leader_elections (
    id TEXT PRIMARY KEY,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    leader TEXT NOT NULL
);

-- Create river_periodic_job_executions table for tracking periodic job executions
CREATE TABLE river_periodic_job_executions (
    id SERIAL PRIMARY KEY,
    job_kind TEXT NOT NULL,
    job_queue TEXT NOT NULL,
    schedule_name TEXT NOT NULL,
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    triggered_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE (job_kind, schedule_name, triggered_at)
);

-- Create river_unique_job_reservations table for unique job tracking
CREATE TABLE river_unique_job_reservations (
    key TEXT PRIMARY KEY,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    job_id BIGINT NOT NULL
);

-- Create river_periodic_schedules table for periodic scheduling
CREATE TABLE river_periodic_schedules (
    id SERIAL PRIMARY KEY,
    job_args JSONB NOT NULL,
    job_kind TEXT NOT NULL,
    job_max_attempts INTEGER NOT NULL,
    job_priority INTEGER NOT NULL,
    job_queue TEXT NOT NULL,
    job_tags TEXT[] NOT NULL DEFAULT '{}',
    last_triggered_at TIMESTAMP WITH TIME ZONE,
    name TEXT NOT NULL,
    period INTERVAL NOT NULL,
    timeout INTERVAL NOT NULL DEFAULT INTERVAL '0',
    enabled BOOLEAN NOT NULL DEFAULT true,
    UNIQUE (job_kind, name)
);