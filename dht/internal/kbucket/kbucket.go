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

// Package kbucket implements a kbucket concept from kademlia.
package kbucket

import (
	"math/big"
	"time"

	n "github.com/LibSEA/mixnet/dht/node"
	"github.com/LibSEA/mixnet/maplist"
)

type KBucket struct {
	nodes           *maplist.MapList[string, *n.Node]
	replacements    *maplist.MapList[string, *n.Node]
	lastUpdated     time.Time
	span            [2]*big.Int
	ksize           int
	maxReplacements int
	self            *n.Node
}

func New(self *n.Node, ksize int, maxReplacements int) *KBucket {
	return &KBucket{
		nodes:        maplist.New[string, *n.Node](),
		replacements: maplist.New[string, *n.Node](),
		span: [2]*big.Int{
			big.NewInt(0),
			big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil),
		},
		ksize:           ksize,
		maxReplacements: maxReplacements,
		self:            self,
	}
}

// AddNode tries to add a node to this KBucket. If the KBucket is full we
// return false and add Node to replacement list. otherwise we add the Node
// and return true.
func (kb *KBucket) AddNode(node *n.Node) bool {
	if elm, ok := kb.nodes.Get(node.GetKey()); ok {
		kb.nodes.MoveToFront(elm)
		return true
	}

	if kb.nodes.Len() < kb.ksize {
		kb.nodes.PushFront(node)
		return true
	}

	kb.replacements.PushFront(node)

	if kb.replacements.Len() > kb.maxReplacements {
		kb.replacements.Remove(kb.replacements.Back())
	}

	return false
}

func (kb *KBucket) Split() {
	var sum big.Int
	var middle big.Int
	var two = big.NewInt(2)

	kb.span[1] = middle.Div(sum.Add(kb.span[0], kb.span[1]), two)

}
