-- MySQL 初始化语句
CREATE TABLE IF NOT EXISTS `file_batches` (
    `id` varchar(36) PRIMARY KEY,
    `pickup_code` varchar(255) UNIQUE NOT NULL,
    `remark` longtext,
    `expire_type` varchar(255),
    `expire_at` datetime(3),
    `max_downloads` bigint,
    `download_count` bigint DEFAULT 0,
    `status` varchar(255) DEFAULT 'active',
    `type` varchar(255) DEFAULT 'file',
    `content` longtext,
    `created_at` datetime(3),
    `updated_at` datetime(3),
    `deleted_at` datetime(3),
    KEY `idx_file_batches_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `file_items` (
    `id` varchar(36) PRIMARY KEY,
    `batch_id` varchar(36) NOT NULL,
    `original_name` longtext,
    `storage_path` longtext,
    `size` bigint,
    `mime_type` longtext,
    `created_at` datetime(3),
    KEY `idx_file_items_batch_id` (`batch_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `api_tokens` (
    `id` bigint unsigned AUTO_INCREMENT PRIMARY KEY,
    `name` longtext,
    `token_hash` varchar(255) UNIQUE NOT NULL,
    `scope` longtext,
    `expire_at` datetime(3),
    `last_used_at` datetime(3),
    `revoked` boolean DEFAULT false,
    `created_at` datetime(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
