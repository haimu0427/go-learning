-- 创建数据库
CREATE DATABASE IF NOT EXISTS gozero DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE gozero;

-- 创建用户表
CREATE TABLE IF NOT EXISTS `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名称',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '用户密码',
  `mobile` varchar(255) NOT NULL DEFAULT '' COMMENT '手机号',
  `gender` char(5) NOT NULL DEFAULT 'male' COMMENT '男｜女｜未公开',
  `nickname` varchar(255) DEFAULT '' COMMENT '用户昵称',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 插入测试数据
INSERT INTO `user` (username, password, mobile, gender, nickname) VALUES
('root', 'root', '13800138000', 'male', '管理员'),
('test', 'test123', '13800138001', 'female', '测试用户')
ON DUPLICATE KEY UPDATE
password = VALUES(password),
mobile = VALUES(mobile),
gender = VALUES(gender),
nickname = VALUES(nickname);