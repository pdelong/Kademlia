package kademlia

import (
	"net"
	"sort"
	"sync"
)

// This file contains the iterative RPCs used for information progagation throughout nodes
// Calls STORE RPC on k Contacts ( Don't call on self?)
func (node *Node) doIterativeStore(key string, value []byte, dest net.TCPAddr) {
	shortlist := node.doIterativeFindNode(dest)

	// get k contacts and send STORE RPC to each
	for _, contact := range shortlist {
		go func() {
			args := StoreArgs{node.addr, key, value}
			var reply StoreReply
			if !node.doRPC("Store", contact.Addr, args, &reply) {
				return
			}
		}()
	}
}

func (node *Node) doIterativeFindValue(key string, dest net.TCPAddr) []byte {
	args := FindValueArgs{node.addr, key}
	var reply FindNodeReply

	if !node.doRPC("FindValue", dest, args, &reply) {
		return nil
	}

	// TODO: If we find the value, STORE RPC sent to closest contact that did not return value
	return nil
}

// Iteratively send a FINDNODE RPC
// Returns a shortlist of k closest nodes
func (node *Node) doIterativeFindNode(dest net.TCPAddr) []Contact {
	//Iterations continue until no contacts returned that are closer or if all contacts in shortlist are active (k contacts have been queried)
	toFind := NewContact(dest)
	contacted := make(map[string]bool)
	shortlist := make([]Contact, 0, 20)

	// Get nearest from own RT - a bit weird to do through RPCs for self though
	shortlist = node.doFindNode(dest)

	// while nearest contacts is not same, keep on iterating
	for {
		pinged := 0
		for i := 0; i < len(shortlist); i++ {
			contactChan := make(chan []Contact)
			var wg sync.WaitGroup
			wg.Add(alpha)

			toPing := shortlist[i].Addr
			_, exists := contacted[toPing.String()]

			if exists {
				continue
			}

			go func(shortlist []Contact) {
				defer wg.Done()
				newShortlist := node.doFindNode(toPing)
				contacted[toPing.String()] = true
				pinged++

				// update the shortlist
				shortlist = append(shortlist, newShortlist...)
				sort.Slice(shortlist, func(i, j int) bool {
					iDist := float64(distanceBetween(toFind.Id, shortlist[i].Id).Uint64())
					jDist := float64(distanceBetween(toFind.Id, shortlist[j].Id).Uint64())
					return iDist < jDist
				})

				contactChan <- shortlist[:k]
			}(shortlist)

			// Wait for all rpcs to return
			wg.Wait()

			// Check if shortlist has any closer nodes
			noCloser := 0
			for s := range contactChan {
				newClosestDist := float64(distanceBetween(toFind.Id, s[0].Id).Uint64())
				currFarthestDist := float64(distanceBetween(toFind.Id, shortlist[len(shortlist)-1].Id).Uint64())
				if newClosestDist >= currFarthestDist {
					noCloser++
				}
			}

			// Stop if already contacted k Contacts
			if len(contacted) >= k {
				shortlist = shortlist[:k]
				return shortlist
			}

			if pinged >= 3 {
				break
			}
		}
	}
}
