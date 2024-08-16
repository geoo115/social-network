-- 000007_create_group_events_table.up.sql
CREATE TABLE group_events (
     id INT PRIMARY KEY AUTOINCREMENT,
    group_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    day_time DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

-- Indexes
CREATE INDEX idx_group_events_group_id ON group_events(group_id);
