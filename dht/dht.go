// Package dht The idea is that the DHT has as its key a public key and it's payload is
// signed by the private key. And we will reject any attempt to write to that
// key that isn't signed and verified.
//
// My hope is this will have the side effect of making it hard to launch some
// of the known attacks against kademlia.
package dht

import (
	"container/list"
	"crypto/ed25519"
	"net"
	"time"

	"github.com/LibSEA/mixnet/dht/node"
	"github.com/LibSEA/mixnet/store"
)

type DHT struct {
	//Self *Contact

	store store.Storage

	routingTable []*list.List
}

const (
	// alpha = 3   // Concurrency
	// b     = 256 // Bits
	// k     = 20
	keySize      = 32
	sigSize      = 64
	maxValueSize = 4096

	expire = 86400 * time.Second

// refresh   = 3600 * time.Second
// replicate = 3600 * time.Second
// republish = 86400 * time.Second
)

type Options struct {
	IP         net.IPAddr
	Port       uint16
	PrivateKey []byte
	PublicKey  []byte
	Store      store.Storage
}

func New(opts Options) *DHT {
	return &DHT{
		store: opts.Store,
	}
}

func (d *DHT) Ping(contact *node.Node) {

}

func splitValue(value []byte) ([]byte, []byte) {
	return value[:len(value)-sigSize], value[len(value)-sigSize:]
}

func (d *DHT) Store(key []byte, value []byte) error {
	if len(value) > maxValueSize {
		return ErrValueTooLarge
	}

	if len(value) < sigSize {
		return ErrValueTooSmall
	}

	if len(key) != keySize {
		return ErrKeyWrongSize
	}

	v, s := splitValue(value)

	if !ed25519.Verify(key, v, s) {
		return ErrSignatureInvalid
	}

	return d.store.Put(key, value, expire)
}

func (d *DHT) FindNode(nodeID node.NodeID) []node.Node {
	return nil
}

func (d *DHT) FindValue(out []byte, key []byte) ([]byte, []node.Node) {
	return nil, nil
}
