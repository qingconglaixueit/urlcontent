# Project Coding Rules (Non-Obvious Only)

## 新功能：方向选择机制

### 前端实现要点
- URL 解析完成后，显示方向选择界面
- 内置方向：技术、成长、生活（硬编码在组件中）
- 自定义方向：允许用户输入任意目录名
- 默认方向：用户未选择时自动使用"自动总结"

### 后端实现要点
- 方向/目录在飞书中是文档类型（obj_type: "docx"）
- 必须先检查文档是否存在于知识库根目录
- 不存在则创建，存在则直接使用
- 写入完成后返回文档名称供前端反馈

### API 扩展
- `/api/write` 请求需添加 `direction` 字段（可选）
- 响应需包含 `documentName` 字段，显示最终写入的文档名称
- 文档检查和创建逻辑在 FeishuService 中实现

## Go 后端开发

### 模块路径
- Go 模块名为 `urltocontent/backend`（不是 `backend`）
- 所有导入路径必须使用完整模块路径：`urltocontent/backend/internal/config`

### 构建要求
- 必须从 `backend/` 目录执行 `go build` 和 `go test`
- 从根目录执行会导致路径错误

### 飞书 API 集成
- 飞书 Wiki API **不支持直接创建文件夹节点**
- 只能创建 `obj_type: "docx"` 的文档类型
- "自动总结"实际上是文档标题，不是真实文件夹
- 文档创建在知识库根目录下

### 服务配置
- ParserService 超时硬编码为 60 秒（services/parser.go:19）
- HTTP Client 用于飞书 API 超时为 30 秒（services/feishu.go:91）
- 使用自定义 User-Agent 模拟浏览器，避免被识别为爬虫

## React 前端开发

### API 配置
- 后端 API 地址硬编码：`const API_BASE_URL = 'http://localhost:8080'`（src/App.jsx:5）
- 修改后端端口需要同步修改此常量

### 状态管理
- 使用 React Hooks（useState, useEffect, useRef）
- 后端健康检查每 30 秒自动执行（src/App.jsx:30）
- 使用 `useRef` 实现自动滚动到消息底部（src/App.jsx:19）

### ESLint 规则
- 未使用变量如果以大写字母或下划线开头则忽略：`'no-unused-vars': ['error', { varsIgnorePattern: '^[A-Z_]' }]`
- React Hooks 和 React Refresh 插件已配置

## 错误处理

### Go 后端
- 所有 HTTP 响应使用 `fmt.Println` 输出日志，而非结构化日志
- 错误使用 `fmt.Errorf` 包装，提供上下文信息
- HTTP 状态码：400（请求错误）、500（服务器错误）

### React 前端
- API 调用错误使用 `try-catch` 捕获
- 错误消息通过 `addBotMessage` 显示在聊天界面
- 后端离线状态自动检测并显示警告

## 测试注意
- 测试文件包含实际飞书凭证（仅用于测试）
- 运行测试会实际写入飞书知识库
- 必须在 backend 目录下执行：`go test ./internal/services`