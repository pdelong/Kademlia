package kademlia

import (
	"container/list"
	"crypto/sha1"
	"fmt"
	"math/big"
)

// An entry in the k-bucket
type NodeEntry struct {
	id   *big.Int
	ip   string
	port int
}

type Node struct {
	id       *big.Int
	address  string
	port     int
	kBuckets []list.List // TODO: Might want array of arrays
}

func (self *Node) String() string {
	return fmt.Sprintf("Node: (id = %v) (address = %s:%d) (kBuckets = %v)",
		self.id,
		self.address,
		self.port,
		self.kBuckets)
}

// Return XOR distance between self and other
func (self *Node) distanceTo(other NodeEntry) *big.Int {
	return big.NewInt(0).Xor(self.id, other.id)
}

func NewNode(address string, port int) *Node {
	node := new(Node)
	node.address = address
	node.port = port

	// #eww
	hash := sha1.Sum([]byte(fmt.Sprintf("%s:%d", address, port)))

	node.id = big.NewInt(0)
	node.id.SetBytes(hash[:])

	return node
}
