-- 000003_create_groups_table.up.sql
CREATE TABLE groups (
      id INT PRIMARY KEY AUTOINCREMENT,
    creator_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

-- Indexes
CREATE INDEX idx_groups_creator_id ON groups(creator_id);
