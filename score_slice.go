package leaderboard

import (
// "golang.org/x/exp/constraints"
)

type ScoreList[K any] struct {
	OnChange func(it K, idx int)
	d        []K

	LessFunc func(a K, b K) bool
}

func (pq ScoreList[K]) Len() int { return len(pq.d) }

func (pq ScoreList[K]) Less(i, j int) bool {
	return pq.LessFunc(pq.d[i], pq.d[j])
}

func (pq *ScoreList[K]) Push(x any) int {
	n := len((*pq).d)
	(*pq).d = append((*pq).d, x.(K))
	return n
}

func (pq *ScoreList[K]) Fix(i int) int {
	for ; i > 0; i-- {
		if pq.Less(i-1, i) {
			pq.OnChange(pq.d[i], i)
			return i
		}

		pq.d[i], pq.d[i-1] = pq.d[i-1], pq.d[i]
		pq.OnChange(pq.d[i], i)

	}

	return 0
}

// update modifies the priority and OrderKey of an Item in the queue.
func (pq *ScoreList[K]) Update(idx int, v K) {
	(*pq).d[idx] = v
	pq.Fix(idx)
}

func (pq ScoreList[K]) Get(idx int) K {
	return pq.d[idx]
}
