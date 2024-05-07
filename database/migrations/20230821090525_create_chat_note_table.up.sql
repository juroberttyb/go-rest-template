CREATE TABLE IF NOT EXISTS public.chat_note
(
    id uuid NOT NULL,
    chat_id uuid NOT NULL,
    target_id uuid NOT NULL,
    text text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    author_id uuid NOT NULL,
    status smallint NOT NULL DEFAULT 0,
    CONSTRAINT chat_note_pkey PRIMARY KEY (id),
    CONSTRAINT chat_note_chat_id_target_id_key UNIQUE (chat_id, target_id)
);
