CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    type TEXT,
    message TEXT,
    is_read BOOLEAN,
    created_at DATETIME,
    details TEXT
);
