package kademlia

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

// Node is an individual Kademlia node
type Node struct {
	id     big.Int
	addr   net.TCPAddr
	ht     map[string][]byte
	rt     *RoutingTable
	logger *log.Logger
	restC  chan CommandMessage
}

// PingArgs contains the arguments for the PING RPC
type PingArgs struct {
	Source net.TCPAddr
}

// PingReply contains the results for the PING RPC
type PingReply struct {
	Source net.TCPAddr
}

// StoreArgs contains the arguments for the STORE RPC
type StoreArgs struct {
	Source net.TCPAddr
	Key    string
	Val    []byte
}

// StoreReply contains the results for the Store RPC
type StoreReply struct {
}

// FindValueArgs contains the arguments for the FINDVALUE RPC
type FindValueArgs struct {
	Source net.TCPAddr
	Key    string
}

// FindValueReply contains the results for the FINDVALUE RPC
type FindValueReply struct {
	Val      []byte
	Contacts []Contact
}

// FindNodeArgs contains the arguments for the FINDNODE RPC
type FindNodeArgs struct {
	Source net.TCPAddr
	//Key    string
}

// FindNodeReply contains the results for the FINDNODE RPC
type FindNodeReply struct {
	Contacts []Contact
}

// Ping is the handler for the PING RPC
func (node *Node) Ping(args PingArgs, reply *PingReply) error {
	contact := NewContact(args.Source)

	if contact == nil {
		return errors.New("Couldn't hash IP address")
	}
	node.rt.add(*contact)

	node.logger.Printf("Ping from %s", args.Source.String())

	// Update k-bucket based on args.Source
	*reply = PingReply{node.addr}
	return nil
}

// Store is the handler for the STORE RPC
func (node *Node) Store(args StoreArgs, reply *StoreReply) error {
	contact := NewContact(args.Source)
	if contact == nil {
		return errors.New("Couldn't hash IP address")
	}
	node.rt.add(*contact)

	node.ht[args.Key] = args.Val
	*reply = StoreReply{}
	return nil
}

// FindValue is the handler for the FINDVALUE RPC
func (node *Node) FindValue(args FindValueArgs, reply *FindValueReply) error {
	contact := NewContact(args.Source)
	if contact == nil {
		return errors.New("Couldn't hash IP address")
	}
	node.rt.add(*contact)
	// If node contains key, returns associated data
	if val, ok := node.ht[args.Key]; ok {
		*reply = FindValueReply{Val: val}
		return nil
	}

	// Otherwise, return set of k triples (equiv. to FindNode)
	kNearest := node.rt.findKNearestContacts(contact.Id)
	*reply = FindValueReply{Contacts: kNearest}
	return nil
}

// FindNode is the handler for the FINDNODE RPC
func (node *Node) FindNode(args FindNodeArgs, reply *FindNodeReply) error {
	contact := NewContact(args.Source)
	if contact == nil {
		return errors.New("Couldn't hash IP address")
	}
	node.rt.add(*contact)

	kNearest := node.rt.findKNearestContacts(contact.Id)
	*reply = FindNodeReply{Contacts: kNearest}
	return nil
}

func (node *Node) String() string {
	return fmt.Sprintf("Node: (id = %s) (address = %s) (kBuckets = %v)",
		node.id.String(),
		node.addr.String(),
		node.rt)
}

// Return XOR distance between node and other
func (node *Node) distanceTo(other *Contact) *big.Int {
	return big.NewInt(0).Xor(&node.id, &other.Id)
}

func distanceBetween(firstId big.Int, secondId big.Int) *big.Int {
	return big.NewInt(0).Xor(&firstId, &secondId)
}

// NewNode returns a new Node struct
func NewNode(address string) *Node {
	node := new(Node)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	node.addr = *addr

	hash := sha1.Sum([]byte(addr.String()))

	node.id = *big.NewInt(0)
	node.id.SetBytes(hash[:])
	// TODO: take in k and tRefresh arguments - for now just hardcoding default
	node.rt = NewRoutingTable(node)
	node.restC = make(chan CommandMessage)

	node.logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	return node
}

