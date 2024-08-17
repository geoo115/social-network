-- 000005_create_notifications_table.up.sql
CREATE TABLE IF NOT EXISTS notifications (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER,
	type TEXT,
	message TEXT,
	is_read BOOLEAN,
	created_at DATETIME,
	details TEXT
)
-- Indexes
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
