create table product
(
    id          serial       not null
        constraint pk_product
            primary key,
    name        varchar(40),
    des         varchar(255),
    status      smallint     not null    default 1, -- 0 下架,1上架
    category_id smallint,
    begin_at    timestamp with time zone,
    expired_at  timestamp with time zone,
    img         varchar(200) not null    default '',
    creator integer not null default 0,
    updated_at  timestamp with time zone default now(),
    created_at  timestamp with time zone default now()
);

create index idx_product_name on product (name);

alter table product
    owner to postgres;

insert into product(id,name)values (1,'流量包月');
