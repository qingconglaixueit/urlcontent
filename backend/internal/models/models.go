package models

// ParseRequest URL 解析请求
type ParseRequest struct {
	URL string `json:"url" binding:"required"`
}

// ParseResponse URL 解析响应
type ParseResponse struct {
	Success   bool   `json:"success"`
	Title     string `json:"title,omitempty"`
	URL       string `json:"url,omitempty"`
	Content   string `json:"content,omitempty"`
	Message   string `json:"message,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// WriteRequest 写入飞书请求
type WriteRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Direction string `json:"direction,omitempty"` // 方向/目录名称
}

// WriteResponse 写入飞书响应
type WriteResponse struct {
	Success     bool   `json:"success"`
	DocumentID  string `json:"document_id,omitempty"`
	DocumentName string `json:"document_name,omitempty"` // 最终写入的文档名称
	Message     string `json:"message,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
