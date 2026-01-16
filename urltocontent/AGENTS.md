# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## 构建和运行命令

### 后端 (Go)
```bash
# 必须从 backend/ 目录执行
cd backend

# 构建
go build -o bin/server ./cmd/server

# 运行
./bin/server

# 测试
go test ./...

# 测试单个文件（必须指定包路径）
go test ./internal/services -run TestFeishuAPICall
```

### 前端 (React + Vite)
```bash
# 开发模式
npm run dev

# 构建
npm run build

# Lint
npm run lint
```

## 重要的非显而易见信息

### 新功能：方向选择机制
- URL 解析后，用户可选择内置方向（技术、成长、生活）或自定义目录名
- 方向/目录在飞书知识库中实际上是文档类型（非真实文件夹）
- 系统会自动检查文档是否存在，不存在则创建
- 未选择时默认写入"自动总结"文档
- 写入完成后会向用户反馈最终写入的目录/文档名称

### 后端架构
- Go 模块名为 `urltocontent/backend`（不是 `backend`）
- 必须从 `backend/` 目录构建和运行，根目录的命令会失败
- 使用标准库 `net/http`，无外部依赖

### 前端配置
- 后端 API 地址硬编码为 `http://localhost:8080`（src/App.jsx:5）
- 使用 React 19 和 Vite，ESM 模块系统

### 飞书集成
- 飞书 API 配置硬编码在 backend/internal/config/config.go
- 可通过环境变量覆盖：`FEISHU_APP_ID`, `FEISHU_APP_SECRET`, `FEISHU_WIKI_ID`, `PORT`
- **限制**：飞书 Wiki API 不支持直接创建文件夹，只能创建 docx 文档类型
- 文档会创建在知识库根目录下，文件名"自动总结"实际上是文档而非文件夹

### 代码风格
- ESLint 规则：未使用变量如果以大写字母或下划线开头则忽略 (`varsIgnorePattern: '^[A-Z_]'`)
- Go 代码：使用 `fmt.Println` 进行日志输出，未使用结构化日志库

### 服务特性
- ParserService 超时设置为 60 秒（处理慢速网页）
- 后端 CORS 允许所有来源（`Access-Control-Allow-Origin: *`）
- 前端每 30 秒自动检查后端健康状态

### 测试说明
- 测试文件包含实际的飞书凭证（仅用于测试）
- 运行测试会实际写入飞书知识库，需谨慎使用
- 测试必须在 backend 目录下运行