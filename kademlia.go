package kademlia

import (
	"container/list"
	"crypto/sha1"
	//	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/rpc"
)

// An entry in the k-bucket
type NodeEntry struct {
	id   *big.Int
	addr *net.IPAddr
}

type Node struct {
	id       *big.Int
	addr     *net.TCPAddr
	kBuckets []list.List // TODO: Might want array of arrays
}

type PingArgs struct {
}

type PingReply struct {
}

type StoreArgs struct {
}

type StoreReply struct {
}

type FindValueArgs struct {
}

type FindValueReply struct {
}

type FindNodeArgs struct {
}

type FindNodeReply struct {
}

func (self *Node) Ping(args PingArgs, reply *PingReply) error {
	return nil
}

func (self *Node) Store(args StoreArgs, reply *StoreReply) error {
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
func (self *Node) distanceTo(other NodeEntry) *big.Int {
	return big.NewInt(0).Xor(self.id, other.id)
}

func NewNode(address string) *Node {
	node := new(Node)
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil
	}

	node.addr = addr

	// #eww
	hash := sha1.Sum([]byte(fmt.Sprintf("%s", address)))

	node.id = big.NewInt(0)
	node.id.SetBytes(hash[:])

	return node
}

func (self *Node) Run() {
	nodeRPC := &NodeRPC{self}
	rpc.Register(nodeRPC)
	rpc.HandleHTTP()

	l, e := net.ListenTCP("tcp", self.addr)
	if e != nil {
		return
	}
	http.Serve(l, nil)
}

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
