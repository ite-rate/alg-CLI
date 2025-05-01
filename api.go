package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// 调用API
func callLLM(prompt string, cfg *Config) (string, error) {
	// 构建请求体
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": cfg.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "你是一个算法专家助手。",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.3, // 低温度以获得更确定性的答案
	})

	if err != nil {
		return "", fmt.Errorf("构建请求失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", cfg.ApiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.ApiKey))

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API返回错误: %s, 状态码: %d", string(body), resp.StatusCode)
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 提取内容
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("无效的API响应结构")
	}

	choice := choices[0].(map[string]interface{})
	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("无效的消息结构")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("无效的内容格式")
	}

	return content, nil
}

// 获取LeetCode题目信息
func fetchProblemInfo(problemID int, cfg *Config) (map[string]interface{}, error) {
	prompt := generateProblemInfoPrompt(problemID)

	response, err := callLLM(prompt, cfg)
	if err != nil {
		return nil, err
	}

	// 尝试解析JSON响应
	var problemInfo map[string]interface{}
	err = json.Unmarshal([]byte(response), &problemInfo)
	if err != nil {
		// 如果JSON解析失败，尝试从文本中提取结构化信息
		problemInfo = extractStructuredInfo(response)
	}

	return problemInfo, nil
}

// 生成代码骨架
func generateCodeSkeleton(problemInfo map[string]interface{}, cfg *Config, category string) (string, error) {
	prompt := generateCodeSkeletonPrompt(problemInfo, cfg, category)

	response, err := callLLM(prompt, cfg)
	if err != nil {
		return "", err
	}

	// 清理代码 - 移除可能的Markdown格式
	cleanedCode := cleanCode(response)

	return cleanedCode, nil
}

// 清理代码 - 移除Markdown代码块标记等
func cleanCode(codeText string) string {
	// 移除可能的Markdown代码块格式
	code := codeText

	// 移除代码块开始标记 (如 ```go 或 ```python 等)
	codeStart := "```"
	if startIdx := bytes.Index([]byte(code), []byte(codeStart)); startIdx != -1 {
		// 找到第一行的换行符
		if nlIdx := bytes.IndexByte([]byte(code)[startIdx+3:], '\n'); nlIdx != -1 {
			code = code[startIdx+3+nlIdx+1:]
		}
	}

	// 移除代码块结束标记 (```)
	codeEnd := "```"
	if endIdx := bytes.LastIndex([]byte(code), []byte(codeEnd)); endIdx != -1 {
		code = code[:endIdx]
	}

	return code
}

// 直接生成代码骨架（一步到位）
func generateDirectCode(problemID int, cfg *Config, category string) (string, map[string]interface{}, error) {
	prompt := generateDirectCodeSkeletonPrompt(problemID, cfg, category)

	response, err := callLLM(prompt, cfg)
	if err != nil {
		return "", nil, err
	}

	// 清理代码 - 移除可能的Markdown格式
	cleanedCode := cleanCode(response)

	// 尝试从代码注释中提取题目名称
	title := extractTitleFromCode(cleanedCode, problemID)

	// 提取基本题目信息用于显示
	info := make(map[string]interface{})
	info["title"] = title
	info["difficulty"] = "未知" // 可以从代码中尝试提取
	info["description"] = "由大模型直接生成"

	return cleanedCode, info, nil
}

// 从代码中提取题目名称
func extractTitleFromCode(code string, defaultID int) string {
	// 默认标题，避免使用特殊字符
	defaultTitle := fmt.Sprintf("题目%d", defaultID)

	// 尝试从头部注释中提取题目名称
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		// 只检查前20行
		if i > 20 {
			break
		}

		line = strings.TrimSpace(line)
		// 查找包含"LeetCode"和题目名称的行
		if strings.Contains(line, "LeetCode") && strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				title := strings.TrimSpace(parts[1])
				// 确保标题不为空且不包含"LeetCode"
				if title != "" && !strings.Contains(title, "LeetCode") {
					return title
				}
			}
		} else if strings.Contains(line, "题目") && strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				title := strings.TrimSpace(parts[1])
				if title != "" {
					return title
				}
			}
		}
	}

	return defaultTitle
}

// 根据题目信息生成代码骨架的Prompt
func generateCodeSkeletonPrompt(problemInfo map[string]interface{}, cfg *Config, category string) string {
	title, _ := problemInfo["title"].(string)
	difficulty, _ := problemInfo["difficulty"].(string)
	description, _ := problemInfo["description"].(string)
	examples, _ := problemInfo["examples"].(string)
	constraints, _ := problemInfo["constraints"].(string)

	levelDescription := "完整代码框架但关键算法实现部分留空，添加TODO注释指导如何实现"
	if cfg.SkeletonLevel < 20 {
		levelDescription = "仅提供基本函数签名和简单注释"
	} else if cfg.SkeletonLevel > 70 {
		levelDescription = "几乎完整的解决方案，只有少量关键部分需要填写"
	}

	algorithmGuidance := ""
	if category != "" {
		algorithmGuidance = fmt.Sprintf(`请特别关注%s算法相关的解题思路，并在注释中提供这种方法的关键步骤。`, category)
	}

	return fmt.Sprintf(`你是一个算法专家，精通LeetCode题库。请为以下LeetCode题目创建一个%s语言的代码骨架：

题目标题: %s
难度: %s
题目描述: %s
示例: %s
约束条件: %s

要求:
1. 代码完整度为%d%%，这意味着%s
2. 添加注释解释算法思路和时间复杂度
3. 对需要学生实现的部分使用TODO注释清晰标记
4. 在注释中提供解题的关键步骤提示，但不给出完整实现
5. 提供至少两种可能的解法框架
6. 所有注释、题目描述和提示必须使用中文
7. 不要使用main函数，而是使用Go语言的测试函数格式 (func TestXxx(t *testing.T))
8. 添加至少2个测试用例，便于使用"go test"命令直接运行和调试
9. 文件应该是一个完整的可直接运行的测试文件，包含必要的import (如"testing"包)
10. 确保测试代码能够直接编译运行，不会有变量声明但未使用的错误
11. 提供比较函数确保测试数据可以正确验证，特别是对于需要忽略顺序的情况

%s

只返回代码，不需要其他解释。`, cfg.Language, title, difficulty, description, examples, constraints, cfg.SkeletonLevel, levelDescription, algorithmGuidance)
}
