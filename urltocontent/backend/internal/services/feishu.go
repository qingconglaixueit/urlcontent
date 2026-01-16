package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	SummaryFolderName = "自动总结"
)

type FeishuService struct {
	AppID            string
	AppSecret        string
	BaseURL          string
	WikiID           string
	httpClient       *http.Client
	summaryFolderToken string
	// 缓存方向文档的 token
	directionTokens map[string]string
}

type TenantAccessTokenResponse struct {
	Code              int    `json:"code"`
	Msg               string `json:"msg"`
	Expire            int    `json:"expire"`
	TenantAccessToken string `json:"tenant_access_token"`
}

type CreateDocumentRequest struct {
	Title       string `json:"title"`
	ParentToken string `json:"parent_node_token,omitempty"`
	ObjType     string `json:"obj_type"`
	NodeType    string `json:"node_type"`
}

type CreateDocumentResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Document struct {
			DocumentID string `json:"document_id"`
		} `json:"document"`
	} `json:"data"`
}

// CreateBlockResponse 飞书批量创建块响应
type CreateBlockResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		BlockIDs []string `json:"block_ids"`
	} `json:"data"`
}

type GetNodesResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Items []struct {
			NodeToken string `json:"node_token"`
			Title     string `json:"title"`
			ObjType   string `json:"obj_type"`
		} `json:"items"`
		HasMore bool   `json:"has_more"`
		Token   string `json:"page_token"`
	} `json:"data"`
}

type NodeResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Node struct {
			NodeToken string `json:"node_token"`
			ObjToken  string `json:"obj_token"`
			Title     string `json:"title"`
			ObjType   string `json:"obj_type"`
		} `json:"node"`
	} `json:"data"`
}

func NewFeishuService(appID, appSecret, wikiID string) *FeishuService {
	return &FeishuService{
		AppID:           appID,
		AppSecret:       appSecret,
		WikiID:          wikiID,
		BaseURL:         "https://open.feishu.cn",
		httpClient:      &http.Client{Timeout: 30 * time.Second},
		directionTokens: make(map[string]string),
	}
}

