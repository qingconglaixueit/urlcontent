package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"urltocontent/backend/internal/config"
	"urltocontent/backend/internal/models"
	"urltocontent/backend/internal/services"
)

type Handler struct {
	Config *config.Config
	Parser *services.ParserService
	Feishu *services.FeishuService
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		Config: cfg,
		Parser: services.NewParserService(),
		Feishu: services.NewFeishuService(cfg.FeishuAppID, cfg.FeishuSecret, cfg.FeishuWikiID),
	}
}

// CORSMiddleware 处理跨域请求
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// ParseURLHandler 处理 URL 解析请求
func (h *Handler) ParseURLHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== 收到 URL 解析请求 ===")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "只支持 POST 请求",
		})
		return
	}

	var req models.ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("❌ 请求体解析失败: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": fmt.Sprintf("请求格式错误: %v", err),
		})
		return
	}

	fmt.Printf("🔗 目标 URL: %s\n", req.URL)

	extracted, err := h.Parser.ParseURL(req.URL)
	if err != nil {
		fmt.Printf("❌ URL 解析失败: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ParseResponse{
			Success: false,
			URL:     req.URL,
			Content: "",
			Message: fmt.Sprintf("URL 解析失败: %v", err),
		})
		return
	}

	fmt.Printf("✅ 解析完成\n")
	fmt.Printf("📝 标题: %s\n", extracted.Title)

	response := models.ParseResponse{
		Success:   true,
		Title:     extracted.Title,
		URL:       req.URL,
		Content:   extracted.Content,
		Timestamp: extracted.Timestamp,
		Message:   "解析成功",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Println("=== URL 解析请求完成 ===\n")
}

// WriteToFeishuHandler 处理写入飞书知识库请求
func (h *Handler) WriteToFeishuHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n=== 收到写入飞书请求 ===")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "只支持 POST 请求",
		})
		return
	}

	var req models.WriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("❌ 请求体解析失败: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": fmt.Sprintf("请求格式错误: %v", err),
		})
		return
	}

	fmt.Printf("📄 文档标题: %s\n", req.Title)
	fmt.Printf("📝 内容长度: %d 字符\n", len(req.Content))
	fmt.Printf("📁 指定方向: %s\n", req.Direction)

	// 查找或创建方向文档
	directionToken, directionName, err := h.Feishu.FindOrCreateDocument(req.Direction)
	if err != nil {
		fmt.Printf("❌ 查找或创建方向文档失败: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     false,
			"message":     fmt.Sprintf("查找或创建方向文档失败: %v", err),
			"documentID":  "",
			"documentName": "",
		})
		return
	}

	fmt.Printf("✅ 使用方向文档: %s (token: %s)\n", directionName, directionToken)

	// 在方向文档下创建子文档
	documentID, err := h.Feishu.CreateDocumentInNode(directionToken, req.Title, req.Content)
	if err != nil {
		fmt.Printf("❌ 写入飞书失败: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     false,
			"message":     fmt.Sprintf("写入飞书失败: %v", err),
			"documentID":  "",
			"documentName": "",
		})
		return
	}

	fmt.Printf("✅ 写入成功，文档ID: %s\n", documentID)
	fmt.Printf("✅ 写入到方向: %s\n", directionName)
	fmt.Println("=== 写入飞书请求完成 ===\n")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"message":     fmt.Sprintf("成功写入飞书知识库的「%s」文档", directionName),
		"documentID":  documentID,
		"documentName": directionName,
	})
}

// HealthCheckHandler 健康检查
func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"service": "urlToContent API",
	})
}
