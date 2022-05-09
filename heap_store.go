package leaderboard

import (
	_ "container/ring"
	"fmt"
	"sync"
	// "golang.org/x/exp/constraints"
)

type HeapStore[K any, V comparable] struct {
	lHeap   ScoreHeap[K, V]
	hHeap   ScoreHeap[K, V]
	pkLheap map[V]int
	pkHheap map[V]int

	mu sync.RWMutex
}

func notFunc[K any](f func(a, b K) bool) func(a, b K) bool {
	return func(a, b K) bool { return !f(a, b) }
}

func NewHeapStore[K any, V comparable](lessFunc func(a, b K) bool) *HeapStore[K, V] {
	res := HeapStore[K, V]{}

	res.lHeap = ScoreHeap[K, V]{LessFunc: lessFunc, OnSplit: res.onLheapSplit}
	res.hHeap = ScoreHeap[K, V]{LessFunc: notFunc[K](lessFunc), OnSplit: res.onHheapSplit}
	res.pkLheap = make(map[V]int)
	res.pkHheap = make(map[V]int)

	return &res
}

func (s *HeapStore[K, V]) onLheapSplit(v Item[K, V]) {
	s.pkLheap[v.Value] = v.Index
}

func (s *HeapStore[K, V]) onHheapSplit(v Item[K, V]) {
	s.pkHheap[v.Value] = v.Index
}

func (s *HeapStore[K, V]) Update(order K, value V) ([]Pair[K, V], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	lidx, ok := s.pkLheap[value]
	if !ok {
		return nil, fmt.Errorf("no such value")
	}

	hidx, ok := s.pkHheap[value]
	if !ok {
		return nil, fmt.Errorf("pkLheap & pkHheap out of sync")
	}

	s.lHeap.Update(lidx, order)
	s.hHeap.Update(hidx, order)

	//TODO
	return []Pair[K, V]{}, nil
}

func (s *HeapStore[K, V]) Insert(order K, value V) []Pair[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pkLheap[value] = s.lHeap.Len()
	s.lHeap.Push(Item[K, V]{OrderKey: order, Value: value})

	s.pkHheap[value] = s.hHeap.Len()
	s.hHeap.Push(Item[K, V]{OrderKey: order, Value: value})

	s.lHeap.Update(s.pkLheap[value], order)
	s.hHeap.Update(s.pkHheap[value], order)

	cnt := 10

	data := make([]Pair[K, V], 0, cnt*2+1)
	for i := s.pkLheap[value]; i > 0 && cnt > 0; {
		cnt--
		item, ok := s.lHeap.GetParent(i)
		if !ok {
			break
		}
		i = item.Index
		data = append(data, Pair[K, V]{item.OrderKey, item.Value})
	}

	data = append(data, Pair[K, V]{order, value})

	cnt = 10
	for i := s.pkHheap[value]; i > 0 && cnt > 0; {
		cnt--
		item, ok := s.hHeap.GetParent(i)
		if !ok {
			break
		}
		i = item.Index
		data = append(data, Pair[K, V]{item.OrderKey, item.Value})
	}

	return data
}
