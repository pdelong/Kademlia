package kademlia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"os"
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

	node.logger.Printf("Performing IP PING of %s", addr)

	if node.doPing(*addr) {
		fmt.Fprintf(w, "Host %s successfully pinged", ipString)
	} else {
		fmt.Fprintf(w, "PING of Host %s unsuccessful", ipString)
	}
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

	node.logger.Printf("Performing ID PING of %s", id.String())

	contact := node.rt.ContactFromID(*id)
	if contact == nil {
		fmt.Fprintf(w, "Could not find %s in routing table", id.String())
		node.logger.Printf("Could not find %s in the routing table", id.String())
		return
	}

	addr := contact.Addr

	node.doPing(addr)

	if node.doPing(addr) {
		fmt.Fprintf(w, "Host %s successfully pinged", id.String())
	} else {
		fmt.Fprintf(w, "PING of Host %s unsuccessful", id.String())
	}
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

	node.logger.Printf("Received STORE for key: (%s), value: (%s)", key, value)

	node.ht.add(key, value, true)

	fmt.Fprintf(w, "Successfully stored key (%s)", key)
}

func (node *Node) handleGetTable(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"GET"}, r, w) {
		return
	}

	ch := node.ht.Iterator()

	a := make([]map[string]interface{}, 0)
	for val := range ch {
		d := make(map[string]interface{})
		d["key"] = val.key
		d["value"] = val.val
		fmt.Println(val.val)
		d["isOrigin"] = val.isOrigin
		a = append(a, d)
	}

	enc := json.NewEncoder(w)
	enc.Encode(a)
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

	node.logger.Println("Shutdown received. Terminating")

	fmt.Fprintf(w, "Called SHUTDOWN")

	os.Exit(0)
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
	// POST /store/<key>
	// Body is raw value
	http.HandleFunc("/store/", func(w http.ResponseWriter, r *http.Request) {
		node.handleStore(w, r)
	})

	http.HandleFunc("/table", func(w http.ResponseWriter, r *http.Request) {
		node.handleGetTable(w, r)
	})

	// Handle oneshot request to find node with specific node id
	// GET /find/<id>
	http.HandleFunc("/oneshot/findnode/", func(w http.ResponseWriter, r *http.Request) {
		node.handleOneshotFindNode(w, r)
	})

	// Handle oneshot request to find specific value
	// GET /findvalue/<key>
	http.HandleFunc("/oneshot/findvalue/", func(w http.ResponseWriter, r *http.Request) {
		node.handleOneshotFindValue(w, r)
	})

	// Handle iterative request to find node with specific node id
	// GET /find/<id>
	http.HandleFunc("/iterative/findnode/", func(w http.ResponseWriter, r *http.Request) {
		node.handleIterativeFindNode(w, r)
	})

	// Handle iterative request to find specific value
	// GET /findvalue/<key>
	http.HandleFunc("/iterative/findvalue/", func(w http.ResponseWriter, r *http.Request) {
		node.handleIterativeFindValue(w, r)
	})

	// Handle request to shutdown server
	// GET /shutdown
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		node.handleShutdown(w, r)
	})
}
