package kademlia

import (
	"math/big"
	"net"
	"sort"
)

// This file contains the iterative RPCs used for information progagation throughout nodes
// Calls STORE RPC on k Contacts ( Don't call on self?)
func (node *Node) doIterativeStore(key string, value []byte) {
	shortlist := node.doIterativeFindNode(key)

	// get k contacts and send STORE RPC to each
	for _, contact := range shortlist {
		go func(contact Contact) {
			args := StoreArgs{node.addr, key, value}
			var reply StoreReply
			if !node.doRPC("Store", contact.Addr, args, &reply) {
				return
			}
		}(contact)
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
func (node *Node) doIterativeFindNode(key string) []Contact {
	//Iterations continue until no contacts returned that are closer or if all contacts in shortlist are active (k contacts have been queried)
	toFindId := new(big.Int)
	toFindId.SetString(key, key_base)
	contacted := make(map[string]bool)
	shortlist := make([]Contact, 0, k)

	// add yourself to contacted
	contacted[node.addr.String()] = true

	shortlist = node.rt.findKNearestContacts(node.id)
	node.logger.Printf("Found %d contacts", len(shortlist))

	contactChan := make(chan []Contact)
	// while nearest contacts is not same, keep on iterating
	for {
		node.logger.Printf("Starting a new round of FindNodes")
		found := 0
		changed := false
		toSend := make([]Contact, 0, alpha)

		// find alpha to contact
		for i := 0; i < len(shortlist); i++ {
			toPing := shortlist[i].Addr
			_, exists := contacted[toPing.String()]

			if exists {
				continue
			}
			toSend = append(toSend, shortlist[i])
			found++
			if (found == alpha) {
				break
			}
		}

		// send alpha (or maybe fewer) RPCs
		for i := 0; i < len(toSend); i++ {
			toPing := toSend[i].Addr
			go func() {
				responseShortlist := node.doFindNode(toPing, key)
				contacted[toPing.String()] = true

				// update the shortlist
				sort.Slice(responseShortlist, func(i, j int) bool {
					iDist := distanceBetween(*toFindId, responseShortlist[i].Id)
					jDist := distanceBetween(*toFindId, responseShortlist[j].Id)
					return (iDist.Cmp(jDist) == -1)
				})
				slice_index := k
				if (len(responseShortlist) < k) {
					slice_index = len(responseShortlist)
				}
				contactChan <- responseShortlist[:slice_index]
			}()
		}

		updatedShortlist := make([]Contact, len(shortlist), k)
		copy(updatedShortlist, shortlist)
		node.logger.Printf("Shortlist length %d", len(updatedShortlist))
		closer := 0
		node.logger.Printf("Going to read from channel")
		for i := 0; i < len(toSend); i++ {
			s := <-contactChan
			newClosestDist := distanceBetween(*toFindId, s[0].Id)
			if (len(updatedShortlist) > 0) {
				currClosestDist := distanceBetween(*toFindId, updatedShortlist[0].Id)
				// if newClosestDist < currClosestDist
				if (newClosestDist.Cmp(currClosestDist) == -1) { 
					closer++
				}
			}

			updatedShortlist = append(updatedShortlist, s...)
			node.logger.Printf("Update: list length: %d\n", len(updatedShortlist))
			// update the shortlist
			sort.Slice(updatedShortlist, func(i, j int) bool {
				iDist := distanceBetween(*toFindId, updatedShortlist[i].Id)
				jDist := distanceBetween(*toFindId, updatedShortlist[j].Id)
				return (iDist.Cmp(jDist) == -1)
			})
			slice_index := k
			if (len(updatedShortlist) < k) {
				slice_index = len(updatedShortlist)
			}
			updatedShortlist = updatedShortlist[:slice_index]
		}

		node.logger.Printf("Finished reading from channel")

		// if we didn't find anything closer in last round, ping the rest of the
		// shortlist that are unseen
		if closer == 0 {
			sendingTo := make([]Contact, 0, k)
			for i := 0; i < len(shortlist); i++ {
				toPing := shortlist[i].Addr
				_, exists := contacted[toPing.String()]

				if !exists {
					sendingTo = append(sendingTo, shortlist[i])
				}
			}
			responseShortlist := node.findNodeToK(toFindId, sendingTo)
			updatedShortlist = append(updatedShortlist, responseShortlist...)
			// update the shortlist
			sort.Slice(updatedShortlist, func(i, j int) bool {
				iDist := distanceBetween(*toFindId, updatedShortlist[i].Id)
				jDist := distanceBetween(*toFindId, updatedShortlist[j].Id)
				return (iDist.Cmp(jDist) == -1)
			})
			slice_index := k
			if (len(updatedShortlist) < k) {
				slice_index = len(updatedShortlist)
			}
			updatedShortlist = updatedShortlist[:slice_index]
		}

		node.logger.Printf("Checking if shortlist has changed")
		// check if the shortlist has changed at all
		// if not, we should terminate
		// comparing shortlist and updatedShortlist
		node.logger.Printf("New shortlist has length %d", len(updatedShortlist))
		changed = false
		loop_index := len(updatedShortlist)
		if (len(shortlist) < loop_index) {
			loop_index = len(shortlist)
		}
		for i := 0; i < loop_index; i++ {
			if updatedShortlist[i].Addr.String() != shortlist[i].Addr.String() {
				changed = true
			}
		}
		node.logger.Println("Shortlist changed this round", changed)
		if !changed {
			return updatedShortlist
		}

		shortlist = updatedShortlist

	}
	return shortlist
}

func (node *Node) findNodeToK(toFindId *big.Int, toSend []Contact) []Contact {
	contactChan := make(chan []Contact)

	for i := 0; i < len(toSend); i++ {
		toPing := toSend[i].Addr
		go func() {
			responseShortlist := node.doFindNode(toPing, toFindId.Text(key_base))

			contactChan <- responseShortlist
		}()
	}

	// Wait for all rpcs to return
	updatedShortlist := make([]Contact, 0)
	for i := 0; i < len(toSend); i++ {
		s := <-contactChan
		updatedShortlist = append(updatedShortlist, s...)
		// update the shortlist
		sort.Slice(updatedShortlist, func(i, j int) bool {
			iDist := distanceBetween(*toFindId, updatedShortlist[i].Id)
			jDist := distanceBetween(*toFindId, updatedShortlist[j].Id)
			return (iDist.Cmp(jDist) == -1)
		})
		slice_index := k
		if (len(updatedShortlist) < k) {
			slice_index = len(updatedShortlist)
		}
		updatedShortlist = updatedShortlist[:slice_index]
	}
	
	return updatedShortlist
}
