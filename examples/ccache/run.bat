@echo off
REM Build the server
go build -o server.exe

REM Start server instances in new windows
start "server8001" server.exe -port=8001
start "server8002" server.exe -port=8002
start "server8003" server.exe -port=8003 -api=1