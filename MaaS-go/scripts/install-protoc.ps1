#Requires -RunAsAdministrator

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Protocol Buffers 自动安装脚本" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 配置
$protocVersion = "24.4"
$installDir = "C:\protoc"
$downloadUrl = "https://github.com/protocolbuffers/protobuf/releases/download/v$protocVersion/protoc-$protocVersion-win64.zip"
$tempZip = "$env:TEMP\protoc.zip"

# 检查是否以管理员身份运行
$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
if (-not $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "请右键点击脚本，选择'使用 PowerShell 运行'并以管理员身份运行！" -ForegroundColor Red
    pause
    exit 1
}

# 检查 Go 是否安装
try {
    $goVersion = go version 2>$null
    Write-Host "✅ 检测到 Go" -ForegroundColor Green
} catch {
    Write-Host "❌ 未检测到 Go，请先安装 Go！" -ForegroundColor Red
    Write-Host "下载地址: https://go.dev/dl/" -ForegroundColor Yellow
    pause
    exit 1
}

# 步骤1: 下载 protoc
Write-Host ""
Write-Host "步骤 1/4: 下载 protoc..." -ForegroundColor Cyan
Write-Host "版本: $protocVersion" -ForegroundColor Gray

if (Test-Path $tempZip) {
    Remove-Item $tempZip -Force
}

try {
    Write-Host "正在下载..." -ForegroundColor Yellow
    Invoke-WebRequest -Uri $downloadUrl -OutFile $tempZip -UseBasicParsing
    Write-Host "✅ 下载完成" -ForegroundColor Green
} catch {
    Write-Host "❌ 下载失败" -ForegroundColor Red
    Write-Host "错误: $_" -ForegroundColor Red
    pause
    exit 1
}

# 步骤2: 解压安装
Write-Host ""
Write-Host "步骤 2/4: 解压安装..." -ForegroundColor Cyan

try {
    if (Test-Path $installDir) {
        Write-Host "清理旧版本..." -ForegroundColor Yellow
        Remove-Item $installDir -Recurse -Force
    }
    
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    
    Write-Host "正在解压..." -ForegroundColor Yellow
    Expand-Archive -Path $tempZip -DestinationPath $installDir -Force
    Remove-Item $tempZip -Force
    
    Write-Host "✅ 解压完成" -ForegroundColor Green
} catch {
    Write-Host "❌ 解压失败" -ForegroundColor Red
    Write-Host "错误: $_" -ForegroundColor Red
    pause
    exit 1
}

# 步骤3: 配置环境变量
Write-Host ""
Write-Host "步骤 3/4: 配置环境变量..." -ForegroundColor Cyan

try {
    $binPath = "$installDir\bin"
    $goBinPath = "$env:USERPROFILE\go\bin"
    
    # 添加 protoc bin 目录
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$binPath*") {
        Write-Host "添加 protoc 到 PATH..." -ForegroundColor Yellow
        $newPath = "$userPath;$binPath"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    } else {
        Write-Host "protoc 已在 PATH 中" -ForegroundColor Gray
    }
    
    # 添加 Go bin 目录
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$goBinPath*") {
        Write-Host "添加 Go bin 到 PATH..." -ForegroundColor Yellow
        $newPath = "$userPath;$goBinPath"
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    } else {
        Write-Host "Go bin 已在 PATH 中" -ForegroundColor Gray
    }
    
    # 更新当前会话
    $env:Path = [Environment]::GetEnvironmentVariable("Path", "User")
    
    Write-Host "✅ 环境变量配置完成" -ForegroundColor Green
} catch {
    Write-Host "❌ 配置环境变量失败" -ForegroundColor Red
    Write-Host "错误: $_" -ForegroundColor Red
    pause
    exit 1
}

# 步骤4: 安装 Go 插件
Write-Host ""
Write-Host "步骤 4/4: 安装 Go 插件..." -ForegroundColor Cyan

try {
    Write-Host "安装 protoc-gen-go..." -ForegroundColor Yellow
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    
    Write-Host "安装 protoc-gen-go-grpc..." -ForegroundColor Yellow
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    
    Write-Host "✅ Go 插件安装完成" -ForegroundColor Green
} catch {
    Write-Host "❌ 安装 Go 插件失败" -ForegroundColor Red
    Write-Host "错误: $_" -ForegroundColor Red
    pause
    exit 1
}

# 验证安装
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "验证安装..." -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

$protocExe = "$installDir\bin\protoc.exe"
if (Test-Path $protocExe) {
    $ver = & $protocExe --version
    Write-Host "✅ protoc: $ver" -ForegroundColor Green
} else {
    Write-Host "❌ protoc 未找到" -ForegroundColor Red
}

Write-Host ""
Write-Host "🎉 安装完成！" -ForegroundColor Green
Write-Host ""
Write-Host "重要提示：" -ForegroundColor Yellow
Write-Host "请关闭当前 PowerShell 窗口，打开新的命令提示符，然后运行:" -ForegroundColor White
Write-Host "  cd D:\code\MaaS\MaaS-go" -ForegroundColor Cyan
Write-Host "  scripts\generate-proto.bat" -ForegroundColor Cyan
Write-Host ""

pause
