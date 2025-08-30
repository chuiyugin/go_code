CREATE TABLE `orders`(
                        `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '主键',
                        `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                        `create_by` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '创建者',
                        `update_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
                        `update_by` VARCHAR(64) NOT NULL DEFAULT '' COMMENT '更新者',
                        `version` SMALLINT(5) UNSIGNED NOT NULL DEFAULT '0' COMMENT '乐观锁版本号',
                        `is_del` tinyint(4) UNSIGNED NOT NULL DEFAULT '0' COMMENT '是否删除：0正常1删除',

                        `user_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '用户id',
                        `order_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '订单id',
                        `trade_id` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '交易单号',
                        `pay_channel` tinyint(4) UNSIGNED NOT NULL DEFAULT '0' COMMENT '支付方式',
                        `status` INT UNSIGNED NOT NULL DEFAULT '0' COMMENT '订单状态:100创建订单/待支付 200已支付 300交易关闭 400完成',
                        `pay_amount` BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '支付金额（分）',
                        `pay_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '支付时间',

                        UNIQUE KEY uniq_order_id (`order_id`),   -- 在 order_id 上加唯一约束
                        INDEX (user_id),
                        INDEX (order_id),
                        INDEX (trade_id),
                        INDEX (is_del)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COMMENT = '订单表';

-- 插入一条数据
-- INSERT INTO `orders` (
--     `create_at`,
--     `create_by`,
--     `update_at`,
--     `update_by`,
--     `version`,
--     `is_del`,
--     `user_id`,
--     `order_id`,
--     `trade_id`,
--     `pay_channel`,
--     `status`,
--     `pay_amount`,
--     `pay_time`
-- ) VALUES (
--     NOW(),             -- `create_at`, 当前时间
--     'admin',           -- `create_by`, 创建者（根据需求修改）
--     NOW(),             -- `update_at`, 当前时间
--     'yugin',           -- `update_by`, 更新者（根据需求修改）
--     0,                 -- `version`, 默认值为0
--     0,                 -- `is_del`, 默认值为0（表示正常）
--     123456,            -- `user_id`, 用户ID（根据实际情况修改）
--     987654321,         -- `order_id`, 订单ID（根据实际情况修改）
--     'TXN1234567890',   -- `trade_id`, 交易单号（根据实际情况修改）
--     1,                 -- `pay_channel`, 支付方式（例如 1：支付宝）
--     200,               -- `status`, 订单状态（200：已支付）
--     1500,              -- `pay_amount`, 支付金额（分为单位，例如 1500 表示 15 元）
--     NOW()              -- `pay_time`, 支付时间
-- );