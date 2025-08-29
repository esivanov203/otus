package hw04lrucache

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
	len  int
	head *ListItem
	back *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.head != nil {
		newItem.Next = l.head
		l.head.Prev = newItem
	}
	l.head = newItem
	if l.len == 0 {
		l.back = newItem
	}
	l.len++

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.head == nil {
		return l.PushFront(v)
	}
	newItem := &ListItem{Value: v, Prev: l.back}
	l.back.Next = newItem
	l.back = newItem
	l.len++

	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i == nil || l.len == 0 {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}
	i.Next = nil
	i.Prev = nil
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.head {
		return
	}
	// вырезаем i из текущего места
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if l.back == i {
		l.back = i.Prev
	}
	// вставляем в начало
	i.Prev = nil
	i.Next = l.head
	if l.head != nil {
		l.head.Prev = i
	}
	l.head = i
}
