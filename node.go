package kademlia

import (
	"container/list"
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

// Data structure for an individual Kademlia node
type Node struct {
	id       big.Int
	addr     net.TCPAddr
	ht       map[string][]byte
	kBuckets []list.List // TODO: Might want array of arrays
	// kBuckets should just be a slice(array?) of kbucket structs
	logger *log.Logger
	// TODO: routing table
}

type PingArgs struct {
	Source net.TCPAddr
}

type PingReply struct{}

type StoreArgs struct {
	Source net.TCPAddr
	Key    string
	Val    []byte
}

type StoreReply struct {
}

type FindValueArgs struct {
	Source net.TCPAddr
	Key    string
}

type FindValueReply struct {
	Val  []byte
	Node Contact
}

type FindNodeArgs struct {
	Source net.TCPAddr
	Key    string
}

type FindNodeReply struct {
}

// Handler for the PING RPC
func (self *Node) Ping(args PingArgs, reply *PingReply) error {
	nodeEntry := NewContact(args.Source)
	if nodeEntry == nil {
		return errors.New("Couldn't hash IP address")
	}

	self.logger.Printf("Ping from %s", args.Source.String())
	// TODO: Update k-bucket based on args.Source
	return nil
}

// Handler for the STORE RPC
func (self *Node) Store(args StoreArgs, reply *StoreReply) error {
	nodeEntry := NewContact(args.Source)
	if nodeEntry == nil {
		return errors.New("Couldn't hash IP address")
	}

	self.ht[args.Key] = args.Val
	// TODO: update kbuckets/ reply?
	return nil
}

// Handler for the FINDVALUE RPC
func (self *Node) FindValue(args FindValueArgs, reply *FindValueReply) error {
	nodeEntry := NewContact(args.Source)
	if nodeEntry == nil {
		return errors.New("Couldn't hash IP address")
	}

	return nil
}

// Handler for the FINDNODE RPC
func (self *Node) FindNode(args FindNodeArgs, reply *FindNodeReply) error {
	nodeEntry := NewContact(args.Source)
	if nodeEntry == nil {
		return errors.New("Couldn't hash IP address")
	}

	return nil
}

func (self *Node) String() string {
	return fmt.Sprintf("Node: (id = %s) (address = %s) (kBuckets = %v)",
		self.id.String(),
		self.addr.String(),
		self.kBuckets)
}

// Return XOR distance between self and other
func (self *Node) distanceTo(other *Contact) *big.Int {
	return big.NewInt(0).Xor(&self.id, &other.Id)
}

// Construct a new Node struct
func NewNode(address string) *Node {
	node := new(Node)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil
	}

	node.addr = *addr

	hash := sha1.Sum([]byte(addr.String()))

	node.id = *big.NewInt(0)
	node.id.SetBytes(hash[:])

	node.logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return node
}

// Called on an initialized Node to begin serving the RPC endpoints
func (self *Node) Run(toPing string) {
	nodeRPC := &NodeRPC{self}
	rpc.Register(nodeRPC)
	rpc.HandleHTTP()

	toPingAddr, err := net.ResolveTCPAddr("", toPing)
	//TODO: handle err
	if err != nil {
		log.Fatal(err)
	}

	l, e := net.ListenTCP("tcp", &self.addr)
	if e != nil {
		return
	}

	// periodically ping
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			self.DoPing(*toPingAddr)
		}
	}()
	http.Serve(l, nil)
}

// Perform the legwork of RPC invocation
func (self *Node) doRPC(method string, dest net.TCPAddr, args interface{}, reply interface{}) bool {
	self.logger.Printf("Sending %s RPC to %s", method, dest.String())

	client, err := rpc.DialHTTP("tcp", dest.String())
	if err != nil {
		self.logger.Printf("Dial to %s failed: %s", dest.String(), err)
		return false
	}

	err = client.Call(fmt.Sprintf("NodeRPC.%s", method), args, reply)
	if err != nil {
		self.logger.Printf("%s RPC to %s failed: %s", method, dest, err)
		return false
	}

	return true
}

// Send a PING RPC to dest
func (self *Node) DoPing(dest net.TCPAddr) {
	args := PingArgs{self.addr}
	var reply PingReply

	if !self.doRPC("Ping", dest, args, &reply) {
		return
	}

	// TODO: Update K-Buckets
}

// Send a STORE RPC for (key, value) to dest
func (self *Node) DoStore(key string, value []byte, dest net.TCPAddr) {
	args := StoreArgs{self.addr, key, value}
	var reply StoreReply

	if !self.doRPC("Store", dest, args, &reply) {
		return
	}

	// TODO: Whatever processing we need to perform afterwards
	// TODO: Update K-Buckets
}

// Send a FINDVALUE RPC for key to dest
func (self *Node) DoFindValue(key string, dest net.TCPAddr) {
	args := FindValueArgs{self.addr, key}
	var reply FindNodeReply

	if !self.doRPC("FindValue", dest, args, &reply) {
		return
	}

	// TODO: Whatever processing we need to perform afterwards
	// TODO: Update K-Buckets
}

// Send a FINDNODE RPC for key to dest
func (self *Node) DoFindNode(key string, dest net.TCPAddr) {
	args := FindNodeArgs{self.addr, key}
	var reply FindNodeReply

	if !self.doRPC("FindNode", dest, args, &reply) {
		return
	}

	// TODO: Whatever processing we need to perform afterwards
	// TODO: Update K-Buckets
}
