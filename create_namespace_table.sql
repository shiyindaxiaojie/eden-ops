CREATE TABLE IF NOT EXISTS `infra_k8s_workload_namespace` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '命名空间ID',
  `config_id` bigint NOT NULL COMMENT 'K8s配置ID',
  `namespace` varchar(255) NOT NULL COMMENT '命名空间名称',
  `workload_count` int DEFAULT 0 COMMENT '工作负载数量',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_namespace` (`config_id`, `namespace`),
  KEY `idx_config_id` (`config_id`),
  KEY `idx_namespace` (`namespace`),
  KEY `idx_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_k8s_namespace_config` FOREIGN KEY (`config_id`) REFERENCES `infra_k8s_config` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='K8s工作负载命名空间表';
