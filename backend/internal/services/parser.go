package services

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ParserService struct {
	httpClient *http.Client
}

func NewParserService() *ParserService {
	return &ParserService{
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // 增加到 60 秒，适应加载较慢的页面
		},
	}
}

// ExtractedContent 提取的内容
type ExtractedContent struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// ProxyResponse All Origins API 响应
type ProxyResponse struct {
	Contents string `json:"contents"`
	Type     string `json:"type"`
}

func (s *ParserService) ParseURL(targetURL string) (*ExtractedContent, error) {
	fmt.Println("═════════════════════════════════════════════════════════════")
	fmt.Println("🔍 开始解析 URL")
	fmt.Println("═════════════════════════════════════════════════════════════")
	fmt.Printf("📡 目标 URL: %s\n", targetURL)
	startTime := time.Now()

	// 验证 URL
	if _, err := url.Parse(targetURL); err != nil {
		return nil, fmt.Errorf("URL 格式错误: %w", err)
	}

	// 直接请求目标 URL 获取内容
	fmt.Println("\n📡 步骤 1: 直接获取网页内容")
	fmt.Println("─────────────────────────────────────────────────────────────")

	// 创建自定义请求，添加浏览器头部
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置浏览器 User-Agent 和其他头部，避免被识别为爬虫
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	fetchStartTime := time.Now()
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	fetchTime := time.Since(fetchStartTime)
	fmt.Printf("✅ 响应状态: %d %s\n", resp.StatusCode, resp.Status)
	fmt.Printf("⏱️  请求耗时: %dms\n", fetchTime.Milliseconds())

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP 错误: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	htmlContent := string(body)
	fmt.Printf("📦 HTML 内容长度: %d 字符\n", len(htmlContent))

	if len(htmlContent) == 0 {
		return nil, fmt.Errorf("响应中没有内容")
	}

	fmt.Println("📝 步骤 2: 解析 HTML 内容")
	fmt.Println("─────────────────────────────────────────────────────────────")

	// 统计 HTML 元素
	contentLower := strings.ToLower(htmlContent)
	paragraphsCount := strings.Count(contentLower, "<p")
	headingsCount := strings.Count(contentLower, "<h") + strings.Count(contentLower, "<H")
	listsCount := strings.Count(contentLower, "<li") + strings.Count(contentLower, "<LI")
	scriptsCount := strings.Count(contentLower, "<script")
	stylesCount := strings.Count(contentLower, "<style")

	fmt.Println("📊 HTML 元素统计:")
	fmt.Printf("   段落 (p): %d\n", paragraphsCount)
	fmt.Printf("   标题 (h1-h6): %d\n", headingsCount)
	fmt.Printf("   列表项 (li): %d\n", listsCount)
	fmt.Printf("   脚本 (script): %d\n", scriptsCount)
	fmt.Printf("   样式 (style): %d\n", stylesCount)

	// 简单的 HTML 解析（提取标题、描述、段落）
	title := extractTitle(htmlContent)
	fmt.Printf("\n📌 提取的标题: %s\n", title)

	metaDesc := extractMetaDescription(htmlContent)
	fmt.Printf("📌 提取的描述: %s...\n", truncateString(metaDesc, 100))
	fmt.Printf("   描述长度: %d 字符\n", len(metaDesc))

	// 清理 HTML 标签
	fmt.Println("\n🧹 步骤 3: 清理无关标签")
	fmt.Println("─────────────────────────────────────────────────────────────")
	cleanedContent := removeTags(htmlContent, []string{"script", "style", "iframe", "noscript"})
	fmt.Printf("移除后内容长度: %d 字符\n", len(cleanedContent))

	// 提取有效文本
	fmt.Println("🔍 步骤 4: 提取有效文本内容")
	fmt.Println("─────────────────────────────────────────────────────────────")

	paragraphs := extractParagraphs(cleanedContent)
	validTexts := filterValidTexts(paragraphs)

	fmt.Printf("找到文本元素总数: %d\n", len(paragraphs))
	fmt.Printf("✅ 有效文本段落数: %d\n", len(validTexts))
	fmt.Printf("✅ 提取的内容长度: %d 字符\n", len(strings.Join(validTexts, "\n\n")))

	// 打印预览
	if len(validTexts) > 0 {
		fmt.Println("📋 前3个有效段落预览:")
		for i, txt := range validTexts {
			if i >= 3 {
				break
			}
			fmt.Printf("   %d. %s%s\n", i+1, truncateString(txt, 80), cond(len(txt) > 80, "...", ""))
		}
	}

	// 构建最终内容并生成总结
	contentText := strings.Join(validTexts, "\n\n")
	if contentText == "" {
		contentText = extractPlainText(htmlContent)
	}

	// 生成内容总结
	fmt.Println("📝 步骤 5: 生成内容总结")
	fmt.Println("─────────────────────────────────────────────────────────────")
	summary := summarizeContent(title, metaDesc, contentText)
	fmt.Printf("✅ 总结生成完成, 长度: %d 字符\n", len(summary))

	finalContent := fmt.Sprintf("标题：%s\n\n描述：%s\n\n来源链接：%s\n\n内容总结：\n%s",
		title,
		cond(metaDesc != "", metaDesc, "无描述"),
		targetURL,
		summary)

	totalTime := time.Since(startTime)
	fmt.Println("\n═════════════════════════════════════════════════════════════")
	fmt.Println("✅ URL 解析完成")
	fmt.Println("═════════════════════════════════════════════════════════════")
	fmt.Printf("⏱️  总耗时: %dms\n", totalTime.Milliseconds())
	fmt.Printf("📊 提取统计:\n")
	fmt.Printf("   - 标题: %s...\n", truncateString(title, 50))
	fmt.Printf("   - 描述: %s...\n", truncateString(metaDesc, 50))
	fmt.Printf("   - 内容: %s... (共 %d 字符)\n", truncateString(contentText, 50), len(contentText))
	fmt.Println("═════════════════════════════════════════════════════════════")

	return &ExtractedContent{
		Title:     title,
		URL:       targetURL,
		Content:   finalContent,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// 辅助函数

func extractTitle(html string) string {
	// 提取 <title> 标签内容
	// 支持 <title> 和 <title anyattr="value"> 格式
	lowerHtml := strings.ToLower(html)

	// 查找 <title 开始位置
	start := strings.Index(lowerHtml, "<title")
	if start == -1 {
		return "无标题"
	}

	// 查找 > 的位置（标签结束）
	tagEnd := strings.Index(html[start:], ">")
	if tagEnd == -1 {
		return "无标题"
	}
	start += tagEnd + 1

	// 查找 </title> 的位置
	end := strings.Index(lowerHtml[start:], "</title>")
	if end == -1 {
		return "无标题"
	}

	title := html[start : start+end]
	// 移除换行和多余空格
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\r", "")
	title = strings.ReplaceAll(title, "\t", " ")
	title = strings.TrimSpace(title)

	// 如果标题为空或太短，尝试从 meta og:title 获取
	if len(title) == 0 || len(title) < 2 {
		ogTitle := extractMetaProperty(html, "og:title")
		if ogTitle != "" {
			return ogTitle
		}
	}

	return title
}

func extractMetaDescription(html string) string {
	// 提取 meta description

	if idx := strings.Index(strings.ToLower(html), "name=\"description\""); idx != -1 {
		start := strings.Index(html[idx:], "content=\"")
		if start != -1 {
			start += idx + 9
			end := strings.Index(html[start:], "\"")
			if end != -1 {
				return html[start : start+end]
			}
		}
	}

	if idx := strings.Index(strings.ToLower(html), "name='description'"); idx != -1 {
		start := strings.Index(html[idx:], "content='")
		if start != -1 {
			start += idx + 9
			end := strings.Index(html[start:], "'")
			if end != -1 {
				return html[start : start+end]
			}
		}
	}

	return ""
}

func extractMetaProperty(html string, property string) string {
	// 提取 meta property，如 og:title
	searchStr := fmt.Sprintf(`property="%s"`, property)
	idx := strings.Index(html, searchStr)
	if idx == -1 {
		searchStr = fmt.Sprintf(`property='%s'`, property)
		idx = strings.Index(html, searchStr)
	}
	if idx == -1 {
		return ""
	}

	contentStart := strings.Index(html[idx:], "content=\"")
	if contentStart == -1 {
		contentStart = strings.Index(html[idx:], "content='")
	}
	if contentStart == -1 {
		return ""
	}
	contentStart += idx + 9

	end := strings.Index(html[contentStart:], `"`)
	if end == -1 {
		return ""
	}

	return strings.TrimSpace(html[contentStart : contentStart+end])
}

func removeTags(html string, tagsToRemove []string) string {
	content := html
	for _, tag := range tagsToRemove {
		// 正确移除整个标签（包括内容）
		lowerContent := strings.ToLower(content)
		maxIterations := 10000 // 防止死循环
		iteration := 0
		for {
			iteration++
			if iteration > maxIterations {
				fmt.Printf("⚠️  警告: 移除标签 <%s> 达到最大迭代次数，可能存在 HTML 结构问题\n", tag)
				break
			}
			
			startTag := "<" + tag
			endTag := "</" + tag + ">"

			startIdx := strings.Index(lowerContent, startTag)
			if startIdx == -1 {
				break
			}

			// 查找结束标签的位置
			endIdx := strings.Index(lowerContent[startIdx:], endTag)
			if endIdx == -1 {
				// 没有找到结束标签，只移除开始标签
				tagEnd := strings.Index(content[startIdx:], ">")
				if tagEnd == -1 {
					break
				}
				content = content[:startIdx] + content[startIdx+tagEnd+1:]
				lowerContent = strings.ToLower(content)
				continue
			}

			// 移除从开始标签到结束标签之间的所有内容
			content = content[:startIdx] + content[startIdx+endIdx+len(endTag):]
			lowerContent = strings.ToLower(content)
		}
	}
	return content
}

func extractParagraphs(html string) []string {
	var paragraphs []string

	// 提取 <p> 标签内容
	content := strings.ToLower(html)
	idx := 0

	for {
		start := strings.Index(content[idx:], "<p")
		if start == -1 {
			break
		}
		start += idx

		// 找到 > 的位置
		endTag := strings.Index(content[start:], ">")
		if endTag == -1 {
			break
		}
		start += endTag + 1

		// 找到 </p> 的位置
		end := strings.Index(content[start:], "</p>")
		if end == -1 {
			break
		}

		text := html[start : start+end]
		paragraphs = append(paragraphs, strings.TrimSpace(text))

		idx = start + end + 4
	}

	// 如果没有 p 标签，尝试其他标签
	if len(paragraphs) == 0 {
		// 提取 h1-h6, li
		for _, tag := range []string{"h1", "h2", "h3", "h4", "h5", "h6", "li"} {
			tagCount := strings.Count(strings.ToLower(html), "<"+tag)
			for i := 0; i < tagCount; i++ {
				tempContent := html
				for j := 0; j <= i; j++ {
					start := strings.Index(strings.ToLower(tempContent), "<"+tag)
					if start == -1 {
						break
					}
					endTag := strings.Index(tempContent[start:], ">")
					if endTag == -1 {
						break
					}
					start += endTag + 1

					end := strings.Index(tempContent[start:], "</"+tag+">")
					if end == -1 {
						break
					}

					text := tempContent[start : start+end]
					paragraphs = append(paragraphs, strings.TrimSpace(text))

					tempContent = tempContent[start+end+len(tag)+3:]
				}
			}
		}
	}

	return paragraphs
}

func filterValidTexts(texts []string) []string {
	var valid []string
	for _, text := range texts {
		// 移除 HTML 标签（简单实现）
		text = removeHTMLTags(text)
		text = strings.TrimSpace(text)
		if len(text) > 10 { // 至少10个字符
			valid = append(valid, text)
		}
	}
	return valid
}

// removeHTMLTags 移除文本中的 HTML 标签
func removeHTMLTags(text string) string {
	// 移除 <a href="...">text</a> 这样的标签
	for {
		start := strings.Index(text, "<")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], ">")
		if end == -1 {
			break
		}
		text = text[:start] + text[start+end+1:]
	}

	// 清理多余的空格
	text = strings.ReplaceAll(text, "  ", " ")
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return strings.TrimSpace(text)
}

