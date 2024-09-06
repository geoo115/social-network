CREATE TABLE IF NOT EXISTS event_rsvps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status TEXT NOT NULL, -- "going" or "not going"
    responded_at DATETIME NOT NULL,
    FOREIGN KEY (event_id) REFERENCES group_events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
