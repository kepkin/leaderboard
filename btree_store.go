package leaderboard

import (
	"fmt"
	"sync"
)

type pkPair[K any, V any] struct {
	btree *Node[K, V]
	idx   int
}

func pkPairEquals[K any, V any](a, b pkPair[K, V]) bool {
	return a.btree == b.btree && a.idx == b.idx
}

type BtreeStore[K any, V any] struct {
	btree   *Node[K, V]
	pkbtree *Node[V, pkPair[K, V]]

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

	res.btree = NewNode[K, V](
		101,
		keyLessFunc,
		keyEqualsFunc,
		valEqualsFunc,
		res.onSplit,
	)
	res.pkbtree = NewNode[V, pkPair[K, V]](
		5,
		valLessFunc,
		valEqualsFunc,
		pkPairEquals[K, V],
		nil)

	return &res
}

func (s *BtreeStore[K, V]) onSplit(value BTreeLeaf[K, V], btree *Node[K, V], idx int) {
	pkbtreeValue := BTreeLeaf[V, pkPair[K, V]]{OrderKey: value.Value, Value: pkPair[K, V]{btree: btree, idx: idx}}
	s.pkbtree, _ = s.pkbtree.Upsert(pkbtreeValue)

	// if s.logSplit {
	// if btree == nil {
	// 	fmt.Printf("was just removed %v\n", value)
	// } else {
	// 	fmt.Printf("was just updated %v %v\n", value, idx)
	// }
	// }
}

func (s *BtreeStore[K, V]) Upsert(value BTreeLeaf[K, V]) (*Iter[K, V], error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	pkbtreeValue := BTreeLeaf[V, pkPair[K, V]]{OrderKey: value.Value, Value: pkPair[K, V]{}}
	localNode, idx := s.pkbtree.Find(pkbtreeValue)
	if localNode != nil {
		node := localNode.At(idx).Value.btree
		nodeIdx := localNode.At(idx).Value.idx
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
			return nil, fmt.Errorf("can not remove %v expecting on %v", value.Value, idx)
		}
	}

	a, b := s.Insert(value)
	s.logSplit = false
	// fmt.Println("------------------------------------")

	return a, b
}

func (s *BtreeStore[K, V]) Insert(value BTreeLeaf[K, V]) (*Iter[K, V], error) {
	// s.mu.Lock()
	// defer s.mu.Unlock()

	var err error
	// if s.logSplit {
	// 	fmt.Printf("inserting %v\n", value)
	// }
	s.btree, _ = s.btree.Insert(value)

	return nil, err
}

func (s *BtreeStore[K, V]) Get(key V) *Iter[K, V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	pkbtreeValue := BTreeLeaf[V, pkPair[K, V]]{OrderKey: key, Value: pkPair[K, V]{}}
	localNode, idx := s.pkbtree.Find(pkbtreeValue)
	if localNode == nil {
		return nil
	}

	node := localNode.Data[idx].Value.btree
	nodeIdx := localNode.Data[idx].Value.idx

	return node.Iter(nodeIdx)
	//TODO:
	return nil
	// return iter.Value().Value.btree.Iter(uint16(iter.Value().Value.idx))
}

func (s *BtreeStore[K, V]) DepthFirstTraverse(visiter func(v BTreeLeaf[K, V])) {
	s.btree.DepthFirstTraverse(visiter)
}

func (s *BtreeStore[K, V]) String() string {
	return s.btree.String()
}

func (s *BtreeStore[K, V]) GetAdjacent(itr *Iter[K, V], visiter func(BTreeLeaf[K, V]), before int, after int) {
	backItr := itr
	for ; before > 0 && backItr != nil; before -= 1 {
		backItr = backItr.Prev()
	}

	for i := 0; i < after && itr != nil; i += 1 {
		visiter(itr.Value())
		itr = itr.Next()
	}
}
