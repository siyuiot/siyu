create table product_attribute
(
    id         serial primary key,
    product_id integer,
    attr_id    integer,
    seq        smallint not null        default 1000, -- 排序
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);
create index idx_pa_pi
    on product_attribute (product_id);

alter table product_attribute
    owner to postgres;

comment on table product_attribute is '货品属性和规格的关联';

-- 初始数据
insert into product_attribute(id,product_id,attr_id)
values (1,1,1),
       (2,1,2),
       (3,1,3)
;
