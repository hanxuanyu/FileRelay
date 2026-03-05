# �����ű�
# �л�����Ŀ��Ŀ¼
Set-Location -Path (Join-Path $PSScriptRoot "..")

$APP_NAME = "filerelay"
$OUTPUT_DIR = "output"

# ����Ŀ��ƽ̨
$PLATFORMS = @(
    "linux/amd64",
    "linux/arm64",
    "windows/amd64",
    "windows/arm64",
    "darwin/amd64",
    "darwin/arm64"
)

Write-Host "��ʼ���� $APP_NAME ��ƽ̨�������ļ�..." -ForegroundColor Cyan

# ���� output Ŀ¼
if (Test-Path $OUTPUT_DIR) {
    Write-Host "�������� $OUTPUT_DIR Ŀ¼..."
    Remove-Item -Path $OUTPUT_DIR -Recurse -Force
}

New-Item -Path $OUTPUT_DIR -ItemType Directory -Force | Out-Null

# ǰ�˹���
Write-Host "���ڹ���ǰ����Ŀ..." -ForegroundColor Cyan
Push-Location webapp
npm install
if ($LASTEXITCODE -ne 0) {
    Write-Host "npm install ʧ�ܣ�ֹͣ���롣" -ForegroundColor Red
    Pop-Location
    exit $LASTEXITCODE
}
npm run build
if ($LASTEXITCODE -ne 0) {
    Write-Host "ǰ�˹���ʧ�ܣ�ֹͣ���롣" -ForegroundColor Red
    Pop-Location
    exit $LASTEXITCODE
}
Pop-Location

# ѭ��������ƽ̨
foreach ($PLATFORM in $PLATFORMS) {
    $parts = $PLATFORM -split "/"
    $os = $parts[0]
    $arch = $parts[1]
    
    $outputName = "$($APP_NAME)-$($os)-$($arch)"
    if ($os -eq "windows") {
        $outputName += ".exe"
    }
    
    Write-Host "���ڱ��� $($os)/$($arch)..."
    $env:CGO_ENABLED = "0"
    $env:GOOS = $os
    $env:GOARCH = $arch
    
    go build -ldflags="-s -w -extldflags=-static" -o (Join-Path $OUTPUT_DIR $outputName) main.go
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  $($os)/$($arch) ����ɹ�" -ForegroundColor Green
        # ѹ��Ϊ tar.gz
        tar -czf (Join-Path $OUTPUT_DIR "$outputName.tar.gz") -C $OUTPUT_DIR $outputName
        # ɾ��ԭʼ�������ļ�
        Remove-Item (Join-Path $OUTPUT_DIR $outputName)
    } else {
        Write-Host "  $($os)/$($arch) ����ʧ��" -ForegroundColor Red
    }
}

# ���û�������
$env:GOOS = $null
$env:GOARCH = $null

Write-Host "----------------------------------------" -ForegroundColor Cyan
Write-Host "��ƽ̨�����ɣ����Ŀ¼: $OUTPUT_DIR" -ForegroundColor Green
