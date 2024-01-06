package main

import (
	"container/heap"
	"fmt"
	"strings"
)

type CharFrequency struct {
	char  byte
	count int
}

type CharFrequencyHeap []CharFrequency

func (h CharFrequencyHeap) Len() int           { return len(h) }
func (h CharFrequencyHeap) Less(i, j int) bool { return h[i].count > h[j].count }
func (h CharFrequencyHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *CharFrequencyHeap) Push(x interface{}) {
	*h = append(*h, x.(CharFrequency))
}

func (h *CharFrequencyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func rearrangeString(s string) string {
	charCount := make(map[byte]int)
	for i := range s {
		charCount[s[i]]++
	}

	maxHeap := &CharFrequencyHeap{}
	heap.Init(maxHeap)

	for char, count := range charCount {
		heap.Push(maxHeap, CharFrequency{char, count})
	}

	var result strings.Builder
	for maxHeap.Len() >= 2 {
		first := heap.Pop(maxHeap).(CharFrequency)
		second := heap.Pop(maxHeap).(CharFrequency)

		result.WriteByte(first.char)
		result.WriteByte(second.char)

		first.count--
		second.count--

		if first.count > 0 {
			heap.Push(maxHeap, first)
		}
		if second.count > 0 {
			heap.Push(maxHeap, second)
		}
	}

	if maxHeap.Len() > 0 {
		last := heap.Pop(maxHeap).(CharFrequency)
		if last.count > 1 {
			return ""
		}
		result.WriteByte(last.char)
	}

	return result.String()
}

func main() {
	s1 := "aab"
	fmt.Println("Input:", s1)
	fmt.Println("Output:", rearrangeString(s1)) // Output: "aba"

	s2 := "aaab"
	fmt.Println("Input:", s2)
	fmt.Println("Output:", rearrangeString(s2)) // Output: ""
}
