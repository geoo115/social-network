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
