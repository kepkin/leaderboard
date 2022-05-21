package leaderboard

import (
	"sync"
)

func (n *Node[K]) Begin() *Iter[K] {
	if n.Childs[0] != nil {
		return n.Childs[0].Begin()
	}

	return newIter[K](n, 0)
}

func newIter[K any](n *Node[K], idx int) *Iter[K] {
	//TODO
	//n.treeRebalanceMu.RLock()
	return &Iter[K]{n, idx, IterStateValid, &n.s.treeRebalanceMu}
}

type iterState int

const (
	IterStateValid iterState = iota
	IterStateEnd
	IterStateStart
)

type Iter[K any] struct {
	n *Node[K]
	i int

	state iterState

	treeRebalanceMu *sync.RWMutex
}

func (it *Iter[K]) Value() K {
	return it.n.Data[it.I()]
}

func (it *Iter[K]) I() int {
	return it.i
}

func (it *Iter[K]) Valid() bool {
	return it.state == IterStateValid
}

func (it *Iter[K]) Close() {
	//TODO
	//it.treeRebalanceMu.RUnlock()
}

func IterEquals[K any, V any](a, b Iter[K]) bool {
	return a.n == b.n && a.i == b.i
}

func (it *Iter[K]) Prev() *Iter[K] {
	if it.state == IterStateStart {
		return it
	} else if it.state == IterStateStart {
		it.state = IterStateValid
		return it
	}

	if it.I() > 0 {
		it.i -= 1
		for it.n.Childs[it.I()+1] != nil {
			it.n = it.n.Childs[it.I()+1]
			it.i = len(it.n.Data) - 1
		}

		return it
	}

	if it.n.Childs[0] != nil {
		it.i = len(it.n.Childs[0].Data) - 1
		it.n = it.n.Childs[0]
		return it
	}

	for it.n.Parent != nil && it.I() == 0 {
		for it.n.Pidx == 0 {
			it.i = it.n.Pidx
			it.n = it.n.Parent

			if it.n.Parent == nil && it.i == 0 {
				it.state = IterStateStart
				return it
			}
		}
		it.i = it.n.Pidx - 1
		it.n = it.n.Parent

		return it
	}

	it.state = IterStateStart
	return it
}

func (it *Iter[K]) Next() *Iter[K] {
	if it.state == IterStateEnd {
		return it
	} else if it.state == IterStateStart {
		it.state = IterStateValid
		return it
	}

	if it.I() < len(it.n.Data)-1 {
		it.i += 1
		for it.n.Childs[it.I()] != nil {
			it.n = it.n.Childs[it.I()]
			it.i = 0
		}

		return it
	}

	if it.n.Childs[len(it.n.Data)] != nil {
		it.i = 0
		it.n = it.n.Childs[len(it.n.Data)]
		return it
	}

	for it.n.Parent != nil && it.I() == len(it.n.Data)-1 {
		for it.n.Pidx == len(it.n.Parent.Data) {
			it.i = it.n.Pidx
			it.n = it.n.Parent

			if it.n.Parent == nil && it.i == len(it.n.Data) {
				it.state = IterStateEnd
				return it
			}

		}
		it.i = it.n.Pidx
		it.n = it.n.Parent

		return it
	}

	it.state = IterStateEnd
	return it
}
