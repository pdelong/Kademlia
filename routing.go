package kademlia

import (
	"math/big"
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

func (node *Node) doIterativeFindValue(key string) []byte {
	value, found := node.ht.get(key)
	if found {
		return value
	}
	//Iterations continue until no contacts returned that are closer or if all contacts in shortlist are active (k contacts have been queried)
	toFindID := new(big.Int)
	toFindID.SetString(key, keyBase)
	contacted := make(map[string]bool)
	shortlist := make([]Contact, 0, k)

	// add yourself to contacted
	contacted[node.addr.String()] = true

	shortlist = node.rt.findKNearestContacts(*toFindID)
	node.logger.Printf("Found %d contacts", len(shortlist))

	contactChan := make(chan []Contact)
	valueChan := make(chan []byte)
	// while nearest contacts is not same, keep on iterating
	for {
		node.logger.Printf("Starting a new round of FindValues")
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
			if found == alpha {
				break
			}
		}

		// send alpha (or maybe fewer) RPCs
		for i := 0; i < len(toSend); i++ {
			toPing := toSend[i].Addr
			go func() {
				response := node.doFindValue(key, toPing)
				if response == nil {
					// Error with performing doFindValue, ignoring for now
					// TODO: Handle error (?)
				} else if response.Val != nil {
					valueChan <- response.Val
					return
				}

				responseShortlist := response.Contacts
				contacted[toPing.String()] = true

				// update the shortlist
				sort.Slice(responseShortlist, func(i, j int) bool {
					iDist := distanceBetween(*toFindID, responseShortlist[i].Id)
					jDist := distanceBetween(*toFindID, responseShortlist[j].Id)
					return (iDist.Cmp(jDist) == -1)
				})
				sliceIndex := k
				if len(responseShortlist) < k {
					sliceIndex = len(responseShortlist)
				}
				contactChan <- responseShortlist[:sliceIndex]
			}()
		}

		updatedShortlist := make([]Contact, len(shortlist), k)
		copy(updatedShortlist, shortlist)
		node.logger.Printf("Shortlist length %d", len(updatedShortlist))
		closer := 0
		node.logger.Printf("Going to read from channel")
		for i := 0; i < len(toSend); i++ {
			var s []Contact
			select {
			case val := <-valueChan:
				return val
			case s = <-contactChan:
			}
			newClosestDist := distanceBetween(*toFindID, s[0].Id)
			if len(updatedShortlist) > 0 {
				currClosestDist := distanceBetween(*toFindID, updatedShortlist[0].Id)
				// if newClosestDist < currClosestDist
				if newClosestDist.Cmp(currClosestDist) == -1 {
					closer++
				}
			}

			updatedShortlist = append(updatedShortlist, s...)
			node.logger.Printf("Update: list length: %d\n", len(updatedShortlist))
			// update the shortlist
			sort.Slice(updatedShortlist, func(i, j int) bool {
				iDist := distanceBetween(*toFindID, updatedShortlist[i].Id)
				jDist := distanceBetween(*toFindID, updatedShortlist[j].Id)
				return (iDist.Cmp(jDist) == -1)
			})
			sliceIndex := k
			if len(updatedShortlist) < k {
				sliceIndex = len(updatedShortlist)
			}
			updatedShortlist = updatedShortlist[:sliceIndex]
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
			value, responseShortlist := node.findValueToK(toFindID, sendingTo)
			if value != nil {
				return value
			}
			updatedShortlist = append(updatedShortlist, responseShortlist...)
			// update the shortlist
			sort.Slice(updatedShortlist, func(i, j int) bool {
				iDist := distanceBetween(*toFindID, updatedShortlist[i].Id)
				jDist := distanceBetween(*toFindID, updatedShortlist[j].Id)
				return (iDist.Cmp(jDist) == -1)
			})
			sliceIndex := k
			if len(updatedShortlist) < k {
				sliceIndex = len(updatedShortlist)
			}
			updatedShortlist = updatedShortlist[:sliceIndex]
		}

		node.logger.Printf("Checking if shortlist has changed")
		// check if the shortlist has changed at all
		// if not, we should terminate
		// comparing shortlist and updatedShortlist
		node.logger.Printf("New shortlist has length %d", len(updatedShortlist))
		changed = false
		loopIndex := len(updatedShortlist)
		if len(shortlist) < loopIndex {
			loopIndex = len(shortlist)
		}
		for i := 0; i < loopIndex; i++ {
			if updatedShortlist[i].Addr.String() != shortlist[i].Addr.String() {
				changed = true
			}
		}
		node.logger.Println("Shortlist changed this round", changed)
		if !changed {
			return nil
		}

		shortlist = updatedShortlist
	}
}

