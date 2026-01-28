-- SQLite 初始化语句
CREATE TABLE IF NOT EXISTS `file_batches` (
    `id` varchar(36) PRIMARY KEY,
    `pickup_code` varchar(255) UNIQUE NOT NULL,
    `remark` text,
    `expire_type` text,
    `expire_at` datetime,
    `max_downloads` integer,
    `download_count` integer DEFAULT 0,
    `status` text DEFAULT 'active',
    `type` text DEFAULT 'file',
    `content` text,
    `created_at` datetime,
    `updated_at` datetime,
    `deleted_at` datetime
);
CREATE INDEX IF NOT EXISTS `idx_file_batches_deleted_at` ON `file_batches` (`deleted_at`);

CREATE TABLE IF NOT EXISTS `file_items` (
    `id` varchar(36) PRIMARY KEY,
    `batch_id` varchar(36) NOT NULL,
    `original_name` text,
    `storage_path` text,
    `size` bigint,
    `mime_type` text,
    `created_at` datetime,
    FOREIGN KEY (`batch_id`) REFERENCES `file_batches`(`id`)
);
CREATE INDEX IF NOT EXISTS `idx_file_items_batch_id` ON `file_items` (`batch_id`);

CREATE TABLE IF NOT EXISTS `api_tokens` (
    `id` integer PRIMARY KEY AUTOINCREMENT,
    `name` text,
    `token_hash` varchar(255) UNIQUE NOT NULL,
    `scope` text,
    `expire_at` datetime,
    `last_used_at` datetime,
    `revoked` boolean DEFAULT 0,
    `created_at` datetime
);
