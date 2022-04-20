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
