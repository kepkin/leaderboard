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

type OnSplitTrigger[K any] func(K, *Node[K], int)
type LessFuncType[K any] func(a, b K) bool
type EqualsFuncType[K any] func(a, b K) bool

func DummyOnSplit[K any](K, *Node[K], int) {}

type NodeSettings[K any] struct {
	Size          int
	LessFunc      LessFuncType[K]
	EqualsFuncKey EqualsFuncType[K]

	OnSplit         OnSplitTrigger[K]
	treeRebalanceMu sync.RWMutex
}

type Node[K any] struct {
	Parent *Node[K]
	Pidx   int
	Data   []K
	Childs []*Node[K]

	mu sync.Mutex

	s *NodeSettings[K]
}

func NewNode[K any](
	size int,
	lessFunc LessFuncType[K],
	equalsFuncKey EqualsFuncType[K],
	onSplit OnSplitTrigger[K],
) *Node[K] {

	if onSplit == nil {
		onSplit = DummyOnSplit[K]
	}

	return &Node[K]{
		Parent: nil,
		Data:   make([]K, 0, size),
		Childs: make([]*Node[K], size+1),

		s: &NodeSettings[K]{
			Size:            size,
			OnSplit:         onSplit,
			LessFunc:        lessFunc,
			EqualsFuncKey:   equalsFuncKey,
			treeRebalanceMu: sync.RWMutex{},
		},
	}
}

func newTreeNode[K any](seed *Node[K], parent *Node[K]) *Node[K] {
	return &Node[K]{
		Parent: parent,
		Data:   make([]K, 0, seed.s.Size),
		Childs: make([]*Node[K], seed.s.Size+1),
		s:      seed.s,
	}
}

// We should be there only if
//  - newValue.OrderKey < Data[idx].OrderKey
//  - Childs[idx] == nil
//  - Childs[:last:] == nil
//  - len(Data) < Size
func (n *Node[K]) insertAt(idx int, newValue K) {
	n.mu.Lock()
	defer n.mu.Unlock()

	var kk K
	n.Data = append(n.Data, kk)
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

func (n *Node[K]) Split() *Node[K] {
	upperNode := n.Parent
	returnNode := n
	leftNode := n

	if upperNode == nil {

		upperNode = newTreeNode(n, nil)
		leftNode.Parent = upperNode
		returnNode = upperNode

		upperNode.Childs[0] = leftNode
	} else {
		if len(upperNode.Data) == n.s.Size {
			upperNode = upperNode.Split()
		}

		// leftNode.Parent = upperNode
	}

	rightNode := newTreeNode(n, upperNode)

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

func (n *Node[K]) Insert(newValue K) (*Node[K], int) {
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
func (n *Node[K]) InsertLocked(newValue K, rebalance bool) (*Node[K], int, bool) {
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

func (n *Node[K]) Upsert(newValue K) (*Node[K], int) {
	localNode, idx := n.Find(newValue)
	if localNode != nil {
		localNode.Data[idx] = newValue
		return n, idx
	}

	return n.Insert(newValue)
}

func (n *Node[K]) Find(value K) (*Node[K], int) {
	if len(n.Data) == 0 {
		return nil, 1
	}

	i := n.FindIdxToInsert(value)
	for ; i < len(n.Data); i++ {
		v := n.Data[i]
		if n.s.LessFunc(v, value) {
			continue
		}

		if n.s.EqualsFuncKey(v, value) {
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

func (n *Node[K]) removeChild(idx int) {
	n.Childs[idx] = nil
}

func (n *Node[K]) RemoveByLocalIdx(idx int) bool {
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

func (n Node[K]) FindIdxToInsert(newValue K) int {
	l := 0
	r := len(n.Data)

	if r == 0 {
		return 0
	}

	for r-l > 1 {
		p := (r-l)/2 + l
		if n.s.LessFunc(newValue, n.Data[p]) {
			r = p
		} else {
			l = p
		}
	}

	if n.s.LessFunc(n.Data[l], newValue) {
		return r
	} else {
		return l
	}
}

func (n *Node[K]) stringHelper() string {
	return fmt.Sprint(n.Data)
}

func (n *Node[K]) stringLoopHelper() string {
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

func (n *Node[K]) DepthFirstTraverse(visiter func(v K)) {
	for i, v := range n.Childs {

		if v != nil {
			v.DepthFirstTraverse(visiter)
		}

		if i < len(n.Data) {
			visiter(n.Data[i])
		}
	}
}

func (n *Node[K]) String() string {

	res := fmt.Sprint(n.Data)
	res += "\n"
	res += n.stringLoopHelper()
	return res
}

func (n *Node[K]) Iter(idx int) *Iter[K] {
	return newIter(n, idx)
}

func (n *Node[K]) At(idx int) K {
	return n.Data[idx]
}
