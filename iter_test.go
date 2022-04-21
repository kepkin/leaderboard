package leaderboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBtreeIter(t *testing.T) {
	tree := buildBtreeWith(t, randomIntData)

	sut := tree.Find(bl(11))
	if sut == nil {
		t.Fatal("can't find 11")
	}

	data := make([]Int, 0)
	c := 10
	for it := sut; it != nil && c > 0; it = it.Next() {
		c -= 1
		data = append(data, it.Value().Value)
	}

	assert.Equal(t, []Int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, data, "ListAtIdx doesn't work")

	sut2 := tree.Find(bl(11))
	data = data[:0]
	for it := sut2; it != nil; it = it.Prev() {
		data = append(data, it.Value().Value)
	}

	assert.Equal(t, []Int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, data, "ListAtIdx doesn't work")
}
