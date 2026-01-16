# Project Debug Rules (Non-Obvious Only)

## 新功能：方向选择机制调试

### 前端调试要点
- 方向选择界面在 URL 解析成功后显示
- 检查内置方向按钮（技术、成长、生活）是否正确触发
- 自定义方向输入框的验证和提交逻辑
- 默认使用"自动总结"的降级逻辑

### 后端调试要点
- 检查文档是否存在：日志应显示查找特定方向的文档
- 文档创建：首次使用某个方向时应显示"创建文档"
- 文档复用：已存在的方向应显示"使用已有文档"
- 响应包含 `documentName` 字段，显示最终写入的文档名称

### 调试场景

#### 方向文档不存在
- 查找文档时日志应显示"未找到 [方向名] 文档"
- 创建文档时日志应显示"创建 [方向名] 文档成功"

#### 方向文档已存在
- 查找文档时日志应显示"找到 [方向名] 文档: [token]"
- 不应出现创建文档的日志

#### 默认方向
- 未选择方向时，日志应显示"使用默认方向: 自动总结"
- 检查"自动总结"文档是否存在

## 后端调试

### 日志位置
- 所有 Go 后端日志输出到标准输出（终端），而非日志文件
- 使用 `fmt.Println` 进行日志输出，无结构化日志库
- 关键操作都有详细的日志输出（带 emoji 标记）

### 日志级别
- 使用 emoji 区分日志类型：
  - ✅ 成功操作
  - ❌ 错误信息
  - ⚠️ 警告信息
  - 🔍 查询/搜索操作
  - 📝 写入操作
  - 📡 网络请求
  - ⏱️ 性能计时

### 常见调试场景

#### 后端无法启动
```bash
# 检查端口占用（Windows）
netstat -ano | findstr :8080

# 修改端口
set PORT=8081
cd backend
go build -o bin/server ./cmd/server
./bin/server
```

#### 飞书 API 失败
- 检查 `backend/internal/config/config.go` 中的凭证
- 查看终端输出的飞书 API 响应（包含详细错误码）
- 飞书 API 错误码非 0 时会输出具体错误信息

#### URL 解析失败
- ParserService 超时设置为 60 秒
- 查看 "📡 步骤 1: 直接获取网页内容" 后的响应状态码
- 检查是否被网站反爬虫（已设置浏览器 User-Agent）

### 调试命令

#### 查看详细 HTTP 请求
- 后端会在每个请求开始时输出 "=== 收到 XXX 请求 ==="
- 请求结束时输出 "=== XXX 请求完成 ==="

#### 性能分析
- URL 解析会输出详细步骤和耗时
- ParserService 输出：
  - 请求耗时（毫秒）
  - HTML 元素统计
  - 提取的内容长度

## 前端调试

### 浏览器开发者工具
- React 应用运行在 `http://localhost:5173`
- 使用浏览器开发者工具（F12）查看网络请求和控制台日志

### API 请求调试
- 所有 API 请求到 `http://localhost:8080`
- 网络标签查看 `/api/parse` 和 `/api/write` 请求
- 控制台查看带前缀的日志：
  - 🔍 发送解析请求到后端
  - ✓ 后端响应
  - ⏱️ 耗时
  - ❌ 解析URL失败
  - 📝 发送写入请求到后端

### 状态调试
- 检查组件状态：
  - `backendStatus`: 'online', 'offline', 'checking'
  - `systemStatus`: 'ready', 'processing', 'success', 'error'
  - `extractedContent`: 解析后的内容
  - `showConfirmButtons`: 是否显示确认按钮

### 后端连接问题
- 前端每 30 秒自动检查后端健康状态
- 后端离线时右上角显示 "后端离线"
- 检查后端是否在 `http://localhost:8080` 运行

## 测试调试

### 运行测试
```bash
cd backend
go test ./internal/services -v
```

### 测试输出
- 测试会实际写入飞书知识库
- 查看终端输出的文档 ID
- 在飞书知识库中验证文档内容

### 测试失败
- 检查飞书 API 凭证是否正确
- 确保网络可以访问 `https://open.feishu.cn`
- 查看测试输出的详细错误信息

## 生产构建调试

### 后端
```bash
cd backend
go build -o bin/server ./cmd/server
./bin/server
```

### 前端
```bash
npm run build
npm run preview
```

### 注意事项
- 生产构建后，CORS 可能导致问题（开发时允许所有来源）
- 确保 `API_BASE_URL` 指向正确的后端地址