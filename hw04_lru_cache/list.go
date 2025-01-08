package hw04lrucache

import "fmt"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	Head *ListItem
	Tail *ListItem
	Size int
}

// create a new list (unexported).
func NewList() *list {
	return &list{}
}

// Length of the list.
func (l *list) Len() int {
	return l.Size
}

// First element of the list.
func (l *list) Front() *ListItem {
	return l.Head
}

// Last element of the list.
func (l *list) Back() *ListItem {
	return l.Tail
}

// Add element to the front of the list.
func (l *list) PushFront(v interface{}) *ListItem {
	newFont := &ListItem{
		Value: v,
	}

	if l.Size == 0 {
		l.Head = newFont
		l.Tail = newFont
	} else {
		newFont.Next = l.Head
		l.Head.Prev = newFont
		l.Head = newFont
	}

	l.Size++
	return newFont
}

// Add element to the back of the list.
func (l *list) PushBack(v interface{}) *ListItem {
	newBack := &ListItem{
		Value: v,
	}

	if l.Size == 0 {
		l.Head = newBack
		l.Tail = newBack
	} else {
		l.Tail.Next = newBack
		newBack.Prev = l.Tail
		l.Tail = newBack
	}

	l.Size++
	return newBack
}

// Remove element from the list.
func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.Head = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.Tail = i.Prev
	}

	i.Prev = nil
	i.Next = nil
	l.Size--
}

// Move element to the front of the list.
func (l *list) MoveToFront(i *ListItem) {
	if i == nil || l.Size == 0 || i == l.Head {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.Tail = i.Prev
	}

	i.Prev = nil
	i.Next = l.Head
	l.Head.Prev = i
	l.Head = i
}

func (l *list) PrintHeadAndTailWithSize() {
	fmt.Printf("HEAD: %v / TAIL: %v / SIZE: %d\n", l.Front().Value, l.Back().Value, l.Size)
}

func (l *list) IterByList() {
	if l.Size < 1 {
		return
	}
	listInSlice := []interface{}{}
	listItem := l.Head
	for listItem != nil {
		listInSlice = append(listInSlice, listItem.Value)
		listItem = listItem.Next
	}
	fmt.Println(l.Len())
	for idx, value := range listInSlice {
		if idx == len(listInSlice)-1 {
			fmt.Printf("| %v |\n", value)
			return
		}
		fmt.Printf("| %v | <-> ", value)
	}
}

func IterByList(l List) {
	if l.Len() < 1 {
		return
	}
	listInSlice := []interface{}{}
	listItem := l.Front()
	for listItem != nil {
		listInSlice = append(listInSlice, listItem.Value)
		listItem = listItem.Next
	}
	fmt.Println(l.Len())
	for idx, value := range listInSlice {
		if idx == len(listInSlice)-1 {
			fmt.Printf("| %v |\n", value)
			return
		}
		fmt.Printf("| %v | <-> ", value)
	}
}
