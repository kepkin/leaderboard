package leaderboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBtreeIter(t *testing.T) {
	tree := buildBtreeWith(t, 5, randomIntData)

	localTree, idx := tree.Find(11)
	sut := localTree.Iter(idx)
	if sut == nil {
		t.Fatal("can't find 11")
	}

	data := make([]int, 0)
	c := 10
	for it := sut; it.Valid() && c > 0; it = it.Next() {
		c -= 1
		data = append(data, it.Value())
	}

	assert.Equal(t, []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, data, "ListAtIdx doesn't work")

	localTree, idx = tree.Find(11)

	sut2 := localTree.Iter(idx)
	data = data[:0]
	for it := sut2; it.Valid(); it = it.Prev() {
		data = append(data, it.Value())
	}

	assert.Equal(t, []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, data, "ListAtIdx doesn't work")
}
