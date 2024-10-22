CREATE TABLE srvs_location_list (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    row_order INTEGER NOT NULL,
    added_timestamp BIGINT NOT NULL DEFAULT (strftime('%s', 'now'))
);

INSERT INTO srvs_location_list (name, row_order) VALUES ('Islander', '10');
INSERT INTO srvs_location_list (name, row_order) VALUES ('Johnny', '20');


CREATE TABLE srvs_employees_list (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    percentage           REAL NOT NULL,
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
    srvs_location_id INTEGER NOT NULL,
    date BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
    srvs_goods_id INTEGER NOT NULL,
    srvs_employees_id INTEGER NOT NULL,
    quantity_start INTEGER NOT NULL,
    quantity_end INTEGER NOT NULL,
    quantity_sell INTEGER NOT NULL,
    quantity_written_off INTEGER NOT NULL
);

CREATE TABLE srvs_shifts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    srvs_location_id INTEGER NOT NULL,
    date BIGINT NOT NULL DEFAULT (strftime('%s', 'now')),
    srvs_employees_id INTEGER NOT NULL,
    salary INTEGER NOT NULL,
    paid INTEGER NOT NULL,
    left_to_pay INTEGER NOT NULL,
    tips INTEGER NOT NULL,
    quantity_post_cards INTEGER NOT NULL,
    quantity_prints INTEGER NOT NULL,
    quantity_feedbacks INTEGER NOT NULL
);
