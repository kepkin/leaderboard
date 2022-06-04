package leaderboard

import (
	"fmt"
	"sync"
)

type Tuple[K any, V any] struct {
	Key K
	Val V
}

type pkPair[K any] struct {
	btree *Node[K]
	idx   int
}

func pkPairEquals[K any, V any](a, b pkPair[Tuple[K, V]]) bool {
	return a.btree == b.btree && a.idx == b.idx
}

type BtreeStore[K any, V any] struct {
	btree   *Node[Tuple[K, V]]
	pkbtree *Node[Tuple[V, pkPair[Tuple[K, V]]]]

	mu   sync.RWMutex
	pkMu sync.RWMutex

	logSplit bool
}

func NewBtreeStore[K any, V any](
	keyLessFunc LessFuncType[K],
	keyEqualsFunc EqualsFuncType[K],
	valLessFunc LessFuncType[V],
	valEqualsFunc EqualsFuncType[V],
) *BtreeStore[K, V] {
	res := BtreeStore[K, V]{}

	res.btree = NewNode(
		101,
		func(a, b Tuple[K, V]) bool { return keyLessFunc(a.Key, b.Key) },
		func(a, b Tuple[K, V]) bool { return keyEqualsFunc(a.Key, b.Key) },
		res.onSplit,
	)
	res.pkbtree = NewNode(
		101,
		func(a, b Tuple[V, pkPair[Tuple[K, V]]]) bool { return valLessFunc(a.Key, b.Key) },
		func(a, b Tuple[V, pkPair[Tuple[K, V]]]) bool { return valEqualsFunc(a.Key, b.Key) },
		nil)

	return &res
}

func (s *BtreeStore[K, V]) onSplit(value Tuple[K, V], btree *Node[Tuple[K, V]], idx int) {
	pkbtreeValue := Tuple[V, pkPair[Tuple[K, V]]]{Key: value.Val, Val: pkPair[Tuple[K, V]]{btree: btree, idx: idx}}
	s.pkbtree, _ = s.pkbtree.Upsert(pkbtreeValue)

	// if s.logSplit {
	// if btree == nil {
	// 	fmt.Printf("was just removed %v\n", value)
	// } else {
	// 	fmt.Printf("was just updated %v %v\n", value, idx)
	// }
	// }
}

func (s *BtreeStore[K, V]) Upsert(value Tuple[K, V]) (*Iter[Tuple[K, V]], error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	pkbtreeValue := Tuple[V, pkPair[Tuple[K, V]]]{Key: value.Val, Val: pkPair[Tuple[K, V]]{}}
	localNode, idx := s.pkbtree.Find(pkbtreeValue)
	if localNode != nil {
		node := localNode.At(idx).Val.btree
		nodeIdx := localNode.At(idx).Val.idx
		// node := iter.Value().Value.btree
		// idx := uint16(iter.Value().Value.idx)
		if node == nil {
			panic(fmt.Errorf("nil to remove %v", value))
		}

		// fmt.Println("------------------------------------")
		// fmt.Printf("tr removing %v\n", value)
		s.logSplit = true
		ok := node.RemoveByLocalIdx(nodeIdx)
		if !ok {
			panic("oh")
			return nil, fmt.Errorf("can not remove %v expecting on %v", value.Val, idx)
		}
	}

	a, b := s.Insert(value)
	s.logSplit = false
	// fmt.Println("------------------------------------")

	return a, b
}

func (s *BtreeStore[K, V]) Insert(value Tuple[K, V]) (*Iter[Tuple[K, V]], error) {
	// s.mu.Lock()
	// defer s.mu.Unlock()

	var err error
	// if s.logSplit {
	// 	fmt.Printf("inserting %v\n", value)
	// }
	s.btree, _ = s.btree.Insert(value)

	return nil, err
}

func (s *BtreeStore[K, V]) Get(key V) *Iter[Tuple[K, V]] {
	s.mu.Lock()
	defer s.mu.Unlock()

	pkbtreeValue := Tuple[V, pkPair[Tuple[K, V]]]{Key: key, Val: pkPair[Tuple[K, V]]{}}
	localNode, idx := s.pkbtree.Find(pkbtreeValue)
	if localNode == nil {
		return nil
	}

	node := localNode.Data[idx].Val.btree
	nodeIdx := localNode.Data[idx].Val.idx

	return node.Iter(nodeIdx)
	//TODO:
	// return nil
	// return iter.Value().Value.btree.Iter(uint16(iter.Value().Value.idx))
}

// func (s *BtreeStore[K, V]) DepthFirstTraverse(visiter func(v Tuple[K, V])) {
// 	s.btree.DepthFirstTraverse(visiter)
// }

func (s *BtreeStore[K, V]) String() string {
	return s.btree.String()
}

// func (s *BtreeStore[K, V]) GetAdjacent(itr *Iter[Tuple[K, V]], visiter func(Tuple[K, V]), before int, after int) {
// 	backItr := itr
// 	for ; before > 0 && backItr != nil; before -= 1 {
// 		backItr = backItr.Prev()
// 	}

// 	for i := 0; i < after && itr != nil; i += 1 {
// 		visiter(itr.Value())
// 		itr = itr.Next()
// 	}
// }