// Iteratively send a FINDNODE RPC
// Returns a shortlist of k closest nodes
func (node *Node) doIterativeFindNode(key string) []Contact {
	//Iterations continue until no contacts returned that are closer or if all contacts in shortlist are active (k contacts have been queried)
	toFindID := new(big.Int)
	toFindID.SetString(key, keyBase)
	contacted := make(map[string]bool)
	shortlist := make([]Contact, 0, k)

	// add yourself to contacted
	contacted[node.addr.String()] = true

	shortlist = node.rt.findKNearestContacts(*toFindID)
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
			if found == alpha {
				break
			}
		}

		// send alpha (or maybe fewer) RPCs
		for i := 0; i < len(toSend); i++ {
			toPing := toSend[i].Addr
			go func() {
				responseShortlist := node.doFindNode(key, toPing)
				contacted[toPing.String()] = true

				// update the shortlist
				sort.Slice(responseShortlist, func(i, j int) bool {
					iDist := distanceBetween(*toFindID, responseShortlist[i].Id)
					jDist := distanceBetween(*toFindID, responseShortlist[j].Id)
					return (iDist.Cmp(jDist) == -1)
				})
				sliceIndex := k
				if len(responseShortlist) < k {
					sliceIndex = len(responseShortlist)
				}
				contactChan <- responseShortlist[:sliceIndex]
			}()
		}

		updatedShortlist := make([]Contact, len(shortlist), k)
		copy(updatedShortlist, shortlist)
		node.logger.Printf("Shortlist length %d", len(updatedShortlist))
		closer := 0
		node.logger.Printf("Going to read from channel")
		for i := 0; i < len(toSend); i++ {
			s := <-contactChan
			newClosestDist := distanceBetween(*toFindID, s[0].Id)
			if len(updatedShortlist) > 0 {
				currClosestDist := distanceBetween(*toFindID, updatedShortlist[0].Id)
				// if newClosestDist < currClosestDist
				if newClosestDist.Cmp(currClosestDist) == -1 {
					closer++
				}
			}

			updatedShortlist = append(updatedShortlist, s...)
			node.logger.Printf("Update: list length: %d\n", len(updatedShortlist))
			// update the shortlist
			sort.Slice(updatedShortlist, func(i, j int) bool {
				iDist := distanceBetween(*toFindID, updatedShortlist[i].Id)
				jDist := distanceBetween(*toFindID, updatedShortlist[j].Id)
				return (iDist.Cmp(jDist) == -1)
			})
			sliceIndex := k
			if len(updatedShortlist) < k {
				sliceIndex = len(updatedShortlist)
			}
			updatedShortlist = updatedShortlist[:sliceIndex]
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
			responseShortlist := node.findNodeToK(toFindID, sendingTo)
			updatedShortlist = append(updatedShortlist, responseShortlist...)
			// update the shortlist
			sort.Slice(updatedShortlist, func(i, j int) bool {
				iDist := distanceBetween(*toFindID, updatedShortlist[i].Id)
				jDist := distanceBetween(*toFindID, updatedShortlist[j].Id)
				return (iDist.Cmp(jDist) == -1)
			})
			sliceIndex := k
			if len(updatedShortlist) < k {
				sliceIndex = len(updatedShortlist)
			}
			updatedShortlist = updatedShortlist[:sliceIndex]
		}

		node.logger.Printf("Checking if shortlist has changed")
		// check if the shortlist has changed at all
		// if not, we should terminate
		// comparing shortlist and updatedShortlist
		node.logger.Printf("New shortlist has length %d", len(updatedShortlist))
		changed = false
		loopIndex := len(updatedShortlist)
		if len(shortlist) < loopIndex {
			loopIndex = len(shortlist)
		}
		for i := 0; i < loopIndex; i++ {
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

	//return shortlist
}

func (node *Node) findNodeToK(toFindID *big.Int, toSend []Contact) []Contact {
	contactChan := make(chan []Contact)

	for i := 0; i < len(toSend); i++ {
		toPing := toSend[i].Addr
		go func() {
			responseShortlist := node.doFindNode(toFindID.Text(keyBase), toPing)

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
			iDist := distanceBetween(*toFindID, updatedShortlist[i].Id)
			jDist := distanceBetween(*toFindID, updatedShortlist[j].Id)
			return (iDist.Cmp(jDist) == -1)
		})
		sliceIndex := k
		if len(updatedShortlist) < k {
			sliceIndex = len(updatedShortlist)
		}
		updatedShortlist = updatedShortlist[:sliceIndex]
	}

	return updatedShortlist
}

func (node *Node) findValueToK(toFindID *big.Int, toSend []Contact) ([]byte, []Contact) {
	contactChan := make(chan []Contact)
	valueChan := make(chan []byte)

	for i := 0; i < len(toSend); i++ {
		toPing := toSend[i].Addr
		go func() {
			response := node.doFindValue(toFindID.Text(keyBase), toPing)
			if response == nil {
				// Error with performing doFindValue, ignoring for now
				// TODO: Handle error (?)
			} else if response.Val != nil {
				valueChan <- response.Val
				return
			}

			responseShortlist := response.Contacts
			contactChan <- responseShortlist
		}()
	}

	// Wait for all rpcs to return
	updatedShortlist := make([]Contact, 0)
	for i := 0; i < len(toSend); i++ {
		var s []Contact
		select {
		case val := <-valueChan:
			return val, nil
		case s = <-contactChan:
		}
		updatedShortlist = append(updatedShortlist, s...)
		// update the shortlist
		sort.Slice(updatedShortlist, func(i, j int) bool {
			iDist := distanceBetween(*toFindID, updatedShortlist[i].Id)
			jDist := distanceBetween(*toFindID, updatedShortlist[j].Id)
			return (iDist.Cmp(jDist) == -1)
		})
		sliceIndex := k
		if len(updatedShortlist) < k {
			sliceIndex = len(updatedShortlist)
		}
		updatedShortlist = updatedShortlist[:sliceIndex]
	}

	return nil, updatedShortlist
}
