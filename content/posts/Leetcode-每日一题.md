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




















