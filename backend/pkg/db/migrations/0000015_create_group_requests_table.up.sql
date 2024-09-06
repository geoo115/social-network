CREATE TABLE IF NOT EXISTS group_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    requester_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    requested_at DATETIME NOT NULL,
    responded_at DATETIME,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (requester_id) REFERENCES users(id)
);
