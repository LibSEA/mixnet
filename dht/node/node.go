package node

import "net"

type NodeID []byte

type Node struct {
	ID   NodeID
	Port uint16
	IP   net.IPAddr
}

func (n *Node) GetKey() string {
	return string(n.ID)
}
