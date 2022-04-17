package main


import (
	"fmt"
)


/*

insert
read
flush



*/


type Lesser interface {
	Less(than Lesser) bool
	Key() string
}

type IntLesser struct {
	Value int64
	KeyValue   string
}

func (s IntLesser) Less(than Lesser) bool {
	return s.Value < than.(IntLesser).Value
}

func (s IntLesser) Key() string {
	return s.KeyValue
}


type OnSplitTrigger func(Lesser, *Node)

type Node struct {
	Parent *Node
	Size uint16
	Data []Lesser
	Childs []*Node

	OnSplit OnSplitTrigger
}

func NewNode(size uint16, parent *Node, onSplit OnSplitTrigger) *Node {
	return &Node{
		Parent: parent,
		Size: size,
		Data: make([]Lesser, 0, size),
		Childs: make([]*Node, size+1),
		OnSplit: onSplit,
	}
}


func (n *Node) insertAt(idx uint16, newValue Lesser) { 
	n.Data = append(n.Data, nil)
	copy(n.Data[idx+1:], n.Data[idx:])
	n.Data[idx] = newValue
}

func (n *Node) Split() *Node {
	
	upperNode := n.Parent
	returnNode := n
	leftNode := n

	if upperNode == nil {

		upperNode = NewNode(n.Size, nil, n.OnSplit)
		leftNode.Parent = upperNode
		returnNode = upperNode

		upperNode.Childs[0] = leftNode
	}

	rightNode := NewNode(n.Size, upperNode, n.OnSplit)

	pivotValue := leftNode.Data[n.Size/2]

	idx := upperNode.FindIdxToInsert(pivotValue)
	upperNode.insertAt(idx, pivotValue)
	n.OnSplit(pivotValue, upperNode)

	upperNode.Childs[idx+1] = rightNode

	rightNode.Data = append(rightNode.Data, leftNode.Data[n.Size/2 + 1:]...)

	for _, v := range(rightNode.Data) {
		n.OnSplit(v, rightNode)
	}

	leftNode.Data = leftNode.Data[:n.Size/2]

	for _, v := range(leftNode.Data) {
		n.OnSplit(v, leftNode)
	}

	return returnNode
}

func (n *Node) Insert(newValue Lesser) (*Node, *Node) {
	insertedNode := n

	if uint16(len(n.Data)) == n.Size {
		n = n.Split()
	}

	idx := n.FindIdxToInsert(newValue)
	if n.Childs[idx] != nil {
		n.Childs[idx], insertedNode = n.Childs[idx].Insert(newValue)
	} else {
		n.insertAt(idx, newValue)	
	}
	
	return n, insertedNode
}

func (n *Node) FindByLocalKey(key string) (int, bool) {
	for i, v := range n.Data {
		if v.Key() == key {
			return i, true
		}
	}

	return 0, false
}


func (n *Node) RemoveByLocalKey(key string) bool {
	idx, ok := n.FindByLocalKey(key)
	if !ok {
		return false
	}

	return n.RemoveByLocalIdx(uint16(idx))
}

func (n *Node) RemoveByLocalIdx(idx uint16) bool {
	if idx >= n.Size {
		return false
	}

	if n.Childs[idx+1] != nil {
		n.Data[idx] = n.Childs[idx+1].Data[0]
		return n.Childs[idx+1].RemoveByLocalIdx(0)
	}

	copy(n.Data[idx:], n.Data[idx+1:])
	return true
}


func (n Node) FindIdxToInsert(newValue Lesser) uint16 {
	idx := uint16(len(n.Data))

	for i, v := range n.Data {
		if newValue.Less(v) {
			return uint16(i)
		}
	}

	return idx
}


func (n *Node) stringHelper() string {
	return fmt.Sprint(n.Data)
}

func (n *Node) stringLoopHelper() string {
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

func (n *Node) DepthFirstTraverse(visiter func(v Lesser)) {
	for i, v := range n.Childs {

		if v != nil {
			v.DepthFirstTraverse(visiter)
		}

		
		if i < len(n.Data) {
			visiter(n.Data[i])	
		}		
	}
}

func (n *Node) String() string {

	res := fmt.Sprint(n.Data)
	res += "\n"
	res += n.stringLoopHelper()
	return res
}

type Store struct {
	btree *Node
	ridx  map[string]*Node
}

func NewStore() Store {
	res := Store{}

	res.btree = NewNode(5, nil, res.onSplit)
	res.ridx = make(map[string]*Node)

	return res
}

func (s *Store) onSplit(value Lesser, n *Node) {
	s.ridx[value.Key()] = n
}

func (s *Store) Insert(value Lesser) error {
	n, ok := s.ridx[value.Key()]
	if ok {
		n.RemoveByLocalKey(value.Key())
	}

	s.btree, n = s.btree.Insert(value)
	s.ridx[value.Key()] = n
	return nil
}

func (s *Store) Get(key string) Lesser {
	n, ok := s.ridx[key]
	if !ok {
		return nil
	}

	if idx, ok := n.FindByLocalKey(key); ok {
		return n.Data[idx]
	}

	return nil
}

func (s *Store) DepthFirstTraverse(visiter func(v Lesser)) {
	s.btree.DepthFirstTraverse(visiter)
}

func (s *Store) String() string {
	return s.btree.String()
}