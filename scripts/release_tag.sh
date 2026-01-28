#!/bin/bash

# 获取脚本所在目录并切换到项目根目录
cd "$(dirname "$0")/.."

# 检查是否提供了 tag 名称
if [ -z "$1" ]; then
    echo "使用方法: scripts/release_tag.sh <tag_name>"
    echo "例如: scripts/release_tag.sh v1.0.0"
    exit 1
fi

TAG_NAME=$1

# 1. 检查当前分支是否为 master
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [ "$CURRENT_BRANCH" != "master" ]; then
    echo "错误：当前分支是 $CURRENT_BRANCH，只能在 master 分支上创建 Tag。"
    exit 1
fi

# 2. 检查本地是否有未提交的内容
if [ -n "$(git status --porcelain)" ]; then
    echo "错误：本地有未提交的更改，请先提交或 stash。"
    exit 1
fi

# 3. 检查本地与远程 master 分支是否一致
echo "正在从远程获取最新 master 分支信息..."
git fetch origin master
LOCAL_HASH=$(git rev-parse HEAD)
REMOTE_HASH=$(git rev-parse origin/master)

if [ "$LOCAL_HASH" != "$REMOTE_HASH" ]; then
    echo "错误：本地 master 分支与远程 origin/master 不一致，请先 pull 或 push。"
    exit 1
fi

echo "安全检查通过，正在处理 Tag: $TAG_NAME"

# 在本地创建或更新 Tag
# -f 表示如果存在则强制更新
git tag -f "$TAG_NAME"

# 强制推送到远程
# --force 会覆盖远程同名 Tag
echo "正在强制推送到远程..."
git push origin "$TAG_NAME" --force

if [ $? -eq 0 ]; then
    echo "成功！Tag $TAG_NAME 已发布并同步至远程。"
else
    echo "失败：推送过程中出现错误。"
    exit 1
fi
