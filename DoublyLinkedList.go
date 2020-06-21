package main

import "strconv"

// Simple x,y pair
type Coord struct {
	x int
	y int
}

// Node to hold Coords
type Node struct {
	val  Coord
	next *Node
	prev *Node
}

// The list itself
type DoublyLinkedList struct {
	head *Node
	tail *Node
}

// Return the front of the list if it exists, otherwise
// return the NilCoord (-999,-999)
func (dll *DoublyLinkedList) Front() Coord {
	if dll.head != nil {
		return dll.head.val
	}
	return NilCoord()
}

// Return the back of the list if it exists, otherwise
// return the NilCoord (-999,-999)
func (dll *DoublyLinkedList) Back() Coord {
	if dll.head != nil {
		return dll.tail.val
	}
	return NilCoord()
}

// Push an item to the front of the list
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

// Push an item to the back of the list
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

// Pop the front of the list and update head
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

// Pop the last item from the list and update the tail
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

// Checks if c is a member of the list
func (dll *DoublyLinkedList) Member(c Coord) bool {
	iter := dll.GetIterator()
	for iter.Next() {
		if Equals(c, iter.Get()) {
			return true
		}
	}
	return false
}

// Checks if two coordinates are equal
func Equals(a Coord, b Coord) bool {
	return a.x == b.x && a.y == b.y
}

// Checks if a coord is equal to the NilCoord
func (c Coord) IsNilCoord() bool {
	return Equals(c, NilCoord())
}

// Returns the NilCoord
func NilCoord() Coord {
	return Coord{-999, -999}
}

// Converts a coordinate to a string
func (c Coord) ToString() string {
	return strconv.Itoa(c.x) + "," + strconv.Itoa(c.y)
}

// Iterator constrution
type DoublyLinkedListIterator struct {
	dll     *DoublyLinkedList
	curNode *Node
}

// Make an iterator for the given list
func (d *DoublyLinkedList) GetIterator() DoublyLinkedListIterator {
	var ret DoublyLinkedListIterator
	ret.dll = d
	ret.Reset()
	return ret
}

// Set the curNode to the head of the list
func (dlli *DoublyLinkedListIterator) Reset() {
	dlli.curNode = &Node{NilCoord(), dlli.dll.head, nil}
}

// Advance through the list by one Node
func (dlli *DoublyLinkedListIterator) Next() bool {
	if dlli.curNode != nil && dlli.curNode.next != nil {
		dlli.curNode = dlli.curNode.next
		return true
	}
	return false
}

// Get the value curNode of our iterator
func (dlli *DoublyLinkedListIterator) Get() Coord {
	if dlli.curNode != nil {
		return dlli.curNode.val
	}
	return NilCoord()
}
