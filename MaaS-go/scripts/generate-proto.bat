@echo off
REM generate-proto.bat - Script to generate Go code from Protocol Buffers on Windows

echo Generating Protocol Buffers Go code...

REM Check if protoc is installed
where protoc >nul 2>nul
if %errorlevel% neq 0 (
    echo Error: protoc is not installed
    echo Please download from: https://github.com/protocolbuffers/protobuf/releases
    echo Download: protoc-<version>-win64.zip
    echo Extract and add 'bin' folder to your PATH environment variable
    exit /b 1
)

REM Check if Go plugins are installed
where protoc-gen-go >nul 2>nul
if %errorlevel% neq 0 (
    echo Installing protoc-gen-go...
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
)

where protoc-gen-go-grpc >nul 2>nul
if %errorlevel% neq 0 (
    echo Installing protoc-gen-go-grpc...
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
)

REM Change to project directory
cd /d "%~dp0\.."

REM Create directory if not exists
if not exist "shared\proto" mkdir "shared\proto"

REM Generate Go code
echo Generating Go code from model.proto...
protoc ^
    --go_out=. ^
    --go_opt=paths=source_relative ^
    --go-grpc_out=. ^
    --go-grpc_opt=paths=source_relative ^
    shared\proto\model.proto

if %errorlevel% equ 0 (
    echo.
    echo ✅ Successfully generated protobuf code!
    echo Generated files:
    echo   - shared/proto/model.pb.go
    echo   - shared/proto/model_grpc.pb.go
) else (
    echo.
    echo ❌ Failed to generate protobuf code
    exit /b 1
)

pause
