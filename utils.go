package kademlia

import (
	"crypto/sha1"
	"fmt"
	"math/big"
	"net"
)

func (node *Node) GetKBucketFromAddr(dest_addr net.TCPAddr) int {
	hash := sha1.Sum([]byte(dest_addr.String()))

	id := *big.NewInt(0)
	id.SetBytes(hash[:])

	return node.GetKBucketFromId(&id)
}

func (node *Node) GetKBucketFromId(dest_id *big.Int) int {
	dest_contact := Contact{*dest_id, net.TCPAddr{}}
	dist := node.distanceTo(&dest_contact)
	node.logger.Printf("Distance is %s", dist)

	// a kludgy hack to the get the floor of log_2 of the distance
	bitstring := fmt.Sprintf("%b", dist)
	log_2_floor := len(bitstring) - 1

	return log_2_floor
}

