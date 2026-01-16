# AI 内容同步机器人

<div align="center">

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![React](https://img.shields.io/badge/React-19.2.0-61DAFB.svg)
![Go](https://img.shields.io/badge/Go-1.23-00ADD8.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

**一个功能强大的全栈应用，可以解析网页内容并将其自动写入飞书知识库**

[快速开始](#快速开始) • [功能特性](#功能特性) • [API 文档](#api-文档) • [常见问题](#常见问题)

</div>

## 📖 项目简介

AI 内容同步机器人是一个全栈 Web 应用，采用 React 前端和 Go 后端架构，提供机器人风格的炫酷界面。它能够智能解析任意网页 URL，自动提取标题、描述和关键内容，并一键写入飞书知识库，帮助用户快速收集和整理网络信息。

## ✨ 功能特性

### 核心功能
- 🤖 **智能 URL 解析** - 自动提取网页标题、描述和关键内容，支持内容总结
- 📚 **飞书知识库集成** - 无缝对接飞书知识库，一键写入文档
- 📁 **方向/目录管理** - 支持自定义方向（目录），智能创建和管理文档结构
- ✅ **用户确认机制** - 写入前预览内容，确保信息准确性
- 🔄 **纯 Go 后端** - 使用标准库实现，零外部依赖，高性能稳定

### UI/UX 特性
- 🎨 **机器人风格 UI** - 炫酷的赛博朋克风格界面
- 💬 **实时状态反馈** - 清晰的操作进度和结果提示
- 📊 **后端健康监测** - 自动检测后端服务状态
- 🔔 **智能消息提示** - 详细的成功/错误信息展示

### 解析能力
- 📝 **HTML 内容提取** - 智能识别标题、段落、列表等元素
- 🎯 **元数据解析** - 提取 meta description、og:title 等信息
- 🧹 **内容清理** - 自动移除脚本、样式等无关标签
- 📋 **内容总结** - 自动生成内容摘要和关键信息

## 🔧 技术栈

### 前端技术
- **React 19.2.0** - 最新版本的 React 框架
- **Vite 7.2.4** - 快速的前端构建工具
- **CSS3** - 现代化样式和动画效果
- **ESLint** - 代码质量检查工具

### 后端技术
- **Go 1.23** - 高性能编程语言
- **标准库 net/http** - 零外部依赖的 HTTP 服务
- **RESTful API** - 标准的 JSON API 接口
- **CORS 支持** - 跨域资源共享
- **飞书开放平台 API** - 集成飞书知识库

### 系统要求
- **Go**: 1.23 或更高版本
- **Node.js**: 18.x 或更高版本
- **npm**: 9.x 或更高版本
- **操作系统**: Windows / macOS / Linux

## 📁 项目结构

```
urltocontent/
├── backend/                      # Go 后端服务
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # 主程序入口
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go       # 配置管理（支持环境变量）
│   │   ├── handlers/
│   │   │   └── handlers.go     # HTTP 处理器和路由
│   │   ├── models/
│   │   │   └── models.go       # 数据模型定义
│   │   └── services/
│   │       ├── feishu.go      # 飞书 API 集成服务
│   │       └── parser.go      # URL 解析服务
│   ├── bin/                     # 编译产物
│   │   ├── server              # Linux/Mac 可执行文件
│   │   └── server.exe          # Windows 可执行文件
│   └── go.mod                   # Go 模块配置
├── src/                         # React 前端应用
│   ├── App.jsx                 # 主应用组件
│   ├── App.css                 # 应用样式
│   ├── main.jsx                # 应用入口
│   ├── index.css               # 全局样式
│   └── assets/                 # 静态资源
├── public/                      # 公共静态文件
│   └── vite.svg
├── .gitignore                   # Git 忽略配置
├── eslint.config.js             # ESLint 配置
├── index.html                   # HTML 模板
├── package.json                 # Node.js 依赖配置
├── vite.config.js               # Vite 构建配置
└── README.md                    # 项目文档
```

## 🚀 快速开始

### 前置条件

在开始之前，请确保您的系统已安装以下软件：

- **Go 1.23+** - [下载地址](https://golang.org/dl/)
- **Node.js 18+** - [下载地址](https://nodejs.org/)
- **npm 9+** - 随 Node.js 一起安装

### 1. 克隆项目

```bash
# 克隆仓库
git clone <repository-url>
cd urltocontent
```

### 2. 配置飞书应用

首次使用需要配置飞书应用凭据：

#### 创建飞书应用
1. 访问 [飞书开放平台](https://open.feishu.cn/)
2. 创建企业自建应用
3. 获取 App ID 和 App Secret
4. 在权限管理中开通以下权限：
   - `wiki:wiki:read` - 读取知识库
   - `wiki:wiki:write` - 写入知识库
   - `docx:document:write` - 创建文档
5. 获取知识库 Wiki ID

#### 配置环境变量（推荐）

```bash
# Windows CMD
set FEISHU_APP_ID=your_app_id
set FEISHU_APP_SECRET=your_app_secret
set FEISHU_WIKI_ID=your_wiki_id

# Windows PowerShell
$env:FEISHU_APP_ID="your_app_id"
$env:FEISHU_APP_SECRET="your_app_secret"
$env:FEISHU_WIKI_ID="your_wiki_id"

# Linux/macOS
export FEISHU_APP_ID="your_app_id"
export FEISHU_APP_SECRET="your_app_secret"
export FEISHU_WIKI_ID="your_wiki_id"
```

### 3. 启动后端服务

```bash
# 进入后端目录
cd backend

# 下载 Go 依赖
go mod download

# 构建项目
go build -o bin/server ./cmd/server

# 启动服务（Windows）
bin\server.exe

# 启动服务（Linux/Mac）
./bin/server
```

后端服务将在 `http://localhost:8080` 启动。

**预期输出：**
```
====================
🤖 URL to Content API
====================

📋 服务器端口: 8080
🚀 飞书 App ID: cli_xxxxxx
📚 飞书 Wiki ID: xxxxxxxx

✅ 服务器启动成功，监听地址: http://localhost:8080

可用端点:
  - GET  /health    - 健康检查
  - POST /api/parse - URL 解析
  - POST /api/write - 写入飞书

按 Ctrl+C 停止服务器
```

### 4. 启动前端服务

打开新的终端窗口：

```bash
# 安装前端依赖（仅首次运行需要）
npm install

# 启动开发服务器
npm run dev
```

前端服务将在 `http://localhost:5173` 启动。

**预期输出：**
```
VITE v7.2.4  ready in xxx ms

➜  Local:   http://localhost:5173/
➜  Network: use --host to expose
```

### 5. 访问应用

在浏览器中打开 `http://localhost:5173`，即可看到机器人风格的界面。

## ⚙️ 配置说明

### 环境变量配置

应用支持通过环境变量进行配置，所有配置项都有默认值。

| 环境变量            | 默认值 | 说明          |
| ------------------- | ------ | ------------- |
| `PORT`              | `8080` | 后端服务端口  |
| `FEISHU_APP_ID`     | `xxx`  | 飞书应用 ID   |
| `FEISHU_APP_SECRET` | `xxx`  | 飞书应用密钥  |
| `FEISHU_WIKI_ID`    | `xxx`  | 飞书知识库 ID |

**配置文件位置：** `backend/internal/config/config.go`

### 飞书知识库配置

#### 权限要求
确保飞书应用拥有以下权限：
- ✅ `wiki:wiki:read` - 读取知识库
- ✅ `wiki:wiki:write` - 写入知识库
- ✅ `docx:document:write` - 创建文档

#### 文档结构
应用会在飞书知识库中自动管理文档结构：

```
知识库根目录
├── 自动总结（默认目录）
│   ├── 网页标题 1
│   ├── 网页标题 2
│   └── ...
├── 自定义方向 1
│   ├── 网页标题 3
│   └── ...
└── 自定义方向 2
    ├── 网页标题 4
    └── ...
```

#### 方向/目录功能
- **默认方向**：如果未指定方向，内容将写入"自动总结"目录
- **自定义方向**：可以输入自定义的目录名称，系统会自动创建
- **方向缓存**：已创建的方向会被缓存，提升性能

## 📡 API 文档

### 基础信息

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **CORS**: 已启用，支持跨域请求

### 1. 健康检查

检查后端服务是否正常运行。

**端点：**
```http
GET /health
```

**响应示例：**
```json
{
  "status": "ok",
  "service": "urlToContent API"
}
```

**状态码：**
- `200` - 服务正常
- `500` - 服务异常

### 2. URL 解析

解析指定 URL 的网页内容，提取标题、描述和关键信息。

**端点：**
```http
POST /api/parse
Content-Type: application/json
```

**请求体：**
```json
{
  "url": "https://example.com"
}
```

**参数说明：**
| 参数 | 类型   | 必填 | 说明             |
| ---- | ------ | ---- | ---------------- |
| url  | string | 是   | 要解析的网页 URL |

**响应示例：**
```json
{
  "success": true,
  "title": "网页标题",
  "url": "https://example.com",
  "content": "标题：网页标题\n\n描述：网页描述...\n\n来源链接：https://example.com\n\n内容总结：\n提取的关键内容...",
  "timestamp": "2026-01-16 10:26:53",
  "message": "解析成功"
}
```

**响应字段说明：**
| 字段      | 类型    | 说明                               |
| --------- | ------- | ---------------------------------- |
| success   | boolean | 请求是否成功                       |
| title     | string  | 网页标题                           |
| url       | string  | 原始 URL                           |
| content   | string  | 提取的内容（包含标题、描述、总结） |
| timestamp | string  | 解析时间戳                         |
| message   | string  | 响应消息                           |

**错误响应：**
```json
{
  "success": false,
  "message": "URL 解析失败: 无法访问该网站"
}
```

### 3. 写入飞书知识库

将解析的内容写入飞书知识库。

**端点：**
```http
POST /api/write
Content-Type: application/json
```

**请求体：**
```json
{
  "title": "文档标题",
  "content": "文档内容",
  "direction": "自定义方向"
}
```

**参数说明：**
| 参数      | 类型   | 必填 | 说明                                  |
| --------- | ------ | ---- | ------------------------------------- |
| title     | string | 是   | 文档标题                              |
| content   | string | 是   | 文档内容                              |
| direction | string | 否   | 方向/目录名称（可选，默认"自动总结"） |

**响应示例：**
```json
{
  "success": true,
  "message": "成功写入飞书知识库的「自定义方向」文档",
  "documentID": "doc_xxxxxxxxxxxxxxxx",
  "documentName": "自定义方向"
}
```

**响应字段说明：**
| 字段         | 类型    | 说明             |
| ------------ | ------- | ---------------- |
| success      | boolean | 请求是否成功     |
| message      | string  | 响应消息         |
| documentID   | string  | 飞书文档 ID      |
| documentName | string  | 文档所在方向名称 |

**错误响应：**
```json
{
  "success": false,
  "message": "查找或创建方向文档失败: 权限不足",
  "documentID": "",
  "documentName": ""
}
```

## 📖 使用说明

### Web 界面使用

1. **启动服务**
   - 确保后端服务运行在 `http://localhost:8080`
   - 前端服务运行在 `http://localhost:5173`
   - 打开浏览器访问 `http://localhost:5173`

2. **解析 URL**
   - 在输入框中输入要解析的网页 URL
   - 点击"解析内容"按钮或按 Enter 键
   - 等待解析完成，查看提取的内容和总结
   - 系统会自动检测后端状态并提示

3. **选择方向**
   - 解析成功后，可以选择文档写入方向
   - 支持选择内置方向或输入自定义方向
   - 如不选择，默认写入"自动总结"目录

4. **写入知识库**
   - 确认预览内容无误后，点击"确认写入知识库"
   - 等待写入完成，查看结果
   - 成功后会显示文档 ID 和写入方向
   - 如需取消，点击"取消写入"

5. **查看日志**
   - 所有操作都会在界面中显示详细日志
   - 包含操作时间、状态和结果信息

### 命令行使用（API 调用）

**使用 curl：**

```bash
# 健康检查
curl http://localhost:8080/health

# 解析 URL
curl -X POST http://localhost:8080/api/parse \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'

# 写入飞书
curl -X POST http://localhost:8080/api/write \
  -H "Content-Type: application/json" \
  -d '{
    "title":"测试文档",
    "content":"这是测试内容",
    "direction":"测试方向"
  }'
```

**使用 PowerShell：**

```powershell
# 健康检查
Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get

# 解析 URL
$body = @{
  url = "https://example.com"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/parse" `
  -Method Post `
  -ContentType "application/json" `
  -Body $body
```

## 💻 开发说明

### 后端开发

```bash
# 进入后端目录
cd backend

# 格式化代码
go fmt ./...

# 运行测试
go test ./...

# 查看测试覆盖率
go test -cover ./...

# 交叉编译
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/server-linux ./cmd/server

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o bin/server-mac ./cmd/server

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o bin/server-mac-arm ./cmd/server

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/server.exe ./cmd/server

# 运行程序（开发模式）
go run ./cmd/server
```

**项目结构说明：**
- `cmd/server/` - 应用入口
- `internal/config/` - 配置管理
- `internal/handlers/` - HTTP 处理器
- `internal/models/` - 数据模型
- `internal/services/` - 业务逻辑服务

### 前端开发

```bash
# 开发模式（热重载）
npm run dev

# 代码检查
npm run lint

# 构建生产版本
npm run build

# 预览生产构建
npm run preview
```

**开发技巧：**
- 使用浏览器开发者工具（F12）查看网络请求和控制台日志
- 前端默认连接 `http://localhost:8080`，可在 [`src/App.jsx`](src/App.jsx:5) 中修改
- 后端每 30 秒自动检测一次健康状态

## ❓ 常见问题

### 1. 后端无法启动

**问题：** 运行 `./bin/server` 时报错 "bind: address already in use"

**解决方案：**

Windows:
```cmd
# 查看占用端口的进程
netstat -ano | findstr :8080

# 杀死进程（替换 PID）
taskkill /PID <进程ID> /F

# 或修改端口号
set PORT=8081
```

Linux/macOS:
```bash
# 查看占用端口的进程
lsof -i :8080

# 杀死进程（替换 PID）
kill -9 <PID>

# 或修改端口号
export PORT=8081
```

### 2. 前端显示"后端离线"

**问题：** 前端界面显示后端服务不可用

**解决方案：**
```bash
# 检查后端是否运行
curl http://localhost:8080/health

# 如果未运行，启动后端
cd backend
go run ./cmd/server
```

**检查清单：**
- ✅ 后端服务是否正常启动
- ✅ 端口是否正确（默认 8080）
- ✅ 防火墙是否阻止连接
- ✅ 浏览器控制台是否有错误信息

### 3. 飞书 API 调用失败

**问题：** 写入飞书时提示权限错误或 API 错误

**解决方案：**

检查飞书应用配置：
1. 登录 [飞书开放平台](https://open.feishu.cn/)
2. 进入应用详情
3. 确认以下权限已开启：
   - `wiki:wiki:read`
   - `wiki:wiki:write`
   - `docx:document:write`
4. 确认知识库已添加应用为成员
5. 验证 API 凭证是否正确

**常见错误码：**
- `99991663` - 权限不足
- `99991400` - 参数错误
- `99991600` - 知识库不存在

### 4. URL 解析失败

**问题：** 解析 URL 时提示错误

**可能原因和解决方案：**

| 错误类型     | 原因                         | 解决方案               |
| ------------ | ---------------------------- | ---------------------- |
| URL 格式错误 | 缺少 `http://` 或 `https://` | 添加协议前缀           |
| 网络错误     | 无法访问目标网站             | 检查网络连接或更换 URL |
| 超时         | 网站响应过慢                 | 增加超时时间或更换 URL |
| 反爬机制     | 网站拒绝爬虫访问             | 尝试手动复制内容       |

**调试技巧：**
```bash
# 手动测试 URL 可访问性
curl -I https://example.com

# 查看后端日志中的详细信息
# 后端会输出完整的解析过程和错误原因
```

### 5. 写入后找不到文档

**问题：** 写入成功但在飞书中找不到文档

**解决方案：**
1. 检查是否选择了正确的知识库（通过 Wiki ID）
2. 刷新飞书知识库页面
3. 查看是否有"自动总结"目录
4. 检查文档是否在其他目录下
5. 确认飞书应用是否有查看权限

## ⚠️ 注意事项

### 使用建议
- 📊 **API 速率限制** - 请勿频繁调用 API，避免触发飞书速率限制（建议间隔 > 1 秒）
- 🔒 **敏感信息安全** - 生产环境务必使用环境变量配置飞书凭证
- 🌐 **网络稳定性** - URL 解析依赖网络，不稳定可能导致解析失败
- 📝 **内容验证** - 写入前请预览内容，确保准确性

### 性能优化
- 后端已实现方向 Token 缓存，提升重复写入性能
- URL 解析默认超时 60 秒，可根据需要调整
- 前端每 30 秒自动检测后端健康状态

### 已知限制
- 部分网站可能有反爬机制，无法正常解析
- 飞书 Wiki API 不支持直接创建文件夹，使用文档替代
- 单次写入内容大小受飞书 API 限制

## 🔍 调试指南

### 浏览器调试

1. 打开浏览器开发者工具（F12）
2. 切换到 "Network" 标签
3. 执行解析或写入操作
4. 查看请求和响应详情
5. 切换到 "Console" 标签查看前端日志

### 后端调试

后端会输出详细的运行日志：

```
═════════════════════════════════════════════════════════════
🔍 开始解析 URL
═════════════════════════════════════════════════════════════
📡 目标 URL: https://example.com
📡 步骤 1: 直接获取网页内容
─────────────────────────────────────────────────────────────
✅ 响应状态: 200 OK
⏱️  请求耗时: 523ms
📦 HTML 内容长度: 45678 字符
...
```

### 日志级别

- `INFO` - 正常操作信息
- `WARN` - 警告信息（如使用缓存）
- `ERROR` - 错误信息（如 API 调用失败）

## 📄 许可证

MIT License

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎贡献代码、报告问题或提出改进建议！

### 如何贡献

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 代码规范

- **Go**: 遵循 [Effective Go](https://go.dev/doc/effective_go) 和 `gofmt` 格式
- **React**: 遵循 [Airbnb React/JSX Style Guide](https://github.com/airbnb/javascript/tree/master/react)
- **提交信息**: 使用清晰的提交信息格式

## 📧 联系方式

- **项目主页**: [GitHub Repository]
- **问题反馈**: [GitHub Issues]
- **邮件**: [项目邮箱]

## 🙏 致谢

感谢以下开源项目和工具：

- [React](https://react.dev/) - 前端框架
- [Vite](https://vitejs.dev/) - 构建工具
- [Go](https://go.dev/) - 后端语言
- [飞书开放平台](https://open.feishu.cn/) - API 支持

---

<div align="center">

**Made with ❤️ by AI Content Sync Bot Team**

[⬆ 回到顶部](#ai-内容同步机器人)

</div>
