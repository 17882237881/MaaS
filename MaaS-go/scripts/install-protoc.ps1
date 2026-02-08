#Requires -RunAsAdministrator

<#
.SYNOPSIS
    自动安装 Protocol Buffers (protoc) 和 Go 插件
.DESCRIPTION
    自动下载、安装 protoc，配置环境变量，并安装 Go 插件
#>

$ErrorActionPreference = "Stop"

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
if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "❌ 请以管理员身份运行此脚本！" -ForegroundColor Red
    Write-Host "   右键点击脚本 -> 以管理员身份运行" -ForegroundColor Yellow
    pause
    exit 1
}

# 检查 Go 是否安装
try {
    $goVersion = go version
    Write-Host "✅ 检测到 Go: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ 未检测到 Go，请先安装 Go！" -ForegroundColor Red
    Write-Host "   下载地址: https://go.dev/dl/" -ForegroundColor Yellow
    pause
    exit 1
}

# 步骤1: 下载 protoc
Write-Host ""
Write-Host "步骤 1/4: 下载 protoc..." -ForegroundColor Cyan
Write-Host "   版本: $protocVersion" -ForegroundColor Gray
Write-Host "   下载地址: $downloadUrl" -ForegroundColor Gray

try {
    if (Test-Path $tempZip) {
        Remove-Item $tempZip -Force
    }
    
    Write-Host "   正在下载..." -ForegroundColor Yellow
    Invoke-WebRequest -Uri $downloadUrl -OutFile $tempZip -UseBasicParsing
    Write-Host "   ✅ 下载完成" -ForegroundColor Green
} catch {
    Write-Host "   ❌ 下载失败: $_" -ForegroundColor Red
    pause
    exit 1
}

# 步骤2: 解压安装
Write-Host ""
Write-Host "步骤 2/4: 解压安装..." -ForegroundColor Cyan

try {
    # 创建安装目录
    if (Test-Path $installDir) {
        Write-Host "   清理旧版本..." -ForegroundColor Yellow
        Remove-Item $installDir -Recurse -Force
    }
    
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    
    # 解压
    Write-Host "   正在解压到 $installDir..." -ForegroundColor Yellow
    Expand-Archive -Path $tempZip -DestinationPath $installDir -Force
    
    # 清理临时文件
    Remove-Item $tempZip -Force
    
    Write-Host "   ✅ 解压完成" -ForegroundColor Green
} catch {
    Write-Host "   ❌ 解压失败: $_" -ForegroundColor Red
    pause
    exit 1
}

# 步骤3: 配置环境变量
Write-Host ""
Write-Host "步骤 3/4: 配置环境变量..." -ForegroundColor Cyan

try {
    $binPath = "$installDir\bin"
    $goBinPath = "$env:USERPROFILE\go\bin"
    
    # 获取当前用户 PATH
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    
    # 添加 protoc bin 目录
    if ($currentPath -notlike "*$binPath*") {
        Write-Host "   添加 protoc 到 PATH..." -ForegroundColor Yellow
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$binPath", "User")
    } else {
        Write-Host "   protoc 已在 PATH 中" -ForegroundColor Gray
    }
    
    # 添加 Go bin 目录
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($currentPath -notlike "*$goBinPath*") {
        Write-Host "   添加 Go bin 到 PATH..." -ForegroundColor Yellow
        [Environment]::SetEnvironmentVariable("Path", "$currentPath;$goBinPath", "User")
    } else {
        Write-Host "   Go bin 已在 PATH 中" -ForegroundColor Gray
    }
    
    # 更新当前会话的 PATH
    $env:Path = [Environment]::GetEnvironmentVariable("Path", "User")
    
    Write-Host "   ✅ 环境变量配置完成" -ForegroundColor Green
} catch {
    Write-Host "   ❌ 配置环境变量失败: $_" -ForegroundColor Red
    pause
    exit 1
}

# 步骤4: 安装 Go 插件
Write-Host ""
Write-Host "步骤 4/4: 安装 Go 插件..." -ForegroundColor Cyan

try {
    Write-Host "   安装 protoc-gen-go..." -ForegroundColor Yellow
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    
    Write-Host "   安装 protoc-gen-go-grpc..." -ForegroundColor Yellow
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    
    Write-Host "   ✅ Go 插件安装完成" -ForegroundColor Green
} catch {
    Write-Host "   ❌ 安装 Go 插件失败: $_" -ForegroundColor Red
    pause
    exit 1
}

# 验证安装
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "验证安装..." -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

$success = $true

try {
    $protocVer = & "$installDir\bin\protoc.exe" --version
    Write-Host "✅ protoc: $protocVer" -ForegroundColor Green
} catch {
    Write-Host "❌ protoc 验证失败" -ForegroundColor Red
    $success = $false
}

try {
    $genGoVer = protoc-gen-go --version 2>&1
    Write-Host "✅ protoc-gen-go: 已安装" -ForegroundColor Green
} catch {
    Write-Host "❌ protoc-gen-go 验证失败" -ForegroundColor Red
    $success = $false
}

try {
    $genGrpcVer = protoc-gen-go-grpc --version 2>&1
    Write-Host "✅ protoc-gen-go-grpc: 已安装" -ForegroundColor Green
} catch {
    Write-Host "❌ protoc-gen-go-grpc 验证失败" -ForegroundColor Red
    $success = $false
}

Write-Host ""
if ($success) {
    Write-Host "🎉 安装成功！" -ForegroundColor Green
    Write-Host ""
    Write-Host "现在可以生成 protobuf 代码了:" -ForegroundColor Cyan
    Write-Host "  1. 关闭当前命令窗口（重要！）" -ForegroundColor Yellow
    Write-Host "  2. 打开新的命令提示符" -ForegroundColor White
    Write-Host "  3. 运行: cd D:\code\MaaS\MaaS-go" -ForegroundColor White
    Write-Host "  4. 运行: scripts\generate-proto.bat" -ForegroundColor White
    Write-Host ""
} else {
    Write-Host "⚠️  部分组件安装失败，请检查错误信息" -ForegroundColor Yellow
}

pause
