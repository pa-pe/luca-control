CREATE TABLE tg_cb_flow_list (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    row_order INTEGER NOT NULL,
    added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

CREATE TABLE tg_cb_flow (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tg_cb_flow_id INTEGER NOT NULL,
    msg TEXT NOT NULL,
    handler_mame VARCHAR(255) NOT NULL,
    row_order INTEGER NOT NULL,
    added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);
