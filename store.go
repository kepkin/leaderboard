package main

import (
	"golang.org/x/exp/constraints"
)

type Store[K constraints.Ordered, V constraints.Ordered] struct {
	btree   *Node[K, V]
	pkbtree *Node[V, *Node[K, V]]
}

func NewStore[K constraints.Ordered, V constraints.Ordered]() *Store[K, V] {
	res := Store[K, V]{}

	res.btree = NewNode[K, V](5, nil, res.onSplit)
	res.pkbtree = NewNode[V, *Node[K, V]](5, nil, nil)

	return &res
}

func (s *Store[K, V]) onSplit(value BTreeLeaf[K, V], n *Node[K, V]) {
	pkbtreeValue := BTreeLeaf[V, *Node[K, V]]{OrderKey: value.Value, Value: n}
	s.pkbtree, _ = s.pkbtree.Upsert(pkbtreeValue)
}

func (s *Store[K, V]) Insert(value BTreeLeaf[K, V]) error {
	var n *Node[K, V] = nil
	s.btree, n = s.btree.Insert(value)

	pkbtreeValue := BTreeLeaf[V, *Node[K, V]]{OrderKey: value.Value, Value: n}
	s.pkbtree, _ = s.pkbtree.Upsert(pkbtreeValue)

	return nil
}

func (s *Store[K, V]) Get(key V) BTreeLeaf[K, V] {
	pkbtreeValue := BTreeLeaf[V, *Node[K, V]]{OrderKey: key, Value: nil}
	node, idx, ok := s.pkbtree.Find(pkbtreeValue)
	if node == nil || !ok {
		return BTreeLeaf[K, V]{}
	}

	res, _ := node.Data[idx].Value.LocalFindByValue(key)
	return res
}

func (s *Store[K, V]) DepthFirstTraverse(visiter func(v BTreeLeaf[K, V])) {
	s.btree.DepthFirstTraverse(visiter)
}

func (s *Store[K, V]) String() string {
	return s.btree.String()
}
