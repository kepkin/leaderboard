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
}

type IntLesser int64

func (s IntLesser) Less(than Lesser) bool {
	return s < than.(IntLesser)
}


type Node struct {
	Parent *Node
	Size uint16
	Data []Lesser
	Childs []*Node
}

func NewNode(size uint16, parent *Node) *Node {
	return &Node{
		Parent: parent,
		Size: size,
		Data: make([]Lesser, 0, size),
		Childs: make([]*Node, size+1),
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

		upperNode = NewNode(n.Size, nil)
		leftNode.Parent = upperNode
		returnNode = upperNode

		upperNode.Childs[0] = leftNode
	}

	rightNode := NewNode(n.Size, upperNode)

	pivotValue := leftNode.Data[n.Size/2]

	idx := upperNode.FindIdxToInsert(pivotValue)
	upperNode.insertAt(idx, pivotValue)
	upperNode.Childs[idx+1] = rightNode

	rightNode.Data = append(rightNode.Data, leftNode.Data[n.Size/2 + 1:]...)
	leftNode.Data = leftNode.Data[:n.Size/2]

	return returnNode
}

func (n *Node) Insert(newValue Lesser) *Node {
	if uint16(len(n.Data)) == n.Size {
		n = n.Split()
	}

	idx := n.FindIdxToInsert(newValue)
	if n.Childs[idx] != nil {
		n.Childs[idx] = n.Childs[idx].Insert(newValue)
	} else {
		n.insertAt(idx, newValue)	
	}
	
	return n
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

func (n *Node) String() string {

	res := fmt.Sprint(n.Data)
	res += "\n"
	res += n.stringLoopHelper()
	return res
}



type Store struct {

}


func (s *Store) Put(a int) error {
	return nil
}