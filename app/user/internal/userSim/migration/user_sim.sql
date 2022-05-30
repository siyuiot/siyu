create table user_sim
(
    uid                   integer                  default 0,
    sid                   integer                  default 0,
    sim_provider          varchar(13)              default null,
    sim_no                varchar(13)              default null,
    iccid                 varchar(20)              default null,
    sim_byte              integer                  default 0,
    sim_available_byte    integer                  default 0,
    bind_ts               integer                  default 0,
    service_end_ts        integer                  default 0,
    service_duration      integer                  default 0,
    remark                varchar(1024)            default null,
    created               integer                  default 0,
    updated               integer                  default 0,
    constraint pk_user_sim primary key (uid,sid)
)with (OIDS = FALSE)
tablespace pg_default;

alter table user_sim owner to postgres;

-- create index idx_id on user_sim(id);

comment on table user_sim is '用户绑定sim卡';

comment on column user_sim.uid is '用户id';
comment on column user_sim.sid is 'sim卡id';
comment on column user_sim.sim_provider is 'sim卡提供商 ChinaTelecom ChinaMobile ChinaUnicom';
comment on column user_sim.sim_no is 'sim卡no';
comment on column user_sim.iccid is 'sim卡iccid';
comment on column user_sim.sim_byte is 'sim卡流量-用户购买流量';
comment on column user_sim.sim_available_byte is 'sim卡可用流量-用户剩余流量';
comment on column user_sim.bind_ts is '绑定时间戳';
comment on column user_sim.service_end_ts is '服务结束时间戳-用户使用截止时间戳';
comment on column user_sim.service_duration is '服务时长 单位:月';
comment on column user_sim.remark is '备注';