package main

import (
	"testing"
)



func TestInsert(t *testing.T) {
	sut := NewNode(5, nil)

	insert := func(val int) {
		sut = sut.Insert(IntLesser(val))
		t.Logf("insert %v: %v", val, sut.String())
	}

	sut = sut.Insert(IntLesser(10))
	sut = sut.Insert(IntLesser(2))
	sut = sut.Insert(IntLesser(5))
	sut = sut.Insert(IntLesser(3))
	sut = sut.Insert(IntLesser(1))
	t.Log(sut.String())
	
	insert(20)
	insert(21)
	insert(18)
	insert(19)
	insert(17)
	insert(6)
	insert(7)
	insert(8)
	insert(9)
	insert(4)
	insert(11)
	insert(12)
	// insert(13)
	// insert(14)
	// insert(15)
	// insert(16)

	t.Fail()
}