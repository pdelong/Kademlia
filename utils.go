package kademlia

import (
	"crypto/sha1"
	"fmt"
	"math/big"
	"net"
)

func RemoveDupesFromShortlist(contacts []Contact) []Contact {
	// make a map from the list
	undupe_map := make(map[string]bool)
	unduped_slice := make([]Contact, 0, len(contacts))
	for i := 0; i < len(contacts); i++ {
		curr := contacts[i]
		_, exist := undupe_map[curr.Addr.String()]
		if !exist {
			unduped_slice = append(unduped_slice, curr)
			undupe_map[curr.Addr.String()] = true
		}
	}

	return unduped_slice
}

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