// Run is called on an initialized Node to begin serving the RPC endpoints
func (node *Node) Run(toPing string) {
	nodeRPC := &NodeRPC{node}
	rpc.Register(nodeRPC)
	rpc.HandleHTTP()
	node.setupControlEndpoints()

	// Removed this from if statement. May need to put it back
	toPingAddr, err := net.ResolveTCPAddr("", toPing)
	//TODO: handle err
	if err != nil {
		node.logger.Printf("%s", err)
	}

	ticker := time.NewTicker(1 * time.Second)
	counter := 0

	go func() {
		for {
			select {
			case msg := <-node.restC:
				switch msg.Command {
				case "PING":
					contact, ok := msg.Arg1.(Contact)
					if !ok {
						fmt.Printf("PING REST argument is not a Contact")
					}

					fmt.Printf("Performing PING for IP: %s (ID: %s)\n",
						contact.Addr.String(),
						contact.Id.String())
				case "STORE":
					key := msg.Arg1
					value := msg.Arg2
					fmt.Printf("Performing STORE of key: %s, value: %s\n", key, value)
				case "FINDNODE":
					id := msg.Arg1
					fmt.Printf("Performing FINDNODE of server id: %s\n", id)
				case "FINDVALUE":
					key := msg.Arg1
					fmt.Printf("Performing FINDVALUE of key: %s\n", key)
				case "SHUTDOWN":
					fmt.Println("Shutting down...")
					os.Exit(0)
				}
			case <-ticker.C:
				if toPing != "" {
					node.doPing(*toPingAddr)
					counter++
					if counter == 5 {
						os.Exit(1)
					}
				}
			}
		}
	}()

	// open our own port for connection
	l, e := net.ListenTCP("tcp", &node.addr)
	if e != nil {
		log.Fatal(e)
		return
	}

	http.Serve(l, nil)
}

// Perform the legwork of RPC invocation
func (node *Node) doRPC(method string, dest net.TCPAddr, args interface{}, reply interface{}) bool {
	node.logger.Printf("Sending %s RPC to %s", method, dest.String())

	client, err := rpc.DialHTTP("tcp", dest.String())
	if err != nil {
		node.logger.Printf("Dial to %s failed: %s", dest.String(), err)
		return false
	}

	err = client.Call(fmt.Sprintf("NodeRPC.%s", method), args, reply)
	if err != nil {
		node.logger.Printf("%s RPC to %s failed: %s", method, dest.String(), err)
		return false
	}

	return true
}

// Send a PING RPC to dest
func (node *Node) doPing(dest net.TCPAddr) {
	args := PingArgs{node.addr}
	var reply PingReply

	if !node.doRPC("Ping", dest, args, &reply) {
		return
	}

	node.logger.Printf("Got ping reply from %s", reply.Source.String())

	// TODO: Update K-Buckets
}

// Send a STORE RPC for (key, value) to dest
func (node *Node) doStore(key string, value []byte, dest net.TCPAddr) {
	args := StoreArgs{node.addr, key, value}
	var reply StoreReply

	if !node.doRPC("Store", dest, args, &reply) {
		return
	}
}

// Send a FINDVALUE RPC for key to dest
func (node *Node) doFindValue(key string, dest net.TCPAddr) {
	args := FindValueArgs{node.addr, key}
	var reply FindNodeReply

	if !node.doRPC("FindValue", dest, args, &reply) {
		return
	}

	// TODO: Whatever processing we need to perform afterwards
	// TODO: Update K-Buckets
}

// Send a FINDNODE RPC for key to dest
func (node *Node) doFindNode(dest net.TCPAddr) {
	args := FindNodeArgs{node.addr}
	var reply FindNodeReply
	if !node.doRPC("FindNode", dest, args, &reply) {
		return
	}
	// TODO: Whatever processing we need to perform afterwards
	// TODO: Update K-Buckets
}
