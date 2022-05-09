package leaderboard

import (
	_ "container/ring"
	"sort"
	"sync"
	// "golang.org/x/exp/constraints"
)

type Pair[K any, V any] struct {
	A K
	B V
}

type ScoreSliceStore[K any, V comparable] struct {
	score ScoreList[Pair[K, V]]
	pk    map[V]int

	lessFunc func(a, b K) bool

	mu sync.RWMutex
}

func NewScoreSliceStore[K any, V comparable](lessFunc func(a, b K) bool) ScoreSliceStore[K, V] {
	res := ScoreSliceStore[K, V]{}

	res.lessFunc = lessFunc
	res.score = ScoreList[Pair[K, V]]{
		LessFunc: func(a, b Pair[K, V]) bool { return lessFunc(a.A, b.A) },
		OnChange: res.onChange,
	}
	res.pk = make(map[V]int)

	return res
}

func (s *ScoreSliceStore[K, V]) PreInitInsert(order K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.score.Push(Pair[K, V]{order, value})
}

func (s *ScoreSliceStore[K, V]) Init() {
	sort.Slice(s.score.d, func(i, j int) bool {
		return s.lessFunc(s.score.d[i].A, s.score.d[j].A)
	})
}

func (s *ScoreSliceStore[K, V]) onChange(v Pair[K, V], idx int) {
	s.pk[v.B] = idx
}

func (s *ScoreSliceStore[K, V]) Insert(order K, value V) []Pair[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := s.score.Push(Pair[K, V]{order, value})
	s.score.Fix(idx)

	cnt := 10
	data := make([]Pair[K, V], 0, cnt*2+1)
	for i := idx - 1; i > 0 && cnt > 0; i-- {
		data = append(data, s.score.Get(i))
		cnt--
	}

	cnt = 11
	for i := idx; i < s.score.Len() && cnt > 0; i++ {
		data = append(data, s.score.Get(i))
		cnt--
	}

	return data
}
