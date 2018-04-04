package kademlia

import (
	"container/list"
	"crypto/sha1"
	"math"
	"math/big"
	"net"
)

// TODO: The naming conventions are atrocious
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

// Returns true if contact Id and Addr are equivalent
// structs can be compared, but structs containing big.Int cannot
func AreEqualContacts(a *Contact, b *Contact) bool {
	return (a.Id.Cmp(&b.Id) == 0 && a.Addr.String() == b.Addr.String())
}

// extra struct because we will want to implement split bucket
type RoutingTable struct {
	owner    *Node
	kBuckets []*KBucket
	k        int
	tRefresh int
}

func NewRoutingTable(owner *Node, k int, tRefresh int) *RoutingTable {
	kBuckets := make([]*KBucket, 160, 160)
	rt := RoutingTable{owner, kBuckets, k, tRefresh}
	return &rt
}

// section 2.4 Kademlia protocol splits bucket when full and range includes own ID
// TODO
func (self *RoutingTable) splitBucket() {
}

//TODO
func (self *RoutingTable) findContact() {

}

func (self *RoutingTable) add(contact Contact) {
	dist := float64(self.owner.distanceTo(&contact).Uint64())
	index := int(math.Ceil(math.Log(dist) / math.Log(2))) // Find which bucket it belongs to
	if self.kBuckets[index] == nil {
		self.kBuckets[index] = NewKBucket(20)
	} else {
		self.kBuckets[index].addContact(contact)
	}
	//TODO: handle failure to add
}

func (self *RoutingTable) remove(contact Contact) {
	dist := float64(self.owner.distanceTo(&contact).Uint64())
	index := int(math.Ceil(math.Log(dist) / math.Log(2))) // Find which bucket it belongs to
	self.kBuckets[index].removeContact(contact)
}

// Not even sure if we will use this
func (self *RoutingTable) clear() {
	// Note that this sets slice capacity to 0
	self.kBuckets = nil
}

type KBucket struct {
	contacts list.List
	k        int       // max number of contacts
	lruCache list.List // not implemented yet but explained in section 4.1
}

func NewKBucket(k int) *KBucket {
	contacts := *list.New()
	lruCache := *list.New()
	kBucket := KBucket{contacts, k, lruCache}
	return &kBucket
}

// If bucket contains contact, returns ptr to element in list. Else, returns nil
func (self *KBucket) getFromList(contact Contact) *list.Element {
	for e := self.contacts.Front(); e != nil; e = e.Next() {
		curr, _ := e.Value.(Contact)
		// TODO: handle error when element can't be cast to Contact
		if AreEqualContacts(&curr, &contact) {
			return e
		}
	}
	return nil
}

// Returns true if contact is added into bucket, false otherwise
func (self *KBucket) addContact(contact Contact) bool {
	// If contact exists, move to tail
	element := self.getFromList(contact)
	if element != nil {
		self.contacts.MoveToBack(element)
		return true
	} else {
		// If bucket isn't full, add to tail
		// list.Len() = O(1)
		if self.contacts.Len() < self.k {
			self.contacts.PushBack(contact)
			return true
		}
		// Otherwise, ping least-recently seen node
		lruNode := self.contacts.Front()
		// TODO: ping node... sigh this is gnna be ugly.
		if true {
			// If no response, node is evicted and new sender is inserted at tail
			self.contacts.Remove(lruNode)
			self.contacts.PushBack(contact)
			return true
		}
		// TODO: implement replacement cache
		return false
	}
}

// Returns true if contact exists, false otherwise
func (self *KBucket) removeContact(contact Contact) bool {
	element := self.getFromList(contact)
	if element != nil {
		self.contacts.Remove(element)
		return true
	} else {
		return false
	}
}
