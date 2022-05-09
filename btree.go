package leaderboard

import (
	"fmt"
	"sync"

	"golang.org/x/exp/constraints"
)

func StdLess[K constraints.Ordered](a, b K) bool {
	return a < b
}

func StdEquals[K constraints.Ordered](a, b K) bool {
	return a == b
}

type OnSplitTrigger[K any, V any] func(BTreeLeaf[K, V], *Node[K, V], int)
type LessFuncType[K any] func(a, b K) bool
type EqualsFuncType[K any] func(a, b K) bool

func DummyOnSplit[K any, V any](BTreeLeaf[K, V], *Node[K, V], int) {}

type BTreeLeaf[K any, V any] struct {
	OrderKey K
	Value    V
}

type NodeSettings[K any, V any] struct {
	Size          int
	LessFunc      LessFuncType[K]
	EqualsFuncKey EqualsFuncType[K]
	EqualsFuncVal EqualsFuncType[V]

	OnSplit         OnSplitTrigger[K, V]
	treeRebalanceMu sync.RWMutex
}

type Node[K any, V any] struct {
	Parent *Node[K, V]
	Pidx   int
	Data   []BTreeLeaf[K, V]
	Childs []*Node[K, V]

	mu sync.Mutex

	s *NodeSettings[K, V]
}

func NewNode[K any, V any](
	size int,
	lessFunc LessFuncType[K],
	equalsFuncKey EqualsFuncType[K],
	equalsFuncVal EqualsFuncType[V],
	onSplit OnSplitTrigger[K, V],
) *Node[K, V] {

	if onSplit == nil {
		onSplit = DummyOnSplit[K, V]
	}

	return &Node[K, V]{
		Parent: nil,
		Data:   make([]BTreeLeaf[K, V], 0, size),
		Childs: make([]*Node[K, V], size+1),

		s: &NodeSettings[K, V]{
			Size:            size,
			OnSplit:         onSplit,
			LessFunc:        lessFunc,
			EqualsFuncKey:   equalsFuncKey,
			EqualsFuncVal:   equalsFuncVal,
			treeRebalanceMu: sync.RWMutex{},
		},
	}
}

func newTreeNode[K any, V any](seed *Node[K, V], parent *Node[K, V]) *Node[K, V] {
	return &Node[K, V]{
		Parent: parent,
		Data:   make([]BTreeLeaf[K, V], 0, seed.s.Size),
		Childs: make([]*Node[K, V], seed.s.Size+1),
		s:      seed.s,
	}
}

// We should be there only if
//  - newValue.OrderKey < Data[idx].OrderKey
//  - Childs[idx] == nil
//  - Childs[:last:] == nil
//  - len(Data) < Size
func (n *Node[K, V]) insertAt(idx int, newValue BTreeLeaf[K, V]) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.Data = append(n.Data, BTreeLeaf[K, V]{})
	copy(n.Data[idx+1:], n.Data[idx:])
	n.Data[idx] = newValue
	for i := idx; i < len(n.Data); i++ {
		n.s.OnSplit(n.Data[i], n, i)
	}

	if n.Childs[idx+1] != nil {
		for i := idx + 1; i < len(n.Childs); i += 1 {
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

		upperNode = newTreeNode[K, V](n, nil)
		leftNode.Parent = upperNode
		returnNode = upperNode

		upperNode.Childs[0] = leftNode
	} else {
		if len(upperNode.Data) == n.s.Size {
			upperNode = upperNode.Split()
		}

		// leftNode.Parent = upperNode
	}

	rightNode := newTreeNode[K, V](n, upperNode)

	pivotValue := leftNode.Data[n.s.Size/2]

	idx := upperNode.FindIdxToInsert(pivotValue)
	leftNode.Pidx = idx
	rightNode.Pidx = idx + 1
	upperNode.insertAt(idx, pivotValue)
	upperNode.Childs[idx+1] = rightNode

	rightNode.Data = append(rightNode.Data, leftNode.Data[n.s.Size/2+1:]...)
	for i, v := range rightNode.Data {
		n.s.OnSplit(v, rightNode, i)
	}

	for i := 0; i < n.s.Size/2+1; i += 1 {
		j := n.s.Size/2 + 1 + i
		if leftNode.Childs[j] != nil {
			rightNode.Childs[i] = leftNode.Childs[j]
			rightNode.Childs[i].Parent = rightNode
			rightNode.Childs[i].Pidx = i
			leftNode.Childs[j] = nil
		}
	}

	leftNode.Data = leftNode.Data[:n.s.Size/2]

	return returnNode
}

