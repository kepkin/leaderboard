package leaderboard

import (
	"container/heap"
	// "golang.org/x/exp/constraints"
)

// An Item is something we manage in a priority queue.
type Item[K any, V any] struct {
	OrderKey K
	Value    V // The OrderKey of the item; arbitrary.
	Index    int
}

// A ScoreHeap implements heap.Interface and holds Items.
type ScoreHeap[K any, V any] struct {
	OnSplit func(it Item[K, V])
	d       []Item[K, V]

	LessFunc func(a K, b K) bool
}

func (pq ScoreHeap[K, V]) Len() int { return len(pq.d) }

func (pq ScoreHeap[K, V]) Less(i, j int) bool {
	// return pq.d[i].OrderKey.Less(pq.d[j].OrderKey)
	return pq.LessFunc(pq.d[i].OrderKey, pq.d[j].OrderKey)
}

func (pq ScoreHeap[K, V]) Swap(i, j int) {
	pq.d[i], pq.d[j] = pq.d[j], pq.d[i]
	pq.d[i].Index = i
	pq.d[j].Index = j
	pq.OnSplit(pq.d[i])
	pq.OnSplit(pq.d[j])
}

func (pq *ScoreHeap[K, V]) Push(x any) {
	n := len((*pq).d)
	(*pq).d = append((*pq).d, x.(Item[K, V]))
	(*pq).d[n].Index = n
}

func (pq *ScoreHeap[K, V]) Pop() any {
	old := (*pq).d
	n := len(old)
	item := old[n-1]
	(*pq).d = old[0 : n-1]
	return item
}

func (pq ScoreHeap[K, V]) GetParent(j int) (Item[K, V], bool) {
	i := (j - 1) / 2 // parent
	if i == j {
		return Item[K, V]{}, false
	}
	return pq.d[i], true
}

// update modifies the OrderKey of an Item in the queue.
func (pq *ScoreHeap[K, V]) Update(idx int, orderKey K) {
	(*pq).d[idx].OrderKey = orderKey
	heap.Fix(pq, idx)
}
