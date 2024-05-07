ALTER TABLE IF EXISTS public."analytics"
    ADD COLUMN url character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying;

ALTER TABLE IF EXISTS public."analytics" DROP CONSTRAINT analytics_pkey;

ALTER TABLE IF EXISTS public."analytics" ADD PRIMARY KEY (broadcast_message_id, url);
