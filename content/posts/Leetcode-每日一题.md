---
title: "Leetcode 每日一题"
date: 2021-07-23T11:50:27+08:00
draft: false
toc: false
images:
tags: 
  - Leetcode
---

## 2021/7/23 1893 Check if All the Integers in a Range Are Covered (EASY)

You are given a 2D integer array `ranges` and two integers `left` and `right`. Each `ranges[i] = [starti, endi]` represents an inclusive interval between `starti` and `endi`.

Return `true` if each integer in the inclusive range `[left, right]` is covered by __at least one__ interval in `ranges`. Return `false` otherwise.

An integer `x` is covered by an interval `ranges[i] = [starti, endi] if starti <= x <= endi.`

#### 暴力
简单题直接暴力破解即可，无需多言
```go
func isCovered(ranges [][]int, left, right int) bool {
    var boolMap = make(map[int]bool)
    for num := left; num <= right; num++ {
        var flag = false
        for _, r := range ranges {
            if num >= r[0] && num <= r[1] {
                flag = true
                break
            }
        }
        boolMap[num] = flag
    }
    for _, v := range boolMap {
        if !v {
            return false
        }
    }
    return true
}
```
执行用时：4 ms, 在所有 Go 提交中击败了21.05%的用户

内存消耗：2.6 MB, 在所有 Go 提交中击败了14.91%的用户

#### 差分数组
另一种简单的方法就是差分数组，先统计每组Range的氛围，然后再看提供的数字是否处于范围之内即可
```go
func isCovered(ranges [][]int, left, right int) bool {
    diff := [52]int{}
    // 1. check ranges
    for _, v := range ranges {
        diff[v[0]]++
        diff[v[1]+1]--
    }
    // 2. check num in left right range
    cnt := 0
    for i := 1; i <= right; i++ {
        cnt += diff[i]
        if i >= left && cnt <= 0 {
            return false
        }
    }
    return true
}
```
执行用时：0 ms, 在所有 Go 提交中击败了100%的用户

内存消耗：2.5 MB, 在所有 Go 提交中击败了80.70%的用户

## 2021/7/26 1713 Minimum Operations to Make a Subsequence (HARD)

You are given an array target that consists of distinct integers and another integer array arr that can have duplicates.

In one operation, you can insert any integer at any position in arr. For example, if arr = [1,4,1,2], you can add 3 in the middle and make it [1,4,3,1,2]. Note that you can insert the integer at the very beginning or end of the array.

Return the minimum number of operations needed to make target a subsequence of arr.

A subsequence of an array is a new array generated from the original array by deleting some elements (possibly none) without changing the remaining elements' relative order. For example, [2,7,4] is a subsequence of [4,2,3,7,2,1,4] (the underlined elements), while [2,4,2] is not.

#### 贪心

题目要求寻找到两个数组中的最少插入元素个数使target成为arr的子集，可以转换思路，寻找到两个数组的最长公共子集，最后要插入的个数就是target数组的长度减去最长公共子集长度

```go
func minOperations(target []int, arr []int) int {
	n := len(target)
	pos := make(map[int]int, n)
	for i, val := range target {
		pos[val] = i
	}
	d := []int{}
	for _, val := range arr {
		if idx, has := pos[val]; has {
			if p := sort.SearchInts(d, idx); p < len(d) {
				d[p] = idx
			} else {
				d = append(d, idx)
			}
		}
	}
	return n - len(d)
}
```

执行用时：168 ms, 在所有 Go 提交中击败了63.89%的用户

内存消耗：12.3 MB, 在所有 Go 提交中击败了91.67%的用户

## 2021/7.26 671 Second Minimum Node In a Binary Tree (EASY)

Given a non-empty special binary tree consisting of nodes with the non-negative value, where each node in this tree has exactly two or zero sub-node. If the node has two sub-nodes, then this node's value is the smaller value among its two sub-nodes. More formally, the property root.val = min(root.left.val, root.right.val) always holds.

Given such a binary tree, you need to output the second minimum value in the set made of all the nodes' value in the whole tree.

If no such second minimum value exists, output -1 instead.

#### DFS暴力
由题意可知，根节点的值是最小的，要找到整棵树中的第二小的值只需要对所有节点进行遍历，寻找到大于根节点的最大值即可

初始化一个`ans`等于-1，代表未找到第二小的值，对树进行遍历，当值大于root的值并大于ans时，更新ans的值。最后返回ans的值，就是第二小的值

