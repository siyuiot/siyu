create table product_sku
(
    id         serial primary key,
    name       varchar(40)  not null    default '',
    price      integer      not null,
    product_id integer      not null    default 0,
    status     smallint     not null    default 0, --0 下架 1 上架
    img        varchar(255) not null    default '',
    des        varchar(2047),
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);
create index idx_ps_pi
    on product_sku (product_id);

alter table product_sku
    owner to postgres;

comment on table product_sku is '商品sku';

alter table product_sku
    add price_origin int default 0 not null;

comment on column product_sku.price_origin is '原价';

-- 初始数据
insert into product_sku(id, name, price, product_id, status)
values (1, 'Robor,秋月白,双电池', 1497600, 1, 1),
       (2, 'Robor,云墨黑,双电池', 1497600, 1, 1),
       (3, 'Robor,绿如蓝,双电池', 1497600, 1, 1),
       (4, 'Robor,秋月白,单电池', 1257600, 1, 1),
       (5, 'Robor,云墨黑,单电池', 1257600, 1, 1),
       (6, 'Robor,绿如蓝,单电池', 1257600, 1, 1),
       (7, 'Robor lite,秋月白,双电池', 1197600, 1, 1),
       (8, 'Robor lite,云墨黑,双电池', 1197600, 1, 1),
       (9, 'Robor lite,秋月白,单电池', 957600, 1, 1),
       (10, 'Robor lite,云墨黑,单电池', 957600, 1, 1)
;

-- 初始数据 sim续费
insert into product_sku(id, name, price,price_origin, product_id, status) values
(11, '整车智能服务', 3900,6900, 2, 1)
;
