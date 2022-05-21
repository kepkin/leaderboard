package leaderboard

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func btreeCompareChild[K any](t *testing.T, n *Node[K], expected [][]K) bool {
	fmt.Print(" | ")

	for _, v := range n.Data {
		fmt.Print(v, ", ")
	}

	fmt.Print(" | ")
	return false
}

func btreeCompare[K any](t *testing.T, n *Node[K], expected [][]K) bool {
	btreeCompareChild(t, n, expected)
	fmt.Println()

	btreeCompareH(t, n, expected)
	return false
}

func btreeCompareH[K any](t *testing.T, n *Node[K], expected [][]K) bool {
	if n == nil {
		fmt.Print("{}, ")
		return false
	}

	for _, v := range n.Childs {
		btreeCompareChild[K](t, v, nil)
	}

	for _, v := range n.Childs {
		btreeCompareH[K](t, v, nil)
	}

	return false
}

func buildNodeWithData(size int, onSplit OnSplitTrigger[int], data ...int) *Node[int] {
	r := NewNode[int](
		size,
		StdLess[int],
		StdEquals[int],
		onSplit,
	)
	for _, v := range data {
		r.Data = append(r.Data, v)
	}

	return r
}

func buildChildWithData(parent *Node[int], pidx int, data ...int) {
	r := newTreeNode[int](
		parent,
		parent,
	)
	r.Pidx = pidx

	parent.Childs[pidx] = r
	for _, v := range data {
		r.Data = append(r.Data, v)
	}
}

// random seq from 1 to 48
var randomIntData = []int{
	10, 32, 45, 34, 26, 16, 4, 40, 22, 21, 29, 20, 24, 12, 6, 15, 27, 1, 43, 44, 17, 46, 3, 8, 30, 35, 41, 18, 47, 42, 13, 36, 7, 9, 28, 25, 48, 5, 14, 19, 31, 23, 11, 38, 33, 37, 2, 39,
}

func buildBtreeWith(t *testing.T, size int, data []int) *Node[int] {
	sut := NewNode[int](
		size,
		StdLess[int],
		StdEquals[int],
		nil,
	)

	logData := make([]int, 0, 49)

	for _, v := range randomIntData {
		var itr *Iter[int]

		sut, _ = sut.Insert(v)

		logData = logData[:0]

		itr = sut.Begin()
		for ; itr.Valid(); itr = itr.Next() {
			logData = append(logData, itr.Value())
		}

		itr.Close()

		// logData = AllocateV(sut, logData)
		t.Logf("insert %v (len: %v): %v", v, len(logData), logData)
	}

	return sut
}

func assertElementsMatch[T any](t *testing.T, listA, listB []T, eqFunc func(a, b T) bool) {
	t.Helper()

	extraA, extraB := diffLists(listA, listB, eqFunc)
	if len(extraA) == 0 && len(extraB) == 0 {
		return
	}

	assert.ElementsMatch(t, listA, listB)
}

// this is a copy paste from assert package with generics variant
func diffLists[T any](listA, listB []T, eqFunc func(a, b T) bool) (extraA, extraB []T) {
	aLen := len(listA)
	bLen := len(listB)

	// Mark indexes in bValue that we already used
	visited := make([]bool, bLen)
	for i := 0; i < aLen; i++ {
		element := listA[i]
		found := false
		for j := 0; j < bLen; j++ {
			if visited[j] {
				continue
			}
			if eqFunc(element, listB[j]) {
				visited[j] = true
				found = true
				break
			}
		}
		if !found {
			extraA = append(extraA, element)
		}
	}

	for j := 0; j < bLen; j++ {
		if visited[j] {
			continue
		}
		extraB = append(extraB, listB[j])
	}

	return
}
