-- user.sql
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名称',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '用户密码',
  `mobile` varchar(255) NOT NULL DEFAULT '' COMMENT '手机号',
  `gender` char(5) NOT NULL DEFAULT 'male' COMMENT '男｜女｜未公开',
  `nickname` varchar(255) DEFAULT '' COMMENT '用户昵称',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`), -- ⚠️注意：唯一索引会生成 FindOneByUsername 方法
  UNIQUE KEY `idx_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;