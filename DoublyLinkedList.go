package main

import "strconv"

type Coord struct {
	x int
	y int
}

type Node struct {
	val  Coord
	next *Node
	prev *Node
}
type DoublyLinkedList struct {
	head *Node
	tail *Node
}

func (dll *DoublyLinkedList) Front() Coord {
	if dll.head != nil {
		return dll.head.val
	}
	return NilCoord()
}

func (dll *DoublyLinkedList) Back() Coord {
	if dll.head != nil {
		return dll.tail.val
	}
	return NilCoord()
}

func (dll *DoublyLinkedList) PushFront(c Coord) {
	newNode := Node{c, nil, nil}
	// Then the head is this node, as is the tail
	if dll.head == nil {
		dll.head = &newNode
		dll.tail = &newNode
	} else {
		newNode.next = dll.head
		dll.head.prev = &newNode
		dll.head = &newNode
	}
}

func (dll *DoublyLinkedList) PushBack(c Coord) {
	// If the tail is null, we know the list is empty
	// So just PushFront the coord lol
	if dll.tail == nil {
		dll.PushFront(c)
	} else {
		newNode := Node{c, nil, dll.tail}
		dll.tail.next = &newNode
		dll.tail = &newNode
	}
}

func (dll *DoublyLinkedList) PopFront() Coord {
	if dll.head != nil {
		retVal := dll.head.val
		// Check if we need to nullify the list
		if dll.head.next == nil {
			dll.head = nil
			dll.tail = nil
		} else {
			dll.head = dll.head.next
			dll.head.prev = nil
		}
		return retVal
	}
	return NilCoord()
}

func (dll *DoublyLinkedList) PopBack() Coord {
	if dll.tail != nil {
		retVal := dll.tail.val
		// Check if we need to nullify the list
		if dll.tail.prev == nil {
			dll.head = nil
			dll.tail = nil
		} else {
			dll.tail = dll.tail.prev
			dll.tail.next = nil
		}
		return retVal
	}
	return NilCoord()
}

func (dll *DoublyLinkedList) Member(c Coord) bool {
	iter := dll.GetIterator()
	for iter.Next() {
		if Equals(c, iter.Get()) {
			return true
		}
	}
	return false
}

func Equals(a Coord, b Coord) bool {
	return a.x == b.x && a.y == b.y
}

func (c Coord) IsNilCoord() bool {
	return Equals(c, NilCoord())
}

func NilCoord() Coord {
	return Coord{-999, -999}
}

func (c Coord) ToString() string {
	return strconv.Itoa(c.x) + "," + strconv.Itoa(c.y)
}

// Iterator-ish portion

type DoublyLinkedListIterator struct {
	dll     *DoublyLinkedList
	curNode *Node
}

func (d *DoublyLinkedList) GetIterator() DoublyLinkedListIterator {
	var ret DoublyLinkedListIterator
	ret.dll = d
	ret.Reset()
	return ret
}

func (dlli *DoublyLinkedListIterator) Reset() {
	dlli.curNode = &Node{NilCoord(), dlli.dll.head, nil}
}

func (dlli *DoublyLinkedListIterator) Next() bool {
	if dlli.curNode != nil && dlli.curNode.next != nil {
		dlli.curNode = dlli.curNode.next
		return true
	}
	return false
}

func (dlli *DoublyLinkedListIterator) Get() Coord {
	if dlli.curNode != nil {
		return dlli.curNode.val
	}
	return NilCoord()
}
