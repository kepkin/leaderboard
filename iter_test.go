package main

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBtreeIter(t *testing.T) {
	tree := buildBtreeWith(t, randomIntData)

	localNode, idx, ok := tree.Find(bl(11))
	if localNode == nil || !ok {
		t.Fatal("can't find 11")
	}

	data := make([]int, 0)
	c := 10
	sut := &Iter[int, int]{localNode, idx}
	for it := sut; it != nil && c > 0; it = it.Next() {
		c -= 1
		data = append(data, it.Value().Value)
	}

	assert.Equal(t, []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, data, "ListAtIdx doesn't work")

	sut2 := &Iter[int, int]{localNode, idx}
	data = data[:0]
	for it := sut2; it != nil; it = it.Prev() {
		data = append(data, it.Value().Value)
	}

	assert.Equal(t, []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, data, "ListAtIdx doesn't work")

}
