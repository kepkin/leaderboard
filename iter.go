package leaderboard

import (
	"sync"
	// "golang.org/x/exp/constraints"
)

func (n *Node[K, V]) Begin() *Iter[K, V] {
	if n.Childs[0] != nil {
		return n.Childs[0].Begin()
	}

	return newIter[K, V](n, 0)
}

func newIter[K Ordered, V Comparable](n *Node[K,V], idx int) *Iter[K, V] {
	n.treeRebalanceMu.RLock()
	return &Iter[K, V]{n, idx, IterStateValid, n.treeRebalanceMu}
}

type iterState int

const (
	IterStateValid iterState = iota
	IterStateEnd
	IterStateStart
)

type Iter[K Ordered, V Comparable] struct {
	n *Node[K, V]
	i int

	state iterState


	treeRebalanceMu *sync.RWMutex
}

func (it *Iter[K, V]) Value() BTreeLeaf[K, V] {
	return it.n.Data[it.I()]
}

func (it *Iter[K, V]) I() int {
	return it.i
}

func (it *Iter[K, V]) Valid() bool {
	return it.state == IterStateValid
}

func (it *Iter[K, V]) Close() {
	it.treeRebalanceMu.RUnlock()
}

func (it *Iter[K, V]) Equals(other Comparable) bool {
	otherItr := other.(*Iter[K,V])
	return it.n == otherItr.n && it.i == otherItr.i
}

func (it *Iter[K, V]) Prev() *Iter[K, V] {
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
		if it.n.Pidx == 0 {
			it.state = IterStateStart
			return it
		}
		it.i = it.n.Pidx - 1
		it.n = it.n.Parent

		return it
	}

	it.state = IterStateStart
	return it
}

func (it *Iter[K, V]) Next() *Iter[K, V] {
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
		if it.n.Pidx == len(it.n.Parent.Data) {
			it.state = IterStateEnd
			return it
		}
		it.i = it.n.Pidx
		it.n = it.n.Parent

		return it
	}

	it.state = IterStateEnd
	return it
}
