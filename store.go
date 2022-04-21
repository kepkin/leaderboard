package main

import (
	"golang.org/x/exp/constraints"
)

type Store[K constraints.Ordered, V constraints.Ordered] struct {
	btree   *Node[K, V]
	pkbtree *Node[V, *Iter[K, V]]
}

func NewStore[K constraints.Ordered, V constraints.Ordered]() *Store[K, V] {
	res := Store[K, V]{}

	res.btree = NewNode[K, V](5, nil, res.onSplit)
	res.pkbtree = NewNode[V, *Iter[K, V]](5, nil, nil)

	return &res
}

func (s *Store[K, V]) onSplit(value BTreeLeaf[K, V], itr *Iter[K, V]) {
	pkbtreeValue := BTreeLeaf[V, *Iter[K, V]]{OrderKey: value.Value, Value: itr}
	s.pkbtree, _ = s.pkbtree.Upsert(pkbtreeValue)
}

func (s *Store[K, V]) Insert(value BTreeLeaf[K, V]) ([]V, error) {
	var iter *Iter[K, V] = nil
	s.btree, iter = s.btree.Insert(value)

	pkbtreeValue := BTreeLeaf[V, *Iter[K, V]]{OrderKey: value.Value, Value: iter}
	s.pkbtree, _ = s.pkbtree.Upsert(pkbtreeValue)

	return []V{}, nil
}

func (s *Store[K, V]) Get(key V) *Iter[K, V] {
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
