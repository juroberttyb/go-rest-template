ALTER TABLE IF EXISTS public."chat_thread"
    ADD COLUMN control_flag smallint NOT NULL DEFAULT 0;
