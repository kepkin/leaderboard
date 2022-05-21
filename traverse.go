package leaderboard

import (
	"fmt"
	// "golang.org/x/exp/constraints"
)

func Allocate[K any](tree *Node[K], dst []K) []K {
	visiter := func(v K) {
		dst = append(dst, v)
	}
	tree.DepthFirstTraverse(visiter)

	return dst
}

func PrintTree[K any](tree *Node[K]) {
	data := make([]K, 0, 100)
	visiter := func(v K) {
		data = append(data, v)
	}
	tree.DepthFirstTraverse(visiter)
	fmt.Println(data)
}
