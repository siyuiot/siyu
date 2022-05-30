create table wechat_pay_mch
(
    mch_id                           varchar(32)              default null,
    mch_certificate_serial_number    varchar(128)             default null,
    mch_API_v3_key                   varchar(32)              default null,
    mch_private_key                  varchar(2048)            default null,
    mch_cert                         varchar(2048)            default null,
    created                          integer                  default 0,
    updated                          integer                  default 0,
    constraint pk_wechat_pay_mch_mch_id primary key (mch_id)
);
comment on table wechat_pay_mch is '微信支付';