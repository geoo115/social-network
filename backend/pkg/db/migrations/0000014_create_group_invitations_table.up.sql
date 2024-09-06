CREATE TABLE IF NOT EXISTS group_invitations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    inviter_id INTEGER NOT NULL,
    invitee_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    invited_at DATETIME NOT NULL,
    responded_at DATETIME,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    FOREIGN KEY (inviter_id) REFERENCES users(id),
    FOREIGN KEY (invitee_id) REFERENCES users(id)
);
