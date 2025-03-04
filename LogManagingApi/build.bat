@echo off

REM Script to build Go executables for multiple operating systems

REM Project directory (automatic detection)
set PROJECT_DIR=%CD%

REM Output directory for Go executables
set GO_OUTPUT_DIR=%PROJECT_DIR%\build\

REM Create output directory if it doesn't exist
if not exist "%GO_OUTPUT_DIR%" mkdir "%GO_OUTPUT_DIR%"

REM Function to build for a specific Go OS and architecture
:build_go_os_arch
set os=%1
set arch=%2
set output_name=%3

echo Building Go for %os% %arch%...

set GOOS=%os%
set GOARCH=%arch%
go build -o "%GO_OUTPUT_DIR%\%output_name%" "%PROJECT_DIR%\logmanagingapi.go"

if %errorlevel% equ 0 (
    echo Build successful: %GO_OUTPUT_DIR%\%output_name%
) else (
    echo Build failed for %os% %arch%
)
goto :eof

REM Build Go executables

REM Build for Windows 64-bit
call :build_go_os_arch windows amd64 "myprogram.exe"

REM Build for Windows 32-bit
call :build_go_os_arch windows 386 "myprogram32.exe"

REM Build for Linux 64-bit
call :build_go_os_arch linux amd64 "myprogram_linux_amd64"

REM Build for Linux 32-bit
call :build_go_os_arch linux 386 "myprogram_linux_386"

REM Build for Linux ARM 64-bit
call :build_go_os_arch linux arm64 "myprogram_linux_arm64"

REM Build for Linux ARM 32-bit
call :build_go_os_arch linux arm "myprogram_linux_arm"

REM Build for macOS 64-bit
call :build_go_os_arch darwin amd64 "myprogram_darwin_amd64"

REM Build for macOS ARM 64-bit
call :build_go_os_arch darwin arm64 "myprogram_darwin_arm64"

echo Go build process completed. Executables are in %GO_OUTPUT_DIR%

pause