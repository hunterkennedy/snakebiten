package main

import "testing"

func teq(t *testing.T, ans Coord, exp Coord) {
	if !Equals(ans, exp) {
		errString := "Expected " + exp.ToString() + " but got " + ans.ToString()
		t.Errorf(errString)
	}
}

func TestEmpty(t *testing.T) {
	var dll DoublyLinkedList
	x1 := dll.PopBack()
	x2 := dll.PopFront()
	if !x1.IsNilCoord() || !x2.IsNilCoord() {
		errString := "Got " + x1.ToString() + " and " + x2.ToString() + " but wanted nil nodes on empty pop"
		t.Errorf(errString)
	}
}

func TestPushBack(t *testing.T) {
	var dll DoublyLinkedList
	dll.PushBack(Coord{1, 1})
	teq(t, dll.Front(), Coord{1, 1})
	teq(t, dll.Back(), Coord{1, 1})
	dll.PushBack(Coord{2, 2})
	teq(t, dll.Front(), Coord{1, 1})
	teq(t, dll.Back(), Coord{2, 2})
	dll.PushBack(Coord{3, 3})
	teq(t, dll.Back(), Coord{3, 3})
}

func TestPushFront(t *testing.T) {
	var dll DoublyLinkedList
	dll.PushFront(Coord{1, 1})
	teq(t, dll.Front(), Coord{1, 1})
	teq(t, dll.Back(), Coord{1, 1})
	dll.PushFront(Coord{2, 2})
	teq(t, dll.Back(), Coord{1, 1})
	teq(t, dll.Front(), Coord{2, 2})
	dll.PushFront(Coord{3, 3})
	teq(t, dll.Back(), Coord{1, 1})
	teq(t, dll.Front(), Coord{3, 3})
}

func TestMixed(t *testing.T) {
	var dll DoublyLinkedList
	dll.PushFront(Coord{1, 1})
	dll.PushBack(Coord{2, 2})
	teq(t, dll.Front(), Coord{1, 1})
	teq(t, dll.Back(), Coord{2, 2})
	dll.PushBack(Coord{3, 3})
	teq(t, dll.Front(), Coord{1, 1})
	teq(t, dll.Back(), Coord{3, 3})
}

func TestPopBack(t *testing.T) {
	var dll DoublyLinkedList
	dll.PushFront(Coord{1, 1})
	dll.PushBack(Coord{2, 2})
	teq(t, dll.Front(), Coord{1, 1})
	teq(t, dll.Back(), Coord{2, 2})
	dll.PushBack(Coord{3, 3})
	teq(t, dll.Front(), Coord{1, 1})
	teq(t, dll.Back(), Coord{3, 3})

	teq(t, dll.PopBack(), Coord{3, 3})
	teq(t, dll.Front(), Coord{1, 1})
	teq(t, dll.Back(), Coord{2, 2})
	teq(t, dll.PopFront(), Coord{1, 1})
	teq(t, dll.Front(), Coord{2, 2})
	teq(t, dll.Back(), Coord{2, 2})
	teq(t, dll.PopFront(), Coord{2, 2})

	teq(t, dll.PopBack(), NilCoord())
	teq(t, dll.PopFront(), NilCoord())
	teq(t, dll.Front(), NilCoord())
	teq(t, dll.Back(), NilCoord())

}
