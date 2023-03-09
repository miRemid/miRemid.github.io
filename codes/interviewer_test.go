package codes

import (
	"fmt"
	"testing"
)

func TestQuickSort(t *testing.T) {
	nums := []int{
		3, 1, 5, 8, 2,
	}
	quick_sort(nums)
	fmt.Println(nums)
	nums = []int{
		3, 1, 5, 8, 2,
	}
	quick_sort_with_stack(nums)
	fmt.Println(nums)
}
