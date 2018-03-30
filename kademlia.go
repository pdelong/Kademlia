package kademlia

import (
	"container/list"
	"crypto/sha1"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

// An entry in the k-bucket
type NodeEntry struct {
	Id   big.Int
	Addr net.TCPAddr
}

type Node struct {
	id       big.Int
	addr     net.TCPAddr
	ht       map[string][]byte
	kBuckets []list.List // TODO: Might want array of arrays
	logger   *log.Logger
	// TODO: routing table
}

type PingArgs struct {
	Source NodeEntry
}

type PingReply struct{}

type StoreArgs struct {
	Source big.Int
	Key    string
	Val    []byte
}

type StoreReply struct {
}

type FindValueArgs struct {
	Source big.Int
	Key    string
}

type FindValueReply struct {
	Val  []byte
	Node NodeEntry
}

type FindNodeArgs struct {
	Source big.Int
	Key    string
}

type FindNodeReply struct {
}

func (self *Node) Ping(args PingArgs, reply *PingReply) error {

	self.logger.Printf("Ping from %s", args.Source.Addr)
	// TODO: Update k-bucket based on args.Source
	return nil
}

func (self *Node) Store(args StoreArgs, reply *StoreReply) error {
	self.ht[args.Key] = args.Val
	// TODO: update kbuckets/ reply?
	return nil
}

func (self *Node) FindValue(args FindValueArgs, reply *FindValueReply) error {
	return nil
}

func (self *Node) FindNode(args FindNodeArgs, reply *FindNodeReply) error {
	return nil
}

func (self *Node) String() string {
	return fmt.Sprintf("Node: (id = %v) (address = %v) (kBuckets = %v)",
		self.id,
		self.addr,
		self.kBuckets)
}

// Return XOR distance between self and other
func (self *Node) distanceTo(other *NodeEntry) *big.Int {
	return big.NewInt(0).Xor(&self.id, &other.Id)
}

func NewNode(address string) *Node {
	node := new(Node)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil
	}

	node.addr = *addr

	// #eww
	hash := sha1.Sum([]byte(fmt.Sprintf("%s", address)))

	node.id = *big.NewInt(0)
	node.id.SetBytes(hash[:])

	node.logger = log.New(os.Stdout, "INFO:  ", log.Ldate|log.Ltime|log.Lshortfile)
	return node
}

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

// DoPing sends RPC, Ping is the one called
func (self *Node) DoPing(dest net.TCPAddr) {
	client, err := rpc.DialHTTP("tcp", dest.String())
	// TODO : handle err
	if err != nil {
		self.logger.Printf("Ping RPC to %s failed %s", dest, err)
		return
	}

	nodeEntry := NodeEntry{Id: self.id, Addr: self.addr}
	args := PingArgs{nodeEntry}
	var reply PingReply
	err = client.Call("NodeRPC.Ping", args, &reply)
	self.logger.Printf("DoPing called on %s", dest.String())
	if err != nil {
		self.logger.Printf("Ping RPC to %s failed %s", dest, err)
	}
}

// do not touch
func (self *NodeRPC) Ping(args PingArgs, reply *PingReply) error {
	self.node.Ping(args, reply)
	return nil
}

func (self *NodeRPC) Store(args StoreArgs, reply *StoreReply) error {
	self.node.Store(args, reply)
	return nil
}

func (self *NodeRPC) FindValue(args FindValueArgs, reply *FindValueReply) error {
	self.node.FindValue(args, reply)
	return nil
}

func (self *NodeRPC) FindNode(args FindNodeArgs, reply *FindNodeReply) error {
	self.node.FindNode(args, reply)
	return nil
}

type NodeRPC struct {
	node *Node
}
