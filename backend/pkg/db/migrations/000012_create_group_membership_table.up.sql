CREATE TABLE group_memberships (
   user_id INT NOT NULL,
    group_id INT NOT NULL,
    joined_at DATETIME NOT NULL,
    left_at DATETIME,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_group_memberships_user ON group_memberships(user_id);
CREATE INDEX idx_group_memberships_group ON group_memberships(group_id);
CREATE INDEX idx_group_memberships_user_group ON group_memberships(user_id, group_id);
