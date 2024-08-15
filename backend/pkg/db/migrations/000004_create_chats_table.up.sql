-- 000004_create_chats_table.up.sql
CREATE TABLE chats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    recipient_id INTEGER,
    group_id INTEGER,
    message TEXT NOT NULL,
    is_group BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (recipient_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

-- Indexes
CREATE INDEX idx_chats_sender_id ON chats(sender_id);
CREATE INDEX idx_chats_recipient_id ON chats(recipient_id);
CREATE INDEX idx_chats_group_id ON chats(group_id);
