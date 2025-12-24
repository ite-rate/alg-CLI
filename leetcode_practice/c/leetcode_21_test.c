/*
 * LeetCode #21: 题目21
 * 难度: 未知
 * 
 * 题目描述:
 * 由大模型直接生成
 * 
 * 代码骨架完整度: 80%
 */

#include <stdio.h>
#include <stdlib.h>
#include <assert.h>

/**
 * LeetCode #21 合并两个有序链表
 * 难度：简单
 * 
 * 题目描述：
 * 将两个升序链表合并为一个新的升序链表并返回。新链表是通过拼接给定的两个链表的所有节点组成的。
 * 
 * 示例：
 * 输入：l1 = [1,2,4], l2 = [1,3,4]
 * 输出：[1,1,2,3,4,4]
 */

// 链表节点定义
struct ListNode {
    int val;
    struct ListNode *next;
};

/**
 * 解法1：迭代法
 * 算法思路：
 * 1. 创建哑节点(dummy)作为新链表的起始点
 * 2. 使用指针遍历两个链表，比较当前节点值
 * 3. 将较小值的节点连接到新链表
 * 4. 当一个链表遍历完后，将另一个链表的剩余部分直接连接
 * 
 * 时间复杂度：O(m+n)，其中m和n分别是两个链表的长度
 * 空间复杂度：O(1)
 */
struct ListNode* mergeTwoLists_iterative(struct ListNode* l1, struct ListNode* l2) {
    // 创建哑节点简化边界处理
    struct ListNode dummy;
    struct ListNode* tail = &dummy;
    dummy.next = NULL;
    
    // TODO: 遍历两个链表，比较节点值并连接
    // 提示：使用while循环同时遍历l1和l2
    // 当l1和l2都不为空时，比较它们的val值
    // 将较小值的节点连接到tail后面，并移动对应指针
    
    while (l1 != NULL && l2 != NULL) {
        if (l1->val <= l2->val) {
            tail->next = l1;
            l1 = l1->next;
        } else {
            tail->next = l2;
            l2 = l2->next;
        }
        tail = tail->next;
    }
    
    // TODO: 处理剩余节点
    // 提示：将非空链表的剩余部分直接连接到tail后面
    if (l1 != NULL) {
        tail->next = l1;
    } else {
        tail->next = l2;
    }
    
    return dummy.next;
}

/**
 * 解法2：递归法
 * 算法思路：
 * 1. 基准情况：如果某个链表为空，直接返回另一个链表
 * 2. 比较两个链表头节点的值
 * 3. 递归地合并剩余部分
 * 
 * 时间复杂度：O(m+n)
 * 空间复杂度：O(m+n)（递归调用栈）
 */
struct ListNode* mergeTwoLists_recursive(struct ListNode* l1, struct ListNode* l2) {
    // TODO: 处理基准情况
    // 提示：如果l1为空返回l2，如果l2为空返回l1
    if (l1 == NULL) return l2;
    if (l2 == NULL) return l1;
    
    // TODO: 比较节点值并递归合并
    // 提示：如果l1的值较小，将l1的next指向递归结果，返回l1
    // 否则将l2的next指向递归结果，返回l2
    if (l1->val <= l2->val) {
        l1->next = mergeTwoLists_recursive(l1->next, l2);
        return l1;
    } else {
        l2->next = mergeTwoLists_recursive(l1, l2->next);
        return l2;
    }
}

// 创建新节点
struct ListNode* createNode(int val) {
    struct ListNode* node = (struct ListNode*)malloc(sizeof(struct ListNode));
    node->val = val;
    node->next = NULL;
    return node;
}

// 打印链表（用于测试）
void printList(struct ListNode* head) {
    while (head != NULL) {
        printf("%d", head->val);
        if (head->next != NULL) printf("->");
        head = head->next;
    }
    printf("\n");
}

// 释放链表内存
void freeList(struct ListNode* head) {
    while (head != NULL) {
        struct ListNode* temp = head;
        head = head->next;
        free(temp);
    }
}

// 比较两个链表是否相等
int compareLists(struct ListNode* l1, struct ListNode* l2) {
    while (l1 != NULL && l2 != NULL) {
        if (l1->val != l2->val) return 0;
        l1 = l1->next;
        l2 = l2->next;
    }
    return l1 == NULL && l2 == NULL;
}

int main() {
    printf("测试LeetCode #21 合并两个有序链表\n");
    
    // 测试用例1：常规情况
    printf("测试用例1: ");
    struct ListNode* l1 = createNode(1);
    l1->next = createNode(2);
    l1->next->next = createNode(4);
    
    struct ListNode* l2 = createNode(1);
    l2->next = createNode(3);
    l2->next->next = createNode(4);
    
    struct ListNode* expected = createNode(1);
    expected->next = createNode(1);
    expected->next->next = createNode(2);
    expected->next->next->next = createNode(3);
    expected->next->next->next->next = createNode(4);
    expected->next->next->next->next->next = createNode(4);
    
    // 测试迭代法
    struct ListNode* result1 = mergeTwoLists_iterative(l1, l2);
    printf("迭代法结果: ");
    printList(result1);
    assert(compareLists(result1, expected));
    printf("迭代法测试通过!\n");
    
    // 重新创建测试数据（因为原链表已被修改）
    l1 = createNode(1);
    l1->next = createNode(2);
    l1->next->next = createNode(4);
    
    l2 = createNode(1);
    l2->next = createNode(3);
    l2->next->next = createNode(4);
    
    // 测试递归法
    struct ListNode* result2 = mergeTwoLists_recursive(l1, l2);
    printf("递归法结果: ");
    printList(result2);
    assert(compareLists(result2, expected));
    printf("递归法测试通过!\n");
    
    // 测试用例2：空链表情况
    printf("\n测试用例2: 空链表\n");
    struct ListNode* empty = NULL;
    struct ListNode* single = createNode(5);
    
    struct ListNode* result3 = mergeTwoLists_iterative(empty, single);
    assert(result3->val == 5 && result3->next == NULL);
    printf("空链表测试通过!\n");
    
    // 释放内存
    freeList(result1);
    freeList(result2);
    freeList(result3);
    freeList(expected);
    freeList(single);
    
    printf("所有测试用例通过!\n");
    return 0;
}
