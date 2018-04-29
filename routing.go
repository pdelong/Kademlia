package kademlia

import (
	//"crypto/sha1"
	//"errors"
	//"fmt"
	//"log"
	//"math/big"
	"net"
	//"net/http"
	//"net/rpc"
	//"os"
	//"time"
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
	return value

}

// Iteratively send a FINDNODE RPC
// Returns a shortlist of k closest nodes
func (node *Node) doIterativeFindNode(dest net.TCPAddr) []Contact {
	//Iterations continue until no contacts returned that are closer or if all contacts in shortlist are active (k contacts have been queried)

	contacted := make(map[string]bool)
	shortlist := make([]Contact, 0, 20)

	// Get nearest from own RT - a bit weird to do through RPCs for self though
	shortlist = node.doFindNode(dest)

	//closest := shortlist[0]

	// while nearest contacts is not same, keep on iterating
	for i := 0; i < alpha; i++ {
		toPing := shortlist[i].Addr
		go func() {
			node.doFindNode(toPing)
			contacted[toPing.String()] = true
		}()
	}
}
