# FileRelay - 文件暂存柜

[![Build Status](https://github.com/hanxuanyu/FileRelay/actions/workflows/build.yml/badge.svg)](https://github.com/hanxuanyu/FileRelay/actions/workflows/build.yml)
[![Release Status](https://github.com/hanxuanyu/FileRelay/actions/workflows/release.yml/badge.svg)](https://github.com/hanxuanyu/FileRelay/actions/workflows/release.yml)
[![Docker Image](https://github.com/hanxuanyu/FileRelay/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/hanxuanyu/FileRelay/actions/workflows/docker-publish.yml)

中文 | [English](README_en.md)

---

FileRelay 是一个简单高效、自托管的文件暂存柜系统，旨在为您提供便捷的临时文件分享和中转服务。

### 🚀 主要特性

- **临时分享**：上传文件或文本，生成唯一取件码，方便在不同设备间快速传输。
- **多数据库支持**：支持 SQLite, MySQL 和 PostgreSQL，默认为 SQLite 开箱即用，生产环境可无缝切换至 MySQL 或 PostgreSQL。
- **多种存储后端**：支持本地磁盘存储、Amazon S3（及兼容 S3 的服务）以及 WebDAV 存储。
- **灵活的过期策略**：支持按时间过期和按下载次数过期（需在管理后台配置或手动清理）。
- **取件保护**：支持自定义取件码长度，内置尝试次数限制防暴力破解。
- **管理后台**：内置美观的管理界面，支持：
    - 动态修改系统配置。
    - 实时查看和管理上传批次。
    - API Token 管理（支持权限细分：上传、取件、管理）。
- **高性能**：采用 Go 语言编写，内存占用极低，支持高并发访问。
- **易于部署**：支持单一二进制文件运行，前端资源已嵌入，无需额外配置 Web 服务器。

### 🛠️ 技术栈

- **后端**: Go 1.24+ (Gin, GORM)
- **数据库**: SQLite, MySQL, PostgreSQL (支持多种数据库，灵活扩展)
- **前端**: Vue 3 + TailwindCSS (位于 `webapp` 目录)，已通过 `embed` 嵌入二进制。
- **文档**: 集成 Swagger UI。

### 🏗️ 开发与构建

#### 环境要求

- **Go**: 1.24 或更高版本
- **Node.js**: 20 或更高版本 (用于前端构建)
- **npm**: 随 Node.js 安装

#### 本地开发

1. **前端开发**:
   ```bash
   cd webapp
   npm install
   npm run dev
   ```
   前端开发服务器默认运行在 `http://localhost:5173`。

2. **后端开发**:
   ```bash
   # 返回项目根目录
   go run main.go -config config/config.yaml
   ```

#### 完整构建

项目提供了自动化的构建脚本，会自动完成前端构建、资源嵌入以及 Go 编译：

- **Windows (CMD)**: `scripts\build.bat`
- **Windows (PowerShell)**: `.\scripts\build.ps1`
- **Linux/macOS**: `chmod +x scripts/build.sh && ./scripts/build.sh`

构建产物将存放在 `output` 目录下。

### 📦 快速开始

#### 1. 获取程序

您可以从 Release 页面下载对应平台的二进制文件，或者按照上述“完整构建”步骤自行编译。

#### 2. 配置文件

在运行之前，请根据需求修改 `config/config.yaml`。详细配置说明请参考 [docs/config_specification.md](docs/config_specification.md)。

主要配置项包括：
- 存储类型 (`local`, `s3`, `webdav`)。
- 管理员密码哈希。
- 数据库配置 (`type`, `host`, `port`, `user`, `password`, `dbname` 等)。系统支持自动迁移，首次连接新数据库时将自动创建表结构。此外，`scripts/sql/` 目录下也提供了各数据库的初始化 SQL 脚本供参考。

#### 3. 运行

```bash
./filerelay -config config/config.yaml
```

### 🐳 Docker 部署

我们提供了 Docker 和 Docker Compose 支持，实现一键部署：

#### 使用 Docker 镜像 (推荐)

如果您不想自行构建，可以直接从 Docker Hub 拉取已构建好的镜像：

```bash
docker run -d \
  --name filerelay \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/config:/app/config \
  hxuanyu521/filerelay:latest
```

**提示：**
- 如果您需要使用自定义的 Web 页面（外置前端资源），可以将资源目录挂载到容器内，并通过环境变量 `FR_WEB_PATH` 指定路径。例如：
  `-v $(pwd)/my-web:/app/my-web -e FR_WEB_PATH=/app/my-web`
- 默认情况下，程序会优先寻找 `FR_WEB_PATH`（或配置文件中的 `web.path`）指定的外部资源，若找不到则使用内置的嵌入资源。

#### 使用 Docker Compose

1. 确保已安装 Docker 和 Docker Compose。
2. 在项目根目录下运行：

```bash
docker-compose up -d
```

程序默认会在当前目录下的 `./data` 文件夹中持久化存储数据、日志、数据库和配置。

#### 环境变量控制

支持通过环境变量覆盖配置文件中的关键项，方便在 Docker 环境下动态调整：

| 环境变量 | 说明 | 默认值 |
| :--- | :--- | :--- |
| **站点设置** | | |
| `FR_SITE_NAME` | 站点名称 | 文件暂存柜 |
| `FR_SITE_BASE_URL` | 站点外部访问地址 (用于生成分享链接) | (空) |
| `FR_SITE_PORT` | 服务监听端口 | 8080 |
| **安全设置** | | |
| `FR_SECURITY_JWT_SECRET` | JWT 签名密钥 | file-relay-secret |
| **上传设置** | | |
| `FR_UPLOAD_MAX_SIZE` | 单个文件最大大小 (MB) | 100 |
| `FR_UPLOAD_RETENTION_DAYS` | 文件最大保留天数 | 30 |
| **数据库设置** | | |
| `FR_DB_TYPE` | 数据库类型 (sqlite, mysql, postgres) | sqlite |
| `FR_DB_PATH` | SQLite 数据库文件路径 | data/file_relay.db |
| `FR_DB_HOST` | 数据库主机地址 (MySQL/Postgres) | (空) |
| `FR_DB_PORT` | 数据库端口 (MySQL/Postgres) | (空) |
| `FR_DB_USER` | 数据库用户名 (MySQL/Postgres) | (空) |
| `FR_DB_PASSWORD` | 数据库密码 (MySQL/Postgres) | (空) |
| `FR_DB_NAME` | 数据库名称 (MySQL/Postgres) | (空) |
| **存储设置** | | |
| `FR_STORAGE_TYPE` | 存储类型 (local, s3, webdav) | local |
| `FR_STORAGE_LOCAL_PATH` | 本地存储路径 | data/storage_data |
| **日志设置** | | |
| `FR_LOG_LEVEL` | 日志级别 (debug, info, warn, error) | info |
| `FR_LOG_FILE_PATH` | 日志文件路径 | data/logs/app.log |
| **Web 设置** | | |
| `FR_WEB_PATH` | 外部 Web 静态资源路径 | web |

程序启动后，您可以通过以下地址访问（默认端口 8080）：
- **Web 界面**: `http://localhost:8080`
- **管理后台**: `http://localhost:8080/admin` (前端路由)
- **API 文档**: `http://localhost:8080/swagger/index.html`

### 📖 接口说明

FileRelay 提供了丰富的 API 接口，支持通过 API Token 进行集成。

- **上传**: `POST /api/batches`
- **取件**: `GET /api/batches/:pickup_code`
- **下载**: `GET /api/files/:file_id/download`

详细接口定义请参考内置的 Swagger 文档。
