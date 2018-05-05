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
		changed := false
		for i := 0; i < len(shortlist); i++ {
			contactChan := make(chan []Contact)
			var wg sync.WaitGroup
			wg.Add(alpha)

			toPing := shortlist[i].Addr
			_, exists := contacted[toPing.String()]

			if exists {
				continue
			}

			pinged++
			go func() {
				defer wg.Done()
				responseShortlist := node.doFindNode(toPing)
				contacted[toPing.String()] = true

				// update the shortlist
				sort.Slice(responseShortlist, func(i, j int) bool {
					iDist := distanceBetween(toFind.Id, responseShortlist[i].Id)
					jDist := distanceBetween(toFind.Id, responseShortlist[j].Id)
					return (iDist.Cmp(jDist) == -1)
				})
				contactChan <- responseShortlist[:k]
			}()

			// Wait for all rpcs to return
			wg.Wait()

			// Check if shortlist has any closer nodes
			noCloser := 0
			updatedShortlist := make([]Contact, 0, 20)
			copy(updatedShortlist, shortlist)

			for s := range contactChan {
				newClosestDist := distanceBetween(toFind.Id, s[0].Id)
				currClosestDist := distanceBetween(toFind.Id, updatedShortlist[0].Id)
				// if newClosestDist >= currClosestDist
				if (newClosestDist.Cmp(currClosestDist) == 0) || 
				   (newClosestDist.Cmp(currClosestDist) == 1) {
					noCloser++
				} 

				updatedShortlist = append(updatedShortlist, s...)
				// update the shortlist
				sort.Slice(updatedShortlist, func(i, j int) bool {
					iDist := distanceBetween(toFind.Id, updatedShortlist[i].Id)
					jDist := distanceBetween(toFind.Id, updatedShortlist[j].Id)
					return (iDist.Cmp(jDist) == -1)
				})
				updatedShortlist = updatedShortlist[:k]
			}

			// if we didn't find anything closer in last round, ping the rest of the 
			// shortlist that are unseen
			if (noCloser == alpha) {
				sendingTo := make([]Contact, 0, 20)
				for i := 0; i < len(shortlist); i++ {
					toPing := shortlist[i].Addr
					_, exists := contacted[toPing.String()]

					if !exists {
						sendingTo = append(sendingTo, shortlist[i])
					}
				}
				responseShortlist := node.findNodeToK(toFind, sendingTo)
				updatedShortlist = append(updatedShortlist, responseShortlist...)
				// update the shortlist
				sort.Slice(updatedShortlist, func(i, j int) bool {
					iDist := distanceBetween(toFind.Id, updatedShortlist[i].Id)
					jDist := distanceBetween(toFind.Id, updatedShortlist[j].Id)
					return (iDist.Cmp(jDist) == -1)
				})
				updatedShortlist = updatedShortlist[:k]
			}

			// check if the shortlist has changed at all
			// if not, we should terminate
			// comparing shortlist and updatedShortlist
			changed = false
			for i := 0; i < len(shortlist); i++ {
				if updatedShortlist[i].Addr.String() != shortlist[i].Addr.String() {
					changed = true
				}
			}
			if !changed {
				return updatedShortlist
			} else {
				shortlist = updatedShortlist
			}

			if pinged >= 3 {
				break
			}
		}
	}
	return shortlist
}

func (node *Node) findNodeToK(toFind *Contact, toSend []Contact) []Contact {
	contactChan := make(chan []Contact)
	var wg sync.WaitGroup
	wg.Add(len(toSend))

	for i := 0; i < len(toSend); i++ {
		toPing := toSend[i].Addr
		go func() {
			defer wg.Done()
			responseShortlist := node.doFindNode(toPing)

			contactChan <- responseShortlist
		}()
	}

	// Wait for all rpcs to return
	wg.Wait()
	updatedShortlist := make([]Contact, 0)
	for s := range contactChan {
		updatedShortlist = append(updatedShortlist, s...)
		// update the shortlist
		sort.Slice(updatedShortlist, func(i, j int) bool {
			iDist := distanceBetween(toFind.Id, updatedShortlist[i].Id)
			jDist := distanceBetween(toFind.Id, updatedShortlist[j].Id)
			return (iDist.Cmp(jDist) == -1)
		})
		updatedShortlist = updatedShortlist[:k]
	}
	return updatedShortlist
}
