@echo off
@rem #"Makefile" for Windows projects.
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
if /I [%1]==[stop] call :stop

@rem # make build
if /I [%1]==[start] call :start

@rem # make build
if /I [%1]==[restart] call :restart

goto :eof

@rem #############
@rem # FUNCTIONS #
@rem #############
:stop
set /p systemPath=<./app.pid
taskkill /PID "%systemPath%"
exit

:start
start ./app.exe -config config.cnf
exit

:restart
set /p systemPath=<./app.pid
taskkill /PID "%systemPath%"
start ./app.exe -config config.cnf
exit


:usage
@echo Usage: %0 ^[ stop ^| start ^| restart ^]
exit /B 1