# Project Documentation Rules (Non-Obvious Only)

## 新功能：方向选择机制

### 功能概述
- URL 解析后，用户可选择内容分类方向
- 提供内置方向：技术、成长、生活
- 支持自定义方向：用户输入任意目录名
- 默认行为：未选择时自动使用"自动总结"

### 工作流程
1. **解析阶段**：URL 解析成功后显示方向选择界面
2. **选择阶段**：
   - 用户选择内置方向按钮（技术/成长/生活）
   - 或用户输入自定义方向名称
   - 或用户不选择（使用默认）
3. **检查阶段**：后端检查知识库根目录下是否存在对应文档
4. **创建阶段**：不存在则创建文档（obj_type: docx）
5. **写入阶段**：将解析内容写入文档
6. **反馈阶段**：显示最终写入的文档名称

### 重要概念
- **方向 = 文档**：在飞书知识库中，方向实际上是文档类型
- **文档检查**：每次写入前检查文档是否存在
- **自动创建**：首次使用某个方向时自动创建对应文档
- **用户反馈**：必须向用户显示最终写入的文档名称

### API 变更
- `/api/write` 请求添加可选字段 `direction`
- `/api/write` 响应添加字段 `documentName`（最终写入的文档名称）

## 项目结构说明

### 目录组织
- `backend/`：Go 后端服务，必须从此目录执行构建和测试
- `src/`：React 前端应用，运行在 `http://localhost:5173`
- 后端 API 硬编码地址为 `http://localhost:8080`（非配置文件）

### 后端架构
- Go 模块名为 `urltocontent/backend`（非 `backend`）
- 使用标准库 `net/http`，无外部 HTTP 框架依赖
- 目录结构：
  - `cmd/server/`：主程序入口
  - `internal/config/`：配置管理（包含硬编码的飞书凭证）
  - `internal/handlers/`：HTTP 处理器（含 CORS 中间件）
  - `internal/services/`：业务逻辑（Parser、Feishu）
  - `internal/models/`：数据模型

### 前端架构
- 单文件应用：主要逻辑在 `src/App.jsx`（410 行）
- 无路由，无状态管理库，纯 React Hooks
- 使用 CSS 模块化：`App.css`、`index.css`

## 关键架构约束

### 飞书集成限制
- **重要**：飞书 Wiki API 不支持直接创建文件夹节点
- 只能创建 `obj_type: "docx"` 的文档类型
- "自动总结"是文档标题，非真实文件夹
- 文档创建在知识库根目录下，无法创建子目录结构

### 通信模式
- 前端 → 后端：RESTful API（POST /api/parse, /api/write）
- 后端 → 飞书：HTTP 请求到 `https://open.feishu.cn`
- 无 WebSocket、无消息队列、无实时推送

### 配置管理
- 后端配置在 `backend/internal/config/config.go`
- 支持环境变量覆盖，但有硬编码默认值
- 前端无配置文件，API 地址硬编码

## 文档位置

### API 文档
- API 端点在 `README.md` 第 123-181 行
- 包含请求/响应示例
- 飞书 API 调用流程在 `backend/internal/services/feishu.go` 注释中

### 测试说明
- 测试在 `backend/internal/services/feishu_test.go`
- 包含实际飞书凭证（仅用于测试）
- 运行测试会实际写入飞书知识库

## 使用流程

### 正常工作流
1. 启动后端：`cd backend && go build -o bin/server ./cmd/server && ./bin/server`
2. 启动前端：`npm run dev`
3. 用户输入 URL
4. 前端调用 `/api/parse` 解析
5. 用户确认后，前端调用 `/api/write` 写入飞书

### 错误处理流程
- 后端离线：前端每 30 秒自动检测，显示警告
- 解析失败：显示错误消息，不中断流程
- 写入失败：显示错误消息，保留解析内容供重试

## 技术债务和限制

### 已知限制
- 飞书无法创建文件夹，所有文档在根目录
- 后端使用 `fmt.Println` 日志，无结构化日志
- 无数据库，无数据持久化
- 无用户认证，无权限控制
- CORS 允许所有来源（生产环境风险）

### 配置硬编码
- 飞书凭证硬编码在 `config.go`
- 前端 API 地址硬编码在 `App.jsx`
- 超时时间硬编码（60 秒解析，30 秒飞书 API）

## 性能特征

### 后端性能
- ParserService 超时：60 秒
- Feishu API 超时：30 秒
- 使用浏览器 User-Agent 避免反爬虫

### 前端性能
- 无服务端渲染（CSR）
- 无代码分割（Vite 自动优化）
- 使用 React 19 并发特性