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
	itemsCount int
	FrontItem  *ListItem
	BackItem   *ListItem
}

func (l *list) Len() int {
	return l.itemsCount
}

func (l *list) Front() *ListItem {
	return l.FrontItem
}

func (l *list) Back() *ListItem {
	return l.BackItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := &ListItem{
		Value: v,
		Next:  l.FrontItem,
		Prev:  nil,
	}
	if l.FrontItem != nil {
		l.FrontItem.Prev = newListItem
	}
	if l.BackItem == nil {
		l.BackItem = newListItem
	}
	l.FrontItem = newListItem
	l.itemsCount++

	return newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.BackItem,
	}
	if l.BackItem != nil {
		l.BackItem.Next = newListItem
	}
	if l.FrontItem == nil {
		l.FrontItem = newListItem
	}
	l.BackItem = newListItem
	l.itemsCount++

	return newListItem
}

func (l *list) Remove(item *ListItem) {
	if item == nil {
		return
	}

	if item.Prev != nil {
		item.Prev.Next = item.Next
	} else {
		l.FrontItem = item.Next
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	} else {
		l.BackItem = item.Prev
	}

	l.itemsCount--
}

func (l *list) MoveToFront(item *ListItem) {
	if l.FrontItem == item {
		return
	}

	l.Remove(item)
	l.PushFront(item.Value)
}

func NewList() List {
	return new(list)
}
