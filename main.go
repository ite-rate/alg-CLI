package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 配置项
type Config struct {
	ApiKey        string
	ApiUrl        string
	Model         string
	Language      string
	SkeletonLevel int // 0-100，代码骨架的完整度
	OutputDir     string
}

// 支持的编程语言
var supportedLanguages = []string{
	"go", "python", "java", "cpp", "javascript", "typescript",
	"rust", "c", "csharp", "php", "ruby", "swift", "kotlin",
}

// 算法分类
var algoCategories = []string{
	"数组", "链表", "栈", "队列", "哈希表",
	"字符串", "二分查找", "排序", "贪心", "动态规划",
	"深度优先搜索", "广度优先搜索", "回溯", "树", "图",
	"数学", "位操作", "并查集", "前缀和", "滑动窗口",
}

func main() {
	// 解析命令行参数
	problemID := flag.Int("id", 0, "LeetCode题目ID")
	language := flag.String("lang", "go", "编程语言")
	level := flag.Int("level", 30, "代码骨架完整度(0-100)")
	category := flag.String("category", "", "算法分类")
	model := flag.String("model", "deepseek-v3-250324", "LLM模型 (可选: deepseek-v3-250324, deepseek-r1-250120)")
	output := flag.String("output", "leetcode_practice", "输出目录")

	flag.Parse()

	// 验证参数
	if *problemID <= 0 {
		fmt.Println("❌ 错误: 必须提供有效的LeetCode题目ID")
		flag.Usage()
		os.Exit(1)
	}

	if *level < 0 || *level > 100 {
		fmt.Println("❌ 错误: 代码骨架完整度必须在0-100之间")
		flag.Usage()
		os.Exit(1)
	}

	// 验证语言
	langValid := false
	for _, lang := range supportedLanguages {
		if lang == *language {
			langValid = true
			break
		}
	}
	if !langValid {
		fmt.Printf("❌ 错误: 不支持的语言: %s\n", *language)
		fmt.Printf("支持的语言: %s\n", strings.Join(supportedLanguages, ", "))
		os.Exit(1)
	}

	// 加载配置
	cfg := &Config{
		ApiKey:        os.Getenv("ARK_API_KEY"),
		ApiUrl:        "https://ark.cn-beijing.volces.com/api/v3/chat/completions",
		Model:         *model,
		Language:      *language,
		SkeletonLevel: *level,
		OutputDir:     *output,
	}

	// 验证API密钥
	if cfg.ApiKey == "" {
		fmt.Println("❌ 错误: 未设置ARK_API_KEY环境变量")
		fmt.Println("请设置环境变量: export ARK_API_KEY='您的API密钥'")
		os.Exit(1)
	}

	// 确保输出目录存在
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		fmt.Printf("❌ 错误: 无法创建输出目录: %v\n", err)
		os.Exit(1)
	}

	// 生成代码文件
	err := generateCodeFile(*problemID, cfg, *category)
	if err != nil {
		fmt.Printf("❌ 错误: %s\n", err)
		os.Exit(1)
	}
}

// 生成代码文件
func generateCodeFile(problemID int, cfg *Config, category string) error {
	fmt.Printf("🔍 处理题目 #%d...\n", problemID)

	// 使用一步到位的方法
	codeContent, problemInfo, err := generateDirectCode(problemID, cfg, category)
	if err != nil {
		// 如果一步到位方法失败，回退到两步法
		fmt.Println("⚠️ 一步到位生成失败，尝试分步获取...")

		// 获取题目信息
		problemInfo, err = fetchProblemInfo(problemID, cfg)
		if err != nil {
			return fmt.Errorf("获取题目信息失败: %v", err)
		}

		fmt.Printf("✅ 题目信息获取成功: %s\n", problemInfo["title"])
		fmt.Printf("🧩 生成代码骨架 (完整度: %d%%)...\n", cfg.SkeletonLevel)

		// 生成代码骨架
		codeContent, err = generateCodeSkeleton(problemInfo, cfg, category)
		if err != nil {
			return fmt.Errorf("生成代码骨架失败: %v", err)
		}
	} else {
		fmt.Println("✅ 代码直接生成成功!")
	}

	// 创建文件名
	title, ok := problemInfo["title"].(string)
	if !ok {
		title = fmt.Sprintf("题目%d", problemID)
	}

	// 使用符合Go测试规范的文件名
	filename := fmt.Sprintf("%s_test.%s",
		fmt.Sprintf("leetcode_%d", problemID),
		getFileExtension(cfg.Language))

	filePath := filepath.Join(cfg.OutputDir, filename)

	// 添加题目信息注释
	difficulty, ok := problemInfo["difficulty"].(string)
	if !ok {
		difficulty = "未知"
	}

	description, ok := problemInfo["description"].(string)
	if !ok {
		description = "由大模型直接生成"
	}

	header := fmt.Sprintf(`/*
 * LeetCode #%d: %s
 * 难度: %s
 * 
 * 题目描述:
 * %s
 * 
 * 代码骨架完整度: %d%%
 */

`, problemID, title, difficulty, description, cfg.SkeletonLevel)

	// 写入文件
	err = os.WriteFile(filePath, []byte(header+codeContent), 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	fmt.Printf("\n✅ 成功生成练习文件: %s\n", filePath)
	fmt.Printf("🔍 题目难度: %s\n", difficulty)
	if category != "" {
		fmt.Printf("📚 算法分类: %s\n", category)
	}
	fmt.Printf("⚙️ 代码完整度: %d%%\n", cfg.SkeletonLevel)
	fmt.Printf("🧪 调试命令: cd %s && go test -v %s\n", cfg.OutputDir, filename)

	return nil
}

// 文件名清理
func sanitizeFilename(name string) string {
	// 移除非法字符
	name = strings.Map(func(r rune) rune {
		if r < 32 || strings.ContainsRune(`<>:"/\|?*`, r) {
			return '_'
		}
		return r
	}, name)

	// 截断过长的文件名
	if len(name) > 50 {
		name = name[:50]
	}

	return name
}

// 获取文件扩展名
func getFileExtension(language string) string {
	switch language {
	case "python":
		return "py"
	case "java":
		return "java"
	case "cpp":
		return "cpp"
	case "javascript":
		return "js"
	case "typescript":
		return "ts"
	case "rust":
		return "rs"
	case "c":
		return "c"
	case "csharp":
		return "cs"
	case "php":
		return "php"
	case "ruby":
		return "rb"
	case "swift":
		return "swift"
	case "kotlin":
		return "kt"
	default:
		return "go"
	}
}
