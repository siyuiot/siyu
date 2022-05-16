create table user_login_token
(
    uid    serial
        constraint user_login_token_pk primary key,
    ts     integer,
    expire integer,
    token  varchar(40)  not null default '',
    des    varchar(255) not null default ''
);

comment on table user_login_token is '用户登录 token';
comment on column user_login_token.ts is '最后一次登录时间';
comment on column user_login_token.token is '登录 token';
comment on column user_login_token.expire is '过期时间';
comment on column user_login_token.des is '最后一次登录时间';
