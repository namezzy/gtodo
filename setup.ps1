# Gtodo 环境初始化脚本 (PowerShell)
# 使用方法: . .\setup.ps1  (注意前面的点和空格，表示在当前会话中执行)

$env:Path = "$PSScriptRoot;$env:Path"
$env:GTODO_STORAGE = "mysql"
$env:GTODO_MYSQL_DSN = "gtodo:gtodo_pass_123@tcp(127.0.0.1:13306)/gtodo?charset=utf8mb4&parseTime=True&loc=Local"

Write-Host "[OK] 环境变量已设置，现在可以直接使用 gtodo 命令：" -ForegroundColor Green
Write-Host ""
Write-Host "  gtodo add '任务描述' -p high"
Write-Host "  gtodo list"
Write-Host "  gtodo done 1"
Write-Host "  gtodo delete 1"
Write-Host ""
