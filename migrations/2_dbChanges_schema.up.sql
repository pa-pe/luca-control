CREATE TABLE db_changes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    web_user_id INTEGER NOT NULL,
    model_name VARCHAR(64) NOT NULL,
    data_from TEXT NOT NULL,
    data_to TEXT NOT NULL,
    added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (web_user_id) REFERENCES web_users(id) ON DELETE CASCADE
);
