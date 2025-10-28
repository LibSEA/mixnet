package dht

import (
	"net"
)

type DHT struct {
	Self *Contact

	RoutingTable [][]*Contact
}

const (
// alpha = 3   // Concurrency
// b     = 256 // Bits
// k     = 20

// expire    = 86400 * time.Second
// refresh   = 3600 * time.Second
// replicate = 3600 * time.Second
// republish = 86400 * time.Second
)

type Contact struct {
	Id   NodeID
	Port uint16
	IP   net.IPAddr
}

// The idea is that the DHT has as its key a public key and it's payload is
// signed by the private key. And we will reject any attempt to write to that
// key that isn't signed and verified.
//
// My hope is this will have the side effect of making it hard to launch some
// of the known attacks against kademlia.
type Options struct {
	IP         net.IPAddr
	Port       uint16
	PrivateKey []byte
	PublicKey  []byte
}

type NodeID [32]byte

func New() {

}

func (d *DHT) Ping() {

}

func (d *DHT) Store() {

}

func (d *DHT) FindNode(nodeId NodeID) []Contact {
	return nil
}

func (d *DHT) FindValue() ([]byte, []NodeID) {
	return nil, nil

}
