package main

import (
	"fmt"
)

// 生成获取LeetCode题目信息的Prompt
func generateProblemInfoPrompt(problemID int) string {
	return fmt.Sprintf(`请提供LeetCode题目#%d的详细信息，包括题目标题、难度级别、完整描述、示例和约束条件。

请以JSON格式返回，格式如下：
{
  "title": "题目标题",
  "difficulty": "简单/中等/困难",
  "description": "题目的详细描述",
  "examples": "所有给出的示例",
  "constraints": "题目的所有约束条件",
  "tags": ["相关标签", "算法分类"]
}

请确保返回的是有效的JSON格式，不要添加额外的说明或解释。`, problemID)
}

// 一步到位生成代码骨架的Prompt
func generateDirectCodeSkeletonPrompt(problemID int, cfg *Config, category string) string {
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

	return fmt.Sprintf(`你是一个算法专家，精通LeetCode题库。请直接为LeetCode题目#%d创建一个%s语言的代码骨架。

要求:
1. 首先简要概述题目（包括题目名称、难度和题目描述）
2. 代码完整度为%d%%，这意味着%s
3. 添加注释解释算法思路和时间复杂度
4. 对需要学生实现的部分使用TODO注释清晰标记
5. 在注释中提供解题的关键步骤提示，但不给出完整实现
6. 提供至少两种可能的解法框架
7. 所有注释、题目描述和提示必须使用中文
8. 不要使用main函数，而是使用Go语言的测试函数格式 (func TestXxx(t *testing.T))
9. 添加至少2个测试用例，便于使用"go test"命令直接运行和调试
10. 文件应该是一个完整的可直接运行的测试文件，包含必要的import (如"testing"包)
11. 确保测试代码能够直接编译运行，不会有变量声明但未使用的错误
12. 提供比较函数确保测试数据可以正确验证，特别是对于需要忽略顺序的情况

%s

只返回代码，不需要其他解释。`, problemID, cfg.Language, cfg.SkeletonLevel, levelDescription, algorithmGuidance)
}

// 提取结构化信息
func extractStructuredInfo(response string) map[string]interface{} {
	// 这个函数用于从非JSON响应中提取信息
	// 简单实现，实际使用时可能需要更复杂的解析逻辑
	result := make(map[string]interface{})

	// 默认值
	result["title"] = "未知题目"
	result["difficulty"] = "中等"
	result["description"] = "题目描述未找到"
	result["examples"] = "示例未找到"
	result["constraints"] = "约束条件未找到"
	result["tags"] = []string{}

	// 这里可以使用正则表达式或其他方法从文本中提取信息
	// 但为了简化，我们这里不做实际提取

	return result
}
