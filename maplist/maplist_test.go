package maplist

import "testing"

type Int struct {
	v int
}

func (i Int) GetKey() int {
	return i.v
}

func TestMapList(t *testing.T) {
	sut := New[int, Int]()

	// []
	if sut.Len() != 0 {
		t.Fatalf("New maplists should be len 0. actual %d", sut.Len())
	}

	e, ok := sut.Get(0)

	if ok || e != nil {
		t.Fatalf(
			"Empty list should return zero element and ok == false. e: %d, ok: %t",
			sut.Val(e).v,
			ok,
		)
	}

	if sut.Front() != nil {
		t.Fatal("Front should return nil for empty list")
	}

	e = sut.PushFront(Int{1})

	// [1]
	if e == nil || sut.Val(e).v != 1 {
		t.Fatal("failed to push 1 to front of list")
	}

	if sut.Front() != sut.Back() || sut.Back() != e {
		t.Fatal("with one element in list front should equal back")
	}

	e = sut.PushFront(Int{1})

	if e == nil || sut.Val(e).v != 1 {
		t.Fatal("failed to push 1 to front of list")
	}

	if sut.Front() != sut.Back() || sut.Back() != e {
		t.Fatal("with one element in list front should equal back")
	}

	e2 := sut.PushBack(Int{2})

	// [1, 2]
	if sut.Front() != e || sut.Back() != e2 || sut.Len() != 2 {
		t.Fatal("failed to PushBack")
	}

	e3 := sut.PushBack(Int{2})

	if e2 != e3 || sut.Len() != 2 || sut.Front() != e || sut.Back() != e2 {
		t.Fatal("PushBack should return original element on duplicate key")
	}

	e3 = sut.PushFront(Int{2})

	// [2, 1]
	if e2 != e3 || sut.Len() != 2 || sut.Front() != e2 || sut.Back() != e {
		t.Fatal(
			"PushFront should return original element on dup key and move to front",
		)
	}

	e3 = sut.PushBack(Int{2})

	// [1, 2]
	if e2 != e3 || sut.Len() != 2 || sut.Front() != e || sut.Back() != e2 {
		t.Fatal("PushBack should move dupe to back")
	}

	e3 = sut.InsertAfter(Int{1}, sut.Back())

	// [2, 1]
	if e != e3 || sut.Back() != e {
		t.Fatal("InsertAfter acts like MoveAfter if already in list")
	}

	e3 = sut.InsertBefore(Int{1}, sut.Front())

	// [1, 2]
	if e != e3 || sut.Front() != e {
		t.Fatal("InsertBefore acts like MoveBefore if already in list")
	}

	e3 = sut.InsertAfter(Int{3}, sut.Front())

	// [1, 3, 2]
	if sut.Len() != 3 || sut.Front().Next() != e3 {
		t.Fatal("InsertAfter should work")
	}

	e4 := sut.InsertBefore(Int{4}, e3)

	// [1, 4, 3, 2]
	if sut.Len() != 4 || sut.Front().Next() != e4 {
		t.Fatal("InsertBefore broken")
	}

	sut.Remove(e3)

	// [1, 4, 2]
	if sut.Len() != 3 {
		t.Fatal("remove broken")
	}

	sut.Remove(e4)

	// [1, 2]
	if sut.Len() != 2 || sut.Front() != e || sut.Back() != e2 {
		t.Fatal("remove broken")
	}

	ml2 := sut.Init()

	if sut != ml2 || sut.Len() != 0 {
		t.Fatal("init should clear list")
	}

	sut.PushBack(Int{1})
	sut.PushBack(Int{2})
	sut.PushBack(Int{5})

	other := New[int, Int]()

	other.PushBack(Int{3})
	other.PushBack(Int{4})
	other.PushBack(Int{5})
	other.PushBack(Int{6})

	sut.PushBackList(other)

	el := sut.Front()
	for i := 0; i < sut.Len(); i++ {
		ev := Int{i + 1}
		if sut.Val(el) != ev {
			t.Fatalf("wrong order after PushBackList %d %d", ev.v, sut.Val(el).v)
		}
		el = el.Next()
	}

	if sut.Len() != 6 {
		t.Fatalf("PushBackList should have length 6. has len: %d", sut.Len())
	}

	sut.Init()
	other.Init()

	sut.PushBack(Int{5})
	sut.PushBack(Int{6})
	sut.PushBack(Int{1})

	other.PushBack(Int{1})
	other.PushBack(Int{2})
	other.PushBack(Int{3})
	other.PushBack(Int{4})

	sut.PushFrontList(other)

	el = sut.Front()
	for i := 0; i < sut.Len(); i++ {
		ev := Int{i + 1}
		if sut.Val(el) != ev {
			t.Fatalf("wrong order after PushFrontList %d %d", ev.v, sut.Val(el).v)
		}
		el = el.Next()
	}

	if sut.Len() != 6 {
		t.Fatalf("PushFrontList should have length 6. has len: %d", sut.Len())
	}
}
