#!/usr/bin/env sh
set -e

# 本地开发脚本：优先使用 air，否则用 go run
if command -v air >/dev/null 2>&1; then
  air -c .air.toml || air
else
  echo "air 未安装，使用 go run"
  go run ./cmd/server
fi