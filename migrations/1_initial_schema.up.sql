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
    web_user_id INTEGER NOT NULL DEFAULT 0,
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
    is_bot numeric default false not null,
    chatbot_permit TINYINT NOT NULL DEFAULT 0,
    srvs_employees_id INTEGER NOT NULL DEFAULT 0,
    tg_cb_flow_step_id INTEGER NOT NULL DEFAULT 0,
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
    date INTEGER NOT NULL DEFAULT 0,
    added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- INSERT INTO web_users (username, password, role) VALUES ('admin', 'admin123', 'admin');

CREATE TABLE db_changes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    web_user_id INTEGER NOT NULL DEFAULT 0,
    model_name VARCHAR(64) NOT NULL,
    data_from TEXT NOT NULL,
    data_to TEXT NOT NULL,
    added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (web_user_id) REFERENCES web_users(id) ON DELETE CASCADE
);

CREATE TABLE tg_cb_flow_list (
                                 id INTEGER PRIMARY KEY AUTOINCREMENT,
                                 name VARCHAR(255) NOT NULL,
                                 row_order INTEGER NOT NULL DEFAULT 0,
                                 added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

CREATE TABLE tg_cb_flow_steps (
                            id INTEGER PRIMARY KEY AUTOINCREMENT,
                            tg_cb_flow_id INTEGER NOT NULL DEFAULT 0,
                            msg TEXT NOT NULL,
                            keyboard VARCHAR(255) NOT NULL,
                            handler_name VARCHAR(255) NOT NULL,
                            row_order INTEGER NOT NULL DEFAULT 0,
                            added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

INSERT INTO tg_cb_flow_list (id, name, row_order, added_timestamp) VALUES (1, 'Shift start', 10, 1729539190);
INSERT INTO tg_cb_flow_list (id, name, row_order, added_timestamp) VALUES (2, 'Sale', 20, 1729539318);
INSERT INTO tg_cb_flow_list (id, name, row_order, added_timestamp) VALUES (3, 'Shift close', 30, 1729539433);
INSERT INTO tg_cb_flow_list (id, name, row_order, added_timestamp) VALUES (4, 'Initial menu', 40, 1729739101);

INSERT INTO tg_cb_flow_steps (id, tg_cb_flow_id, msg, keyboard, handler_name, row_order, added_timestamp) VALUES (1, 1, 'Choose your location', 'func:getLocationsKeyboard', 'handleUserChooseLocation', 10, 1729727117);
INSERT INTO tg_cb_flow_steps (id, tg_cb_flow_id, msg, keyboard, handler_name, row_order, added_timestamp) VALUES (2, 1, 'Enter the remainder of Frame A', '', 'handleRemainderProduct(FrameA)', 20, 1729732871);
INSERT INTO tg_cb_flow_steps (id, tg_cb_flow_id, msg, keyboard, handler_name, row_order, added_timestamp) VALUES (3, 1, 'Enter the remainder of Frame B', '', 'handleRemainderProduct(FrameB)', 30, 1729733731);
INSERT INTO tg_cb_flow_steps (id, tg_cb_flow_id, msg, keyboard, handler_name, row_order, added_timestamp) VALUES (4, 1, 'Enter the remainder of Paper', '', 'handleRemainderProduct(Paper)', 40, 1729735056);
INSERT INTO tg_cb_flow_steps (id, tg_cb_flow_id, msg, keyboard, handler_name, row_order, added_timestamp) VALUES (5, 4, 'Please tap menu button', 'Start shift', '', 10, 1729739392);
