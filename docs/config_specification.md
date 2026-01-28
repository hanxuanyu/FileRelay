# FileRelay 配置项详细说明文档

本文档整理了 FileRelay 系统 `config.yaml` 配置文件中各字段的含义、类型及示例，供前端配置页面开发参考。

## 1. 站点设置 (site)
用于定义前端展示的站点基本信息。

| 字段名 | 类型 | 含义 | 示例 |
| :--- | :--- | :--- | :--- |
| `name` | string | 站点名称，显示在网页标题和页头 | `文件暂存柜` |
| `description` | string | 站点描述，显示在首页或元标签中 | `临时文件中转服务` |
| `logo` | string | 站点 Logo 的 URL 地址 | `/logo.png` |
| `base_url` | string | 站点外部访问地址。若配置则固定使用该地址拼接直链；若留空，系统将尝试从请求头（如 `X-Forwarded-Proto`, `:scheme`, `Forwarded` 等）或 `Referer` 中自动检测协议及主机名，以确保在 HTTPS 代理环境下链接正确。 | `https://file.example.com` |
| `port` | int | 后端服务监听端口 | `8080` |

## 2. 安全设置 (security)
涉及系统鉴权、取件保护相关的核心安全配置。

| 字段名 | 类型 | 含义 | 示例 |
| :--- | :--- | :--- | :--- |
| `admin_password_hash` | string | 管理员密码的 bcrypt 哈希值。可以通过更新配置接口修改，修改后立即生效，且不再依赖数据库存储。 | `$2a$10$...` |
| `pickup_code_length` | int | 自动生成的取件码长度。变更后系统将自动对存量取件码进行右侧补零或截取以适配新长度。 | `6` |
| `pickup_fail_limit` | int | 单个 IP 对单个取件码尝试失败的最大次数，超过后将被临时封禁 | `5` |
| `jwt_secret` | string | 用于签发管理端 JWT Token 的密钥，建议设置为复杂随机字符串 | `file-relay-secret` |

## 3. 上传设置 (upload)
控制文件上传的限制和策略。

| 字段名 | 类型 | 含义 | 示例 |
| :--- | :--- | :--- | :--- |
| `max_file_size_mb` | int64 | 单个文件的最大允许大小（单位：MB） | `100` |
| `max_batch_files` | int | 一个取件批次中允许包含的最大文件数量 | `20` |
| `max_retention_days` | int | 文件在服务器上的最长保留天数（针对 time 类型的过期策略） | `30` |
| `require_token` | bool | 是否强制要求提供 API Token 才能进行上传操作 | `false` |

## 4. 存储设置 (storage)
定义文件的实际物理存储方式。系统支持多种存储后端。

| 字段名 | 类型 | 含义 | 示例 |
| :--- | :--- | :--- | :--- |
| `type` | string | 当前激活的存储类型。可选值：`local`, `webdav`, `s3` | `local` |

### 4.1 本地存储 (local)
当 `type` 为 `local` 时生效。

| 字段名 | 类型 | 含义 | 示例 |
| :--- | :--- | :--- | :--- |
| `path` | string | 文件存储在服务器本地的相对或绝对路径 | `data/storage_data` |

### 4.2 WebDAV 存储 (webdav)
当 `type` 为 `webdav` 时生效。

| 字段名 | 类型 | 含义 | 示例 |
| :--- | :--- | :--- | :--- |
| `url` | string | WebDAV 服务器的 API 地址 | `https://dav.example.com` |
| `username` | string | WebDAV 登录用户名 | `user` |
| `password` | string | WebDAV 登录密码 | `pass` |
| `root` | string | WebDAV 上的基础存储根目录 | `/file-relay` |

### 4.3 S3 存储 (s3)
当 `type` 为 `s3` 时生效。

| 字段名 | 类型 | 含义 | 示例 |
| :--- | :--- | :--- | :--- |
| `endpoint` | string | S3 服务端点 | `s3.amazonaws.com` |
| `region` | string | S3 区域 | `us-east-1` |
| `access_key` | string | S3 Access Key | `your-access-key` |
| `secret_key` | string | S3 Secret Key | `your-secret-key` |
| `bucket` | string | S3 存储桶名称 | `file-relay-bucket` |
| `use_ssl` | bool | 是否强制使用 SSL (HTTPS) 连接 | `false` |

