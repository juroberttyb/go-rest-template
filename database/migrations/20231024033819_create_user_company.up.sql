CREATE TABLE IF NOT EXISTS public.company
(
    id uuid NOT NULL,
    name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    picture character varying(100) COLLATE pg_catalog."default",
    email character varying(100) COLLATE pg_catalog."default",
    phone character varying(100) COLLATE pg_catalog."default",
    address character varying(100) COLLATE pg_catalog."default",
    site character varying(100) COLLATE pg_catalog."default",
    profile text COLLATE pg_catalog."default",
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT company_pkey PRIMARY KEY (id),
    CONSTRAINT company_email_key UNIQUE (email)
);

ALTER TABLE IF EXISTS public."user"
    ADD COLUMN company_id uuid;

ALTER TABLE IF EXISTS public."user"
    ADD FOREIGN KEY (company_id)
    REFERENCES public.company (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

