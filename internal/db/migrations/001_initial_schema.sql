CREATE TABLE metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    value REAL NOT NULL,
    timestamp DATETIME NOT NULL,
    tags TEXT NOT NULL  -- JSON string of tags
);

CREATE INDEX idx_metrics_name_timestamp ON metrics(name, timestamp);

CREATE TABLE alerts (
    id TEXT PRIMARY KEY,
    metric_name TEXT NOT NULL,
    threshold REAL NOT NULL,
    operator TEXT NOT NULL,
    duration INTEGER NOT NULL,  -- in seconds
    tags TEXT NOT NULL,        -- JSON string of tags
    is_triggered BOOLEAN NOT NULL DEFAULT FALSE,
    last_check DATETIME
);

CREATE TABLE alert_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alert_id TEXT NOT NULL,
    metric_name TEXT NOT NULL,
    value REAL NOT NULL,
    threshold REAL NOT NULL,
    timestamp DATETIME NOT NULL,
    FOREIGN KEY (alert_id) REFERENCES alerts(id)
);

CREATE INDEX idx_alert_history_timestamp ON alert_history(timestamp);
CREATE INDEX idx_alert_history_alert_id ON alert_history(alert_id);