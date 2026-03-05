@echo off
chcp 65001 > nul
setlocal enabledelayedexpansion

:: �л�����Ŀ��Ŀ¼
cd /d "%~dp0.."

set APP_NAME=filerelay
set OUTPUT_DIR=output

echo ��ʼ���� %APP_NAME% ��ƽ̨�������ļ�...

:: ���� output Ŀ¼
if exist "%OUTPUT_DIR%" (
    echo �������� %OUTPUT_DIR% Ŀ¼...
    rd /s /q "%OUTPUT_DIR%"
)

mkdir "%OUTPUT_DIR%"

:: ǰ�˹���
echo ���ڹ���ǰ����Ŀ...
pushd webapp
call npm install
if %ERRORLEVEL% neq 0 (
    echo npm install ʧ�ܣ�ֹͣ���롣
    popd
    exit /b %ERRORLEVEL%
)
call npm run build
if %ERRORLEVEL% neq 0 (
    echo ǰ�˹���ʧ�ܣ�ֹͣ���롣
    popd
    exit /b %ERRORLEVEL%
)
popd

:: ����Ŀ��ƽ̨ (OS/Arch)
set PLATFORMS=linux/amd64 linux/arm64 windows/amd64 windows/arm64 darwin/amd64 darwin/arm64

for %%P in (%PLATFORMS%) do (
    for /f "tokens=1,2 delims=/" %%A in ("%%P") do (
        set CGO_ENABLED=0
        set GOOS=%%A
        set GOARCH=%%B
        
        set OUTPUT_NAME=%APP_NAME%-%%A-%%B
        if "%%A"=="windows" set OUTPUT_NAME=!OUTPUT_NAME!.exe
        
        echo ���ڱ��� %%A/%%B...
        
        go build -ldflags="-s -w -extldflags=-static" -o "%OUTPUT_DIR%\!OUTPUT_NAME!" main.go
        
        if !ERRORLEVEL! equ 0 (
            echo   %%A/%%B ����ɹ�
            :: ѹ��Ϊ tar.gz (Windows 10+ �Դ� tar)
            tar -czf "%OUTPUT_DIR%\!OUTPUT_NAME!.tar.gz" -C "%OUTPUT_DIR%" "!OUTPUT_NAME!"
            :: ɾ��ԭʼ�������ļ�
            del "%OUTPUT_DIR%\!OUTPUT_NAME!"
        ) else (
            echo   %%A/%%B ����ʧ��
        )
    )
)

echo ----------------------------------------
echo ��ƽ̨�����ɣ����Ŀ¼: %OUTPUT_DIR%
pause
