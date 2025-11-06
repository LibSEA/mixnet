package routingtable

import (
	kb "github.com/LibSEA/mixnet/dht/internal/kbucket"
	n "github.com/LibSEA/mixnet/dht/node"
)

type RoutingTable struct {
	self    *n.Node
	buckets []kb.KBucket
}


