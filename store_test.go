package main

import (
	"testing"
)



func TestInsert(t *testing.T) {
	sut := NewStore()

	insert := func(key string, val int, log bool) {
		/*sut, _ =*/ sut.Insert(IntLesser{KeyValue: key, Value: int64(val)})


		data := make([]IntLesser, 0, 100)
		visiter := func(v Lesser) {
			data = append(data, v.(IntLesser))
		}

		if log {
			sut.DepthFirstTraverse(visiter)
			t.Logf("insert %v: %v", val, data)
			//t.Logf("insert %v: %v", val, sut.String())
		}
	}

	insert("10", 10, false)
	insert("2", 2, false)
	insert("5", 5, false)
	insert("3", 3, false)
	insert("1", 1, false)
	t.Log(sut.String())
	
	insert("20", 20, true)
	insert("21", 21, true)
	insert("18", 18, true)
	insert("19", 19, true)
	insert("17", 17, true)
	insert("6", 6, true)
	insert("7", 7, true)
	insert("8", 8, true)
	insert("9", 9, true)
	insert("4", 4, true)
	insert("11", 11, true)
	insert("12", 12, true)
	insert("13", 13, true)
	insert("14", 14, true)
	insert("15", 15, true)
	insert("16", 16, true)

	t.Log("---------------------------")

	v := sut.Get("11")
	t.Logf("Get 11 %v", v)

	t.Fail()
}