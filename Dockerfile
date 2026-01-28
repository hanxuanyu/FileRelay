# Frontend build stage
FROM node:20-alpine AS frontend-builder
ARG NPM_REGISTRY
# 如果设置了 NPM_REGISTRY，则配置 npm 镜像
RUN if [ -n "$NPM_REGISTRY" ]; then npm config set registry $NPM_REGISTRY; fi

WORKDIR /webapp
# 仅拷贝包定义文件以利用缓存
COPY webapp/package*.json ./
RUN npm ci
# 拷贝前端源码并构建
COPY webapp/ ./
RUN npm run build

# Backend build stage
FROM golang:1.24-alpine AS builder
ARG GOPROXY
ENV GOPROXY=$GOPROXY
WORKDIR /app
# 仅拷贝依赖文件以利用缓存
COPY go.mod go.sum ./
RUN go mod download
# 拷贝后端源码（不包含 webapp）
COPY internal/ ./internal/
COPY main.go .
COPY docs/ ./docs/
# 从前端构建阶段拷贝产物到 web 目录（Go embed 需要）
COPY --from=frontend-builder /web ./web
# 编译二进制，移除调试信息并进行静态链接优化
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o filerelay main.go

# Final stage
FROM alpine:3.21
WORKDIR /app
# 合并指令以减少层数，安装必要包并创建所需目录
RUN apk add --no-cache ca-certificates mailcap && \
    mkdir -p config data/storage_data data/logs
# 仅拷贝编译后的二进制文件，前端资源已通过 Go embed 包含在内
COPY --from=builder /app/filerelay .

EXPOSE 8080
# 合并环境变量定义
ENV FR_SITE_PORT=8080 \
    FR_STORAGE_LOCAL_PATH=data/storage_data \
    FR_DB_TYPE=sqlite \
    FR_DB_PATH=data/file_relay.db \
    FR_LOG_LEVEL=info \
    FR_LOG_FILE_PATH=data/logs/app.log

CMD ["./filerelay"]
