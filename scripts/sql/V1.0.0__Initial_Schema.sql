-- 删除已存在的表（如果存在）
DROP TABLE IF EXISTS `sys_role_menu`;
DROP TABLE IF EXISTS `sys_user_role`;
DROP TABLE IF EXISTS `sys_menu`;
DROP TABLE IF EXISTS `sys_role`;
DROP TABLE IF EXISTS `sys_user`;

-- 创建用户表
CREATE TABLE IF NOT EXISTS `sys_user` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0=禁用，1=启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 创建角色表
CREATE TABLE IF NOT EXISTS `sys_role` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(50) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `description` varchar(255) DEFAULT NULL COMMENT '角色描述',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0=禁用，1=启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 创建用户角色关联表
CREATE TABLE IF NOT EXISTS `sys_user_role` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`,`role_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 创建菜单表
CREATE TABLE IF NOT EXISTS `sys_menu` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `parent_id` bigint DEFAULT NULL COMMENT '父菜单ID',
  `name` varchar(50) NOT NULL COMMENT '菜单名称',
  `path` varchar(200) DEFAULT NULL COMMENT '路由路径',
  `component` varchar(255) DEFAULT NULL COMMENT '组件路径',
  `perms` varchar(100) DEFAULT NULL COMMENT '权限标识',
  `type` tinyint NOT NULL DEFAULT '1' COMMENT '类型：0=目录，1=菜单，2=按钮',
  `icon` varchar(100) DEFAULT NULL COMMENT '图标',
  `sort_order` int DEFAULT '0' COMMENT '排序',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0=禁用，1=启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

-- 创建角色菜单关联表
CREATE TABLE IF NOT EXISTS `sys_role_menu` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `menu_id` bigint NOT NULL COMMENT '菜单ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_menu` (`role_id`,`menu_id`),
  KEY `idx_menu_id` (`menu_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关联表';

-- 初始化角色数据
INSERT INTO `sys_role` (`name`, `code`, `status`, `description`) 
SELECT '系统管理员', 'admin', 1, '系统管理员，拥有所有权限'
WHERE NOT EXISTS (SELECT 1 FROM `sys_role` WHERE `code` = 'admin');

INSERT INTO `sys_role` (`name`, `code`, `status`, `description`) 
SELECT '普通用户', 'user', 1, '普通用户'
WHERE NOT EXISTS (SELECT 1 FROM `sys_role` WHERE `code` = 'user');

-- 获取管理员角色ID
SET @admin_role_id = NULL;
SELECT @admin_role_id := id FROM `sys_role` WHERE `code` = 'admin';

-- 创建管理员用户
-- 密码: admin123 -> SHA256: 240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9 -> bcrypt
INSERT INTO `sys_user` (
  `username`,
  `password`,
  `nickname`,
  `email`,
  `status`
) SELECT
  'admin',
  '$2a$10$RBVyo4ojLZEG.j.uJnpcEuhY4DRC0TChUAdo1XBnlKB/mKnZDEbUS',
  '系统管理员',
  'admin@example.com',
  1
WHERE NOT EXISTS (SELECT 1 FROM `sys_user` WHERE `username` = 'admin');

-- 获取管理员用户ID
SET @admin_user_id = NULL;
SELECT @admin_user_id := id FROM `sys_user` WHERE `username` = 'admin';

-- 关联管理员用户和角色
INSERT IGNORE INTO `sys_user_role` (`user_id`, `role_id`) 
SELECT @admin_user_id, @admin_role_id 
WHERE @admin_user_id IS NOT NULL AND @admin_role_id IS NOT NULL;

-- 系统管理菜单
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `perms`, `type`, `icon`, `sort_order`, `status`) 
VALUES (NULL, '系统管理', '/system', 'Layout', NULL, 0, 'setting', 99, 1);

-- 获取系统管理菜单ID
SET @system_id = NULL;
SELECT @system_id := LAST_INSERT_ID();

-- 系统管理子菜单
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `perms`, `type`, `icon`, `sort_order`, `status`) VALUES 
(@system_id, '用户管理', 'user', 'system/user/index', 'system:user:list', 1, 'user', 1, 1),
(@system_id, '角色管理', 'role', 'system/role/index', 'system:role:list', 1, 'role', 2, 1),
(@system_id, '菜单管理', 'menu', 'system/menu/index', 'system:menu:list', 1, 'menu', 3, 1);

-- 授权所有菜单给管理员角色
INSERT IGNORE INTO `sys_role_menu` (`role_id`, `menu_id`)
SELECT @admin_role_id, id FROM `sys_menu` 
WHERE @admin_role_id IS NOT NULL AND deleted_at IS NULL;
