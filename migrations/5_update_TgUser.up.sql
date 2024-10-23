ALTER TABLE tg_users DROP COLUMN chatbot_state;
ALTER TABLE tg_users ADD COLUMN tg_cb_flow_step_id INTEGER NOT NULL DEFAULT 0;
