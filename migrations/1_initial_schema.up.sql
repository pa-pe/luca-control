-- 1_initial_schema.up.sql
CREATE TABLE web_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL,
    added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))

);

CREATE TABLE web_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    web_user_id INTEGER NOT NULL,
    session_key VARCHAR(64) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (web_user_id) REFERENCES web_users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tg_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name VARCHAR(50) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    language_code CHAR(2) NOT NULL,
    chatbot_permit TINYINT NOT NULL DEFAULT 0,
    chatbot_state VARCHAR(100) NOT NULL,
    shift_state TINYINT NOT NULL DEFAULT 0,
    added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

CREATE TABLE IF NOT EXISTS tg_msgs (
    internal_id INTEGER PRIMARY KEY AUTOINCREMENT,
    tg_id INTEGER NOT NULL DEFAULT 0,
    tg_user_id INTEGER NOT NULL DEFAULT 0,
    chat_id INTEGER NOT NULL DEFAULT 0,
    reply_to_message_id INTEGER NOT NULL DEFAULT 0,
    is_outgoing TINYINT NOT NULL DEFAULT 0,
    text TEXT NOT NULL,
    date INTEGER DEFAULT 0,
    added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- INSERT INTO web_users (username, password, role) VALUES ('admin', 'admin123', 'admin');
