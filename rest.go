package kademlia

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"strings"
)

func checkMethod(methods []string, request *http.Request, w http.ResponseWriter) bool {
	for _, method := range methods {
		if method == request.Method {
			return true
		}
	}

	fmt.Fprintf(w, "This endpoint only works with %s", strings.Join(methods[:], " "))
	return false
}

func (node *Node) handlePingIP(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"GET"}, r, w) {
		return
	}

	ipString := r.URL.Path[len("/ping/ip/"):]

	addr, err := net.ResolveTCPAddr("", ipString)
	if err != nil {
		fmt.Fprintf(w, "Couldn't resolve IP address %s: %s", ipString, err)
		return
	}

	node.doPing(*addr)

	// TODO: Return diagnostic information from doPing
	// TODO: Logging

	fmt.Fprintf(w, "Sent PING (IP) to %s", addr.String())
}

func (node *Node) handlePingID(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"GET"}, r, w) {
		return
	}

	idString := r.URL.Path[len("/ping/id/"):]

	id := new(big.Int)
	_, success := id.SetString(idString, 16)
	if !success {
		fmt.Fprintf(w, "Invalid id: %s", idString)
		return
	}

	contact := node.rt.ContactFromID(*id)
	if contact == nil {
		fmt.Fprintf(w, "Could not find %s in routing table", id.String())
		return
	}

	addr := contact.Addr

	node.doPing(addr)

	// TODO: Return diagnostic information from doPing
	// TODO: Logging

	fmt.Fprintf(w, "Sent PING (ID) to server %s", addr.String())
}

func (node *Node) handleStore(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"POST"}, r, w) {
		return
	}

	key := r.URL.Path[len("/store/"):]
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error reading value")
	}

	// TODO: Store the key
	// TODO: Mark this node as the originator
	// TODO: Send back response
	fmt.Fprintf(w, "Called STORE for key %s with value %s", key, value)
}

func (node *Node) handleOneshotFindNode(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"GET"}, r, w) {
		return
	}

	id := r.URL.Path[len("/iterative/findnode/"):]

	// TODO: Check for valid id
	// TODO: Perform necessary stuff
	// TODO: Send back response
	fmt.Fprintf(w, "Called FINDNODE for server %s", id)
}

func (node *Node) handleOneshotFindValue(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"GET"}, r, w) {
		return
	}

	key := r.URL.Path[len("/oneshot/findvalue/"):]

	// TODO: Check for valid key
	// TODO: Perform necessary stuff
	// TODO: Send back response
	fmt.Fprintf(w, "Called FINDVALUE for value %s", key)
}

func (node *Node) handleIterativeFindNode(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"GET"}, r, w) {
		return
	}

	id := r.URL.Path[len("/iterative/findnode/"):]

	// TODO: Check for valid id
	// TODO: Perform necessary stuff
	// TODO: Send back response
	fmt.Fprintf(w, "Called FINDNODE for server %s", id)
}

func (node *Node) handleIterativeFindValue(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"GET"}, r, w) {
		return
	}

	key := r.URL.Path[len("/iterative/findvalue/"):]

	// TODO: Check for valid key
	// TODO: Perform necessary stuff
	// TODO: Send back response
	fmt.Fprintf(w, "Called FINDVALUE for value %s", key)
}

func (node *Node) handleShutdown(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"GET"}, r, w) {
		return
	}

	// TODO: Perform necessary stuff

	fmt.Fprintf(w, "Called SHUTDOWN")
}

// setupControlEndpoints registers handlers for the remote control REST API
func (node *Node) setupControlEndpoints() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "")
	})

	// Handle request to ping a specific server by IP address
	// GET /ping/ip/<ip addr>
	http.HandleFunc("/ping/ip/", func(w http.ResponseWriter, r *http.Request) {
		node.handlePingIP(w, r)
	})

	// Handle request to ping a specific server by ID
	// GET /ping/id/<id>
	http.HandleFunc("/ping/id/", func(w http.ResponseWriter, r *http.Request) {
		node.handlePingID(w, r)
	})

	// Handle request to store (key,value) in the DHT
	// This node becomes the originator
	// POST /store/<key hash>
	// Body is raw value
	http.HandleFunc("/store/", func(w http.ResponseWriter, r *http.Request) {
		node.handleStore(w, r)
	})

	// Handle oneshot request to find node with specific node id
	// GET /find/<id>
	http.HandleFunc("/oneshot/findnode/", func(w http.ResponseWriter, r *http.Request) {
		node.handleOneshotFindNode(w, r)
	})

	// Handle oneshot request to find specific value
	// GET /findvalue/<key hash>
	http.HandleFunc("/oneshot/findvalue/", func(w http.ResponseWriter, r *http.Request) {
		node.handleOneshotFindValue(w, r)
	})

	// Handle iterative request to find node with specific node id
	// GET /find/<id>
	http.HandleFunc("/iterative/findnode/", func(w http.ResponseWriter, r *http.Request) {
		node.handleIterativeFindNode(w, r)
	})

	// Handle iterative request to find specific value
	// GET /findvalue/<key hash>
	http.HandleFunc("/iterative/findvalue/", func(w http.ResponseWriter, r *http.Request) {
		node.handleIterativeFindValue(w, r)
	})

	// Handle request to shutdown server
	// GET /shutdown
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		node.handleShutdown(w, r)
	})
}
