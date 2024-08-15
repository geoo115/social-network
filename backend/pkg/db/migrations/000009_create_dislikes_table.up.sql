-- 000009_create_dislikes_table.up.sql
CREATE TABLE dislikes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

-- Indexes
CREATE INDEX idx_dislikes_user_id ON dislikes(user_id);
CREATE INDEX idx_dislikes_post_id ON dislikes(post_id);
