CREATE TABLE IF NOT EXISTS monitors (
    monitor_id TEXT PRIMARY KEY,
    url TEXT NOT NULL,
    error_condition TEXT NOT NULL
);
