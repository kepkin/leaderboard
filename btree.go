package leaderboard

import (
	"fmt"

	// "golang.org/x/exp/constraints"
)


type Comparable interface {
	Equals(other Comparable) bool
}

type Ordered interface {
	Comparable
	Less(than Ordered) bool
}


type OnSplitTrigger[K Ordered, V Comparable] func(BTreeLeaf[K, V], *Iter[K, V])

func DummyOnSplit[K Ordered, V Comparable](BTreeLeaf[K, V], *Iter[K, V]) {}

type BTreeLeaf[K Ordered, V Comparable] struct {
	OrderKey K
	Value    V
}

type Node[K Ordered, V Comparable] struct {
	Parent *Node[K, V]
	Pidx   int
	Size   uint16
	Data   []BTreeLeaf[K, V]
	Childs []*Node[K, V]

	OnSplit OnSplitTrigger[K, V]
}

func NewNode[K Ordered, V Comparable](size uint16, parent *Node[K, V], onSplit OnSplitTrigger[K, V]) *Node[K, V] {

	if onSplit == nil {
		onSplit = DummyOnSplit[K, V]
	}

	return &Node[K, V]{
		Parent:  parent,
		Size:    size,
		Data:    make([]BTreeLeaf[K, V], 0, size),
		Childs:  make([]*Node[K, V], size+1),
		OnSplit: onSplit,
	}
}

// We should be there only if
//  - newValue.OrderKey < Data[idx].OrderKey
//  - Childs[idx] == nil
//  - Childs[:last:] == nil
//  - len(Data) < Size
func (n *Node[K, V]) insertAt(idx uint16, newValue BTreeLeaf[K, V]) {
	n.Data = append(n.Data, BTreeLeaf[K, V]{})
	copy(n.Data[idx+1:], n.Data[idx:])
	n.Data[idx] = newValue

	if n.Childs[idx+1] != nil {
		for i := int(idx) + 1; i < len(n.Childs); i += 1 {
			if n.Childs[i] != nil {
				n.Childs[i].Pidx += 1
			}
		}

		copy(n.Childs[idx+2:], n.Childs[idx+1:])
		n.Childs[idx+1] = nil

	}
}

func (n *Node[K, V]) Split() *Node[K, V] {

	upperNode := n.Parent
	returnNode := n
	leftNode := n

	if upperNode == nil {

		upperNode = NewNode(n.Size, nil, n.OnSplit)
		leftNode.Parent = upperNode
		returnNode = upperNode

		upperNode.Childs[0] = leftNode
	} else {
		if uint16(len(upperNode.Data)) == n.Size {
			upperNode = upperNode.Split()
		}

		// leftNode.Parent = upperNode
	}

	rightNode := NewNode(n.Size, upperNode, n.OnSplit)

	pivotValue := leftNode.Data[n.Size/2]

	idx := upperNode.FindIdxToInsert(pivotValue)
	leftNode.Pidx = int(idx)
	rightNode.Pidx = int(idx) + 1
	upperNode.insertAt(idx, pivotValue)
	upperNode.Childs[idx+1] = rightNode

	n.OnSplit(pivotValue, &Iter[K,V]{upperNode, int(idx)})

	rightNode.Data = append(rightNode.Data, leftNode.Data[n.Size/2+1:]...)
	for i, v := range rightNode.Data {
		n.OnSplit(v, &Iter[K,V]{rightNode, i})
	}

	for i := uint16(0); i < n.Size/2+1; i += 1 {
		j := n.Size/2 + 1 + i
		if leftNode.Childs[j] != nil {
			rightNode.Childs[i] = leftNode.Childs[j]
			rightNode.Childs[i].Parent = rightNode
			rightNode.Childs[i].Pidx = int(i)
			leftNode.Childs[j] = nil
		}
	}

	leftNode.Data = leftNode.Data[:n.Size/2]
	for i, v := range leftNode.Data {
		n.OnSplit(v, &Iter[K,V]{leftNode, i})
	}

	return returnNode
}

