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
  `description` varchar(200) DEFAULT NULL COMMENT '描述',
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
    `sync_interval` INT DEFAULT 30 COMMENT '同步间隔(秒)',
    `version` VARCHAR(20) COMMENT 'Kubernetes版本',
    `context` VARCHAR(100) COMMENT '上下文',
    `cluster_id` VARCHAR(100) COMMENT '集群ID',
    `node_count` INT DEFAULT 0 COMMENT '节点数量',
    `pod_count` INT DEFAULT 0 COMMENT 'Pod数量',
    `cpu_total` VARCHAR(20) COMMENT 'CPU总量',
    `cpu_used` VARCHAR(20) COMMENT 'CPU使用量',
    `memory_total` VARCHAR(20) COMMENT '内存总量',
    `memory_used` VARCHAR(20) COMMENT '内存使用量',
    `workload_count` INT DEFAULT 0 COMMENT '工作负载数量',
    `workload_running` INT DEFAULT 0 COMMENT '运行中工作负载数量',
    `workload_idle` INT DEFAULT 0 COMMENT '闲置工作负载数量',
    `pod_total` INT DEFAULT 0 COMMENT 'Pod总数',
    `pod_running` INT DEFAULT 0 COMMENT '运行中Pod数量',
    `pod_error` INT DEFAULT 0 COMMENT '异常Pod数量',
    `node_total` INT DEFAULT 0 COMMENT '节点总数',
    `node_running` INT DEFAULT 0 COMMENT '运行中节点数量',
    `node_error` INT DEFAULT 0 COMMENT '异常节点数量',
    `workload_destroyed_count` INT DEFAULT 0 COMMENT '工作负载销毁数量',
    `pod_destroyed_count` INT DEFAULT 0 COMMENT 'Pod销毁数量',
    `node_destroyed_count` INT DEFAULT 0 COMMENT 'Node销毁数量',
    `last_sync_time` TIMESTAMP NULL COMMENT '最后同步时间',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Kubernetes集群表';

-- Create database configuration table
CREATE TABLE IF NOT EXISTS `infra_database_config` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '数据库ID',
  `name` varchar(50) NOT NULL COMMENT '数据库名称',
  `provider_id` bigint DEFAULT NULL COMMENT '云厂商ID',
  `host` varchar(100) NOT NULL COMMENT '主机地址',
  `port` int NOT NULL COMMENT '端口',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) NOT NULL COMMENT '密码密文',
  `type` varchar(20) NOT NULL COMMENT '数据库类型：MySQL、PostgreSQL等',
  `database` VARCHAR(100) NOT NULL COMMENT '数据库名称',
  `parameters` VARCHAR(500) COMMENT '连接参数',
  `description` varchar(200) DEFAULT NULL COMMENT '描述',
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
  `host` varchar(100) NOT NULL COMMENT '主机地址',
  `port` int NOT NULL DEFAULT '22' COMMENT 'SSH端口',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(100) DEFAULT NULL COMMENT '密码密文',
  `private_key` text DEFAULT NULL COMMENT '私钥',
  `os` VARCHAR(50) NOT NULL COMMENT '操作系统',
  `arch` VARCHAR(20) COMMENT '架构',
  `cpu` INT COMMENT 'CPU核数',
  `memory` INT COMMENT '内存大小(MB)',
  `disk` INT COMMENT '磁盘大小(GB)',
  `description` varchar(200) DEFAULT NULL COMMENT '描述',
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
  `description` varchar(200) DEFAULT NULL COMMENT '描述',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：0=禁用，1=启用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_provider_id` (`provider_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='云账号表';

