CREATE TABLE IF NOT EXISTS group_memberships (
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    joined_at DATETIME NOT NULL,
    left_at DATETIME,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);
