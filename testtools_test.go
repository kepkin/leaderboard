package leaderboard

import (
	"testing"
)

type Int int

func (this Int) Less(than Ordered) bool {
	return this < than.(Int)
}

func (this Int) Equals(than Comparable) bool {
	return this == than.(Int)
}

func bl(v int) BTreeLeaf[Int, Int] {
	return BTreeLeaf[Int, Int]{Int(v), Int(v)}
}

func buildNodeWithData(size int, data ...int) *Node[Int, Int] {
	r := NewNode[Int, Int](uint16(size), nil, nil)
	for _, v := range data {
		r.Data = append(r.Data, bl(v))
	}

	return r
}

func buildNodeDataToCmp(data ...int) []BTreeLeaf[Int, Int] {

	r := make([]BTreeLeaf[Int, Int], 0, len(data))
	for _, v := range data {
		r = append(r, bl(v))
	}
	return r
}

func leafArrToIntArr(src []BTreeLeaf[Int, Int]) []Int {
	r := make([]Int, 0, len(src))
	for _, v := range src {
		r = append(r, v.OrderKey)
	}

	return r
}

// random seq from 1 to 48
var randomIntData = []int{
	10, 32, 45, 34, 26, 16, 4, 40, 22, 21, 29, 20, 24, 12, 6, 15, 27, 1, 43, 44, 17, 46, 3, 8, 30, 35, 41, 18, 47, 42, 13, 36, 7, 9, 28, 25, 48, 5, 14, 19, 31, 23, 11, 38, 33, 37, 2, 39,
}

func buildBtreeWith(t *testing.T, data []int) *Node[Int, Int] {
	sut := NewNode[Int, Int](5, nil, nil)

	logData := make([]Int, 0, 49)

	for _, v := range randomIntData {
		sut, _ = sut.Insert(bl(v))

		logData = logData[:0]

		itr := sut.Begin()
		for ; itr != nil; itr = itr.Next() {
			logData = append(logData, itr.Value().OrderKey)
		}

		// logData = AllocateV(sut, logData)
		t.Logf("insert %v (len: %v): %v", v, len(logData), logData)
	}

	return sut
}
