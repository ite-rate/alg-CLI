/*
 * LeetCode #15: 题目15
 * 难度: 未知
 * 
 * 题目描述:
 * 由大模型直接生成
 * 
 * 代码骨架完整度: 30%
 */

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * LeetCode #15 - 三数之和
 * 难度：中等
 * 题目描述：给定一个包含n个整数的数组nums，判断nums中是否存在三个元素a、b、c，使得a + b + c = 0？
 * 找出所有满足条件且不重复的三元组。
 * 注意：答案中不可以包含重复的三元组。
 *
 * 示例：
 * 输入：nums = [-1,0,1,2,-1,-4]
 * 输出：[[-1,-1,2],[-1,0,1]]
 */
class Solution {
    /**
     * 解法1：排序 + 双指针
     * 思路：
     * 1. 先对数组进行排序
     * 2. 遍历数组，固定一个数nums[i]，然后在剩余部分使用双指针寻找两数之和等于-nums[i]
     * 3. 注意跳过重复元素以避免重复解
     * 时间复杂度：O(n^2)
     */
    public List<List<Integer>> threeSum(int[] nums) {
        List<List<Integer>> res = new ArrayList<>();
        Arrays.sort(nums); // 先排序
        
        for (int i = 0; i < nums.length - 2; i++) {
            if (i > 0 && nums[i] == nums[i - 1]) continue; // 跳过重复元素
            
            int left = i + 1;
            int right = nums.length - 1;
            int target = -nums[i];
            
            while (left < right) {
                int sum = nums[left] + nums[right];
                
                if (sum < target) {
                    left++;
                } else if (sum > target) {
                    right--;
                } else {
                    // 找到一组解
                    res.add(Arrays.asList(nums[i], nums[left], nums[right]));
                    
                    // 跳过重复元素
                    while (left < right && nums[left] == nums[left + 1]) left++;
                    while (left < right && nums[right] == nums[right - 1]) right--;
                    
                    // 移动指针
                    left++;
                    right--;
                }
            }
        }
        
        return res;
    }

    /**
     * 解法2：哈希表法（不推荐，仅作为备选方案）
     * 思路：
     * 1. 先对数组进行排序
     * 2. 双重循环固定两个数，用哈希表查找第三个数
     * 3. 注意处理重复解的问题
     * 时间复杂度：O(n^2)，但实际效率不如双指针
     */
    public List<List<Integer>> threeSum2(int[] nums) {
        List<List<Integer>> res = new ArrayList<>();
        Arrays.sort(nums);
        
        // TODO: 实现哈希表解法
        // 提示：
        // 1. 使用双重循环遍历所有两数组合
        // 2. 使用HashSet记录已经见过的数字
        // 3. 检查是否存在需要的补数
        // 4. 注意处理重复解
        
        return res;
    }
}

// 测试代码
class SolutionTest {
    /**
     * 比较两个二维列表是否相等（忽略顺序）
     */
    private boolean isEqual(List<List<Integer>> a, List<List<Integer>> b) {
        if (a.size() != b.size()) return false;
        
        List<String> aStrs = new ArrayList<>();
        for (List<Integer> list : a) {
            Collections.sort(list);
            aStrs.add(list.toString());
        }
        Collections.sort(aStrs);
        
        List<String> bStrs = new ArrayList<>();
        for (List<Integer> list : b) {
            Collections.sort(list);
            bStrs.add(list.toString());
        }
        Collections.sort(bStrs);
        
        return aStrs.equals(bStrs);
    }

    @Test
    public void testThreeSum() {
        Solution solution = new Solution();
        
        // 测试用例1
        int[] nums1 = {-1, 0, 1, 2, -1, -4};
        List<List<Integer>> expected1 = Arrays.asList(
            Arrays.asList(-1, -1, 2),
            Arrays.asList(-1, 0, 1)
        );
        List<List<Integer>> result1 = solution.threeSum(nums1);
        assertTrue(isEqual(expected1, result1));
        
        // 测试用例2
        int[] nums2 = {0, 0, 0, 0};
        List<List<Integer>> expected2 = Arrays.asList(
            Arrays.asList(0, 0, 0)
        );
        List<List<Integer>> result2 = solution.threeSum(nums2);
        assertTrue(isEqual(expected2, result2));
    }

    @Test
    public void testThreeSum2() {
        Solution solution = new Solution();
        
        // 测试用例1
        int[] nums1 = {-1, 0, 1, 2, -1, -4};
        List<List<Integer>> expected1 = Arrays.asList(
            Arrays.asList(-1, -1, 2),
            Arrays.asList(-1, 0, 1)
        );
        List<List<Integer>> result1 = solution.threeSum2(nums1);
        assertTrue(isEqual(expected1, result1));
        
        // 测试用例2
        int[] nums2 = {0, 0, 0, 0};
        List<List<Integer>> expected2 = Arrays.asList(
            Arrays.asList(0, 0, 0)
        );
        List<List<Integer>> result2 = solution.threeSum2(nums2);
        assertTrue(isEqual(expected2, result2));
    }
}
