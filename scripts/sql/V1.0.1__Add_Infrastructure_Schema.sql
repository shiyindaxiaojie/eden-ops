-- 删除已存在的表（如果存在）
DROP TABLE IF EXISTS `infra_k8s_workload`;
DROP TABLE IF EXISTS `infra_k8s_config`;
DROP TABLE IF EXISTS `infra_database_config`;
DROP TABLE IF EXISTS `infra_server_config`;
DROP TABLE IF EXISTS `infra_cloud_account`;
DROP TABLE IF EXISTS `infra_cloud_provider`;

-- 创建云厂商表
CREATE TABLE IF NOT EXISTS `infra_cloud_provider` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '云厂商ID',
  `name` varchar(50) NOT NULL COMMENT '云厂商名称',
  `code` varchar(50) NOT NULL COMMENT '云厂商代码',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0=禁用，1=启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='云厂商表';

-- Create Kubernetes configuration table
CREATE TABLE IF NOT EXISTS `infra_k8s_config` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(100) NOT NULL,
    `provider_id` BIGINT,
    `kubeconfig` TEXT NOT NULL,
    `description` VARCHAR(500),
    `status` TINYINT DEFAULT 1,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Kubernetes集群表';

-- Create database configuration table
CREATE TABLE IF NOT EXISTS `infra_database_config` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '数据库ID',
  `name` varchar(50) NOT NULL COMMENT '数据库名称',
  `provider_id` bigint DEFAULT NULL COMMENT '云厂商ID',
  `host` varchar(255) NOT NULL COMMENT '主机地址',
  `port` int NOT NULL COMMENT '端口',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT '密码密文',
  `type` varchar(20) NOT NULL COMMENT '数据库类型：MySQL、PostgreSQL等',
  `database` VARCHAR(100) NOT NULL COMMENT '数据库名称',
  `parameters` VARCHAR(500) COMMENT '连接参数',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0=禁用，1=启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_provider_id` (`provider_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据库配置表';

-- Create server configuration table
CREATE TABLE IF NOT EXISTS `infra_server_config` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '服务器ID',
  `name` varchar(50) NOT NULL COMMENT '服务器名称',
  `provider_id` bigint DEFAULT NULL COMMENT '云厂商ID',
  `host` varchar(255) NOT NULL COMMENT '主机地址',
  `port` int NOT NULL DEFAULT '22' COMMENT 'SSH端口',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) DEFAULT NULL COMMENT '密码密文',
  `private_key` text DEFAULT NULL COMMENT '私钥',
  `os` VARCHAR(50) NOT NULL COMMENT '操作系统',
  `arch` VARCHAR(20) COMMENT '架构',
  `cpu` INT COMMENT 'CPU核数',
  `memory` INT COMMENT '内存大小(MB)',
  `disk` INT COMMENT '磁盘大小(GB)',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0=禁用，1=启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_provider_id` (`provider_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务器配置表';

-- 创建云账号表
CREATE TABLE IF NOT EXISTS `infra_cloud_account` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '云账号ID',
  `name` varchar(50) NOT NULL COMMENT '账号名称',
  `provider_id` bigint DEFAULT NULL COMMENT '云厂商ID',
  `access_key` varchar(100) NOT NULL COMMENT '访问密钥',
  `secret_key` varchar(100) NOT NULL COMMENT '访问密钥密文',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0=禁用，1=启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_provider_id` (`provider_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='云账号表';

-- 创建 Kubernetes 工作负载表
CREATE TABLE IF NOT EXISTS `infra_k8s_workload` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `config_id` bigint NOT NULL COMMENT '集群配置ID',
  `name` varchar(200) NOT NULL COMMENT '工作负载名称',
  `namespace` varchar(100) NOT NULL COMMENT '命名空间',
  `kind` varchar(20) NOT NULL COMMENT '工作负载类型',
  `replicas` int DEFAULT '0' COMMENT '副本数',
  `ready_replicas` int DEFAULT '0' COMMENT '就绪副本数',
  `status` varchar(20) DEFAULT NULL COMMENT '状态',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_config_id` (`config_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Kubernetes工作负载';

-- 初始化云厂商数据
INSERT INTO `infra_cloud_provider` (`name`, `code`, `description`, `status`) VALUES 
('本地机房', 'local', '自建机房/本地数据中心', 1),
('阿里云', 'aliyun', '阿里云', 1),
('腾讯云', 'tencent', '腾讯云', 1),
('华为云', 'huawei', '华为云', 1),
('AWS', 'aws', 'Amazon Web Services', 1),
('Azure', 'azure', 'Microsoft Azure', 1),
('Google Cloud', 'gcp', 'Google Cloud Platform', 1);

-- 获取管理员角色ID
SET @admin_role_id = NULL;
SELECT @admin_role_id := id FROM `sys_role` WHERE `code` = 'admin';

-- 基础设施菜单
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `perms`, `type`, `icon`, `sort_order`, `status`) 
VALUES (NULL, '基础设施', '/infrastructure', 'Layout', NULL, 0, 'server', 1, 1);

