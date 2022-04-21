package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildStoreWithData(data ...int) *Store[int, int] {
	r := NewStore[int, int]()
	for _, v := range data {
		r .Insert(bl(v))
	}

	return r
}

func TestStoreGettingElementByValue(t *testing.T) {
	sut := buildStoreWithData(randomIntData...)

	v := sut.Get(11)
	assert.Equal(t, v.Value().Value, 11, "")
}

func TestStoreGettingElementWithLeaderTabel(t *testing.T) {
	sut := buildStoreWithData(randomIntData...)

	itr := sut.Get(11)

	tableSize := 7
	leaderTable := make([]int, 0, tableSize)
	
	leftIdx := tableSize / 2
	backItr := itr
	for ; leftIdx > 0 && backItr != nil; leftIdx -= 1 {
		backItr.Prev()
	}


	for i := 0; i < tableSize - leftIdx && itr != nil; i+=1 {
		leaderTable = append(leaderTable, itr.Value().Value)
		itr.Next()
	}

	assert.Equal(t, leaderTable, 11, "")

}