-- 创建 Kubernetes 节点表
CREATE TABLE IF NOT EXISTS `infra_k8s_node` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '节点ID',
  `config_id` bigint NOT NULL COMMENT 'K8s配置ID',
  `name` varchar(100) NOT NULL COMMENT '节点名称',
  `internal_ip` varchar(45) DEFAULT NULL COMMENT '内部IP',
  `external_ip` varchar(45) DEFAULT NULL COMMENT '外部IP',
  `hostname` varchar(100) DEFAULT NULL COMMENT '主机名',
  `os_image` varchar(100) DEFAULT NULL COMMENT '操作系统镜像',
  `kernel_version` varchar(100) DEFAULT NULL COMMENT '内核版本',
  `container_runtime` varchar(100) DEFAULT NULL COMMENT '容器运行时',
  `kubelet_version` varchar(50) DEFAULT NULL COMMENT 'Kubelet版本',
  `kube_proxy_version` varchar(50) DEFAULT NULL COMMENT 'KubeProxy版本',
  `cpu_capacity` varchar(20) DEFAULT NULL COMMENT 'CPU容量',
  `memory_capacity` varchar(20) DEFAULT NULL COMMENT '内存容量',
  `pods_capacity` varchar(20) DEFAULT NULL COMMENT 'Pod容量',
  `cpu_allocatable` varchar(20) DEFAULT NULL COMMENT 'CPU可分配',
  `memory_allocatable` varchar(20) DEFAULT NULL COMMENT '内存可分配',
  `pods_allocatable` varchar(20) DEFAULT NULL COMMENT 'Pod可分配',
  `cpu_usage` varchar(20) DEFAULT NULL COMMENT 'CPU使用量',
  `memory_usage` varchar(20) DEFAULT NULL COMMENT '内存使用量',
  `pods_usage` int DEFAULT 0 COMMENT 'Pod使用量',
  `labels` text DEFAULT NULL COMMENT '标签(JSON格式)',
  `annotations` text DEFAULT NULL COMMENT '注解(JSON格式)',
  `taints` text DEFAULT NULL COMMENT '污点(JSON格式)',
  `conditions` text DEFAULT NULL COMMENT '状态条件(JSON格式)',
  `status` varchar(50) DEFAULT 'Unknown' COMMENT '节点状态',
  `ready` tinyint(1) DEFAULT 0 COMMENT '是否就绪',
  `schedulable` tinyint(1) DEFAULT 1 COMMENT '是否可调度',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_name` (`config_id`, `name`),
  KEY `idx_config_id` (`config_id`),
  KEY `idx_status` (`status`),
  KEY `idx_ready` (`ready`),
  KEY `idx_internal_ip` (`internal_ip`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Kubernetes节点表';

-- 创建 Kubernetes 工作负载表
CREATE TABLE IF NOT EXISTS `infra_k8s_workload` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `config_id` bigint NOT NULL COMMENT '集群配置ID',
  `name` varchar(100) NOT NULL COMMENT '工作负载名称',
  `namespace` varchar(63) NOT NULL COMMENT '命名空间',
  `kind` varchar(20) NOT NULL COMMENT '工作负载类型',
  `replicas` int DEFAULT '0' COMMENT '副本数',
  `ready_replicas` int DEFAULT '0' COMMENT '就绪副本数',
  `status` varchar(20) DEFAULT NULL COMMENT '状态',
  `labels` text DEFAULT NULL COMMENT '标签(JSON格式)',
  `selector` text DEFAULT NULL COMMENT '选择器(JSON格式)',
  `images` text DEFAULT NULL COMMENT '容器镜像列表(JSON格式)',
  `cpu_request` varchar(20) DEFAULT NULL COMMENT 'CPU请求',
  `cpu_limit` varchar(20) DEFAULT NULL COMMENT 'CPU限制',
  `memory_request` varchar(20) DEFAULT NULL COMMENT '内存请求',
  `memory_limit` varchar(20) DEFAULT NULL COMMENT '内存限制',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_name_namespace_kind` (`config_id`, `name`, `namespace`, `kind`),
  KEY `idx_config_id` (`config_id`),
  KEY `idx_namespace` (`namespace`),
  KEY `idx_kind` (`kind`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Kubernetes工作负载';

-- 创建K8s命名空间表
CREATE TABLE IF NOT EXISTS `infra_k8s_namespace` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '命名空间ID',
  `config_id` bigint NOT NULL COMMENT 'K8s配置ID',
  `namespace` varchar(63) NOT NULL COMMENT '命名空间名称',
  `workload_count` int DEFAULT 0 COMMENT '工作负载数量',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_namespace` (`config_id`, `namespace`),
  KEY `idx_config_id` (`config_id`),
  KEY `idx_namespace` (`namespace`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='K8s命名空间表';

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
(@infra_id, '云厂商管理', 'cloud-provider', 'infrastructure/cloud-provider/index', 'infrastructure:cloud-provider:list', 1, 'cloud', 4, 1),
(@infra_id, '云账号管理', 'cloud-account', 'infrastructure/cloud-account/index', 'infrastructure:cloud-account:list', 1, 'cloudy', 5, 1);

-- 获取各个子菜单的ID
SET @k8s_menu_id = NULL;
SET @server_menu_id = NULL;
SET @db_menu_id = NULL;
SET @cloud_provider_menu_id = NULL;
SET @cloud_account_menu_id = NULL;

SELECT @k8s_menu_id := id FROM `sys_menu` WHERE `path` = 'kubernetes' AND parent_id = @infra_id;
SELECT @server_menu_id := id FROM `sys_menu` WHERE `path` = 'server' AND parent_id = @infra_id;
SELECT @db_menu_id := id FROM `sys_menu` WHERE `path` = 'database' AND parent_id = @infra_id;
SELECT @cloud_provider_menu_id := id FROM `sys_menu` WHERE `path` = 'cloud-provider' AND parent_id = @infra_id;
SELECT @cloud_account_menu_id := id FROM `sys_menu` WHERE `path` = 'cloud-account' AND parent_id = @infra_id;

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
(@cloud_provider_menu_id, '云厂商列表', 'infrastructure:cloud-provider:list', 2, NULL, 1, 1),
(@cloud_provider_menu_id, '云厂商创建', 'infrastructure:cloud-provider:create', 2, NULL, 2, 1),
(@cloud_provider_menu_id, '云厂商更新', 'infrastructure:cloud-provider:update', 2, NULL, 3, 1),
(@cloud_provider_menu_id, '云厂商删除', 'infrastructure:cloud-provider:delete', 2, NULL, 4, 1);

-- Add CRUD permissions for Cloud Account
INSERT INTO `sys_menu` (`parent_id`, `name`, `perms`, `type`, `icon`, `sort_order`, `status`)
VALUES 
(@cloud_account_menu_id, '云账号列表', 'infrastructure:cloud-account:list', 2, NULL, 1, 1),
(@cloud_account_menu_id, '云账号创建', 'infrastructure:cloud-account:create', 2, NULL, 2, 1),
(@cloud_account_menu_id, '云账号更新', 'infrastructure:cloud-account:update', 2, NULL, 3, 1),
(@cloud_account_menu_id, '云账号删除', 'infrastructure:cloud-account:delete', 2, NULL, 4, 1);

-- Add CRUD permissions for Kubernetes
INSERT INTO `sys_menu` (`parent_id`, `name`, `perms`, `type`, `icon`, `sort_order`, `status`)
VALUES 
(@k8s_menu_id, 'Kubernetes列表', 'infrastructure:kubernetes:list', 2, NULL, 1, 1),
(@k8s_menu_id, 'Kubernetes创建', 'infrastructure:kubernetes:create', 2, NULL, 2, 1),
(@k8s_menu_id, 'Kubernetes更新', 'infrastructure:kubernetes:update', 2, NULL, 3, 1),
(@k8s_menu_id, 'Kubernetes删除', 'infrastructure:kubernetes:delete', 2, NULL, 4, 1);

-- Add CRUD permissions for Database
INSERT INTO `sys_menu` (`parent_id`, `name`, `perms`, `type`, `icon`, `sort_order`, `status`)
VALUES 
(@db_menu_id, '数据库列表', 'infrastructure:database:list', 2, NULL, 1, 1),
(@db_menu_id, '数据库创建', 'infrastructure:database:create', 2, NULL, 2, 1),
(@db_menu_id, '数据库更新', 'infrastructure:database:update', 2, NULL, 3, 1),
(@db_menu_id, '数据库删除', 'infrastructure:database:delete', 2, NULL, 4, 1);

-- Add CRUD permissions for Server
INSERT INTO `sys_menu` (`parent_id`, `name`, `perms`, `type`, `icon`, `sort_order`, `status`)
VALUES
(@server_menu_id, '服务器列表', 'infrastructure:server:list', 2, NULL, 1, 1),
(@server_menu_id, '服务器创建', 'infrastructure:server:create', 2, NULL, 2, 1),
(@server_menu_id, '服务器更新', 'infrastructure:server:update', 2, NULL, 3, 1),
(@server_menu_id, '服务器删除', 'infrastructure:server:delete', 2, NULL, 4, 1);

-- 创建K8s Pod表
CREATE TABLE IF NOT EXISTS `infra_k8s_pod` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'Pod ID',
  `config_id` bigint NOT NULL COMMENT 'K8s配置ID',
  `workload_id` bigint DEFAULT NULL COMMENT '工作负载ID',
  `name` varchar(63) NOT NULL COMMENT 'Pod名称',
  `namespace` varchar(63) NOT NULL COMMENT '命名空间',
  `workload_name` varchar(100) DEFAULT NULL COMMENT '工作负载名称',
  `workload_kind` varchar(50) DEFAULT NULL COMMENT '工作负载类型',
  `status` varchar(50) NOT NULL COMMENT 'Pod状态',
  `phase` varchar(50) DEFAULT NULL COMMENT 'Pod阶段',
  `node_name` varchar(100) DEFAULT NULL COMMENT '节点名称',
  `pod_ip` varchar(45) DEFAULT NULL COMMENT 'Pod IP',
  `host_ip` varchar(45) DEFAULT NULL COMMENT '主机IP',
  `instance_ip` varchar(45) DEFAULT NULL COMMENT '实例IP',
  `cpu_request` varchar(20) DEFAULT NULL COMMENT 'CPU请求',
  `cpu_limit` varchar(20) DEFAULT NULL COMMENT 'CPU限制',
  `memory_request` varchar(20) DEFAULT NULL COMMENT '内存请求',
  `memory_limit` varchar(20) DEFAULT NULL COMMENT '内存限制',
  `restart_count` int DEFAULT 0 COMMENT '重启次数',
  `start_time` datetime DEFAULT NULL COMMENT '启动时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_name_namespace_instance_ip` (`config_id`, `name`, `namespace`, `instance_ip`),
  KEY `idx_config_id` (`config_id`),
  KEY `idx_namespace` (`namespace`),
  KEY `idx_status` (`status`),
  KEY `idx_workload_name` (`workload_name`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='K8s Pod表';

-- ==================== 历史表定义 ====================

-- 创建K8s Pod历史表
CREATE TABLE IF NOT EXISTS `infra_k8s_pod_history` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '历史记录ID',
  `original_id` bigint NOT NULL COMMENT '原始Pod ID',
  `config_id` bigint NOT NULL COMMENT 'K8s配置ID',
  `workload_id` bigint DEFAULT NULL COMMENT '工作负载ID',
  `name` varchar(100) NOT NULL COMMENT 'Pod名称',
  `namespace` varchar(100) NOT NULL COMMENT '命名空间',
  `workload_name` varchar(100) DEFAULT NULL COMMENT '工作负载名称',
  `workload_kind` varchar(50) DEFAULT NULL COMMENT '工作负载类型',
  `status` varchar(50) NOT NULL COMMENT 'Pod状态',
  `phase` varchar(50) DEFAULT NULL COMMENT 'Pod阶段',
  `node_name` varchar(100) DEFAULT NULL COMMENT '节点名称',
  `pod_ip` varchar(45) DEFAULT NULL COMMENT 'Pod IP',
  `host_ip` varchar(45) DEFAULT NULL COMMENT '主机IP',
  `instance_ip` varchar(45) DEFAULT NULL COMMENT '实例IP',
  `cpu_request` varchar(20) DEFAULT NULL COMMENT 'CPU请求',
  `cpu_limit` varchar(20) DEFAULT NULL COMMENT 'CPU限制',
  `memory_request` varchar(20) DEFAULT NULL COMMENT '内存请求',
  `memory_limit` varchar(20) DEFAULT NULL COMMENT '内存限制',
  `restart_count` int DEFAULT 0 COMMENT '重启次数',
  `start_time` datetime DEFAULT NULL COMMENT '启动时间',
  `created_at` datetime NOT NULL COMMENT '原始创建时间',
  `updated_at` datetime NOT NULL COMMENT '原始更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '原始删除时间',
  `archived_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '归档时间',
  `archive_reason` varchar(100) DEFAULT 'sync_cleanup' COMMENT '归档原因：sync_cleanup=同步清理, manual=手动归档',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_name_namespace_created_at` (`config_id`, `name`, `namespace`, `created_at`),
  KEY `idx_original_id` (`original_id`),
  KEY `idx_config_id` (`config_id`),
  KEY `idx_namespace` (`namespace`),
  KEY `idx_archived_at` (`archived_at`),
  KEY `idx_config_name` (`config_id`, `name`(100)),
  KEY `idx_archive_reason` (`archive_reason`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='K8s Pod历史表';

-- 创建K8s Node历史表
CREATE TABLE IF NOT EXISTS `infra_k8s_node_history` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '历史记录ID',
  `original_id` bigint NOT NULL COMMENT '原始节点ID',
  `config_id` bigint NOT NULL COMMENT 'K8s配置ID',
  `name` varchar(100) NOT NULL COMMENT '节点名称',
  `internal_ip` varchar(45) DEFAULT NULL COMMENT '内部IP',
  `external_ip` varchar(45) DEFAULT NULL COMMENT '外部IP',
  `hostname` varchar(100) DEFAULT NULL COMMENT '主机名',
  `os_image` varchar(100) DEFAULT NULL COMMENT '操作系统镜像',
  `kernel_version` varchar(100) DEFAULT NULL COMMENT '内核版本',
  `container_runtime` varchar(100) DEFAULT NULL COMMENT '容器运行时',
  `kubelet_version` varchar(50) DEFAULT NULL COMMENT 'Kubelet版本',
  `kube_proxy_version` varchar(50) DEFAULT NULL COMMENT 'KubeProxy版本',
  `cpu_capacity` varchar(20) DEFAULT NULL COMMENT 'CPU容量',
  `memory_capacity` varchar(20) DEFAULT NULL COMMENT '内存容量',
  `pods_capacity` varchar(20) DEFAULT NULL COMMENT 'Pod容量',
  `cpu_allocatable` varchar(20) DEFAULT NULL COMMENT 'CPU可分配',
  `memory_allocatable` varchar(20) DEFAULT NULL COMMENT '内存可分配',
  `pods_allocatable` varchar(20) DEFAULT NULL COMMENT 'Pod可分配',
  `cpu_usage` varchar(20) DEFAULT NULL COMMENT 'CPU使用量',
  `memory_usage` varchar(20) DEFAULT NULL COMMENT '内存使用量',
  `pods_usage` int DEFAULT 0 COMMENT 'Pod使用量',
  `labels` text DEFAULT NULL COMMENT '标签(JSON格式)',
  `annotations` text DEFAULT NULL COMMENT '注解(JSON格式)',
  `taints` text DEFAULT NULL COMMENT '污点(JSON格式)',
  `conditions` text DEFAULT NULL COMMENT '状态条件(JSON格式)',
  `status` varchar(50) DEFAULT 'Unknown' COMMENT '节点状态',
  `ready` tinyint(1) DEFAULT 0 COMMENT '是否就绪',
  `schedulable` tinyint(1) DEFAULT 1 COMMENT '是否可调度',
  `created_at` datetime NOT NULL COMMENT '原始创建时间',
  `updated_at` datetime NOT NULL COMMENT '原始更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '原始删除时间',
  `archived_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '归档时间',
  `archive_reason` varchar(100) DEFAULT 'sync_cleanup' COMMENT '归档原因：sync_cleanup=同步清理, manual=手动归档',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_name_created_at` (`config_id`, `name`, `created_at`),
  KEY `idx_original_id` (`original_id`),
  KEY `idx_config_id` (`config_id`),
  KEY `idx_name` (`name`(100)),
  KEY `idx_archived_at` (`archived_at`),
  KEY `idx_config_name` (`config_id`, `name`(100)),
  KEY `idx_internal_ip` (`internal_ip`),
  KEY `idx_status` (`status`),
  KEY `idx_archive_reason` (`archive_reason`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='K8s Node历史表';

-- 创建K8s Workload历史表
CREATE TABLE IF NOT EXISTS `infra_k8s_workload_history` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '历史记录ID',
  `original_id` bigint NOT NULL COMMENT '原始工作负载ID',
  `config_id` bigint NOT NULL COMMENT '集群配置ID',
  `name` varchar(100) NOT NULL COMMENT '工作负载名称',
  `namespace` varchar(100) NOT NULL COMMENT '命名空间',
  `kind` varchar(20) NOT NULL COMMENT '工作负载类型',
  `replicas` int DEFAULT '0' COMMENT '副本数',
  `ready_replicas` int DEFAULT '0' COMMENT '就绪副本数',
  `status` varchar(20) DEFAULT NULL COMMENT '状态',
  `labels` text DEFAULT NULL COMMENT '标签(JSON格式)',
  `selector` text DEFAULT NULL COMMENT '选择器(JSON格式)',
  `images` text DEFAULT NULL COMMENT '容器镜像列表(JSON格式)',
  `cpu_request` varchar(20) DEFAULT NULL COMMENT 'CPU请求',
  `cpu_limit` varchar(20) DEFAULT NULL COMMENT 'CPU限制',
  `memory_request` varchar(20) DEFAULT NULL COMMENT '内存请求',
  `memory_limit` varchar(20) DEFAULT NULL COMMENT '内存限制',
  `created_at` datetime NOT NULL COMMENT '原始创建时间',
  `updated_at` datetime NOT NULL COMMENT '原始更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '原始删除时间',
  `archived_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '归档时间',
  `archive_reason` varchar(100) DEFAULT 'sync_cleanup' COMMENT '归档原因：sync_cleanup=同步清理, manual=手动归档',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_name_namespace_kind_created_at` (`config_id`, `name`, `namespace`, `kind`, `created_at`),
  KEY `idx_original_id` (`original_id`),
  KEY `idx_config_id` (`config_id`),
  KEY `idx_namespace` (`namespace`),
  KEY `idx_kind` (`kind`),
  KEY `idx_archived_at` (`archived_at`),
  KEY `idx_config_name` (`config_id`, `name`(100)),
  KEY `idx_name_namespace` (`name`(100), `namespace`),
  KEY `idx_archive_reason` (`archive_reason`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='K8s Workload历史表';