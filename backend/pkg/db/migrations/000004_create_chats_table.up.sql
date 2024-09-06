CREATE TABLE IF NOT EXISTS chats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    recipient_id INTEGER,
    group_id INTEGER,
    message TEXT NOT NULL,
    is_group BOOLEAN NOT NULL,
    created_at DATETIME,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (recipient_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);
