package maplist

import "container/list"

type GetKey[K comparable] interface {
	GetKey() K
}

type MapList[K comparable, V GetKey[K]] struct {
	list *list.List
	dict map[K]*list.Element
}

func New[K comparable, V GetKey[K]]() *MapList[K, V] {
	return &MapList[K, V]{
		list: list.New(),
		dict: make(map[K]*list.Element),
	}
}

func (ml *MapList[K, V]) Get(key K) (*list.Element, bool) {
	e, ok := ml.dict[key]
	return e, ok
}

func (ml *MapList[K, V]) Remove(elm *list.Element) V {
	delete(ml.dict, ml.Val(elm).GetKey())
	return ml.list.Remove(elm).(V)
}

func (ml *MapList[K, V]) Back() *list.Element {
	return ml.list.Back()
}

func (ml *MapList[K, V]) Front() *list.Element {
	return ml.list.Front()
}

func (ml *MapList[K, V]) Init() *MapList[K, V] {
	ml.list.Init()
	clear(ml.dict)

	return ml
}

func (ml *MapList[K, V]) InsertAfter(v V, mark *list.Element) *list.Element {
	if el, ok := ml.dict[v.GetKey()]; ok {
		ml.MoveAfter(el, mark)
		return el
	}

	el := ml.list.InsertAfter(v, mark)

	ml.dict[v.GetKey()] = el

	return el
}

func (ml *MapList[K, V]) InsertBefore(v V, mark *list.Element) *list.Element {
	if el, ok := ml.dict[v.GetKey()]; ok {
		ml.MoveBefore(el, mark)
		return el
	}

	el := ml.list.InsertBefore(v, mark)

	ml.dict[v.GetKey()] = el

	return el
}

func (ml *MapList[K, V]) Val(el *list.Element) V {
	return el.Value.(V)
}

func (ml *MapList[K, V]) Len() int {
	return ml.list.Len()
}

func (ml *MapList[K, V]) MoveAfter(e, mark *list.Element) {
	ml.list.MoveAfter(e, mark)
}

func (ml *MapList[K, V]) MoveBefore(e, mark *list.Element) {
	ml.list.MoveBefore(e, mark)
}

func (ml *MapList[K, V]) MoveToBack(e *list.Element) {
	ml.list.MoveToBack(e)
}

func (ml *MapList[K, V]) MoveToFront(e *list.Element) {
	ml.list.MoveToFront(e)
}

func (ml *MapList[K, V]) PushBack(v V) *list.Element {
	if el, ok := ml.dict[v.GetKey()]; ok {
		ml.MoveToBack(el)
		return el
	}

	el := ml.list.PushBack(v)

	ml.dict[v.GetKey()] = el

	return el
}

func (ml *MapList[K, V]) PushBackList(other *MapList[K, V]) {
	e := other.Front()
	for i := other.list.Len(); i > 0; i-- {
		ml.PushBack(ml.Val(e))
		e = e.Next()
	}
}

func (ml *MapList[K, V]) PushFront(v V) *list.Element {
	if el, ok := ml.dict[v.GetKey()]; ok {
		ml.MoveToFront(el)
		return el
	}

	el := ml.list.PushFront(v)

	ml.dict[v.GetKey()] = el

	return el
}

func (ml *MapList[K, V]) PushFrontList(other *MapList[K, V]) {
	e := other.Back()
	for i := other.list.Len(); i > 0; i-- {
		ml.PushFront(ml.Val(e))
		e = e.Prev()
	}
}
