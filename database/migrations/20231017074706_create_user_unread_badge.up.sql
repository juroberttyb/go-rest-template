ALTER TABLE IF EXISTS public."user"
    ADD COLUMN unread_count bigint NOT NULL DEFAULT 0;
