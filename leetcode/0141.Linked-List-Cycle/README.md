# [141. Linked List Cycle](https://leetcode.cn/problems/linked-list-cycle/)

> 难度：Easy

## 题目

Given `head`, the head of a linked list, determine if the linked list has a cycle in it.

There is a cycle in a linked list if there is some node in the list that can be reached again by continuously following the `next` pointer. Internally, `pos` is used to denote the index of the node that tail's `next` pointer is connected to. **Note that `pos` is not passed as a parameter**.

Return `true` *if there is a cycle in the linked list*. Otherwise, return `false`.

**Example 1:**

![](https://assets.leetcode.com/uploads/2018/12/07/circularlinkedlist.png)

```
Input: head = [3,2,0,-4], pos = 1
Output: true
Explanation: There is a cycle in the linked list, where the tail connects to the 1st node (0-indexed).
```

**Example 2:**

![](https://assets.leetcode.com/uploads/2018/12/07/circularlinkedlist_test2.png)

```
Input: head = [1,2], pos = 0
Output: true
Explanation: There is a cycle in the linked list, where the tail connects to the 0th node.
```

**Example 3:**

![](https://assets.leetcode.com/uploads/2018/12/07/circularlinkedlist_test3.png)

```
Input: head = [1], pos = -1
Output: false
Explanation: There is no cycle in the linked list.
```

**Constraints:**

- The number of the nodes in the list is in the range `[0, 10^4]`.
- `-10^5 <= Node.val <= 10^5`
- `pos` is `-1` or a **valid index** in the linked-list.

**Follow up:** Can you solve it using `O(1)` (i.e. constant) memory?

<details>
<summary>官方中文题目翻译</summary>

给你一个链表的头节点 `head` ，判断链表中是否有环。

如果链表中有某个节点，可以通过连续跟踪 `next` 指针再次到达，则链表中存在环。 为了表示给定链表中的环，评测系统内部使用整数 `pos` 来表示链表尾连接到链表中的位置（索引从 0 开始）。** 注意：`pos` 不作为参数进行传递** 。仅仅是为了标识链表的实际情况。

*如果链表中存在环* ，则返回 `true` 。 否则，返回 `false` 。

**示例 1：**

![](https://assets.leetcode.cn/aliyun-lc-upload/uploads/2018/12/07/circularlinkedlist.png)

```
输入：head = [3,2,0,-4], pos = 1
输出：true
解释：链表中有一个环，其尾部连接到第二个节点。
```

**示例 2：**

![](https://assets.leetcode.cn/aliyun-lc-upload/uploads/2018/12/07/circularlinkedlist_test2.png)

```
输入：head = [1,2], pos = 0
输出：true
解释：链表中有一个环，其尾部连接到第一个节点。
```

**示例 3：**

![](https://assets.leetcode.cn/aliyun-lc-upload/uploads/2018/12/07/circularlinkedlist_test3.png)

```
输入：head = [1], pos = -1
输出：false
解释：链表中没有环。
```

**提示：**

- 链表中节点的数目范围是 `[0, 10^4]`
- `-10^5 <= Node.val <= 10^5`
- `pos` 为 `-1` 或者链表中的一个 **有效索引** 。

**进阶：** 你能用 `O(1)`（即，常量）内存解决此问题吗？

</details>

## 题目大意

判断链表中是否有环

## 解题思路

1. 通过两个指针标记节点,快指针一次走两步,慢指针走一步,存在环的情况快指针一定会追上慢指针(参考了[halfrost](https://github.com/halfrost/LeetCode-Go/blob/master/leetcode/0141.Linked-List-Cycle/README.md))
2. 通过map记录节点地址

```
claude code 给的其他思路:

1. Brent 算法（不修改链表，O(1) 空间）
和快慢指针一样，也是两个指针。区别是慢指针不跟着走，而是隔一段"传送"一次：
- 快指针每次走一步，同时计数。
- 走到 1、2、4、8…… 步时，把慢指针传送到快指针当前的位置，计数清零。
- 快指针撞上慢指针，说明有环；走到 nil，说明无环。

复杂度同样是 O(n) 时间、O(1) 空间，但每步只移动一个指针，实际走的步数通常比 Floyd 少。它的另一个好处是顺带就能得到环长。

2. 利用题目的数据范围
- 数步数：题目给了节点数上限 10⁴。如果走了 10⁴+1 步还没遇到 nil，就一定有环。写法最简单，但完全依赖这个上限，换一题就不能用了。
- 改值做标记：节点值的范围是 [-10⁵, 10⁵]。走过的节点把 val 改成范围外的值，比如 10⁵+1；再遇到这个值，就说明有环。这也依赖数据范围，而且会改掉节点值。

3. 修改链表结构
- 断链或指向哨兵：走过一个节点，就把它的 Next 指向一个固定的哨兵节点（或者指回 head）。之后走到哨兵，就说明这个节点以前来过。
- 反转链表：边走边把链表反转。有环的话，指针会绕回原来的 head；无环的话，会走到 nil。所以只要看最后停在哪里：回到 head 就有环，前提是链表至少有两个节点。
```

### 同类题

## 复杂度

- 时间复杂度： hasCycle = O(n); hasCycleMap = O(n);
- 空间复杂度： hasCycle = O(1); hasCycleMap = O(n);
