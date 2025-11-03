# Daily-Spark Backend

本后端采用 Go 项目标准布局（简化版），适配前后端分离：

```
backend/
  ├── cmd/               # 入口程序（可包含多个可执行）
  │   └── server/        # HTTP 服务入口
  ├── internal/          # 私有包（仅限本项目使用）
  ├── pkg/               # 公共包（可被外部引用，如有需要）
  ├── api/               # OpenAPI/Swagger/协议定义
  ├── configs/           # 配置模板或默认配置
  ├── docs/              # 设计/用户文档
  ├── scripts/           # 构建、运行、开发辅助脚本
  ├── test/              # 集成/端到端测试及测试数据
  └── go.mod             # Go 模块定义
```

## 开发运行
- 使用 `docker compose up -d` 启动依赖与后端（容器内通过 Air 热重载）。
- 或本机运行：在 `backend` 目录执行 `go run ./cmd/server`。

## 模块路径
- 当前 `go.mod` 模块名为占位 `github.com/cin/daily-spark`，如需发布到 GitHub/GitLab，请将其修改为你的仓库路径。