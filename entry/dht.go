package entry

import (
	"net"
	"sync"
	"time"

	"github.com/flynn/noise"
)

type DHT struct {
	Self *Contact

	RoutingTable [][]*Contact
}

const (
	ALPHA = 3   // Concurrency
	B     = 256 // Bits
	K     = 20

	EXPIRE    = 86400 * time.Second
	REFRESH   = 3600 * time.Second
	REPLICATE = 3600 * time.Second
	REPUBLISH = 86400 * time.Second
)

type Contact struct {
	Id   NodeID
	Port uint16
	IP   net.IPAddr
}

type DhtOptions struct {
	IP         net.IPAddr
	Port       uint16
	PrivateKey []byte
	PublicKey  []byte
	Key noise.DHKey
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
