create table user_sim_log
(
    ts            integer                  default 0,
    uid           integer                  default 0,
    sid           integer                  default 0,
    phone_num     varchar(32)              default null,
    sim_no        varchar(32)              default null,
    imsi          varchar(32)              default null,
    iccid         varchar(32)              default null,
    remark        varchar(1024)            default null,
    created       integer                  default 0,
    constraint pk_user_sim_log primary key (ts)
)with (OIDS = FALSE)
tablespace pg_default;

alter table user_sim_log owner to postgres;

-- create index idx_id on user_sim_log(id);

comment on table user_sim_log is '用户绑定sim卡log';

comment on column user_sim_log.ts is '绑定/解绑时间';
comment on column user_sim_log.uid is '用户id';
comment on column user_sim_log.sid is 'sim卡id';
comment on column user_sim_log.phone_num is 'phone_num';
comment on column user_sim_log.sim_no is 'sim_no';
comment on column user_sim_log.imsi is 'imsi';
comment on column user_sim_log.iccid is 'iccid';
comment on column user_sim_log.remark is '备注';