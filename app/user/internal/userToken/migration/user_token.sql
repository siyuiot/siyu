create table user_token
(
    uid    serial  constraint pk_user_token primary key,
    ts     integer,
    expire integer,
    token  varchar(40)  not null default '',
    des    varchar(255) not null default ''
);

comment on table user_token is '用户登录token';
comment on column user_token.ts is '最后一次登录时间';
comment on column user_token.token is '登录token';
comment on column user_token.expire is '过期时间';
comment on column user_token.des is '描述';
