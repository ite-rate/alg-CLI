/*
 * LeetCode #15: 题目15
 * 难度: 未知
 *
 * 题目描述:
 * 由大模型直接生成
 *
 * 代码骨架完整度: 30%
 */

package leetcode

import (
	"sort"
	"strconv"
	"testing"
)

/*
题目：#15 三数之和
给你一个包含n个整数的数组nums，判断nums中是否存在三个元素a，b，c，使得a + b + c = 0？
请你找出所有和为0且不重复的三元组。

注意：答案中不可以包含重复的三元组。

示例1：
输入：nums = [-1,0,1,2,-1,-4]
输出：[[-1,-1,2],[-1,0,1]]

示例2：
输入：nums = []
输出：[]

解题思路：
解法1：排序+双指针法
1. 先对数组进行排序
2. 固定一个数nums[i]，然后用双指针在剩余数组中寻找两数之和等于-nums[i]
3. 注意跳过重复元素以避免重复解
时间复杂度：O(n^2)，其中排序O(nlogn)，双指针遍历O(n^2)

解法2：哈希表法（略复杂，需要处理重复问题）
1. 先对数组进行排序
2. 使用哈希表记录每个数的出现次数
3. 两层循环遍历所有可能的两个数组合，在哈希表中查找第三个数
时间复杂度：O(n^2)，但需要额外空间
*/

// 解法1：排序+双指针法
func threeSum(nums []int) [][]int {
	var res [][]int
	n := len(nums)

	// 先排序
	sort.Ints(nums)

	for i := 0; i < n-2; i++ {
		// 跳过重复元素
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}

		left, right := i+1, n-1
		target := -nums[i]

		// TODO: 使用双指针寻找两数之和等于target
		// 提示：
		// 1. 当sum < target时，左指针右移
		// 2. 当sum > target时，右指针左移
		// 3. 当sum == target时，记录结果并跳过重复元素
		for left < right {
			sum := nums[left] + nums[right]
			if sum == target {
				res = append(res, []int{nums[i], nums[left], nums[right]})
				left++
				right--
				for left < right && nums[left] == nums[left-1] {
					left++
				}
				for left < right && nums[right] == nums[right+1] {
					right--
				}
			} else if sum < target {
				left++
			} else {
				right--
			}
		}
	}
	return res
}

// 解法2：哈希表法（框架）
// func threeSumHash(nums []int) [][]int {
// 	var result [][]int
// 	n := len(nums)
// 	if n < 3 {
// 		return result
// 	}

// 	sort.Ints(nums)
// 	numCount := make(map[int]int)

// 	// TODO: 实现哈希表解法
// 	// 提示：
// 	// 1. 先统计每个数字的出现次数存入哈希表
// 	// 2. 双层循环遍历所有可能的两个数组合
// 	// 3. 在哈希表中查找第三个数-(nums[i]+nums[j])
// 	// 4. 注意处理重复组合和数字重用问题

// 	return result
// }

// 比较函数，用于测试结果是否正确（忽略顺序）
func compareSlices(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}

	// 将每个子切片排序后转换为字符串进行比较
	m := make(map[string]bool)
	for _, sub := range a {
		sort.Ints(sub)
		key := ""
		for _, num := range sub {
			key += strconv.Itoa(num) + ","
		}
		m[key] = true
	}

	for _, sub := range b {
		sort.Ints(sub)
		key := ""
		for _, num := range sub {
			key += strconv.Itoa(num) + ","
		}
		if !m[key] {
			return false
		}
	}

	return true
}

func TestThreeSum(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected [][]int
	}{
		{
			name:     "example1",
			nums:     []int{-1, 0, 1, 2, -1, -4},
			expected: [][]int{{-1, -1, 2}, {-1, 0, 1}},
		},
		{
			name:     "empty",
			nums:     []int{},
			expected: [][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := threeSum(tt.nums)
			if !compareSlices(result, tt.expected) {
				t.Errorf("threeSum(%v) = %v, want %v", tt.nums, result, tt.expected)
			}
		})
	}
}

// func TestThreeSumHash(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		nums     []int
// 		expected [][]int
// 	}{
// 		{
// 			name:     "example1",
// 			nums:     []int{-1, 0, 1, 2, -1, -4},
// 			expected: [][]int{{-1, -1, 2}, {-1, 0, 1}},
// 		},
// 		{
// 			name:     "empty",
// 			nums:     []int{},
// 			expected: [][]int{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result := threeSumHash(tt.nums)
// 			if !compareSlices(result, tt.expected) {
// 				t.Errorf("threeSumHash(%v) = %v, want %v", tt.nums, result, tt.expected)
// 			}
// 		})
// 	}
// }
