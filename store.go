package leaderboard

import (
	"fmt"
	"sync"
	// "golang.org/x/exp/constraints"
)

type Store[K Ordered, V Ordered] struct {
	btree   *Node[K, V]
	pkbtree *Node[V, *Iter[K, V]]

	mu sync.RWMutex
	pkMu sync.RWMutex
}

func NewStore[K Ordered, V Ordered]() *Store[K, V] {
	res := Store[K, V]{}

	res.btree = NewNode[K, V](101, nil, res.onSplit)
	res.pkbtree = NewNode[V, *Iter[K, V]](5, nil, nil)

	return &res
}

func (s *Store[K, V]) onSplit(value BTreeLeaf[K, V], itr *Iter[K, V]) {
	// pkbtreeValue := BTreeLeaf[V, *Iter[K, V]]{OrderKey: value.Value, Value: itr}
	// s.pkbtree, _ = s.pkbtree.Upsert(pkbtreeValue)
}

func (s *Store[K, V]) Insert(value BTreeLeaf[K, V]) (*Iter[K, V], error) {
    needRebalance := false
	var iter *Iter[K, V] = nil
	var err error



	for iter == nil {
		s.btree, iter, needRebalance = s.btree.Insert(value, false)
		if needRebalance {
			s.btree, iter, needRebalance = s.btree.Insert(value, true)	
		}
	}

	// iter, err, needRebalance = s.insertLocking(value, false)
	// for needRebalance {
	// 	if iter != nil {
	// 		iter.Close()	
	// 	}
		
	// 	needRebalance = false
	// 	iter, err, needRebalance = s.insertLocking(value, true)
	// }

	// err, needRebalance = s.insertPKLocking(value.Value, iter, false)
	// for needRebalance {
	// 	needRebalance = false
	// 	err, needRebalance = s.insertPKLocking(value.Value, iter, true)	
	// }

	return iter, err
}


func (s *Store[K, V]) insertLocking(value BTreeLeaf[K, V], rebalance bool) (*Iter[K, V], error, bool) {
	// if rebalance {
	// 	s.mu.Lock()
	// 	defer s.mu.Unlock()
	// } else {
	// 	s.mu.RLock()
	//     defer s.mu.RUnlock()
	// }

    needRebalance := false
	var iter *Iter[K, V] = nil

	s.btree, iter, needRebalance = s.btree.Insert(value, rebalance)
	if iter == nil && needRebalance {
		return nil, nil, needRebalance
	} else if iter == nil {
		return nil, fmt.Errorf("failed to insert value"), needRebalance
	}

	return iter, nil, needRebalance
}

func (s *Store[K, V]) insertPKLocking(pkey V, itr *Iter[K, V], rebalance bool) (error, bool) {
	if rebalance {
		s.pkMu.Lock()
		defer s.pkMu.Unlock()
	} else {
		s.pkMu.RLock()
	    defer s.pkMu.RUnlock()
	}

    needRebalance := false
	pkbtreeValue := BTreeLeaf[V, *Iter[K, V]]{OrderKey: pkey, Value: itr}
	s.pkbtree, _, needRebalance = s.pkbtree.Insert(pkbtreeValue, rebalance)

	return nil, needRebalance
}

func (s *Store[K, V]) Get(key V) *Iter[K, V] {
	s.mu.Lock()
    defer s.mu.Unlock()

	pkbtreeValue := BTreeLeaf[V, *Iter[K, V]]{OrderKey: key, Value: nil}
	iter := s.pkbtree.Find(pkbtreeValue)
	if iter == nil {
		return nil
	}

	return iter.Value().Value
}

func (s *Store[K, V]) DepthFirstTraverse(visiter func(v BTreeLeaf[K, V])) {
	s.btree.DepthFirstTraverse(visiter)
}

func (s *Store[K, V]) String() string {
	return s.btree.String()
}

func (s *Store[K, V]) GetAdjacent(itr *Iter[K,V], visiter func(BTreeLeaf[K, V]), before int, after int) {
	backItr := itr
	for ; before > 0 && backItr != nil; before -= 1 {
		backItr = backItr.Prev()
	}

	for i := 0; i < after && itr != nil; i+=1 {
		visiter(itr.Value())
		itr = itr.Next()
	}
}