create table sim_order
(
    oid          serial,
    uid          integer       default 0,
    sid          integer       default 0,
    name         varchar(40)   default null,
    no           varchar(40)   default null,
    typ          integer       default 0,
    sku_id       integer       default 0,
    sku_count    integer       default 0,
    status       integer       default 0,
    amount_price integer       default 0,
    coupon_price integer       default 0,
    logistic_fee integer       default 0,
    due_price    integer       default 0,
    pay_price    integer       default 0,
    pay_channel  varchar(20)   default null,
    payment_id   integer       default 0,
    paid         integer       default 0,
    refund       integer       default 0,
    expired      integer       default 0,
    remark       varchar(2048) default null,
    created      integer       default 0,
    updated      integer       default 0,
    constraint pk_sim_order primary key(oid)
);
alter table sim_order owner to postgres;
comment on table sim_order is '订单';
comment on column sim_order.oid is '订单id';
comment on column sim_order.uid is '用户 id';
comment on column sim_order.sid is 'sim卡id';
comment on column sim_order.no is '订单编号';
comment on column sim_order.typ is '订单类型 1 商品';
comment on column sim_order.sku_id is '商品id';
comment on column sim_order.status is '订单状态';
comment on column sim_order.amount_price is '物品总价';
comment on column sim_order.coupon_price is '优惠价格';
comment on column sim_order.logistic_fee is '运费';
comment on column sim_order.due_price is '应付价格';
comment on column sim_order.pay_price is '实付价格';
comment on column sim_order.pay_channel is '支付渠道';
comment on column sim_order.payment_id is '支付 id payment 表关联';
comment on column sim_order.paid is '支付时间';
comment on column sim_order.refund is '退款时间';
comment on column sim_order.expired is '过期时间';
comment on column sim_order.remark is '备注';