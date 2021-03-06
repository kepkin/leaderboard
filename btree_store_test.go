package leaderboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func bl(v int) Tuple[int, int] {
	return Tuple[int, int]{v, v}
}

func buildStoreWithData(data ...int) *BtreeStore[int, int] {
	r := NewBtreeStore[int, int](
		StdLess[int],
		StdEquals[int],
		StdLess[int],
		StdEquals[int],
	)
	for _, v := range data {
		it, _ := r.Insert(bl(v))
		it.Close()
	}

	return r
}

func TestStoreGettingElementByValue(t *testing.T) {
	t.Skip()
	sut := buildStoreWithData(randomIntData...)

	v := sut.Get(11)
	assert.Equal(t, 11, v.Value().Val, "")
}

func TestStoreGettingElementWithLeaderTabel(t *testing.T) {
	sut := buildStoreWithData(randomIntData...)

	itr := sut.Get(11)

	tableSize := 7
	leaderTable := make([]int, 0, tableSize)

	leftIdx := tableSize / 2
	backItr := itr
	for ; leftIdx > 0 && backItr.Valid(); leftIdx -= 1 {
		backItr.Prev()
	}

	for i := 0; i < tableSize-leftIdx && itr.Valid(); i += 1 {
		leaderTable = append(leaderTable, itr.Value().Val)
		itr.Next()
	}
	itr.Close()

	assert.Equal(t, []int{8, 9, 10, 11, 12, 13, 14}, leaderTable, "")

	//Test that pkbtree has correct number of items
	it := sut.pkbtree.Begin()
	defer it.Close()
	for ; it.Valid(); it.Next() {
		t.Log(it.Value())
	}
	// t.Fatal()
}
