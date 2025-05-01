
# alg 使用说明手册

## 工具简介

alg是一个命令行工具，用于生成LeetCode算法题目的代码骨架，帮助您更高效地练习算法题。工具特点：

- 从LeetCode获取题目并生成可直接运行的测试文件
- 支持调整代码完整度（10%-80%），适合不同学习阶段
- 提供多种算法思路注释和提示
- 生成的是测试文件，可直接运行和调试
- 支持多种编程语言和算法分类

## 安装方法

### 方法一：放置到用户bin目录（推荐）

```bash
# 克隆仓库
git clone https://github.com/yourusername/alg.git
cd alg

# 构建
go build -o alg

# 安装到用户bin目录
mkdir -p ~/bin
cp alg ~/bin/

# 确保~/bin在PATH中
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc  # 或者 source ~/.zshrc（如果使用zsh）
```

安装后，您可以在任何目录下直接使用`alg`命令。

### 方法二：本地使用

```bash
# 克隆仓库
git clone https://github.com/yourusername/alg.git
cd alg

# 构建
go build -o alg

# 本地运行
./alg -id 1
```

## 使用前准备

在使用前，请设置API密钥环境变量：

```bash
export ARK_API_KEY='您的API密钥'
```

您可以将此行添加到`.bashrc`或`.zshrc`文件中以永久保存。

## 基本用法

```bash
# 生成题目#1的代码骨架（默认Go语言）
alg -id 1

# 指定编程语言
alg -id 15 -lang python

# 调整代码完整度(0-100)
alg -id 15 -level 10  # 最小骨架
alg -id 15 -level 50  # 中等骨架
alg -id 15 -level 80  # 几乎完整的解决方案

# 指定算法分类
alg -id 33 -category "二分查找"

# 调试生成的代码（以Go为例）
cd leetcode_practice
go test -v leetcode_15_test.go           # 测试所有解法
go test -v -run=TestThreeSum leetcode_15_test.go  # 只测试特定解法
```

## 参数说明

```
  -id int
    	LeetCode题目ID（必需）
  -lang string
    	编程语言 (默认 "go")
  -level int
    	代码骨架完整度(0-100) (默认 30)
  -category string
    	算法分类，如"二分查找"、"动态规划"等
  -model string
    	LLM模型 (默认 "deepseek-v3-250324")
  -output string
    	输出目录 (默认 "leetcode_practice")
```

## 支持的编程语言

go, python, java, cpp, javascript, typescript, rust, c, csharp, php, ruby, swift, kotlin

## 支持的算法分类

数组, 链表, 栈, 队列, 哈希表, 字符串, 二分查找, 排序, 贪心, 动态规划, 深度优先搜索, 广度优先搜索, 回溯, 树, 图, 数学, 位操作, 并查集, 前缀和, 滑动窗口

## 使用技巧

1. 不同的完整度适合不同学习阶段：
   - 初学者可以使用较高完整度(60-80)来学习算法思路
   - 熟练者可以选择中等完整度(30-50)练习核心逻辑
   - 挑战自己可以选择低完整度(10-20)几乎从零开始编写

2. 代码调试技巧：
   - 使用`-run`参数测试特定方法: `go test -v -run=TestThreeSum leetcode_15_test.go`
   - 测试特定用例: `go test -v -run=TestThreeSum/example1 leetcode_15_test.go`
