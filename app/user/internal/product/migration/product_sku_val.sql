create table product_sku_val
(
    id         serial primary key,
    sku_id     integer,
    val_id     integer,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

create index idx_psv_si on product_sku_val (sku_id);
create index idx_psv_val_id on product_sku_val (val_id);

comment on table product_sku_val is 'sku 属性值关联表';

alter table product_sku_val
    owner to postgres;



insert into product_sku_val(sku_id, val_id)
values
------ 蓝鲨Robor'
    -- 双电池
    (1, 1),
    (1, 3),--秋月白
    (1, 6),
    (2, 1),
    (2, 4),--云墨黑
    (2, 6),
    (3, 1),
    (3, 5),--绿如蓝
    (3, 6),
    -- 单电池
    (4, 1),
    (4, 3),--秋月白
    (4, 7),--单电池
    (5, 1),
    (5, 4),--云墨黑
    (5, 7),--单电池
    (6, 1),
    (6, 5),--绿如蓝
    (6, 7),--单电池
---------------蓝鲨 Robor lite
    ---- 双电池
    (7, 2),
    (7, 3),--秋月白
    (7, 6),--双电池
    (8, 2),
    (8, 4),--云墨黑
    (8, 6),--双电池
    -- 单电池
    (9, 2),
    (9, 3),--秋月白
    (9, 7),
    (10, 2),
    (10, 4),--云墨黑
    (10, 7)
;


insert into product_sku_val(sku_id, val_id)
values (11,8);
