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

INSERT INTO tg_cb_flow_list (id, name, row_order, added_timestamp) VALUES (1, 'Shift start', 10, 1729539190);
INSERT INTO tg_cb_flow_list (id, name, row_order, added_timestamp) VALUES (2, 'Sale', 20, 1729539318);
INSERT INTO tg_cb_flow_list (id, name, row_order, added_timestamp) VALUES (3, 'Shift close', 30, 1729539433);
