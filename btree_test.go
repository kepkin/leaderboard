package leaderboard

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const treeSize = 5

type splitRecord struct {
	value int
	idx   int
	n     *Node[int]
}

func TestBtreeSizeof(t *testing.T) {
	res := buildNodeWithData(treeSize, nil, 1, 2, 3, 4, 5)

	leafSize := unsafe.Sizeof(0)

	if unsafe.Sizeof(*res) > 8 {

		t.Log("btree leaf has sizeof ", leafSize)
		t.Log("btree has sizeof ", unsafe.Sizeof(*res))
		t.Fail()
	}
}

func TestSplit(t *testing.T) {

	splitRecorder := func(dst *[]splitRecord) OnSplitTrigger[int] {
		return func(v int, n *Node[int], idx int) {
			*dst = append(*dst, splitRecord{value: v, idx: idx, n: n})
		}
	}

	// left corner case
	{
		var splitRecords []splitRecord

		p := buildNodeWithData(treeSize, splitRecorder(&splitRecords), 1, 2, 3, 4, 5)
		p, _ = p.Insert(0)

		assert.Equal(t, p.Data, []int{3}, "")
		assert.Equal(t, p.Childs[0].Data, []int{0, 1, 2}, "")
		assert.Equal(t, p.Childs[1].Data, []int{4, 5}, "")
		assert.Equal(t, splitRecords, []splitRecord{{3, 0, p}, {4, 0, p.Childs[1]}, {5, 1, p.Childs[1]}, {0, 0, p.Childs[0]}, {1, 1, p.Childs[0]}, {2, 2, p.Childs[0]}}, "")

		t.Log(splitRecords)
		// btreeCompare[int](t, p, nil)
	}

	// right corner case
	{
		var splitRecords []splitRecord

		p := buildNodeWithData(treeSize, splitRecorder(&splitRecords), 1, 2, 3, 4, 5)
		p, _ = p.Insert(6)

		assert.Equal(t, p.Data, []int{3}, "")
		assert.Equal(t, p.Childs[0].Data, []int{1, 2}, "")
		assert.Equal(t, p.Childs[1].Data, []int{4, 5, 6}, "")
		// left node with 1, 2 actually stays the same
		assert.Equal(t, splitRecords, []splitRecord{{3, 0, p}, {4, 0, p.Childs[1]}, {5, 1, p.Childs[1]}, {6, 2, p.Childs[1]}}, "")
	}

	// middle corner case
	{
		var splitRecords []splitRecord

		p := buildNodeWithData(treeSize, splitRecorder(&splitRecords), 1, 2, 4, 5, 6)
		p, _ = p.Insert(3)

		assert.Equal(t, p.Data, []int{4}, "")
		assert.Equal(t, p.Childs[0].Data, []int{1, 2, 3}, "")
		assert.Equal(t, p.Childs[1].Data, []int{5, 6}, "")
		assert.Equal(t, splitRecords, []splitRecord{{4, 0, p}, {5, 0, p.Childs[1]}, {6, 1, p.Childs[1]}, {3, 2, p.Childs[0]}}, "")
	}

	// case with preserving childs

	{
		var splitRecords []splitRecord

		p := buildNodeWithData(treeSize, splitRecorder(&splitRecords), 3, 6, 9, 16, 19)
		buildChildWithData(p, 0, 0, 1, 2)
		buildChildWithData(p, 1, 4, 5)
		buildChildWithData(p, 2, 7, 8)
		buildChildWithData(p, 3, 11, 12, 13, 14, 15)
		buildChildWithData(p, 4, 17, 18)
		buildChildWithData(p, 5, 20, 21)

		p, _ = p.Insert(10)
		assert.Equal(t, p.Data, []int{9}, "")
		assert.Equal(t, p.Childs[0].Data, []int{3, 6}, "")
		assert.Equal(t, p.Childs[1].Data, []int{13, 16, 19}, "")
		assert.Equal(t, p.Childs[0].Childs[0].Data, []int{0, 1, 2}, "")
		assert.Equal(t, p.Childs[0].Childs[1].Data, []int{4, 5}, "")
		assert.Equal(t, p.Childs[0].Childs[2].Data, []int{7, 8}, "")
		assert.Equal(t, p.Childs[1].Childs[0].Data, []int{10, 11, 12}, "")
		assert.Equal(t, p.Childs[1].Childs[1].Data, []int{14, 15}, "")
		assert.Equal(t, p.Childs[1].Childs[2].Data, []int{17, 18}, "")
		assert.Equal(t, p.Childs[1].Childs[3].Data, []int{20, 21}, "")

		assert.Equal(t, splitRecords, []splitRecord{
			// first split
			{9, 0, p},                                  //pivot
			{16, 0, p.Childs[1]}, {19, 1, p.Childs[1]}, // right node
			// second split
			{13, 0, p.Childs[1]}, //pivot
			{16, 1, p.Childs[1]},
			{19, 2, p.Childs[1]},

			{14, 0, p.Childs[1].Childs[1]}, {15, 1, p.Childs[1].Childs[1]}, // right node

			// insertion
			{10, 0, p.Childs[1].Childs[0]},
			{11, 1, p.Childs[1].Childs[0]},
			{12, 2, p.Childs[1].Childs[0]},
		}, "")
	}
}

func TestInsertKeepsOrder(t *testing.T) {
	sut := buildBtreeWith(t, 9, randomIntData)

	res := make([]int, 0, len(randomIntData))
	res = Allocate(sut, res)

	assert.Equal(t, len(randomIntData), len(res), "Some elements were lost")
	for i, v := range res {
		assert.Equal(t, i+1, int(v), "Elements not in order")
		if i+1 != int(v) {
			break
		}
	}

	t.Log("----------------")
	t.Log(res)
}
