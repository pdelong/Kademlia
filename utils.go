package kademlia

import (
	"crypto/sha1"
	"fmt"
	"math/big"
	"net"
)



// GetKBucketFromAddr returns the KBucket that would contain destAddr
func (node *Node) GetKBucketFromAddr(destAddr net.TCPAddr) int {
	hash := sha1.Sum([]byte(destAddr.String()))

	id := *big.NewInt(0)
	id.SetBytes(hash[:])

	return node.GetKBucketFromID(&id)
}

// GetKBucketFromID returns the KBucket that would contain destID
func (node *Node) GetKBucketFromID(destID *big.Int) int {
	destContact := Contact{*destID, net.TCPAddr{}}
	dist := node.distanceTo(&destContact)
	node.logger.Printf("Distance is %s", dist)

	// a kludgy hack to the get the floor of log_2 of the distance
	bitstring := fmt.Sprintf("%b", dist)
	log2Floor := len(bitstring) - 1

	return log2Floor
}
