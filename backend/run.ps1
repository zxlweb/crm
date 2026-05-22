# 本地启动 API（不依赖终端是否已刷新系统 PATH）
$env:Path = "C:\Program Files\Go\bin;" + $env:Path
Set-Location $PSScriptRoot
go run ./cmd/api/
