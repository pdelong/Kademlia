package kademlia

// The following definitions are to present the proper RPC interface.
// They immediately delegate functionality to the corresponding functions on the
// Node struct. THESE SHOULD NOT BE EDITED!!!
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
