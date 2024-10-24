CREATE TABLE srvs_location_list (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    row_order INTEGER NOT NULL DEFAULT 0,
    added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

INSERT INTO srvs_location_list (name, row_order) VALUES ('Islander', '10');
INSERT INTO srvs_location_list (name, row_order) VALUES ('Johnny', '20');


CREATE TABLE srvs_employees_list (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    percentage           REAL NOT NULL,
    srvs_shift_id INTEGER NOT NULL DEFAULT 0,
    added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

INSERT INTO srvs_employees_list (name, percentage) VALUES ('Test', 0.3);


CREATE TABLE srvs_goods_list (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    price           REAL NOT NULL,
    added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

INSERT INTO srvs_goods_list (name, price) VALUES ('Frame A', 15);
INSERT INTO srvs_goods_list (name, price) VALUES ('Frame B', 15);
INSERT INTO srvs_goods_list (name, price) VALUES ('Paper', 110);
INSERT INTO srvs_goods_list (name, price) VALUES ('Single A', 20);
INSERT INTO srvs_goods_list (name, price) VALUES ('Single B', 20);

CREATE TABLE srvs_leftovers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    srvs_location_id INTEGER NOT NULL DEFAULT 0,
    date BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
    srvs_goods_id INTEGER NOT NULL DEFAULT 0,
    srvs_employees_id INTEGER NOT NULL DEFAULT 0,
    quantity_start INTEGER NOT NULL DEFAULT 0,
    quantity_end INTEGER NOT NULL DEFAULT 0,
    quantity_sell INTEGER NOT NULL DEFAULT 0,
    quantity_written_off INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE srvs_shifts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    srvs_location_id INTEGER NOT NULL DEFAULT 0,
    date BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
    srvs_employees_id INTEGER NOT NULL DEFAULT 0,
    salary INTEGER NOT NULL DEFAULT 0,
    paid INTEGER NOT NULL DEFAULT 0,
    left_to_pay INTEGER NOT NULL DEFAULT 0,
    tips INTEGER NOT NULL DEFAULT 0,
    quantity_post_cards INTEGER NOT NULL DEFAULT 0,
    quantity_prints INTEGER NOT NULL DEFAULT 0,
    quantity_feedbacks INTEGER NOT NULL DEFAULT 0
);
