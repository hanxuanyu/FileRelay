# FileRelay - File Relay Station

[![Build Status](https://github.com/hanxuanyu/FileRelay/actions/workflows/build.yml/badge.svg)](https://github.com/hanxuanyu/FileRelay/actions/workflows/build.yml)
[![Release Status](https://github.com/hanxuanyu/FileRelay/actions/workflows/release.yml/badge.svg)](https://github.com/hanxuanyu/FileRelay/actions/workflows/release.yml)
[![Docker Image](https://github.com/hanxuanyu/FileRelay/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/hanxuanyu/FileRelay/actions/workflows/docker-publish.yml)

[中文](README.md) | English

---

FileRelay is a simple, efficient, self-hosted file relay/temporary storage system designed to provide convenient temporary file sharing and transfer services.

### 🚀 Key Features

- **Temporary Sharing**: Upload files or text and generate a unique pickup code for quick transfer across different devices.
- **Multi-database Support**: Supports SQLite, MySQL, and PostgreSQL. Defaults to SQLite for out-of-the-box use, with seamless switching to MySQL or PostgreSQL for production.
- **Multiple Storage Backends**: Supports local disk storage, Amazon S3 (and S3-compatible services), and WebDAV.
- **Flexible Expiration**: Supports expiration by time and download counts (configurable in admin or via manual cleanup).
- **Pickup Protection**: Customizable pickup code length and built-in rate limiting to prevent brute-force attacks.
- **Admin Dashboard**: Built-in management interface for:
    - Dynamic system configuration updates.
    - Real-time batch management (view, delete, clean).
    - API Token management with granular scopes (upload, pickup, admin).
- **High Performance**: Written in Go with low memory footprint and high concurrency support.
- **Easy Deployment**: Single binary execution with embedded frontend assets, no extra Web server needed.

### 🛠️ Tech Stack

- **Backend**: Go 1.24+ (Gin, GORM)
- **Database**: SQLite, MySQL, PostgreSQL (supports multiple databases for scalability)
- **Frontend**: Vue 3 + TailwindCSS (located in `webapp` directory), embedded into binary via `embed`.
- **Documentation**: Integrated Swagger UI.

### 🏗️ Development & Build

#### Prerequisites

- **Go**: 1.24 or higher
- **Node.js**: 20 or higher (for frontend build)
- **npm**: Installed with Node.js

#### Local Development

1. **Frontend**:
   ```bash
   cd webapp
   npm install
   npm run dev
   ```
   The frontend dev server runs at `http://localhost:5173` by default.

2. **Backend**:
   ```bash
   # Go back to the project root
   go run main.go -config config/config.yaml
   ```

#### Full Build

The project provides automated build scripts that handle frontend building, asset embedding, and Go compilation:

- **Windows (CMD)**: `scripts\build.bat`
- **Windows (PowerShell)**: `.\scripts\build.ps1`
- **Linux/macOS**: `chmod +x scripts/build.sh && ./scripts/build.sh`

The build artifacts will be stored in the `output` directory.

### 📦 Quick Start

#### 1. Get the Program

Download the binary for your platform from the Release page, or build it yourself using the "Full Build" steps above.

#### 2. Configuration

Modify `config/config.yaml` before running. See [docs/config_specification.md](docs/config_specification.md) for detailed field descriptions.

Key settings:
- Storage type (`local`, `s3`, `webdav`).
- Admin password hash.
- Database configuration (`type`, `host`, `port`, `user`, `password`, `dbname`, etc.). The system supports auto-migration and will create tables automatically on the first connection. Manual SQL scripts are also available in the `scripts/sql/` directory.

#### 3. Running

```bash
./filerelay -config config/config.yaml
```

### 🐳 Docker Deployment

We provide Docker and Docker Compose support for one-click deployment:

#### Using Docker Image (Recommended)

If you don't want to build it yourself, you can pull the pre-built image from Docker Hub:

```bash
docker run -d \
  --name filerelay \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/config:/app/config \
  hxuanyu521/filerelay:latest
```

**Tips:**
- If you need to use custom Web pages (external frontend assets), you can mount the assets directory into the container and specify the path via the `FR_WEB_PATH` environment variable. For example:
  `-v $(pwd)/my-web:/app/my-web -e FR_WEB_PATH=/app/my-web`
- By default, the program will first look for external assets specified by `FR_WEB_PATH` (or `web.path` in the config file), and fallback to embedded assets if not found.

#### Using Docker Compose

1. Ensure Docker and Docker Compose are installed.
2. Run in the project root:

```bash
docker-compose up -d
```

By default, data, logs, database, and configuration will be persisted in the `./data` folder in the current directory.

#### Environment Variables

You can override key configuration items using environment variables, which is convenient for dynamic adjustments in Docker:

| Env Var | Description | Default |
| :--- | :--- | :--- |
| **Site Settings** | | |
| `FR_SITE_NAME` | Site Name | 文件暂存柜 |
| `FR_SITE_BASE_URL` | Site Base URL (for sharing links) | (empty) |
| `FR_SITE_PORT` | Service Port | 8080 |
| **Security Settings** | | |
| `FR_SECURITY_JWT_SECRET` | JWT Secret | file-relay-secret |
| **Upload Settings** | | |
| `FR_UPLOAD_MAX_SIZE` | Max File Size (MB) | 100 |
| `FR_UPLOAD_RETENTION_DAYS` | Max Retention Days | 30 |
| **Database Settings** | | |
| `FR_DB_TYPE` | DB Type (sqlite, mysql, postgres) | sqlite |
| `FR_DB_PATH` | SQLite DB Path | data/file_relay.db |
| `FR_DB_HOST` | DB Host (MySQL/Postgres) | (empty) |
| `FR_DB_PORT` | DB Port (MySQL/Postgres) | (empty) |
| `FR_DB_USER` | DB User (MySQL/Postgres) | (empty) |
| `FR_DB_PASSWORD` | DB Password (MySQL/Postgres) | (empty) |
| `FR_DB_NAME` | DB Name (MySQL/Postgres) | (empty) |
| **Storage Settings** | | |
| `FR_STORAGE_TYPE` | Storage Type (local, s3, webdav) | local |
| `FR_STORAGE_LOCAL_PATH` | Local Storage Path | data/storage_data |
| **Log Settings** | | |
| `FR_LOG_LEVEL` | Log Level (debug, info, warn, error) | info |
| `FR_LOG_FILE_PATH` | Log File Path | data/logs/app.log |
| **Web Settings** | | |
| `FR_WEB_PATH` | External Web Assets Path | web |

Access via:
- **Web UI**: `http://localhost:8080`
- **Admin Panel**: `http://localhost:8080/admin`
- **API Docs**: `http://localhost:8080/swagger/index.html`

### 📖 API Reference

FileRelay provides a rich set of APIs for integration via API Tokens.

- **Upload**: `POST /api/batches`
- **Pickup**: `GET /api/batches/:pickup_code`
- **Download**: `GET /api/files/:file_id/download`

Refer to the built-in Swagger documentation for details.
