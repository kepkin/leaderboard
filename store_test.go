package leaderboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildStoreWithData(data ...int) *Store[Int, Int] {
	r := NewStore[Int, Int]()
	for _, v := range data {
		r .Insert(bl(v))
	}

	return r
}

func TestStoreGettingElementByValue(t *testing.T) {
	sut := buildStoreWithData(randomIntData...)

	v := sut.Get(11)
	assert.Equal(t, Int(11), v.Value().Value, "")
}


func TestStoreGettingElementWithLeaderTabel(t *testing.T) {
	sut := buildStoreWithData(randomIntData...)

	itr := sut.Get(11)

	tableSize := 7
	leaderTable := make([]Int, 0, tableSize)
	
	leftIdx := tableSize / 2
	backItr := itr
	for ; leftIdx > 0 && backItr != nil; leftIdx -= 1 {
		backItr.Prev()
	}


	for i := 0; i < tableSize - leftIdx && itr != nil; i+=1 {
		leaderTable = append(leaderTable, itr.Value().Value)
		itr.Next()
	}

	assert.Equal(t, []Int{8, 9, 10, 11, 12, 13, 14}, leaderTable, "")

}
