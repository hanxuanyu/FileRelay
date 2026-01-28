param (
    [Parameter(Mandatory=$true, Position=0)]
    [string]$TagName
)

# 切换到项目根目录
Set-Location -Path (Join-Path $PSScriptRoot "..")

# 1. 检查当前分支是否为 master
$currentBranch = git rev-parse --abbrev-ref HEAD
if ($currentBranch -ne "master") {
    Write-Host "错误：当前分支是 $currentBranch，只能在 master 分支上创建 Tag。" -ForegroundColor Red
    exit 1
}

# 2. 检查本地是否有未提交的内容
$status = git status --porcelain
if ($null -ne $status -and $status.Length -gt 0) {
    Write-Host "错误：本地有未提交的更改，请先提交或 stash。" -ForegroundColor Red
    exit 1
}

# 3. 检查本地与远程 master 分支是否一致
Write-Host "正在从远程获取最新 master 分支信息..." -ForegroundColor Cyan
git fetch origin master
$localHash = git rev-parse HEAD
$remoteHash = git rev-parse origin/master

if ($localHash -ne $remoteHash) {
    Write-Host "错误：本地 master 分支与远程 origin/master 不一致，请先 pull 或 push。" -ForegroundColor Red
    exit 1
}

Write-Host "安全检查通过，正在处理 Tag: $TagName" -ForegroundColor Cyan

# 在本地创建或更新 Tag
git tag -f $TagName

if ($LASTEXITCODE -ne 0) {
    Write-Host "本地创建 Tag 失败。" -ForegroundColor Red
    exit 1
}

# 强制推送到远程
Write-Host "正在强制推送到远程..." -ForegroundColor Cyan
git push origin $TagName --force

if ($LASTEXITCODE -eq 0) {
    Write-Host "成功！Tag $TagName 已发布并同步至远程。" -ForegroundColor Green
} else {
    Write-Host "失败：推送过程中出现错误。" -ForegroundColor Red
    exit 1
}
