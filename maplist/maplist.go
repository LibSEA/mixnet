/*
mixnet - tool to create and manage LibSEA mixnets
Copyright (C) 2025  Liberatory Sofware Engineering Association

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// Package maplist is a combination of a [pkg/container/list] and a map.
//
// Any operation that would have added/inserted a list element whose key
// already is used will move the original list element to the position the
// new element would have been added/inserted.
//
// For a value to be usable in a maplist the value must implement the
// [GetKey] interface.
package maplist

import "container/list"

// GetKey is an interface that must be implemented by a value intended
// to be used in a maplist.
type GetKey[K comparable] interface {
	GetKey() K
}

type MapList[K comparable, V GetKey[K]] struct {
	list *list.List
	dict map[K]*list.Element
}

// New creates a new maplist. To reset a maplist see [MapList.Init].
func New[K comparable, V GetKey[K]]() *MapList[K, V] {
	return &MapList[K, V]{
		list: list.New(),
		dict: make(map[K]*list.Element),
	}
}

// Get attempts to act like indexing into a map. The biggest difference
// is the optional ok return value is not optional.
func (ml *MapList[K, V]) Get(key K) (*list.Element, bool) {
	e, ok := ml.dict[key]
	return e, ok
}

func (ml *MapList[K, V]) Remove(elm *list.Element) V {
	delete(ml.dict, ml.Val(elm).GetKey())
	return ml.list.Remove(elm).(V)
}

// Back is the same as [pkg/container/list.List.Back]
func (ml *MapList[K, V]) Back() *list.Element {
	return ml.list.Back()
}

// Front is the same as [pkg/container/list.List.Front]
func (ml *MapList[K, V]) Front() *list.Element {
	return ml.list.Front()
}

// Init clears all items from the maplist
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