func (s *FeishuService) getTenantAccessToken() (string, error) {
	reqBody := map[string]string{
		"app_id":     s.AppID,
		"app_secret": s.AppSecret,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("请求体序列化失败: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		s.BaseURL+"/open-apis/auth/v3/tenant_access_token/internal",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result TenantAccessTokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("获取令牌失败: %s (code: %d)", result.Msg, result.Code)
	}

	return result.TenantAccessToken, nil
}

// GetSummaryFolderToken 获取或创建"自动总结"文件夹的token
func (s *FeishuService) GetSummaryFolderToken() (string, error) {
	// 如果已经缓存，直接返回
	if s.summaryFolderToken != "" {
		fmt.Println("✅ 使用缓存的自动总结文件夹 token")
		return s.summaryFolderToken, nil
	}

	fmt.Println("🔍 开始查找自动总结文件夹...")

	token, err := s.getTenantAccessToken()
	if err != nil {
		return "", fmt.Errorf("获取访问令牌失败: %w", err)
	}

	// 获取知识库节点列表（page_size 最大为 50）
	nodesURL := fmt.Sprintf("%s/open-apis/wiki/v2/spaces/%s/nodes?page_size=50", s.BaseURL, s.WikiID)
	req, err := http.NewRequest("GET", nodesURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("获取节点列表失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result GetNodesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("获取节点列表失败: %s (code: %d)", result.Msg, result.Code)
	}

	// 查找"自动总结"文件夹
	fmt.Printf("📋 查找到 %d 个节点\n", len(result.Data.Items))
	for _, node := range result.Data.Items {
		fmt.Printf("   - 标题: %s, Token: %s, 类型: %s\n", node.Title, node.NodeToken, node.ObjType)
		if node.Title == SummaryFolderName {
			fmt.Printf("✅ 找到自动总结文件夹: %s\n", node.NodeToken)
			s.summaryFolderToken = node.NodeToken
			return node.NodeToken, nil
		}
	}

	fmt.Println("⚠️  未找到自动总结文件夹，正在创建...")

	// 如果没找到，创建"自动总结"文件夹
	folderToken, err := s.createSummaryFolder(token)
	if err != nil {
		return "", fmt.Errorf("创建自动总结文件夹失败: %w", err)
	}

	fmt.Printf("✅ 自动总结文件夹创建成功: %s\n", folderToken)
	s.summaryFolderToken = folderToken
	return folderToken, nil
}

// createSummaryFolder 创建"自动总结"文件夹
// createSummaryFolder 创建"自动总结"文件夹
// 注意：飞书 Wiki API 不支持直接创建文件夹节点
// 建议：使用现有的节点作为父节点，或者手动在知识库中创建
func (s *FeishuService) createSummaryFolder(token string) (string, error) {
	// 飞书 Wiki API 创建节点时不支持普通文件夹类型
	// 我们需要使用一个已有的节点作为父节点
	// 或者直接在知识库根目录下创建页面

	// 这里使用一个简单的策略：直接在根目录下创建一个名为"自动总结"的文档
	// 用户可以在知识库中手动将其转换为文件夹或整理结构

	fmt.Println("📝 注意：飞书 Wiki API 不支持直接创建文件夹")
	fmt.Println("💡 将直接在知识库根目录下创建页面")

	// 根据飞书 API 文档，创建节点需要必需字段
	// node_type: origin (原始节点) 或 shortcut (快捷方式)
	// obj_type: docx (文档), bitable (多维表格) 等
	createReq := CreateDocumentRequest{
		Title:    SummaryFolderName,
		ObjType:  "docx",
		NodeType: "origin",
	}

	jsonData, err := json.Marshal(createReq)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	fmt.Printf("📝 创建页面的请求体: %s\n", string(jsonData))

	req, err := http.NewRequest(
		"POST",
		s.BaseURL+"/open-apis/wiki/v2/spaces/"+s.WikiID+"/nodes",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("创建页面请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	fmt.Printf("📦 响应内容: %s\n", string(body))

	var result NodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("创建页面失败: %s (code: %d)", result.Msg, result.Code)
	}

	return result.Data.Node.NodeToken, nil
}

func (s *FeishuService) CreateDocument(title, content string) (string, error) {
	fmt.Println("\n📝 开始创建飞书文档...")

	// 获取访问令牌
	token, err := s.getTenantAccessToken()
	if err != nil {
		return "", fmt.Errorf("获取访问令牌失败: %w", err)
	}

	// 获取或创建"自动总结"文件夹
	parentToken, err := s.GetSummaryFolderToken()
	if err != nil {
		return "", fmt.Errorf("获取自动总结文件夹失败: %w", err)
	}

	fmt.Printf("📁 使用父节点: %s\n", parentToken)

	// 先在 Wiki 空间中创建一个文档节点
	// 注意：在现有节点下创建文档需要正确的 API 调用
	createWikiReq := map[string]interface{}{
		"title":             title,
		"parent_node_token": parentToken,
		"obj_type":          "docx",
		"node_type":         "origin",
	}

	jsonWikiData, err := json.Marshal(createWikiReq)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	fmt.Printf("📝 创建 Wiki 节点请求体: %s\n", string(jsonWikiData))

	req, err := http.NewRequest(
		"POST",
		s.BaseURL+"/open-apis/wiki/v2/spaces/"+s.WikiID+"/nodes",
		bytes.NewReader(jsonWikiData),
	)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("创建 Wiki 节点请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("📦 创建 Wiki 节点响应: %s\n", string(body))

	var wikiResult NodeResponse
	if err := json.Unmarshal(body, &wikiResult); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if wikiResult.Code != 0 {
		return "", fmt.Errorf("创建 Wiki 节点失败: %s (code: %d)", wikiResult.Msg, wikiResult.Code)
	}

	wikiNodeToken := wikiResult.Data.Node.NodeToken
	objToken := wikiResult.Data.Node.ObjToken

	fmt.Printf("✅ Wiki 节点创建成功: %s\n", wikiNodeToken)
	fmt.Printf("✅ 文档对象 Token: %s\n", objToken)

	// 等待文档初始化并写入内容
	fmt.Println("⏳ 检查文档初始化状态...")
	maxRetries := 10
	retryDelay := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		if s.isDocumentReady(objToken, token) {
			fmt.Printf("✅ 文档已就绪 (尝试 %d/%d)\n", i+1, maxRetries)
			break
		}
		if i < maxRetries-1 {
			fmt.Printf("⏳ 文档未就绪，等待 %d 后重试...\n", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	// 写入文档内容使用 objToken（文档 ID），而不是 nodeToken（节点 token）
	if err := s.createDocumentContent(objToken, content, token); err != nil {
		return objToken, fmt.Errorf("文档内容写入失败: %w", err)
	}

	fmt.Println("✅ 文档内容已写入")
	return objToken, nil
}

func (s *FeishuService) createDocumentContent(documentID, content, token string) error {
	fmt.Printf("📝 开始写入文档内容，文档 ID: %s\n", documentID)
	fmt.Printf("📝 内容长度: %d 字符\n", len(content))

	// 如果内容为空，记录警告但不返回错误
	if content == "" {
		fmt.Println("⚠️  警告：内容为空，跳过写入")
		return nil
	}

	// 步骤 1: 获取文档的根 block_id
	rootBlockID, err := s.getRootBlockID(documentID, token)
	if err != nil {
		return fmt.Errorf("获取文档根 block_id 失败: %w", err)
	}

	fmt.Printf("✅ 获取到根 block_id: %s\n", rootBlockID)

	// 步骤 2: 在根块下创建子块
	// 根据飞书官方文档，使用 /blocks/:block_id/children 端点
	// block_type: 1=page, 2=text, 3=heading1, 4=heading2, 5=heading3 等
	// 使用 2 表示文本块
	createBlockReq := map[string]interface{}{
		"children": []map[string]interface{}{
			{
				"block_type": 2,
				"text": map[string]interface{}{
					"elements": []map[string]interface{}{
						{
							"text_run": map[string]interface{}{
								"content": content,
								"style":   map[string]interface{}{},
							},
						},
					},
				},
			},
		},
		"index": -1,
	}

	// 打印内容预览（只显示前 200 字符）
	contentPreview := content
	if len(contentPreview) > 200 {
		contentPreview = contentPreview[:200] + "..."
	}
	fmt.Printf("📝 内容预览: %s\n", contentPreview)

	jsonData, err := json.Marshal(createBlockReq)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	// 使用正确的 API 端点：/documents/:document_id/blocks/:block_id/children
	// 这里使用根 block_id 作为父块 ID
	fmt.Printf("🌐 调用飞书 API: POST /open-apis/docx/v1/documents/%s/blocks/%s/children\n", documentID, rootBlockID)

	req, err := http.NewRequest(
		"POST",
		s.BaseURL+"/open-apis/docx/v1/documents/"+documentID+"/blocks/"+rootBlockID+"/children",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("写入内容请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	fmt.Printf("📝 HTTP 状态码: %d\n", resp.StatusCode)
	fmt.Printf("📝 API 响应: %s\n", string(body))

	// 检查 HTTP 状态码（飞书 API 可能返回 200 或其他成功状态）
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API 返回错误: HTTP %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result CreateBlockResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %w, 响应内容: %s", err, string(body))
	}

	fmt.Printf("📝 响应 code: %d, msg: %s\n", result.Code, result.Msg)

	if result.Code != 0 {
		return fmt.Errorf("写入内容失败: %s (code: %d)", result.Msg, result.Code)
	}

	fmt.Printf("✅ 文档内容写入成功，block_id: %v\n\n", result.Data.BlockIDs)

	return nil
}

// getRootBlockID 获取文档的根 block_id
func (s *FeishuService) getRootBlockID(documentID, token string) (string, error) {
	req, err := http.NewRequest(
		"GET",
		s.BaseURL+"/open-apis/docx/v1/documents/"+documentID+"/blocks",
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("获取文档块信息失败: HTTP %d, 响应: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("📋 文档块信息响应: %s\n", string(body))

	var result struct {
		Code int `json:"code"`
		Data struct {
			Items []struct {
				BlockID string `json:"block_id"`
			} `json:"items"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("获取文档块失败: code %d", result.Code)
	}

	// 取第一个块作为根块
	if len(result.Data.Items) == 0 {
		return "", fmt.Errorf("文档中没有找到任何块")
	}

	return result.Data.Items[0].BlockID, nil
}

func (s *FeishuService) getDocumentID(objToken, token string) (string, error) {
	req, err := http.NewRequest(
		"GET",
		s.BaseURL+"/open-apis/docx/v1/documents/"+objToken,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("获取文档信息失败: HTTP %d, 响应: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	fmt.Printf("📋 文档信息响应: %s\n", string(body))

	var result struct {
		Code int `json:"code"`
		Data struct {
			Document struct {
				DocumentID string `json:"document_id"`
				Token      string `json:"token"` // Wiki节点可能有不同的token
			} `json:"document"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("获取文档失败: code %d", result.Code)
	}

	// 检查是否有token字段（对于Wiki文档）
	if result.Data.Document.Token != "" {
		fmt.Printf("✅ 使用 Wiki 文档 Token: %s\n", result.Data.Document.Token)
		return result.Data.Document.Token, nil
	}

	return result.Data.Document.DocumentID, nil
}

func (s *FeishuService) isDocumentReady(documentID, token string) bool {
	req, err := http.NewRequest(
		"GET",
		s.BaseURL+"/open-apis/docx/v1/documents/"+documentID,
		nil,
	)
	if err != nil {
		return false
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// HTTP 200 表示文档已初始化并可用
	return resp.StatusCode == 200
}

// FindOrCreateDocument 查找或创建方向文档
// 如果文档不存在则创建，存在则返回缓存的 token
func (s *FeishuService) FindOrCreateDocument(direction string) (string, string, error) {
	// 如果方向为空，使用默认方向
	if direction == "" {
		direction = SummaryFolderName
		fmt.Println("ℹ️  未指定方向，使用默认方向: " + direction)
	}

	fmt.Printf("🔍 查找方向文档: %s\n", direction)

	// 检查缓存
	if token, ok := s.directionTokens[direction]; ok {
		fmt.Printf("✅ 使用缓存的方向文档 token: %s\n", token)
		return token, direction, nil
	}

	// 获取访问令牌
	token, err := s.getTenantAccessToken()
	if err != nil {
		return "", "", fmt.Errorf("获取访问令牌失败: %w", err)
	}

	// 获取知识库节点列表（page_size 最大为 50）
	nodesURL := fmt.Sprintf("%s/open-apis/wiki/v2/spaces/%s/nodes?page_size=50", s.BaseURL, s.WikiID)
	req, err := http.NewRequest("GET", nodesURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("获取节点列表失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result GetNodesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 0 {
		return "", "", fmt.Errorf("获取节点列表失败: %s (code: %d)", result.Msg, result.Code)
	}

	// 查找方向文档
	fmt.Printf("📋 查找到 %d 个节点\n", len(result.Data.Items))
	for _, node := range result.Data.Items {
		if node.Title == direction {
			fmt.Printf("✅ 找到方向文档: %s, Token: %s\n", direction, node.NodeToken)
			s.directionTokens[direction] = node.NodeToken
			return node.NodeToken, direction, nil
		}
	}

	// 未找到，创建文档
	fmt.Printf("⚠️  未找到方向文档 %s，正在创建...\n", direction)
	documentToken, err := s.createDirectionDocument(token, direction)
	if err != nil {
		return "", "", fmt.Errorf("创建方向文档失败: %w", err)
	}

	fmt.Printf("✅ 方向文档创建成功: %s, Token: %s\n", direction, documentToken)
	s.directionTokens[direction] = documentToken
	return documentToken, direction, nil
}

// createDirectionDocument 创建方向文档
func (s *FeishuService) createDirectionDocument(token, directionName string) (string, error) {
	fmt.Printf("📝 创建方向文档: %s\n", directionName)

	createReq := CreateDocumentRequest{
		Title:    directionName,
		ObjType:  "docx",
		NodeType: "origin",
	}

	jsonData, err := json.Marshal(createReq)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	fmt.Printf("📝 创建文档的请求体: %s\n", string(jsonData))

	req, err := http.NewRequest(
		"POST",
		s.BaseURL+"/open-apis/wiki/v2/spaces/"+s.WikiID+"/nodes",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("创建文档请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	fmt.Printf("📦 响应内容: %s\n", string(body))

	var result NodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("创建文档失败: %s (code: %d)", result.Msg, result.Code)
	}

	return result.Data.Node.NodeToken, nil
}

// CreateDocumentInNode 在指定节点下创建文档
func (s *FeishuService) CreateDocumentInNode(parentToken, title, content string) (string, error) {
	fmt.Println("\n📝 开始创建飞书文档...")
	fmt.Printf("📁 父节点 Token: %s\n", parentToken)

	// 获取访问令牌
	token, err := s.getTenantAccessToken()
	if err != nil {
		return "", fmt.Errorf("获取访问令牌失败: %w", err)
	}

	// 在指定节点下创建文档
	createWikiReq := map[string]interface{}{
		"title":             title,
		"parent_node_token": parentToken,
		"obj_type":          "docx",
		"node_type":         "origin",
	}

	jsonWikiData, err := json.Marshal(createWikiReq)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	fmt.Printf("📝 创建 Wiki 节点请求体: %s\n", string(jsonWikiData))

	req, err := http.NewRequest(
		"POST",
		s.BaseURL+"/open-apis/wiki/v2/spaces/"+s.WikiID+"/nodes",
		bytes.NewReader(jsonWikiData),
	)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("创建 Wiki 节点请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("📦 创建 Wiki 节点响应: %s\n", string(body))

	var wikiResult NodeResponse
	if err := json.Unmarshal(body, &wikiResult); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if wikiResult.Code != 0 {
		return "", fmt.Errorf("创建 Wiki 节点失败: %s (code: %d)", wikiResult.Msg, wikiResult.Code)
	}

	wikiNodeToken := wikiResult.Data.Node.NodeToken
	objToken := wikiResult.Data.Node.ObjToken

	fmt.Printf("✅ Wiki 节点创建成功: %s\n", wikiNodeToken)
	fmt.Printf("✅ 文档对象 Token: %s\n", objToken)

	// 等待文档初始化并写入内容
	fmt.Println("⏳ 检查文档初始化状态...")
	maxRetries := 10
	retryDelay := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		if s.isDocumentReady(objToken, token) {
			fmt.Printf("✅ 文档已就绪 (尝试 %d/%d)\n", i+1, maxRetries)
			break
		}
		if i < maxRetries-1 {
			fmt.Printf("⏳ 文档未就绪，等待 %d 后重试...\n", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	// 写入文档内容使用 objToken（文档 ID），而不是 nodeToken（节点 token）
	if err := s.createDocumentContent(objToken, content, token); err != nil {
		return objToken, fmt.Errorf("文档内容写入失败: %w", err)
	}

	fmt.Println("✅ 文档内容已写入")
	return objToken, nil
}

// 简单的 UUID 生成
func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
