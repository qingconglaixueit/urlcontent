import { useEffect, useRef, useState } from 'react'
import './App.css'

// 后端 API 地址
const API_BASE_URL = 'http://localhost:8080'

function App() {
  const [messages, setMessages] = useState([])
  const [urlInput, setUrlInput] = useState('')
  const [isProcessing, setIsProcessing] = useState(false)
  const [extractedContent, setExtractedContent] = useState(null)
  const [showDirectionSelection, setShowDirectionSelection] = useState(false)
  const [selectedDirection, setSelectedDirection] = useState('')
  const [customDirection, setCustomDirection] = useState('')
  const [useCustomDirection, setUseCustomDirection] = useState(false)
  const [systemStatus, setSystemStatus] = useState('ready')
  const [backendStatus, setBackendStatus] = useState('checking')

  const messagesEndRef = useRef(null)

  // 自动滚动到最新消息
  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  // 检查后端服务状态
  useEffect(() => {
    checkBackendHealth()
    const interval = setInterval(checkBackendHealth, 30000) // 每30秒检查一次
    return () => clearInterval(interval)
  }, [])

  const checkBackendHealth = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/health`)
      if (response.ok) {
        setBackendStatus('online')
      } else {
        setBackendStatus('offline')
      }
    } catch (error) {
      setBackendStatus('offline')
      console.log('后端服务不可用:', error.message)
    }
  }

  // 解析URL内容（调用后端API）
  const parseUrlContent = async (url) => {
    setSystemStatus('processing')
    addBotMessage(`正在解析 URL: ${url}...`)

    console.log('🔍 发送解析请求到后端:', url)

    try {
      const startTime = Date.now()

      const response = await fetch(`${API_BASE_URL}/api/parse`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ url })
      })

      const data = await response.json()
      const duration = Date.now() - startTime

      console.log('✓ 后端响应:', data)
      console.log(`⏱️  耗时: ${duration}ms`)

      if (!data.success) {
        throw new Error(data.message || '解析失败')
      }

      setExtractedContent({
        title: data.title,
        url: data.url,
        content: data.content,
        timestamp: data.timestamp
      })

      addBotMessageWithContent('✅ 内容解析成功！以下是提取的关键信息：', {
        title: data.title,
        url: data.url,
        content: data.content.substring(0, 500) + (data.content.length > 500 ? '...' : '')
      })

      setShowDirectionSelection(true)
      setSystemStatus('success')

    } catch (error) {
      console.error('❌ 解析URL失败:', error)
      addBotMessage(`解析失败: ${error.message}`)
      setSystemStatus('error')
    } finally {
      setIsProcessing(false)
    }
  }

  // 写入飞书知识库（调用后端API）
  const writeToFeishuWiki = async () => {
    setSystemStatus('processing')

    // 保存最终使用的方向，用于后续提示
    let finalDirection = ''

    // 确定最终方向：优先使用自定义方向，如果没有则使用选中的内置方向
    if (customDirection.trim()) {
      finalDirection = customDirection.trim()
      addBotMessage(`📁 使用自定义方向: ${finalDirection}`)
    } else if (selectedDirection) {
      finalDirection = selectedDirection
      addBotMessage(`📁 使用内置方向: ${finalDirection}`)
    } else {
      finalDirection = ''
      addBotMessage('📁 未选择方向，使用默认: 自动总结')
    }

    try {
      addBotMessage('🔄 正在写入飞书知识库...')

      console.log('📝 发送写入请求到后端:', {
        title: extractedContent.title,
        contentLength: extractedContent.content.length,
        direction: finalDirection
      })

      const startTime = Date.now()

      const response = await fetch(`${API_BASE_URL}/api/write`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          title: extractedContent.title,
          content: extractedContent.content,
          direction: finalDirection
        })
      })

      const data = await response.json()
      const duration = Date.now() - startTime

      console.log('✓ 后端响应:', data)
      console.log(`⏱️  耗时: ${duration}ms`)

      if (!data.success) {
        throw new Error(data.message || '写入失败')
      }

      addBotMessage('✅ 内容已成功写入飞书知识库！')

      // 明确提示最终写入的方向
      if (finalDirection) {
        addBotMessage(`📁 最终写入方向: 「${finalDirection}」`)
      } else {
        addBotMessage(`📁 最终写入方向: 自动总结`)
      }

      if (data.documentID) {
        addBotMessage(`📄 文档 ID: ${data.documentID}`)
      }
      if (data.documentName) {
        addBotMessage(`📁 文档名称: 「${data.documentName}」`)
      }
      addBotMessage(`⏰ 写入时间: ${new Date().toLocaleString('zh-CN')}`)

      setSystemStatus('success')
      setExtractedContent(null)
      setShowDirectionSelection(false)
      setSelectedDirection('')
      setCustomDirection('')
      setUseCustomDirection(false)

    } catch (error) {
      console.error('❌ 写入飞书知识库失败:', error)
      addBotMessage(`写入失败: ${error.message}`)
      setSystemStatus('error')
    }
  }

  // 添加机器人消息
  const addBotMessage = (content) => {
    const message = {
      id: Date.now(),
      sender: 'bot',
      content,
      timestamp: new Date().toLocaleTimeString('zh-CN')
    }
    setMessages(prev => [...prev, message])
  }

  // 添加带内容的机器人消息
  const addBotMessageWithContent = (content, extractedData) => {
    const message = {
      id: Date.now(),
      sender: 'bot',
      content,
      extractedData,
      timestamp: new Date().toLocaleTimeString('zh-CN')
    }
    setMessages(prev => [...prev, message])
  }

  // 处理用户输入
  const handleUrlSubmit = (e) => {
    e.preventDefault()

    if (backendStatus !== 'online') {
      addBotMessage('⚠️ 后端服务未连接，请确保后端服务已启动')
      return
    }

    if (!urlInput.trim()) {
      addBotMessage('请输入有效的 URL 地址')
      return
    }

    // 简单的URL验证
    try {
      new URL(urlInput)
    } catch {
      addBotMessage('请输入有效的 URL 格式（例如：https://example.com）')
      return
    }

    setIsProcessing(true)
    setExtractedContent(null)

    // 添加用户消息
    const userMessage = {
      id: Date.now(),
      sender: 'user',
      content: `解析 URL: ${urlInput}`,
      timestamp: new Date().toLocaleTimeString('zh-CN')
    }
    setMessages(prev => [...prev, userMessage])

    // 解析URL
    parseUrlContent(urlInput)

    setUrlInput('')
  }

  // 选择方向
  const handleDirectionSelect = (direction) => {
    setSelectedDirection(direction)
    // 如果已经输入了自定义方向，不清除它，只是取消使用自定义方向的标记
    if (!customDirection.trim()) {
      setUseCustomDirection(false)
    }
  }

  // 选择自定义方向
  const handleCustomDirectionChange = (value) => {
    setCustomDirection(value)
    // 输入自定义方向时，保留已选择的内置方向作为备用
    setUseCustomDirection(true)
  }

  // 拒绝写入
  const handleReject = () => {
    addBotMessage('❌ 已取消写入操作')
    setShowDirectionSelection(false)
    setExtractedContent(null)
    setSelectedDirection('')
    setCustomDirection('')
    setUseCustomDirection(false)
    setSystemStatus('ready')
  }

  // 清除对话历史
  const clearMessages = () => {
    setMessages([])
  }

  return (
    <div className="app-container">
      <header className="header">
        <div className="robot-icon">🤖</div>
        <div className="header-content">
          <h1 className="title">AI 内容同步机器人</h1>
          <p className="subtitle">URL 内容解析 · 飞书知识库集成</p>
        </div>
        <div className={`backend-status ${backendStatus}`}>
          <span className="status-indicator"></span>
          <span>
            {backendStatus === 'online' && '后端在线'}
            {backendStatus === 'offline' && '后端离线'}
            {backendStatus === 'checking' && '检查中...'}
          </span>
        </div>
      </header>

      <div className="main-panel">
        <div className="chat-panel">
          <h2 className="section-title">💬 对话交互</h2>

          <div className="system-status">
            <span className="system-status-item">
              <span className={`status-dot ${systemStatus === 'ready' ? '' : systemStatus}`}></span>
              <span>系统状态</span>
            </span>
            <span className={`status-text status-${systemStatus}`}>
              {systemStatus === 'ready' && '就绪'}
              {systemStatus === 'processing' && '处理中'}
              {systemStatus === 'success' && '成功'}
              {systemStatus === 'error' && '错误'}
            </span>
          </div>

          <div className="chat-messages">
            {messages.length === 0 && (
              <div className="message bot">
                <div className="message-content">
                  👋 欢迎使用 AI 内容同步机器人！\n\n
                  📌 功能说明：\n
                  1. 输入 URL 链接，自动解析内容\n
                  2. 提取关键信息并预览\n
                  3. 确认后自动写入飞书知识库的"自动总结"目录\n\n
                  🚀 后端服务状态：{backendStatus === 'online' ? '✅ 在线' : '⚠️ 离线'}\n\n
                  💡 如果后端显示离线，请先启动后端服务：\n
                  cd backend && ./bin/server
                </div>
              </div>
            )}

            {messages.map(msg => (
              <div key={msg.id} className={`message ${msg.sender}`}>
                <div className="message-header">
                  <span className="message-sender">
                    {msg.sender === 'bot' ? '🤖 机器人' : '👤 用户'}
                  </span>
                  <span className="message-time">{msg.timestamp}</span>
                </div>
                <div className="message-content">{msg.content}</div>
                {msg.extractedData && (
                  <div className="extracted-content">
                    <div className="extracted-content-title">📄 提取的内容</div>
                    <div className="extracted-content-text">
                      <strong>标题：</strong>{msg.extractedData.title}\n\n
                      <strong>链接：</strong>{msg.extractedData.url}\n\n
                      <strong>内容预览：</strong>\n{msg.extractedData.content}
                    </div>
                  </div>
                )}
              </div>
            ))}

            {isProcessing && (
              <div className="message bot">
                <div className="message-content">
                  <span className="loading"></span>
                  正在处理中...
                </div>
              </div>
            )}

            <div ref={messagesEndRef} />
          </div>

          {showDirectionSelection && extractedContent && (
            <div className="direction-selection">
              <div className="direction-selection-title">📁 请选择内容分类方向：</div>

              <div className="direction-buttons">
                <button
                  onClick={() => handleDirectionSelect('技术')}
                  className={`direction-button ${selectedDirection === '技术' ? 'selected' : ''}`}
                >
                  💻 技术
                </button>
                <button
                  onClick={() => handleDirectionSelect('成长')}
                  className={`direction-button ${selectedDirection === '成长' ? 'selected' : ''}`}
                >
                  🌱 成长
                </button>
                <button
                  onClick={() => handleDirectionSelect('生活')}
                  className={`direction-button ${selectedDirection === '生活' ? 'selected' : ''}`}
                >
                  🎯 生活
                </button>
              </div>

              <div className="custom-direction-input">
                <div className="custom-direction-label">或自定义方向：</div>
                <input
                  type="text"
                  className="custom-input"
                  placeholder="输入自定义方向名称"
                  value={customDirection}
                  onChange={(e) => handleCustomDirectionChange(e.target.value)}
                  disabled={isProcessing}
                />
              </div>

              <div className="direction-confirm-buttons">
                <button
                  onClick={handleReject}
                  className="reject"
                >
                  ❌ 取消
                </button>
                <button
                  onClick={writeToFeishuWiki}
                  className="approve"
                  disabled={!selectedDirection && !customDirection.trim()}
                >
                  ✅ 确认写入知识库
                </button>
              </div>
            </div>
          )}

          <form onSubmit={handleUrlSubmit} className="input-area">
            <input
              type="text"
              className="url-input"
              placeholder={backendStatus === 'online' ? "输入要解析的 URL（例如：https://example.com）" : "等待后端服务启动..."}
              value={urlInput}
              onChange={(e) => setUrlInput(e.target.value)}
              disabled={isProcessing || backendStatus !== 'online'}
            />
            <div className="action-buttons">
              <button
                type="submit"
                disabled={isProcessing || !urlInput.trim() || backendStatus !== 'online'}
              >
                {isProcessing ? '解析中...' : '🔍 解析内容'}
              </button>
              <button
                type="button"
                onClick={clearMessages}
                disabled={isProcessing || messages.length === 0}
                className="secondary"
              >
                🗑️ 清除
              </button>
            </div>
          </form>
        </div>

        <div className="info-panel">
          <h2 className="section-title">📊 系统信息</h2>

          <div className="info-item">
            <div className="info-label">后端 API</div>
            <div className="info-value">
              {API_BASE_URL}
              <span className={`status-badge ${backendStatus}`}>
                {backendStatus === 'online' && '✓ 运行中'}
                {backendStatus === 'offline' && '✗ 停止'}
                {backendStatus === 'checking' && '⟳ 检查中'}
              </span>
            </div>
          </div>

          <div className="info-item">
            <div className="info-label">可用端点</div>
            <div className="info-value">
              <div>GET  /health - 健康检查</div>
              <div>POST /api/parse - URL 解析</div>
              <div>POST /api/write - 写入飞书</div>
            </div>
          </div>

          <div className="info-item">
            <div className="info-label">功能特性</div>
            <div className="info-value">
              <div>✅ 智能 URL 内容解析</div>
              <div>✅ 自动创建"自动总结"目录</div>
              <div>✅ 飞书知识库集成</div>
              <div>✅ 用户确认机制</div>
              <div>✅ 实时状态反馈</div>
            </div>
          </div>

          <div className="info-item">
            <div className="info-label">技术栈</div>
            <div className="info-value">
              <div>🎨 前端: React + Vite</div>
              <div>⚙️  后端: Go (net/http)</div>
              <div>📚 存储: 飞书知识库</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default App