-- 获取基础设施菜单ID
SET @infra_id = NULL;
SELECT @infra_id := LAST_INSERT_ID();

-- 基础设施子菜单
INSERT INTO `sys_menu` (`parent_id`, `name`, `path`, `component`, `perms`, `type`, `icon`, `sort_order`, `status`) VALUES 
(@infra_id, 'Kubernetes', 'kubernetes', 'infrastructure/kubernetes/index', 'infrastructure:kubernetes:list', 1, 'kubernetes', 1, 1),
(@infra_id, '服务器管理', 'server', 'infrastructure/server/index', 'infrastructure:server:list', 1, 'server', 2, 1),
(@infra_id, '数据库管理', 'database', 'infrastructure/database/index', 'infrastructure:database:list', 1, 'database', 3, 1),
(@infra_id, '云账号管理', 'cloud-provider', 'infrastructure/cloud-provider/index', 'infrastructure:cloud:list', 1, 'cloud', 4, 1);

-- 获取各个子菜单的ID
SET @k8s_menu_id = NULL;
SET @server_menu_id = NULL;
SET @db_menu_id = NULL;
SET @cloud_menu_id = NULL;

SELECT @k8s_menu_id := id FROM `sys_menu` WHERE `path` = 'kubernetes' AND parent_id = @infra_id;
SELECT @server_menu_id := id FROM `sys_menu` WHERE `path` = 'server' AND parent_id = @infra_id;
SELECT @db_menu_id := id FROM `sys_menu` WHERE `path` = 'database' AND parent_id = @infra_id;
SELECT @cloud_menu_id := id FROM `sys_menu` WHERE `path` = 'cloud-provider' AND parent_id = @infra_id;

-- 授权基础设施菜单给管理员角色
INSERT IGNORE INTO `sys_role_menu` (`role_id`, `menu_id`)
SELECT @admin_role_id, id FROM `sys_menu` 
WHERE @admin_role_id IS NOT NULL AND id = @infra_id;

INSERT IGNORE INTO `sys_role_menu` (`role_id`, `menu_id`)
SELECT @admin_role_id, id FROM `sys_menu` 
WHERE @admin_role_id IS NOT NULL AND parent_id = @infra_id;

-- Add CRUD permissions for Cloud Provider
INSERT INTO `sys_menu` (`parent_id`, `name`, `perms`, `type`, `icon`, `sort_order`, `status`)
VALUES 
(@cloud_menu_id, 'CloudProviderList', 'infrastructure:cloud-provider:list', 2, NULL, 1, 1),
(@cloud_menu_id, 'CloudProviderCreate', 'infrastructure:cloud-provider:create', 2, NULL, 2, 1),
(@cloud_menu_id, 'CloudProviderUpdate', 'infrastructure:cloud-provider:update', 2, NULL, 3, 1),
(@cloud_menu_id, 'CloudProviderDelete', 'infrastructure:cloud-provider:delete', 2, NULL, 4, 1);

-- Add CRUD permissions for Kubernetes
INSERT INTO `sys_menu` (`parent_id`, `name`, `perms`, `type`, `icon`, `sort_order`, `status`)
VALUES 
(@k8s_menu_id, 'KubernetesList', 'infrastructure:kubernetes:list', 2, NULL, 1, 1),
(@k8s_menu_id, 'KubernetesCreate', 'infrastructure:kubernetes:create', 2, NULL, 2, 1),
(@k8s_menu_id, 'KubernetesUpdate', 'infrastructure:kubernetes:update', 2, NULL, 3, 1),
(@k8s_menu_id, 'KubernetesDelete', 'infrastructure:kubernetes:delete', 2, NULL, 4, 1);

-- Add CRUD permissions for Database
INSERT INTO `sys_menu` (`parent_id`, `name`, `perms`, `type`, `icon`, `sort_order`, `status`)
VALUES 
(@db_menu_id, 'DatabaseList', 'infrastructure:database:list', 2, NULL, 1, 1),
(@db_menu_id, 'DatabaseCreate', 'infrastructure:database:create', 2, NULL, 2, 1),
(@db_menu_id, 'DatabaseUpdate', 'infrastructure:database:update', 2, NULL, 3, 1),
(@db_menu_id, 'DatabaseDelete', 'infrastructure:database:delete', 2, NULL, 4, 1);

-- Add CRUD permissions for Server
INSERT INTO `sys_menu` (`parent_id`, `name`, `perms`, `type`, `icon`, `sort_order`, `status`)
VALUES 
(@server_menu_id, 'ServerList', 'infrastructure:server:list', 2, NULL, 1, 1),
(@server_menu_id, 'ServerCreate', 'infrastructure:server:create', 2, NULL, 2, 1),
(@server_menu_id, 'ServerUpdate', 'infrastructure:server:update', 2, NULL, 3, 1),
(@server_menu_id, 'ServerDelete', 'infrastructure:server:delete', 2, NULL, 4, 1); 