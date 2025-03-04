@echo off

REM Script to build Go executables for multiple operating systems

REM Function to build for a specific Go OS and architecture
:build_go_os_arch
set os=%1
set arch=%2
set output_name=%3
echo OS: %os%
echo ARCH: %arch%
echo OUTPUT_NAME: %output_name%

cd ../
echo Current Directory after cd ..: %cd%

REM Project directory (automatic detection)
for /f "delims=" %%i in ("%cd%") do (
    set "PROJECT_DIR=%%~fi"
)
echo PROJECT_DIR: %PROJECT_DIR%

REM Output directory for Go executables
set GO_OUTPUT_DIR=%PROJECT_DIR%\build\api

REM Create output directory if it doesn't exist
echo Creating directory: "%GO_OUTPUT_DIR%"
if not exist "%GO_OUTPUT_DIR%" mkdir "%GO_OUTPUT_DIR%"

cd LogManagingApi

echo Building Go for %os% %arch%...

set GOOS=%os%
set GOARCH=%arch%
go build -o "%GO_OUTPUT_DIR%\%output_name%" "%PROJECT_DIR%\LogManagingApi\logmanagingapi.go"
echo Errorlevel: %errorlevel%
exit /b %errorlevel%
pause

if errorlevel 1 (
    echo Build failed for %os% %arch%
    pause
) else (
    echo Build successful: %GO_OUTPUT_DIR%\%output_name%
    pause
)

cd "../"
goto :eof

REM Main execution
call :build_go_os_arch %1 %2 %3
echo Go build process completed. Executables are in %GO_OUTPUT_DIR%
exit /b %errorlevel%