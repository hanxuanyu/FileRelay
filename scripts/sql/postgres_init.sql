-- PostgreSQL 初始化语句
CREATE TABLE IF NOT EXISTS "file_batches" (
    "id" varchar(36) PRIMARY KEY,
    "pickup_code" varchar(255) UNIQUE NOT NULL,
    "remark" text,
    "expire_type" text,
    "expire_at" timestamptz,
    "max_downloads" bigint,
    "download_count" bigint DEFAULT 0,
    "status" text DEFAULT 'active',
    "type" text DEFAULT 'file',
    "content" text,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz
);
CREATE INDEX IF NOT EXISTS "idx_file_batches_deleted_at" ON "file_batches" ("deleted_at");

CREATE TABLE IF NOT EXISTS "file_items" (
    "id" varchar(36) PRIMARY KEY,
    "batch_id" varchar(36) NOT NULL,
    "original_name" text,
    "storage_path" text,
    "size" bigint,
    "mime_type" text,
    "created_at" timestamptz,
    CONSTRAINT "fk_file_batches_file_items" FOREIGN KEY ("batch_id") REFERENCES "file_batches"("id")
);
CREATE INDEX IF NOT EXISTS "idx_file_items_batch_id" ON "file_items" ("batch_id");

CREATE TABLE IF NOT EXISTS "api_tokens" (
    "id" bigserial PRIMARY KEY,
    "name" text,
    "token_hash" varchar(255) UNIQUE NOT NULL,
    "scope" text,
    "expire_at" timestamptz,
    "last_used_at" timestamptz,
    "revoked" boolean DEFAULT false,
    "created_at" timestamptz
);
