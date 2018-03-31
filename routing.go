package kademlia

import (
	"crypto/sha1"
	"math/big"
	"net"
)

// We should do all kbucket-related work here
// Preferably a struct with functions we can use to get/store easily in kbuckets

// An entry in the k-bucket
type Contact struct {
	Id   big.Int
	Addr net.TCPAddr
}

// Create a new Contact struct based on addr by taking the hash
func NewContact(addr net.TCPAddr) *Contact {
	hash := sha1.Sum([]byte(addr.String()))

	id := *big.NewInt(0)
	id.SetBytes(hash[:])

	nodeEntry := Contact{id, addr}
	return &nodeEntry
}

type KBucket struct {
}

func NewKBucket() *KBucket {
	kBucket := KBucket{}
	return &kBucket
}

func (self *KBucket) addContact(contact Contact) {
}

func (self *KBucket) removeContact(contact Contact) {
}

//func (self *KBucket) findNeighbor() {
//}
// yada yada
// Will most likely include helper functions to search for neighbors
