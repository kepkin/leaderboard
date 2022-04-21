package main

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	sut := buildBtreeWith(t, randomIntData)

	res := make([]int, 0, len(randomIntData))
	res = AllocateV(sut, res)

	assert.Equal(t, len(randomIntData), len(res), "Some elements were lost")
	for i, v := range res {
		assert.Equal(t, i+1, v, "Elements not in order")
		if i+1 != v {
			break
		}
	}

	t.Log("----------------")
	t.Log(res)
}
