-- Table: public.users

-- DROP TABLE public.users;

CREATE TABLE public.users
(
    user_id serial not null,
    phone_num character varying(20) COLLATE pg_catalog."default",
    account character varying(20) COLLATE pg_catalog."default",
    email character varying(100) COLLATE pg_catalog."default",
    nick_name character varying(32) COLLATE pg_catalog."default",
    real_name character varying(40) COLLATE pg_catalog."default",
    gender smallint,
    birthday date,
    id_no character varying(20) COLLATE pg_catalog."default",
    icon character varying(255) COLLATE pg_catalog."default",
    password character varying(40) COLLATE pg_catalog."default",
    location character varying(40) COLLATE pg_catalog."default",
    created_time timestamp with time zone NOT NULL DEFAULT now(),
    updated_time timestamp with time zone NOT NULL DEFAULT now(),
    pwd_salt character varying(64) COLLATE pg_catalog."default" NOT NULL,
    per_sign character varying(100) COLLATE pg_catalog."default",
    completion smallint DEFAULT 30,
    home jsonb DEFAULT '{}'::jsonb,
    mile_remind smallint DEFAULT 0,
    reg_type smallint DEFAULT 0,
    tz character varying(40) COLLATE pg_catalog."default" NOT NULL DEFAULT 'Asia/Shanghai'::character varying,
    app character varying(20) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    phone_area character varying(10) COLLATE pg_catalog."default",
    general_setup jsonb NOT NULL DEFAULT '{}'::jsonb,
    CONSTRAINT pk_users PRIMARY KEY (user_id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.users
    OWNER to postgres;
-- Index: users_union_pk

-- DROP INDEX public.users_union_pk;

CREATE UNIQUE INDEX users_union_pk
    ON public.users USING btree
    (phone_num COLLATE pg_catalog."default" ASC NULLS LAST, phone_area COLLATE pg_catalog."default" ASC NULLS LAST, app COLLATE pg_catalog."default" ASC NULLS LAST, email COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;