func (n *Node[K, V]) Insert(newValue BTreeLeaf[K, V]) (*Node[K, V], int) {
	var needRebalance = false
	var idx int

	n.s.treeRebalanceMu.RLock()
	n, idx, needRebalance = n.InsertLocked(newValue, false)
	n.s.treeRebalanceMu.RUnlock()

	if !needRebalance {
		return n, idx
	}

	n.s.treeRebalanceMu.Lock()
	defer n.s.treeRebalanceMu.Unlock()

	for needRebalance {
		n, idx, needRebalance = n.InsertLocked(newValue, true)
		n, idx, needRebalance = n.InsertLocked(newValue, false)
	}

	return n, idx
}

// bool - returns if rebalancing is needed
func (n *Node[K, V]) InsertLocked(newValue BTreeLeaf[K, V], rebalance bool) (*Node[K, V], int, bool) {
	lookupNode := n
	var needRebalance bool = false

	if len(n.Data) == n.s.Size {
		if !rebalance {
			return lookupNode, -1, true
		}

		n = n.Split()

		lookupNode = n
		if n.Parent != nil {
			lookupNode = n.Parent
		}

		return n, -1, false
	}

	idx := lookupNode.FindIdxToInsert(newValue)
	if lookupNode.Childs[idx] != nil {
		lookupNode.Childs[idx], idx, needRebalance = lookupNode.Childs[idx].InsertLocked(newValue, rebalance)
	} else {
		if rebalance {
			return n, -1, false
		}

		lookupNode.insertAt(idx, newValue)
	}

	return n, idx, needRebalance
}

func (n *Node[K, V]) Upsert(newValue BTreeLeaf[K, V]) (*Node[K, V], int) {
	localNode, idx := n.Find(newValue)
	if localNode != nil {
		localNode.Data[idx] = newValue
		return n, idx
	}

	return n.Insert(newValue)
}

func (n *Node[K, V]) LocalFindByValue(value V) (BTreeLeaf[K, V], bool) {
	for _, v := range n.Data {
		if n.s.EqualsFuncVal(v.Value, value) {
			return v, true
		}
	}

	return BTreeLeaf[K, V]{}, false
}

func (n *Node[K, V]) Find(value BTreeLeaf[K, V]) (*Node[K, V], int) {
	for i, v := range n.Data {
		if n.s.LessFunc(v.OrderKey, value.OrderKey) {
			continue
		}

		if n.s.EqualsFuncKey(value.OrderKey, v.OrderKey) {
			return n, i
		} else { // Means we need to search left node
			if n.Childs[i] == nil {
				return nil, -1
			}

			return n.Childs[i].Find(value)
		}
	}

	if n.Childs[len(n.Data)] == nil {
		return nil, -1
	}

	return n.Childs[len(n.Data)].Find(value)
}

//TODO: write Remove

func (n *Node[K, V]) removeChild(idx int) {
	n.Childs[idx] = nil
}

func (n *Node[K, V]) RemoveByLocalIdx(idx int) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	if idx >= len(n.Data) {
		return false
	}

	if n.Childs[idx+1] != nil {
		n.Data[idx] = n.Childs[idx+1].Data[0]
		res := n.Childs[idx+1].RemoveByLocalIdx(0)

		n.s.OnSplit(n.Data[idx], n, idx)
		return res
	} else if n.Childs[idx] != nil {
		chidx := len(n.Childs[idx].Data) - 1
		n.Data[idx] = n.Childs[idx].Data[chidx]
		res := n.Childs[idx].RemoveByLocalIdx(chidx)

		n.s.OnSplit(n.Data[idx], n, idx)
		return res
	}

	// rememberValue := n.Data[idx]
	n.s.OnSplit(n.Data[idx], nil, idx)
	copy(n.Data[idx:], n.Data[idx+1:])
	n.Data = n.Data[:len(n.Data)-1]

	for i := idx; i < len(n.Data); i++ {
		n.s.OnSplit(n.Data[i], n, i)
	}

	for i := idx + 1; i < len(n.Childs); i += 1 {
		if n.Childs[i] != nil {
			n.Childs[i].Pidx -= 1
		}
	}

	copy(n.Childs[idx+1:], n.Childs[idx+2:])
	n.Childs[len(n.Data)+1] = nil

	if len(n.Data) == 0 {
		// fmt.Printf("removed last child %v\n", rememberValue)
		n.Data = nil
		n.Childs = nil
		n.Parent.removeChild(n.Pidx)
		n.Parent = nil
	}
	return true
}

func (n Node[K, V]) FindIdxToInsert(newValue BTreeLeaf[K, V]) int {
	for i, v := range n.Data {
		if n.s.LessFunc(newValue.OrderKey, v.OrderKey) {
			return i
		}
	}

	return len(n.Data)
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

func (n *Node[K, V]) Iter(idx int) *Iter[K, V] {
	return newIter[K, V](n, idx)
}

func (n *Node[K, V]) At(idx int) BTreeLeaf[K, V] {
	return n.Data[idx]
}
