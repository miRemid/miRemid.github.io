package codes

type ListNode struct {
	Val  int
	Next *ListNode
}

func quick_sort(nums []int) {
	_quick_sort(nums, 0, len(nums)-1)
}

func _quick_sort(nums []int, left, right int) {
	if left >= right {
		return
	}
	k := nums[right]
	l, r := left, right-1
	for l < r {
		for l < r && nums[l] < k {
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
