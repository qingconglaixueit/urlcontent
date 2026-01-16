package services

import (
	"fmt"
	"testing"
)

// TestFeishuAPICall 测试飞书API调用流程
func TestFeishuAPICall(t *testing.T) {
	// 使用实际配置
	appID := "cli_a9d27bd8db78dbb4"
	appSecret := "swcvzxSrgtxMQsSr4YMyLfPdTnbbAibe"
	wikiID := "7102436789893267458"

	service := NewFeishuService(appID, appSecret, wikiID)

	// 测试内容
	testTitle := "测试文档标题-单元测试"
	testContent := "这是测试内容。第一段内容：用于验证飞书API写入功能正常。第二段内容：确保内容能够正确写入飞书知识库。"

	fmt.Println("=======================================")
	fmt.Println("开始测试飞书API写入功能")
	fmt.Println("=======================================")

	// 创建文档并写入内容
	documentID, err := service.CreateDocument(testTitle, testContent)
	if err != nil {
		t.Errorf("❌ 创建文档失败: %v", err)
		return
	}

	fmt.Println("=======================================")
	fmt.Printf("✅ 测试成功！文档ID: %s\n", documentID)
	fmt.Println("✅ 请在飞书知识库中验证文档内容是否正确写入")
	fmt.Println("=======================================")
}

// TestFeishuWriteEmptyContent 测试写入空内容
func TestFeishuWriteEmptyContent(t *testing.T) {
	appID := "cli_a9d27bd8db78dbb4"
	appSecret := "swcvzxSrgtxMQsSr4YMyLfPdTnbbAibe"
	wikiID := "7102436789893267458"

	service := NewFeishuService(appID, appSecret, wikiID)

	fmt.Println("=======================================")
	fmt.Println("开始测试写入空内容")
	fmt.Println("=======================================")

	// 创建文档并写入空内容
	documentID, err := service.CreateDocument("测试空内容", "")
	if err != nil {
		t.Errorf("❌ 创建文档失败: %v", err)
		return
	}

	fmt.Printf("✅ 空内容测试通过！文档ID: %s\n", documentID)
	fmt.Println("=======================================")
}

// TestWeChatArticleWrite 测试微信公众号文章写入
func TestWeChatArticleWrite(t *testing.T) {
	appID := "cli_a9d27bd8db78dbb4"
	appSecret := "swcvzxSrgtxMQsSr4YMyLfPdTnbbAibe"
	wikiID := "7102436789893267458"

	// 创建 Parser 服务来解析微信文章
	parser := NewParserService()

	fmt.Println("=======================================")
	fmt.Println("开始测试微信公众号文章写入")
	fmt.Println("=======================================")

	// 解析微信公众号文章
	wechatURL := "https://mp.weixin.qq.com/s/zCOiWZPAdNTsA5EzXGbWlA"
	fmt.Printf("🔗 解析URL: %s\n", wechatURL)

	extracted, err := parser.ParseURL(wechatURL)
	if err != nil {
		t.Errorf("❌ 解析微信公众号文章失败: %v", err)
		return
	}

	fmt.Printf("✅ 解析成功: %s\n", extracted.Title)
	fmt.Printf("📝 内容长度: %d 字符\n\n", len(extracted.Content))

	// 写入飞书
	feishuService := NewFeishuService(appID, appSecret, wikiID)
	documentID, err := feishuService.CreateDocument(extracted.Title, extracted.Content)
	if err != nil {
		t.Errorf("❌ 写入飞书失败: %v", err)
		return
	}

	fmt.Println("=======================================")
	fmt.Printf("✅ 微信公众号文章测试成功！\n")
	fmt.Printf("📄 标题: %s\n", extracted.Title)
	fmt.Printf("📄 文档ID: %s\n", documentID)
	fmt.Println("✅ 请在飞书知识库中验证文章内容是否正确写入")
	fmt.Println("=======================================")
}

// TestWeChatArticleWriteWithLongContent 测试写入长内容（微信公众号文章通常较长）
func TestWeChatArticleWriteWithLongContent(t *testing.T) {
	appID := "cli_a9d27bd8db78dbb4"
	appSecret := "swcvzxSrgtxMQsSr4YMyLfPdTnbbAibe"
	wikiID := "7102436789893267458"

	fmt.Println("=======================================")
	fmt.Println("开始测试长内容写入（模拟微信公众号文章）")
	fmt.Println("=======================================")

	// 创建较长的测试内容（模拟微信公众号文章）
	testTitle := "测试长内容写入-微信公众号文章模拟"
	testContent := `这是第一段内容：模拟微信公众号的文章格式和长度。微信公众号文章通常包含多个段落，每段可能达到200-300字。

这是第二段内容：在实际使用中，解析器会从HTML中提取有效内容，去除脚本、样式等无关标签。

这是第三段内容：我们需要确保系统能够处理较长的内容并将其正确写入飞书知识库。长内容写入是测试系统健壮性的重要环节。

这是第四段内容：飞书API应该能够接受并存储较大篇幅的文字内容，而不应该因为内容长度问题导致写入失败。

这是第五段内容：通过这个测试，我们可以验证系统对实际生产环境中微信公众号文章的处理能力。这是结束段落。`

	service := NewFeishuService(appID, appSecret, wikiID)

	documentID, err := service.CreateDocument(testTitle, testContent)
	if err != nil {
		t.Errorf("❌ 写入长内容失败: %v", err)
		return
	}

	fmt.Println("=======================================")
	fmt.Printf("✅ 长内容测试成功！\n")
	fmt.Printf("📝 内容长度: %d 字符\n", len(testContent))
	fmt.Printf("📄 文档ID: %s\n", documentID)
	fmt.Println("=======================================")
}
