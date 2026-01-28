# 构建脚本
# 切换到项目根目录
Set-Location -Path (Join-Path $PSScriptRoot "..")

$APP_NAME = "filerelay"
$OUTPUT_DIR = "output"

# 定义目标平台
$PLATFORMS = @(
    "linux/amd64",
    "linux/arm64",
    "windows/amd64",
    "windows/arm64",
    "darwin/amd64",
    "darwin/arm64"
)

Write-Host "开始构建 $APP_NAME 多平台二进制文件..." -ForegroundColor Cyan

# 清理 output 目录
if (Test-Path $OUTPUT_DIR) {
    Write-Host "正在清理 $OUTPUT_DIR 目录..."
    Remove-Item -Path $OUTPUT_DIR -Recurse -Force
}

New-Item -Path $OUTPUT_DIR -ItemType Directory -Force | Out-Null

# 前端构建
Write-Host "正在构建前端项目..." -ForegroundColor Cyan
Push-Location webapp
npm install
if ($LASTEXITCODE -ne 0) {
    Write-Host "npm install 失败，停止编译。" -ForegroundColor Red
    Pop-Location
    exit $LASTEXITCODE
}
npm run build
if ($LASTEXITCODE -ne 0) {
    Write-Host "前端构建失败，停止编译。" -ForegroundColor Red
    Pop-Location
    exit $LASTEXITCODE
}
Pop-Location

# 循环构建各平台
foreach ($PLATFORM in $PLATFORMS) {
    $parts = $PLATFORM -split "/"
    $os = $parts[0]
    $arch = $parts[1]
    
    $outputName = "$($APP_NAME)-$($os)-$($arch)"
    if ($os -eq "windows") {
        $outputName += ".exe"
    }
    
    Write-Host "正在编译 $($os)/$($arch)..."
    $env:GOOS = $os
    $env:GOARCH = $arch
    
    go build -o (Join-Path $OUTPUT_DIR $outputName) main.go
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  $($os)/$($arch) 编译成功" -ForegroundColor Green
        # 压缩为 tar.gz
        tar -czf (Join-Path $OUTPUT_DIR "$outputName.tar.gz") -C $OUTPUT_DIR $outputName
        # 删除原始二进制文件
        Remove-Item (Join-Path $OUTPUT_DIR $outputName)
    } else {
        Write-Host "  $($os)/$($arch) 编译失败" -ForegroundColor Red
    }
}

# 重置环境变量
$env:GOOS = $null
$env:GOARCH = $null

Write-Host "----------------------------------------" -ForegroundColor Cyan
Write-Host "多平台打包完成！输出目录: $OUTPUT_DIR" -ForegroundColor Green
