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

	languageReq := getLanguageSpecificRequirements(cfg.Language)

	return fmt.Sprintf(`你是一个算法专家，精通LeetCode题库。请直接为LeetCode题目#%d创建一个%s语言的代码骨架。

要求:
1. 首先简要概述题目（包括题目名称、难度和题目描述）
2. 代码完整度为%d%%，这意味着%s
3. 添加注释解释算法思路和时间复杂度
4. 对需要学生实现的部分使用TODO注释清晰标记
5. 在注释中提供解题的关键步骤提示，但不给出完整实现
6. 提供至少两种可能的解法框架
7. 所有注释、题目描述和提示必须使用中文
%s

%s

只返回代码，不需要其他解释。`, problemID, cfg.Language, cfg.SkeletonLevel, levelDescription, languageReq, algorithmGuidance)
}

// 根据语言返回单文件可运行/可测试的要求，避免不同语言约束互相污染
func getLanguageSpecificRequirements(language string) string {
	switch language {
	case "go":
		return `8. 不要使用main函数，而是使用Go语言的测试函数格式 (func TestXxx(t *testing.T))
9. 添加至少2个测试用例，便于使用"go test"命令直接运行和调试
10. 文件应该是一个完整的可直接运行的测试文件，包含必要的import (如"testing"包)
11. 确保测试代码能够直接编译运行，不会有变量声明但未使用的错误
12. 提供比较函数确保测试数据可以正确验证，特别是对于需要忽略顺序的情况`
	case "c":
		return `8. 输出必须是标准C代码（C11），不得出现任何Go/Python/Java等其他语言语法
9. 允许使用main函数，并在main中编写至少2个自测用例（可用assert或手写检查+打印）
10. 代码应可直接用gcc/clang编译运行（无未声明函数、无语法错误）
11. 为核心解法提供清晰的函数签名与必要的数据结构定义（按题目需要）
12. 代码风格保持一致、可读（命名清晰、边界检查明确，可参考Google C风格的精神）`
	case "cpp":
		return `8. 输出必须是标准C++17代码，不得出现任何Go/Python/Java等其他语言语法
9. 允许使用main函数，并在main中编写至少2个自测用例（建议使用assert/打印对比）
10. 代码应可直接用g++/clang++编译运行
11. 为核心解法提供清晰的函数签名与必要的数据结构定义（按题目需要）
12. 代码风格保持一致、可读（可参考Google C++风格的精神）`
	case "python":
		return `8. 输出必须是Python 3代码，不得出现任何Go/Java/C等其他语言语法
9. 使用标准库unittest编写至少2个测试用例，确保可直接运行
10. 文件结构为：核心函数/类 + unittest.TestCase + if __name__ == "__main__": unittest.main()
11. 不使用任何第三方依赖`
	case "java":
		return `8. 输出必须是Java代码，不得出现任何Go/Python/C等其他语言语法
9. 不引入第三方依赖：使用public static void main(String[] args)编写至少2个自测用例（断言/打印对比）
10. 代码应可直接用javac编译、java运行
11. 类/方法命名清晰、结构清楚（可参考Google Java风格的精神）`
	case "javascript":
		return `8. 输出必须是JavaScript（Node.js）代码，不得出现其他语言语法
9. 使用Node.js标准能力编写至少2个自测用例（可用console.assert或手写比较函数+打印）
10. 代码应可直接用node运行，不使用第三方依赖`
	case "typescript":
		return `8. 输出必须是TypeScript代码，不得出现其他语言语法
9. 编写至少2个自测用例（可用console.assert或手写比较函数+打印）
10. 不使用第三方依赖；代码应可在存在tsc的环境下通过命令行编译后运行（例如：tsc file.ts && node file.js）`
	case "rust":
		return `8. 输出必须是Rust代码，不得出现其他语言语法
9. 允许使用fn main()，并在main中编写至少2个自测用例（assert!/打印对比）
10. 代码应可直接用rustc编译运行（不依赖Cargo项目结构，不使用第三方crate）`
	case "csharp":
		return `8. 输出必须是C#代码，不得出现其他语言语法
9. 使用Main方法编写至少2个自测用例（断言/打印对比）
10. 代码尽量保持单文件可编译运行（不引入第三方依赖）`
	case "php":
		return `8. 输出必须是PHP代码，不得出现其他语言语法
9. 编写至少2个自测用例（手写比较函数+打印/断言）
10. 代码应可直接用php运行，不使用第三方依赖`
	case "ruby":
		return `8. 输出必须是Ruby代码，不得出现其他语言语法
9. 编写至少2个自测用例（可用raise/打印对比）
10. 代码应可直接用ruby运行，不使用第三方依赖`
	case "swift":
		return `8. 输出必须是Swift代码，不得出现其他语言语法
9. 编写至少2个自测用例（assert/打印对比）
10. 代码应可直接用swiftc编译运行（单文件），不使用第三方依赖`
	case "kotlin":
		return `8. 输出必须是Kotlin代码，不得出现其他语言语法
9. 使用fun main()编写至少2个自测用例（check/打印对比）
10. 代码应可直接用kotlinc编译后运行（单文件），不使用第三方依赖`
	default:
		return fmt.Sprintf(`8. 输出必须只包含%s语言代码，不得混入其他语言语法
9. 编写至少2个可运行的自测用例（使用该语言的标准方式），避免第三方依赖
10. 保证代码可直接编译/运行`, language)
	}
}

// 提取结构化信息
func extractStructuredInfo(_ string) map[string]interface{} {
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
