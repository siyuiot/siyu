create table mall_value
(
    id         serial primary key,
    attr_id    integer not null,
    name       varchar(20),
    seq        smallint not null        default 1000,
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now()
);

create index idx_mv_ai on mall_value (attr_id);

alter table mall_value
    owner to postgres;
comment on table mall_value is '货品属性值';


-- 初始数据
insert into mall_value (id, attr_id, name)
values (1, 1, '蓝鲨 Robor'),
       (2, 1, '蓝鲨 Robor lite'),
       (3, 2, '秋月白'),
       (4, 2, '云墨黑'),
       (5, 2, '绿如蓝'),
       (6, 3, '锂电池包*2'),
       (7, 3, '锂电池包*1')
;

insert into mall_value (id, attr_id, name)
values (8, 4, '1年')
;
