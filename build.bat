@echo off
setlocal enabledelayedexpansion

echo [DEBUG] Script started
echo [DEBUG] Current directory: %CD%

REM Script to build both Go API and Kotlin library

REM Check if Go is installed
echo [DEBUG] Checking if Go is installed
where "go" >nul 2>nul
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH
    echo.
    echo To install Go:
    echo 1. Download Go from https://golang.org/dl/
    echo 2. Run the installer
    echo 3. Add Go to your PATH if not done automatically
    echo 4. Open a new terminal and verify with 'go version'
    echo.
    exit /b 1
)

REM Project directory (automatic detection)
for /f "delims=" %%i in ("%cd%") do (
    set "PROJECT_DIR=%%~fi"
)
echo [DEBUG] Project directory: %PROJECT_DIR%

REM Output directory for Go executables
set "GO_OUTPUT_DIR=%PROJECT_DIR%\build\api"
echo [DEBUG] Output directory: %GO_OUTPUT_DIR%

REM Create output directory if it doesn't exist
if not exist "%GO_OUTPUT_DIR%" (
    echo [DEBUG] Creating output directory
    mkdir "%GO_OUTPUT_DIR%"
)

REM Build Go API
echo Building Go API...
echo [DEBUG] Changing to LogManagingApi directory
cd LogManagingApi
if errorlevel 1 (
    echo Error: Failed to change to LogManagingApi directory
    exit /b 1
)

set "GOOS=windows"
set "GOARCH=amd64"
echo [DEBUG] Building with GOOS=%GOOS% GOARCH=%GOARCH%
go build -o "%GO_OUTPUT_DIR%\LogManagingApi-windows-64.exe" "%PROJECT_DIR%\LogManagingApi\logmanagingapi.go"

if errorlevel 1 (
    echo Build failed for Windows AMD64
    cd ..
    exit /b 1
) else (
    echo Build successful: %GO_OUTPUT_DIR%\LogManagingApi-windows-64.exe
)

cd ..

REM Build Kotlin library
echo Building Kotlin library...
echo [DEBUG] Changing to LogManagingKotlinLib directory
cd LogManagingKotlinLib
if errorlevel 1 (
    echo Error: Failed to change to LogManagingKotlinLib directory
    exit /b 1
)

echo [DEBUG] Running Gradle build using wrapper
REM Build the library using Gradle wrapper
call ..\gradlew.bat clean build

if errorlevel 1 (
    echo Kotlin library build failed
    cd ..
    exit /b 1
) else (
    echo Kotlin library build successful
)

cd ..

echo Build process completed successfully
echo Go API executables are in %GO_OUTPUT_DIR%
echo Kotlin library JAR is in LogManagingKotlinLib/build/libs/

exit /b 0