## 5. API Token 设置 (api_token)
控制系统对外开放的 API Token 管理功能。

| 字段名 | 类型 | 含义 | 示例 |
| :--- | :--- | :--- | :--- |
| `enabled` | bool | 是否启用 API Token 功能模块 | `true` |
| `allow_admin_api` | bool | 是否允许具备 `admin` 权限的 API Token 访问管理接口 | `true` |
| `max_tokens` | int | 系统允许创建的 API Token 最大总数限制 | `20` |

### 5.1 API Token 权限说明 (Scopes)
在创建 API Token 时，可以通过 `scope` 字段赋予以下一种或多种权限（多个权限用逗号分隔，如 `upload,pickup`）：

| 权限值 | 含义 | 说明 |
| :--- | :--- | :--- |
| `upload` | 上传权限 | 允许调用文件和长文本上传接口 |
| `pickup` | 取件权限 | 允许获取批次详情、下载文件及查询下载次数 |
| `admin` | 管理权限 | 允许访问管理端（Admin）所有接口。需开启 `allow_admin_api` 且 Token 功能已启用 |

## 6. Web 前端设置 (web)
定义前端静态资源的加载方式。

| 字段名 | 类型 | 含义 | 示例 |
| :--- | :--- | :--- | :--- |
| `path` | string | 外部前端资源目录路径。若该路径存在且包含 `index.html`，系统将优先使用此目录；否则回退使用内置前端资源。 | `web` |

## 7. 数据库设置 (database)
系统元数据存储配置。支持 SQLite, MySQL 和 PostgreSQL。

**特性说明**：
- **动态切换**：通过管理接口更新数据库配置后，系统会自动尝试连接新数据库，无需重启应用。
- **自动迁移**：切换数据库时，系统会自动将原数据库中的元数据（如存量批次、文件记录、API Token 等）同步迁移至新数据库中。

| 字段名 | 类型 | 含义 | 示例 |
| :--- | :--- | :--- | :--- |
| `type` | string | 数据库类型。可选值：`sqlite`, `mysql`, `postgres` | `sqlite` |
| `path` | string | SQLite 数据库文件的路径 (仅在 `type` 为 `sqlite` 时生效) | `data/file_relay.db` |
| `host` | string | 数据库地址 (MySQL/PostgreSQL) | `127.0.0.1` |
| `port` | int | 数据库端口 (MySQL/PostgreSQL) | `3306` 或 `5432` |
| `user` | string | 数据库用户名 (MySQL/PostgreSQL) | `root` |
| `password` | string | 数据库密码 (MySQL/PostgreSQL) | `password` |
| `dbname` | string | 数据库名称 (MySQL/PostgreSQL) | `file_relay` |
| `config` | string | 额外 DSN 配置参数。MySQL 如 `charset=utf8mb4&parseTime=True&loc=Local`；PostgreSQL 如 `sslmode=disable` | `charset=utf8mb4&parseTime=True&loc=Local` |

## 8. 日志设置 (log)
控制系统日志的输出级别和目的地。

| 字段名 | 类型 | 含义 | 示例 |
| :--- | :--- | :--- | :--- |
| `level` | string | 日志级别。可选值：`debug`, `info`, `warn`, `error` | `info` |
| `file_path` | string | 日志文件路径。如果设置，日志将同时输出到控制台和该文件；若为空则仅输出到控制台。 | `data/logs/app.log` |

---

## 附录：公共配置接口 (/api/config)

为了方便前端展示和交互约束，系统提供了 `/api/config` 接口，该接口不需要鉴权，返回以下非敏感字段（结构与完整配置保持一致）：

- **site**: 完整内容（`name`, `description`, `logo`, `base_url`）
- **security**: 仅包含 `pickup_code_length`
- **upload**: 完整内容（`max_file_size_mb`, `max_batch_files`, `max_retention_days`, `require_token`）
- **api_token**: 仅包含 `enabled` 开关
- **storage**: 仅包含 `type`（存储类型）

## 附录：其他关键接口

### 查询下载次数 (GET /api/batches/:pickup_code/count)
- **用途**：供前端实时刷新当前取件批次的下载次数。
- **特性**：支持查询已过期的文件。

### 恢复 API Token (POST /api/admin/api-tokens/:id/recover)
- **权限**：需要管理员权限。
- **用途**：将状态为“已撤销”的 Token 重新恢复为有效状态。
