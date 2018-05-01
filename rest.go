package kademlia

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"strings"
)

// CommandMessage is a container for messages sent from the REST api to the node
type CommandMessage struct {
	Command string
	Arg1    interface{}
	Arg2    interface{}
	Resp    chan interface{} //TODO: There might be a better option. string?
}

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

	contact := NewContact(*addr)

	c := make(chan interface{})
	node.restC <- CommandMessage{"PING", contact, nil, c}

	<-c

	// TODO: Send back response
	fmt.Fprintf(w, "REST: Received PING (IP) for server %s", addr.String())
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

	// TODO: See if id in kbuckets, else return error

	c := make(chan interface{})
	node.restC <- CommandMessage{"PING", nil, id, c}

	<-c

	// TODO: Send back response
	fmt.Fprintf(w, "Called PING (ID) for server %s", id)
}

func (node *Node) handleOneshotStore(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"POST"}, r, w) {
		return
	}

	key := r.URL.Path[len("/oneshot/store/"):]
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error reading value")
	}

	c := make(chan interface{})
	node.restC <- CommandMessage{"STORE", key, value, c}

	<-c

	// TODO: Send back response
	fmt.Fprintf(w, "Called STORE for key %s with value %s", key, value)
}

func (node *Node) handleOneshotFindNode(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"GET"}, r, w) {
		return
	}

	id := r.URL.Path[len("/oneshot/findnode/"):]

	// TODO: Check for valid id
	// TODO: Get response

	c := make(chan interface{})
	node.restC <- CommandMessage{"FINDNODE", id, nil, c}

	fmt.Fprintf(w, "Called FINDNODE for server %s", id)
}

func (node *Node) handleOneshotFindValue(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"GET"}, r, w) {
		return
	}

	key := r.URL.Path[len("/oneshot/findvalue/"):]

	// TODO: Check for valid key
	// TODO: Get response

	c := make(chan interface{})
	node.restC <- CommandMessage{"FINDVALUE", key, nil, c}

	fmt.Fprintf(w, "Called FINDVALUE for value %s", key)
}

func (node *Node) handleIterativeStore(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"POST"}, r, w) {
		return
	}

	key := r.URL.Path[len("/iterative/store/"):]
	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error reading value")
	}

	c := make(chan interface{})
	node.restC <- CommandMessage{"STORE", key, value, c}

	<-c

	// TODO: Send back response
	fmt.Fprintf(w, "Called STORE for key %s with value %s", key, value)
}

func (node *Node) handleIterativeFindNode(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"GET"}, r, w) {
		return
	}

	id := r.URL.Path[len("/iterative/findnode/"):]

	// TODO: Check for valid id
	// TODO: Get response

	c := make(chan interface{})
	node.restC <- CommandMessage{"FINDNODE", id, nil, c}

	fmt.Fprintf(w, "Called FINDNODE for server %s", id)
}

func (node *Node) handleIterativeFindValue(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"GET"}, r, w) {
		return
	}

	key := r.URL.Path[len("/iterative/findvalue/"):]

	// TODO: Check for valid key
	// TODO: Get response

	c := make(chan interface{})
	node.restC <- CommandMessage{"FINDVALUE", key, nil, c}

	fmt.Fprintf(w, "Called FINDVALUE for value %s", key)
}

func (node *Node) handleShutdown(w http.ResponseWriter, r *http.Request) {
	if !checkMethod([]string{"GET"}, r, w) {
		return
	}

	node.restC <- CommandMessage{"SHUTDOWN", nil, nil, nil}

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

	// Handle oneshot request to store (key,value) in the DHT
	// POST /store/<key hash>
	// Body is raw value
	http.HandleFunc("/oneshot/store/", func(w http.ResponseWriter, r *http.Request) {
		node.handleOneshotStore(w, r)
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

	// Handle iterative request to store (key,value) in the DHT
	// POST /store/<key hash>
	// Body is raw value
	http.HandleFunc("/iterative/store/", func(w http.ResponseWriter, r *http.Request) {
		node.handleIterativeStore(w, r)
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
