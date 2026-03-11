@echo off
REM Gtodo 环境初始化脚本
REM 使用方法: 双击运行，或在 cmd 中执行 setup.bat

REM 将 gtodo.exe 所在目录加入当前会话 PATH
set "PATH=%~dp0;%PATH%"

REM 设置 MySQL 存储
set GTODO_STORAGE=mysql
set GTODO_MYSQL_DSN=gtodo:gtodo_pass_123@tcp(127.0.0.1:13306)/gtodo?charset=utf8mb4^&parseTime=True^&loc=Local

echo [OK] 环境变量已设置，现在可以直接使用 gtodo 命令：
echo.
echo   gtodo add "任务描述" -p high
echo   gtodo list
echo   gtodo done 1
echo   gtodo delete 1
echo.

cmd /k
