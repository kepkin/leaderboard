package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

func Allocate[K constraints.Ordered, V comparable](tree *Node[K, V], dst []BTreeLeaf[K, V]) []BTreeLeaf[K, V] {
	visiter := func(v BTreeLeaf[K, V]) {
		dst = append(dst, v)
	}
	tree.DepthFirstTraverse(visiter)

	return dst
}

type AllocateVV[K constraints.Ordered, V comparable] struct {
	Data []V
}

func (a *AllocateVV[K, V]) visit(v BTreeLeaf[K, V]) {
	a.Data = append(a.Data, v.Value)
}

func AllocateV[K constraints.Ordered, V comparable](tree *Node[K, V], dst []V) []V {
	visiter := func(v BTreeLeaf[K, V]) {
		dst = append(dst, v.Value)
	}
	tree.DepthFirstTraverse(visiter)
	return dst
}

func PrintTree[K constraints.Ordered, V comparable](tree *Node[K, V]) {
	data := make([]K, 0, 100)
	visiter := func(v BTreeLeaf[K, V]) {
		data = append(data, v.OrderKey)
	}
	tree.DepthFirstTraverse(visiter)
	fmt.Println(data)
}
