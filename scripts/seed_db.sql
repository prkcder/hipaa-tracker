CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    sanitized BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS web_tracker_scans (
    id SERIAL PRIMARY KEY,
    scanned_url TEXT NOT NULL,
    tracker_name TEXT NOT NULL,
    tracker_domain TEXT NOT NULL,
    risk_level TEXT NOT NULL,
    page_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);