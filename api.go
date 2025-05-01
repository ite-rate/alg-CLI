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