func (n *Node[K, V]) Insert(newValue BTreeLeaf[K, V]) (*Node[K, V], *Iter[K, V]) {
	lookupNode := n
	var iter *Iter[K,V] = nil

	if uint16(len(n.Data)) == n.Size {
		n = n.Split()
		lookupNode = n
		if n.Parent != nil {
			lookupNode = n.Parent
		}
	}

	idx := lookupNode.FindIdxToInsert(newValue)
	if lookupNode.Childs[idx] != nil {
		lookupNode.Childs[idx], iter = lookupNode.Childs[idx].Insert(newValue)
	} else {
		lookupNode.insertAt(idx, newValue)
		iter = &Iter[K,V]{lookupNode, int(idx)}
	}

	return n, iter
}

func (n *Node[K, V]) Upsert(newValue BTreeLeaf[K, V]) (*Node[K, V], *Iter[K, V]) {
	itr := n.Find(newValue)
	if itr != nil {
		localNode := itr.n
		localNode.Data[itr.i] = newValue
		return n, itr
	}

	return n.Insert(newValue)
}

func (n *Node[K, V]) LocalFindByValue(value V) (BTreeLeaf[K, V], bool) {
	for _, v := range n.Data {
		if v.Value.Equals(value) {
			return v, true
		}
	}

	return BTreeLeaf[K, V]{}, false
}

func (n *Node[K, V]) Find(value BTreeLeaf[K, V]) (*Iter[K, V]) {
	for i, v := range n.Data {
		if v.OrderKey.Less(value.OrderKey) {
			continue
		}

		if value.OrderKey.Equals(v.OrderKey) {
			return &Iter[K,V]{n, i}
		} else { // Means we need to search left node
			if n.Childs[i] == nil {
				return nil
			}

			return n.Childs[i].Find(value)
		}
	}

	if n.Childs[n.Size] == nil {
		return nil
	}

	return n.Childs[n.Size].Find(value)
}

//TODO: write Remove

func (n *Node[K, V]) RemoveByLocalIdx(idx uint16) bool {
	if idx >= n.Size {
		return false
	}

	if n.Childs[idx+1] != nil {
		//TODO: we are having childs with no Data
		n.Data[idx] = n.Childs[idx+1].Data[0]
		return n.Childs[idx+1].RemoveByLocalIdx(0)
	}

	copy(n.Data[idx:], n.Data[idx+1:])
	n.Data = n.Data[:len(n.Data)-1]

	// TODO: somehow we need to remove this Node after it has no Data
	// if len(n.Data) == 0 {
	// 	n.Parent
	// }
	return true
}

func (n Node[K, V]) FindIdxToInsert(newValue BTreeLeaf[K, V]) uint16 {
	idx := uint16(len(n.Data))

	for i, v := range n.Data {
		if newValue.OrderKey.Less(v.OrderKey) {
			return uint16(i)
		}
	}

	return idx
}

func (n *Node[K, V]) stringHelper() string {
	return fmt.Sprint(n.Data)
}

func (n *Node[K, V]) stringLoopHelper() string {
	res := fmt.Sprintf(" {%v}", n.Data)

	for i, v := range n.Childs {

		if v == nil {
			res += " nil "
		} else {
			if i != 0 && i-1 < len(n.Data) {
				res += fmt.Sprintf(" < %v < ", n.Data[i-1])
			}

			res += v.stringHelper()
		}
	}

	res += "\n"

	for _, v := range n.Childs {
		if v != nil {
			res += v.stringLoopHelper()
		}
	}

	return res
}

func (n *Node[K, V]) DepthFirstTraverse(visiter func(v BTreeLeaf[K, V])) {
	for i, v := range n.Childs {

		if v != nil {
			v.DepthFirstTraverse(visiter)
		}

		if i < len(n.Data) {
			visiter(n.Data[i])
		}
	}
}

func (n *Node[K, V]) String() string {

	res := fmt.Sprint(n.Data)
	res += "\n"
	res += n.stringLoopHelper()
	return res
}
