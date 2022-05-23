create table mall_attribute
(
    id         serial primary key,
    name       varchar(20),
    seq        smallint not null        default 1000,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

alter table mall_attribute
    owner to postgres;
comment
on table mall_attribute is '货品属性';


-- 初始数据
insert into mall_attribute(id, name)
values (1, '车型'),
       (2, '颜色'),
       (3, '电池数量')
;
insert into mall_attribute(id, name)
values (4, '年限');
