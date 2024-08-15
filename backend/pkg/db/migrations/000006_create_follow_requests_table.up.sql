-- 000006_create_follow_requests_table.up.sql
CREATE TABLE follow_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    recipient_id INTEGER NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (recipient_id) REFERENCES users(id)
);

-- Indexes
CREATE INDEX idx_follow_requests_sender_id ON follow_requests(sender_id);
CREATE INDEX idx_follow_requests_recipient_id ON follow_requests(recipient_id);
CREATE INDEX idx_follow_requests_status ON follow_requests(status);