```go
func findSecondMinimumValue(root *TreeNode) int {
    ans := -1
    rootVal := root.Val
    var dfs func(*TreeNode)
    dfs = func(node *TreeNode) {
        if node == nil || ans != -1 && node.Val >= ans {
            return
        }
        if node.Val > rootVal {
            ans = node.Val
        }
        dfs(node.Left)
        dfs(node.Right)
    }
    dfs(root)
    return ans
}
```

## 2021/7.28 863 
Given the root of a binary tree, the value of a target node target, and an integer k, return an array of the values of all nodes that have a distance k from the target node.

You can return the answer in any order.

#### 哈系+DFS

首先寻找到所有节点的父节点，然后从Target节点出发，寻找距离K的节点即可

```go
func distanceK(root *TreeNode, target *TreeNode, k int) []int {
    // 1. get all father node
	var fathers = make(map[int]*TreeNode)
	var getFather func(root *TreeNode)
	getFather = func(root *TreeNode) {
		if root == nil {
			return
		}
		if root.Left != nil {
			fathers[root.Left.Val] = root
			getFather(root.Left)
		}
		if root.Right != nil {
			fathers[root.Right.Val] = root
			getFather(root.Right)
		}
	}
	getFather(root)
	var res = make([]int, 0)
	// 2. find the target node
	var getNode func(*TreeNode, *TreeNode, int)
	getNode = func(node, from *TreeNode, dis int) {
		if node == nil {
			return
		}
		if dis == k {
			res = append(res, node.Val)
			return
		}
		if node.Left != from {
			getNode(node.Left, node, dis+1)
		}
		if node.Right != from {
			getNode(node.Right, node, dis+1)
		}
		if fathers[node.Val] != from {
			getNode(fathers[node.Val], node, dis+1)
		}
	}
	getNode(target, nil, 0)
	return res
}
```

执行用时：0 ms, 在所有 Go 提交中击败了100%的用户

内存消耗：3.2 MB, 在所有 Go 提交中击败了52.08%的用户

## 2021/7.30 171 Excel Sheet Column Number (EASY)
Given a string columnTitle that represents the column title as appear in an Excel sheet, return its corresponding column number.

    A -> 1
    B -> 2
    ...
    Z -> 26
    AA -> 27
    ...

#### 遍历
进制转换题，非常简单
```go
func titleToNumber(columnTitle string) int {
    var num int = 0
    for i, mul := len(columnTitle) - 1, 1; i>= 0; i-- {
        k := columnTitle[i] - 'A' + 1
        num += int(k) * mul
        mul *= 26
    }
    return num
}
```
## 987 2021/7/31 Vertical Order Traversall of a Binary Tree

Give the root of a binary tree, calculate the vertical order traversal of the binary tree.

For each node at position (row, col), its left and right children will be at positions (row + 1, col - 1) and (row + 1, col + 1) respectively. The rootof the tree is at (0, 0)

The vertical order traversal of a binary tree is a list of top-to-bottom orderings for each column index starting from the leftmost column and ending on the rightmost column. There may be multiple nodes in the same row and same column. In such a case, sort these nodes by their valeus.

Return the vertical order traversal of the binary tree.

#### DFS暴力+Sort+Hash
Map直接存同col的节点，然后从小到大遍历并排序输出即可

```go
func verticalTraversal(root *TreeNode) (res [][]int) {
    // 1. hash map, find the same coloumn nodes
    var hashMap = make(map[int][]int)
    var left, right = 999, -999
    var dfs func(root *TreeNode, row, col int)
    dfs = func(root *TreeNode, row, col int) {
        if root == nil {
            return
        }
        if col < left {
            left = col
        }
        if col > right {
            right = col
        }
        if _, ok := hashMap[col]; ok {
            hashMap[col] = append(hashMap[col], root.Val)
        } else {
            hashMap[col] = make([]int, 0)
            hashMap[col] = append(hashMap[col], root.Val)
        }
        dfs(root.Left, row+1, col-1)
        dfs(root.Right, row+1, col+1)
    }
    dfs(root, 0, 0)
    for i := left; i <= right; i++ {
        if arr, ok := hashMap[i]; ok {
            sort.Ints(arr)
            res = append(res, arr)
        }
    }
    return
}
```

Leetcode-cn测试集有问题

    Input:
    [3, 1, 4, 0, 2, 2]
    Output:
    [[0], [1], [2, 2, 3], [4]]
    Expect:
    [[0], [1], [3, 2, 2], [4]] # 题意从小到大，但预期为从大到小

