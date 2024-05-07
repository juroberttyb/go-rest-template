CREATE TABLE IF NOT EXISTS public.chat_sticker
(
    id character varying(20) COLLATE pg_catalog."default" NOT NULL,
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    folder character varying(20) COLLATE pg_catalog."default" NOT NULL,
    line_url character varying(320) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone NOT NULL,
    CONSTRAINT chat_sticker_pkey PRIMARY KEY (id)
);
