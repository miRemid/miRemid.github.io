---
title: "TopK面试题1"
date: 2023-03-09T14:42:35+08:00
draft: false
toc: true
images:
tags: 
  - 面试
---

> [TopK来源](https://osjobs.net/topk/all/)

## 使用递归及非递归两种方式实现快速排序

快排思想，选取一个哨兵节点，从左边寻找到第一个比哨兵大的节点，从右边寻找到第一个比哨兵小的节点，交换两者位置。
如果寻找过程中没有找到对应的元素，则说明排序已经完成

首先是递归的方式
```golang

func quick_sort(nums []int){
    _quick_sort(nums, 0, len(nums)-1)
}

func _quick_sort(nums []int, left, right int) {
    if left >= right {
        return 
    }
    k := nums[right]
    l, r := left, right - 1
    for l < r {
        for l < r && nums[i] < k {
            l++
        }
        for r > l && nums[r] > k {
            r--
        }
        nums[l], nums[r] = nums[r], nums[l]
    }
    nums[right], nums[l] = nums[l], nums[right]
    _quick_sort(nums, left, l-1)
    _quick_sort(nums, l+1, right)
}
```

非递归版本，其实就是用栈模拟递归
```golang
func quick_sort_with_stack(nums []int) {
    type pair struct {
        x, y int
    }
    var s = make([]pair, 0)
    if len(nums) < 2 {
        return
    }
    s = append(s, pair{0, len(nums) - 1})
    for len(s) != 0 {
        p := s[0]
        s = s[1:]
        left, right := p.x, p.y-1
        if left >= right {
            continue
        }
        k := nums[p.y]
        for left < right {
            for left < right && nums[left] < k {
                left++
            }
            for left < right && nums[right] > k {
                right--
            }
            nums[left], nums[right] = nums[right], nums[left]
        }
        nums[p.y], nums[left] = nums[left], nums[p.y]
        s = append(s, pair{p.x, left - 1})
        s = append(s, pair{left + 1, p.y})
    }
}
```

## 环形链表

> 给你一个链表的头节点 head ，判断链表中是否有环。
如果链表中有某个节点，可以通过连续跟踪 next 指针再次到达，则链表中存在环。 为了表示给定链表中的环，评测系统内部使用整数 pos 来表示链表尾连接到链表中的位置（索引从 0 开始）。注意：pos 不作为参数进行传递 。仅仅是为了标识链表的实际情况。
如果链表中存在环 ，则返回 true 。 否则，返回 false 。

经典双指针问题，可以用hashmap但没必要

```golang
type ListNode struct {
    Val int
    Next *ListNode
}
func hasCycle(head *ListNode) bool {
    if head == nil || head.Next == nil {
        return false
    } 
    slow, fast := head, head.Next.Next
    for fast != nil && fast.Next != nil {
        if slow == fast {
            return true
        }
        slow = slow.Next
        fast = fast.Next.Next
    }
    return false
}
```
但还是加一个HashMap的方法吧，万一面试官要求用呢
```golang
func hasCycle(head *ListNode) bool {
    var hashMap = make(map[*ListNode]struct{})
    for head != nil {
        if _, ok := hashMap[head]; ok {
            return true
        }
        hashMap[head] = struct{}{}
        head = head.Next
    }
    return false
}
```

## LRU 缓存机制

经典LRU缓存机制，简单复习一下LRU

LRU缓存简单来讲是一个栈入结构，保存着使用最为频繁的缓存数据，使得可以快速在内存中查询到数据而不需要进行一次额外的数据查询例如数据库查询以提高效率，在操作系统的页表中也有使用

为了实现最热的缓存数据，我们首先想到的是用数组栈来模拟，但为了更方便我们选择双向链表来实现，当然在Leetcode官方题解中也使用了双向链表

在这里，为了体验Golang的泛型并且充分了解双向链表的使用，决定从头编写双向链表而非使用官方的container包

```golang



```

## 反转链表

经典面试题，分递归和非递归版本，首先是递归版本
```golang
func reverseLink(head *ListNode) *ListNode {
    if head == nil || head.Next == nil {
        return head
    }
    newHead := reverseLink(head.Next)
    head.Next.Next = head
    head.Next = nil
    return newHead
}
```
然后是非递归版本
```golang
func reverseLink(head *ListNode) *ListNode {
    var prev *ListNode
    cur := head
    for cur != nil {
        next := cur.Next
        cur.Next = prev
        prev = cur
        cur = next
    }
    return prev
}
```

## 删除链表的倒数第 N 个结点
经典双指针题目，先让快指针走N次，随后生成新的指针直到快指针走完，慢指针指向的节点就是第N个节点，删除即可

```golang
func removeNthFromEnd(head *ListNode, n int) *ListNode {
    dummyNode := &ListNode{
        Next: head,
    }
    fast := dummyNode
    for i := 0; i < n && fast != nil; i++ {
        fast = (fast.Next)
    }
    ptr := dummyNode
    for fast.Next != nil {
        ptr = ptr.Next
        fast = fast.Next
    }
    ptr.Next = ptr.Next.Next
    return dummyNode.Next
}
```

## 快速排序的空间复杂度是多少？时间复杂度的最好最坏的情况是多少，有哪些优化方案？

1. 快速排序的空间复杂度可以为$O(logn)$
2. 时间复杂度最好的情况下为O(nlogn)，最坏情况下为O(n^2)
3. 快速排序的优化点如下
    - 采用随机哨兵策略，如[5,4,3,2,1]这种数据的影响。另外也可以选择三数的中位数（左、中、右，中间数据，如8, 0, 6，选择6作为哨兵）
    - 快排针对数据量小并且部分有序的数组效率并不高，可以在切分为一定大小后转换为插入排序
    - 将与哨兵相同的元素放在分割点附近，减少分割后的数组长度