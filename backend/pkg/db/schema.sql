CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				email TEXT UNIQUE NOT NULL,
				password VARCHAR(255) NULL,
				first_name TEXT,
				last_name TEXT,
				date_of_birth TEXT,
				avatar TEXT,
				nickname TEXT,
				about_me TEXT,
				provider TEXT,
				is_private BOOLEAN DEFAULT FALSE,
				created_at DATETIME NOT NULL,
				updated_at DATETIME NOT NULL
			);


CREATE TABLE IF NOT EXISTS posts (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				content TEXT NOT NULL,
				image TEXT,
				privacy TEXT NOT NULL,
				created_at DATETIME NOT NULL,
				updated_at DATETIME NOT NULL,
				FOREIGN KEY (user_id) REFERENCES users(id)
			);


CREATE TABLE IF NOT EXISTS followers (
				follower_id INTEGER NOT NULL,
				followed_id INTEGER NOT NULL,
				PRIMARY KEY (follower_id, followed_id),
				FOREIGN KEY (follower_id) REFERENCES users(id),
				FOREIGN KEY (followed_id) REFERENCES users(id)
			);


CREATE TABLE IF NOT EXISTS groups (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				creator_id INTEGER NOT NULL,
				title TEXT NOT NULL,
				description TEXT,
				created_at DATETIME NOT NULL,
				updated_at DATETIME NOT NULL
			);


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

CREATE TABLE IF NOT EXISTS group_requests (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				group_id INTEGER NOT NULL,
				requester_id INTEGER NOT NULL,
				status TEXT NOT NULL DEFAULT 'pending',
				requested_at DATETIME NOT NULL,
				responded_at DATETIME,
				FOREIGN KEY (group_id) REFERENCES groups(id),
				FOREIGN KEY (requester_id) REFERENCES users(id)
			);

CREATE TABLE IF NOT EXISTS group_events (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				group_id INTEGER NOT NULL,
				title TEXT NOT NULL,
				description TEXT,
				day_time DATETIME NOT NULL,
				created_at DATETIME NOT NULL,
				updated_at DATETIME NOT NULL,
				FOREIGN KEY (group_id) REFERENCES groups(id)
			);


CREATE TABLE IF NOT EXISTS event_rsvps (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				event_id INTEGER NOT NULL,
				user_id INTEGER NOT NULL,
				status TEXT NOT NULL, -- "going" or "not going"
				responded_at DATETIME NOT NULL,
				FOREIGN KEY (event_id) REFERENCES group_events(id),
				FOREIGN KEY (user_id) REFERENCES users(id)
			);


CREATE TABLE IF NOT EXISTS group_memberships (
				user_id INTEGER NOT NULL,
				group_id INTEGER NOT NULL,
				joined_at DATETIME NOT NULL,
				left_at DATETIME,
				PRIMARY KEY (user_id, group_id),
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (group_id) REFERENCES groups(id)
			);


CREATE TABLE IF NOT EXISTS chats (
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
			;


CREATE TABLE IF NOT EXISTS notifications (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER,
				type TEXT,
				message TEXT,
				is_read BOOLEAN,
				created_at DATETIME,
				details TEXT
			);


CREATE TABLE IF NOT EXISTS follow_requests (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				sender_id INTEGER NOT NULL,
				recipient_id INTEGER NOT NULL,
				status TEXT NOT NULL,
				created_at DATETIME,
				FOREIGN KEY (sender_id) REFERENCES users(id),
				FOREIGN KEY (recipient_id) REFERENCES users(id)
			);


CREATE TABLE IF NOT EXISTS likes (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				post_id INTEGER NOT NULL,
				created_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (post_id) REFERENCES posts(id)
			);


CREATE TABLE IF NOT EXISTS dislikes (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				post_id INTEGER NOT NULL,
				created_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (post_id) REFERENCES posts(id)
			);


CREATE TABLE IF NOT EXISTS sessions (
				session_id TEXT PRIMARY KEY,
				user_id INTEGER NOT NULL,
				expires_at DATETIME NOT NULL,
				FOREIGN KEY (user_id) REFERENCES users(id)
			);


CREATE TABLE IF NOT EXISTS comments (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				user_id INTEGER NOT NULL,
				post_id INTEGER NOT NULL,
				content TEXT NOT NULL,
				created_at DATETIME,
				updated_at DATETIME,
				FOREIGN KEY (user_id) REFERENCES users(id),
				FOREIGN KEY (post_id) REFERENCES posts(id)
			);
