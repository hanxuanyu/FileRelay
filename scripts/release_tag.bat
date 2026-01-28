@echo off
chcp 65001 > nul
setlocal

:: 切换到项目根目录
cd /d "%~dp0.."

if "%~1"=="" (
    echo 使用方法: release_tag.bat <tag_name>
    echo 例如: release_tag.bat v1.0.0
    exit /b 1
)

set TAG_NAME=%~1

REM 1. 检查当前分支是否为 master
for /f "tokens=*" %%i in ('git rev-parse --abbrev-ref HEAD') do set CURRENT_BRANCH=%%i
if not "%CURRENT_BRANCH%"=="master" (
    echo 错误：当前分支是 %CURRENT_BRANCH%，只能在 master 分支上创建 Tag。
    exit /b 1
)

REM 2. 检查本地是否有未提交的内容
for /f "tokens=*" %%i in ('git status --porcelain') do (
    echo 错误：本地有未提交的更改，请先提交或 stash。
    exit /b 1
)

REM 3. 检查本地与远程 master 分支是否一致
echo 正在从远程获取最新 master 分支信息...
git fetch origin master
if %ERRORLEVEL% neq 0 (
    echo 错误：获取远程分支信息失败。
    exit /b 1
)

for /f "tokens=*" %%i in ('git rev-parse HEAD') do set LOCAL_HASH=%%i
for /f "tokens=*" %%i in ('git rev-parse origin/master') do set REMOTE_HASH=%%i

if not "%LOCAL_HASH%"=="%REMOTE_HASH%" (
    echo 错误：本地 master 分支与远程 origin/master 不一致，请先 pull 或 push。
    exit /b 1
)

echo 安全检查通过，正在处理 Tag: %TAG_NAME%

:: 在本地创建或更新 Tag
git tag -f %TAG_NAME%

:: 强制推送到远程
echo 正在强制推送到远程...
git push origin %TAG_NAME% --force

if %ERRORLEVEL% equ 0 (
    echo 成功！Tag %TAG_NAME% 已发布并同步至远程。
) else (
    echo 失败：推送过程中出现错误。
    exit /b 1
)
