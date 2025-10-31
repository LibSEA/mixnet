package routingtable

import (
	"container/list"

	"github.com/LibSEA/mixnet/dht/node"
)

type RoutingTable struct {
	self    *node.Node
	buckets *list.List
}
