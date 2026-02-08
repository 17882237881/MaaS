@echo off
REM fix-proto.bat - Fix protobuf import issues

echo ========================================
echo  Fixing Protobuf Import Issues
echo ========================================
echo.

cd /d "%~dp0\.."

echo Step 1: Cleaning Go cache...
go clean -cache 2>nul
go clean -modcache 2>nul
echo Done.
echo.

echo Step 2: Checking generated files...
if not exist "shared\proto\model.pb.go" (
    echo ERROR: model.pb.go not found!
    echo Please run: scripts\generate-proto.bat
    pause
    exit /b 1
)

if not exist "shared\proto\model_grpc.pb.go" (
    echo ERROR: model_grpc.pb.go not found!
    echo Please run: scripts\generate-proto.bat
    pause
    exit /b 1
)

echo Found generated files.
echo.

echo Step 3: Checking package name...
findstr /B "package" shared\proto\model.pb.go | findstr "modelpb" >nul
if errorlevel 1 (
    echo ERROR: Package name in model.pb.go is not 'modelpb'!
    echo Please regenerate the protobuf files.
    pause
    exit /b 1
)

echo Package name is correct: modelpb
echo.

echo Step 4: Running go mod tidy...
go mod tidy
if errorlevel 1 (
    echo ERROR: go mod tidy failed!
    pause
    exit /b 1
)
echo Done.
echo.

echo Step 5: Verifying build...
go build ./shared/proto/... 2>build_error.log
if errorlevel 1 (
    echo ERROR: Build failed!
    echo Error details:
    type build_error.log
    del build_error.log
    pause
    exit /b 1
)
del build_error.log 2>nul

echo.
echo ========================================
echo  ✅ Fix completed successfully!
echo ========================================
echo.
echo You can now build the project:
echo   go build ./...
echo.

pause