func extractPlainText(html string) string {
	// 简单的文本提取
	content := html
	// 移除所有 HTML 标签
	for {
		start := strings.Index(content, "<")
		if start == -1 {
			break
		}
		end := strings.Index(content, ">")
		if end == -1 {
			break
		}
		content = content[:start] + content[end+1:]
	}

	// 清理空白
	lines := strings.Split(content, "\n")
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Count(line, "") > 10 {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

func cond(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

// summarizeContent 生成内容总结
func summarizeContent(title, metaDesc, content string) string {
	var summaryParts []string

	// 总是包含标题
	titleLine := fmt.Sprintf("【%s】", title)
	summaryParts = append(summaryParts, titleLine)

	// 如果有描述，优先使用描述
	if metaDesc != "" {
		summaryParts = append(summaryParts, metaDesc)
	}

	// 提取关键段落
	lines := strings.Split(content, "\n")
	var keyParagraphs []string

	// 查找包含重要关键词的段落
	keywords := []string{"重要", "关键", "注意", "总结", "结论", "因此", "所以", "首先", "其次", "最后"}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 10 {
			continue
		}

		// 检查是否包含关键词
		containsKeyword := false
		lineLower := strings.ToLower(line)
		for _, keyword := range keywords {
			if strings.Contains(lineLower, strings.ToLower(keyword)) {
				containsKeyword = true
				break
			}
		}

		// 限制总结长度
		if containsKeyword || len(keyParagraphs) < 3 {
			keyParagraphs = append(keyParagraphs, line)
			if len(keyParagraphs) >= 5 {
				break
			}
		}
	}

	// 添加关键段落
	if len(keyParagraphs) > 0 {
		summaryParts = append(summaryParts, "\n【关键要点】")
		for i, para := range keyParagraphs {
			if i >= 3 {
				break
			}
			summaryParts = append(summaryParts, fmt.Sprintf("• %s", truncateString(para, 150)))
		}
	}

	// 添加原文链接提示
	summaryParts = append(summaryParts, "\n（此内容为自动生成的总结，详细信息请查看原文）")

	return strings.Join(summaryParts, "\n")
}
