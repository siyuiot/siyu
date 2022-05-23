-- Table: public.sim

-- DROP TABLE public.sim;

CREATE TABLE IF NOT EXISTS public.sim
(
    id serial not null,
    sim_no character varying(20) COLLATE pg_catalog."default" NOT NULL,
    imsi character varying(40) COLLATE pg_catalog."default",
    iccid character varying(40) COLLATE pg_catalog."default",
    test_period smallint DEFAULT 2,
    quiet_period smallint DEFAULT 6,
    default_service_duration smallint DEFAULT 6,
    status smallint DEFAULT 0,
    remark character varying(1000) COLLATE pg_catalog."default" DEFAULT ''::character varying,
    created integer NOT NULL DEFAULT 0,
    innet_date date NOT NULL DEFAULT '1970-01-01'::date,
    activated_time timestamp without time zone NOT NULL DEFAULT '0001-01-01 00:00:00'::timestamp without time zone,
    last_activated date DEFAULT '0001-01-01'::date,
    service integer NOT NULL DEFAULT 0,
    country character varying(20)[] COLLATE pg_catalog."default" DEFAULT ARRAY[]::character varying[],
    extra jsonb DEFAULT '{}'::jsonb,
    updated integer DEFAULT 0,
    CONSTRAINT pk_sim PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.sim OWNER to postgres;

COMMENT ON TABLE public.sim IS 'sim 卡';

COMMENT ON COLUMN public.sim.test_period IS '测试期';

COMMENT ON COLUMN public.sim.quiet_period IS '静默期';

COMMENT ON COLUMN public.sim.default_service_duration IS '默认服务时间';

COMMENT ON COLUMN public.sim.status IS '1正常 2停卡 3欠费';

COMMENT ON COLUMN public.sim.last_activated IS '最后允许激活日期';

COMMENT ON COLUMN public.sim.service IS '服务属性';
