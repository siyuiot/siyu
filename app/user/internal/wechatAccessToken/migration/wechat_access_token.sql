create table wechat_access_token
(
    app_id        varchar(32)              default null,
    secret        varchar(128)             default null,
    access_token  varchar(512)             default null,
    expires_in    integer                  default 0,
    expires_at    integer                  default 0,
    remark        varchar(1024)            default null,
    created       integer                  default 0,
    updated       integer                  default 0,
    constraint pk_id primary key (app_id)
)with (OIDS = FALSE)
tablespace pg_default;

alter table wechat_access_token owner to postgres;

-- create index idx_id on wechat_access_token(id);

comment on table wechat_access_token is '微信accessToken';

comment on column wechat_access_token.app_id is '第三方用户唯一凭证';
comment on column wechat_access_token.secret is '第三方用户唯一凭证密钥，即appsecret';
comment on column wechat_access_token.access_token is '获取到的凭证';
comment on column wechat_access_token.expires_in is '凭证有效时间，单位：秒';
comment on column wechat_access_token.expires_at is '凭证过期时间戳';
comment on column wechat_access_token.remark is '备注';