@echo off
chcp 65001 > nul
setlocal enabledelayedexpansion

:: 切换到项目根目录
cd /d "%~dp0.."

set APP_NAME=filerelay
set OUTPUT_DIR=output

echo 开始构建 %APP_NAME% 多平台二进制文件...

:: 清理 output 目录
if exist "%OUTPUT_DIR%" (
    echo 正在清理 %OUTPUT_DIR% 目录...
    rd /s /q "%OUTPUT_DIR%"
)

mkdir "%OUTPUT_DIR%"

:: 前端构建
echo 正在构建前端项目...
pushd webapp
call npm install
if %ERRORLEVEL% neq 0 (
    echo npm install 失败，停止编译。
    popd
    exit /b %ERRORLEVEL%
)
call npm run build
if %ERRORLEVEL% neq 0 (
    echo 前端构建失败，停止编译。
    popd
    exit /b %ERRORLEVEL%
)
popd

:: 定义目标平台 (OS/Arch)
set PLATFORMS=linux/amd64 linux/arm64 windows/amd64 windows/arm64 darwin/amd64 darwin/arm64

for %%P in (%PLATFORMS%) do (
    for /f "tokens=1,2 delims=/" %%A in ("%%P") do (
        set GOOS=%%A
        set GOARCH=%%B
        
        set OUTPUT_NAME=%APP_NAME%-%%A-%%B
        if "%%A"=="windows" set OUTPUT_NAME=!OUTPUT_NAME!.exe
        
        echo 正在编译 %%A/%%B...
        
        go build -o "%OUTPUT_DIR%\!OUTPUT_NAME!" main.go
        
        if !ERRORLEVEL! equ 0 (
            echo   %%A/%%B 编译成功
            :: 压缩为 tar.gz (Windows 10+ 自带 tar)
            tar -czf "%OUTPUT_DIR%\!OUTPUT_NAME!.tar.gz" -C "%OUTPUT_DIR%" "!OUTPUT_NAME!"
            :: 删除原始二进制文件
            del "%OUTPUT_DIR%\!OUTPUT_NAME!"
        ) else (
            echo   %%A/%%B 编译失败
        )
    )
)

echo ----------------------------------------
echo 多平台打包完成！输出目录: %OUTPUT_DIR%
pause
