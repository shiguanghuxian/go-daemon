@echo off
@rem #"Makefile" for Windows 25jd_jxk projects.
@rem #Copyright (c) http://www.atisafe.com/, 2015. All rights reserved.
@rem
SETLOCAL

@rem ###############
@rem # PRINT USAGE #
@rem ###############

if [%1]==[] goto usage

@rem ################
@rem # SWITCH BLOCK #
@rem ################

@rem # make build
if /I [%1]==[build] call :build

@rem # make clean
if /I [%1]==[clean] call :clean

goto :eof

@rem #############
@rem # FUNCTIONS #
@rem #############
:build
for /f "delims=" %%i in ('go version') do (set go_version=%%i)
for /f "delims=" %%i in ('git rev-parse HEAD') do (set git_hash=%%i)
@go build -ldflags "-X main.VERSION=1.0.0 -X 'main.BUILD_TIME=%DATE% %TIME%' -X 'main.GO_VERSION=%go_version%' -X main.GIT_HASH=%git_hash%" -o ./build/godaemon.exe ./
exit /B %ERRORLEVEL%

:windows_build
for /f "delims=" %%i in ('go version') do (set go_version=%%i)
for /f "delims=" %%i in ('git rev-parse HEAD') do (set git_hash=%%i)
SET CGO_ENABLED=0
set GOARCH=amd64
set GOOS=windows
@go build -ldflags "-X main.VERSION=1.0.0 -X 'main.BUILD_TIME=%DATE% %TIME%' -X 'main.GO_VERSION=%go_version%' -X main.GIT_HASH=%git_hash%" -o ./build/godaemon.exe ./
exit /B %ERRORLEVEL%

:windows_build_386
for /f "delims=" %%i in ('go version') do (set go_version=%%i)
for /f "delims=" %%i in ('git rev-parse HEAD') do (set git_hash=%%i)
SET CGO_ENABLED=0
set GOARCH=386
set GOOS=windows
@go build -ldflags "-X main.VERSION=1.0.0 -X 'main.BUILD_TIME=%DATE% %TIME%' -X 'main.GO_VERSION=%go_version%' -X main.GIT_HASH=%git_hash%" -o ./build/godaemon.exe ./
exit /B %ERRORLEVEL%

:linux_build
for /f "delims=" %%i in ('go version') do (set go_version=%%i)
for /f "delims=" %%i in ('git rev-parse HEAD') do (set git_hash=%%i)
SET CGO_ENABLED=0
set GOARCH=amd64
set GOOS=linux
@go build -ldflags "-X main.VERSION=1.0.0 -X 'main.BUILD_TIME=%DATE% %TIME%' -X 'main.GO_VERSION=%go_version%' -X main.GIT_HASH=%git_hash%" -o ./build/godaemon ./
exit /B %ERRORLEVEL%

:clean
@del /S /F /Q "build\godaemon*"
@del /S /F /Q "build\logs\*.log"
exit /B 0

:usage
@echo Usage: %0 ^[ build ^| clean ^| linux_build ^| windows_build_386 ^| windows_build  ^]
exit /B 1