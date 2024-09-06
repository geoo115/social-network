CREATE TABLE IF NOT EXISTS group_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    day_time DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);
