package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Initialize the SQLite connection and apply migrations
func Initialize() error {
	var err error
	DB, err = sql.Open("sqlite3", "./socialNetwork.db")
	if err != nil {
		return err
	}

	// Apply migrations (create tables)
	err = applyMigrations(DB)
	if err != nil {
		return err
	}

	return nil
}

func applyMigrations(db *sql.DB) error {
	tables := []struct {
		name   string
		create string
	}{
		{
			name: "users",
			create: `CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				email TEXT UNIQUE NOT NULL,
				password TEXT NOT NULL,
				first_name TEXT NOT NULL,
				last_name TEXT NOT NULL,
				date_of_birth TEXT NOT NULL,
				avatar TEXT,
				nickname TEXT,
				about_me TEXT,
				created_at DATETIME,
				updated_at DATETIME
			)`,
		},
		{
			name: "posts",
			create: `CREATE TABLE IF NOT EXISTS posts (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				content TEXT NOT NULL,
				image TEXT,
				privacy TEXT NOT NULL,
				created_at DATETIME,
				updated_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id)
			)`,
		},
		{
			name: "groups",
			create: `CREATE TABLE IF NOT EXISTS groups (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				creator_id INTEGER NOT NULL,
				title TEXT NOT NULL,
				description TEXT,
				created_at DATETIME,
				updated_at DATETIME,
				FOREIGN KEY (creator_id) REFERENCES users(id)
			)`,
		},
		{
			name: "chats",
			create: `CREATE TABLE IF NOT EXISTS chats (
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
			)`,
		},
		{
			name: "notifications",
			create: `CREATE TABLE IF NOT EXISTS notifications (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				message TEXT NOT NULL,
				is_read BOOLEAN NOT NULL,
				created_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id)
			)`,
		},
		{
			name: "follow_requests",
			create: `CREATE TABLE IF NOT EXISTS follow_requests (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				sender_id INTEGER NOT NULL,
				recipient_id INTEGER NOT NULL,
				status TEXT NOT NULL,
				created_at DATETIME,
				FOREIGN KEY (sender_id) REFERENCES users(id),
				FOREIGN KEY (recipient_id) REFERENCES users(id)
			)`,
		},
		{
			name: "group_events",
			create: `CREATE TABLE IF NOT EXISTS group_events (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				group_id INTEGER NOT NULL,
				title TEXT NOT NULL,
				description TEXT,
				day_time DATETIME NOT NULL,
				created_at DATETIME,
				updated_at DATETIME,
				FOREIGN KEY (group_id) REFERENCES groups(id)
			)`,
		},
		{
			name: "group_memberships",
			create: `CREATE TABLE IF NOT EXISTS group_memberships (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				group_id INTEGER NOT NULL,
				joined_at DATETIME,
				left_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (group_id) REFERENCES groups(id)
			)`,
		},
		{
			name: "likes",
			create: `CREATE TABLE IF NOT EXISTS likes (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				post_id INTEGER NOT NULL,
				created_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (post_id) REFERENCES posts(id)
			)`,
		},
		{
			name: "dislikes",
			create: `CREATE TABLE IF NOT EXISTS dislikes (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				post_id INTEGER NOT NULL,
				created_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (post_id) REFERENCES posts(id)
			)`,
		},
		{
			name: "sessions",
			create: `CREATE TABLE IF NOT EXISTS sessions (
				session_id TEXT PRIMARY KEY,
				user_id INTEGER NOT NULL,
				expires_at DATETIME NOT NULL,
				FOREIGN KEY (user_id) REFERENCES users(id)
			)`,
		},
		{
			name: "comments",
			create: `CREATE TABLE IF NOT EXISTS comments (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				post_id INTEGER NOT NULL,
				content TEXT NOT NULL,
				created_at DATETIME,
				updated_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (post_id) REFERENCES posts(id)
			)`,
		},
	}

	for _, table := range tables {
		if _, err := db.Exec(table.create); err != nil {
			return err
		}
	}

	return nil
}
