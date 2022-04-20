package main

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func bl(v int) BTreeLeaf[int, int] {
	return BTreeLeaf[int, int]{v, v}
}

func buildNodeWithData(size int, data ...int) *Node[int, int] {
	r := NewNode[int, int](uint16(size), nil, nil)
	for _, v := range data {
		r.Data = append(r.Data, bl(v))
	}

	return r
}

func buildNodeDataToCmp(data ...int) []BTreeLeaf[int, int] {

	r := make([]BTreeLeaf[int, int], 0, len(data))
	for _, v := range data {
		r = append(r, bl(v))
	}
	return r
}

func leafArrToIntArr(src []BTreeLeaf[int, int]) []int {
	r := make([]int, 0, len(src))
	for _, v := range src {
		r = append(r, v.OrderKey)
	}

	return r
}

func TestSplit(t *testing.T) {

	// left corner case
	{
		p := buildNodeWithData(5, 1, 2, 3, 4, 5)
		p, _ = p.Insert(bl(0))

		assert.Equal(t, p.Data, []BTreeLeaf[int, int]{bl(3)}, "")
		assert.Equal(t, p.Childs[0].Data, []BTreeLeaf[int, int]{bl(0), bl(1), bl(2)}, "")
		assert.Equal(t, p.Childs[1].Data, []BTreeLeaf[int, int]{bl(4), bl(5)}, "")
	}

	// right corner case
	{
		p := buildNodeWithData(5, 1, 2, 3, 4, 5)
		p, _ = p.Insert(bl(6))

		assert.Equal(t, p.Data, []BTreeLeaf[int, int]{bl(3)}, "")
		assert.Equal(t, p.Childs[0].Data, []BTreeLeaf[int, int]{bl(1), bl(2)}, "")
		assert.Equal(t, p.Childs[1].Data, []BTreeLeaf[int, int]{bl(4), bl(5), bl(6)}, "")
	}

	// middle corner case
	{
		p := buildNodeWithData(5, 1, 2, 4, 5, 6)
		p, _ = p.Insert(bl(3))

		assert.Equal(t, p.Data, []BTreeLeaf[int, int]{bl(4)}, "")
		assert.Equal(t, p.Childs[0].Data, []BTreeLeaf[int, int]{bl(1), bl(2), bl(3)}, "")
		assert.Equal(t, p.Childs[1].Data, []BTreeLeaf[int, int]{bl(5), bl(6)}, "")
	}

	// case with preserving childs

	// {
	// 	r  := buildNodeWithData(2, 3, 6)
	// 	r.Childs[0] = buildNodeWithData(2, 1, 2)
	// 	r.Childs[1] = buildNodeWithData(2, 4, 5)
	// 	r.Childs[2] = buildNodeWithData(2, 7, 8)

	// 	r, _ = r.Insert(bl(0))

	// 	t.Log(r.Data)
	// 	t.Log(r.Childs[0].Data)
	// 	t.Log(r.Childs[0].Childs[0].Data)
	// 	t.Log(r.Childs[1].Data)
	// 	// assert.Equal(t, leafArrToIntArr(r.Data), []int{6}, "")
	// 	// assert.Equal(t, leafArrToIntArr(r.Childs[0].Data), []int{2, 3}, "")
	// 	// assert.Equal(t, leafArrToIntArr(r.Childs[0].Childs[0].Data), []int{0, 1}, "")
	// 	// assert.Equal(t, leafArrToIntArr(r.Childs[0].Childs[1].Data), []int{}, "")

	// 	// assert.Equal(t, leafArrToIntArr(r.Childs[0].Childs[2].Data), []int{4, 5}, "")

	// 	// assert.Equal(t, leafArrToIntArr(r.Childs[1].Data), []int{7, 8}, "")

	// 	data := make([]int, 0, 49)
	// 	data = AllocateV(r, data)
	// 	t.Log(data)
	// 	t.Fail()

	// }

}

func TestInsertKeepsOrder(t *testing.T) {
	sut := NewNode[int, int](5, nil, nil)

	logData := make([]int, 0, 49)

	insert := func(key int, log bool) {
		sut, _ = sut.Insert(BTreeLeaf[int, int]{OrderKey: key, Value: key})

		if log {
			logData = logData[:0]

			logData = AllocateV(sut, logData)
			t.Logf("insert %v (len: %v): %v", key, len(logData), logData)
		}
	}

	insert(10, true)
	insert(32, true)
	insert(45, true)
	insert(34, true)
	insert(26, true)
	insert(16, true)
	insert(4, true)
	insert(40, true)
	insert(22, true)
	insert(21, true)
	insert(29, true)
	insert(20, true)
	insert(24, true)
	insert(12, true)
	insert(6, true)
	insert(15, true)
	insert(27, true)
	insert(1, true)
	insert(43, true)
	insert(44, true)
	insert(17, true)
	insert(46, true)
	insert(3, true)
	insert(8, true)
	insert(30, true)
	insert(35, true)
	insert(41, true)
	insert(18, true)
	insert(47, true)
	insert(42, true)
	insert(13, true)
	insert(36, true)
	insert(7, true)
	insert(9, true)
	insert(28, true)
	insert(25, true)
	insert(48, true)
	insert(5, true)
	insert(14, true)
	insert(19, true)
	insert(31, true)
	insert(23, true)
	insert(11, true)
	insert(38, true)
	insert(33, true)
	insert(37, true)
	insert(2, true)
	insert(39, true)

	res := make([]int, 0, 48)
	res = AllocateV(sut, res)

	assert.Equal(t, 48, len(res), "Some elements were lost")
	for i, v := range res {
		assert.Equal(t, i+1, v, "Elements not in order")
		if i+1 != v {
			break
		}
	}

	t.Log("----------------")
	t.Log(res)
}
