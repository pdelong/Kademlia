package kademlia

// The following definitions are to present the proper RPC interface.
// They immediately delegate functionality to the corresponding functions on the
// Node struct. THESE SHOULD NOT BE EDITED!!!

// Ping is a stub function that exposes the PING RPC
func (fakeNode *NodeRPC) Ping(args PingArgs, reply *PingReply) error {
	fakeNode.node.Ping(args, reply)
	return nil
}

// Store is a stub function that exposes the STORE RPC
func (fakeNode *NodeRPC) Store(args StoreArgs, reply *StoreReply) error {
	fakeNode.node.Store(args, reply)
	return nil
}

// FindValue is a stub function that exposes the FINDVALUE RPC
func (fakeNode *NodeRPC) FindValue(args FindValueArgs, reply *FindValueReply) error {
	fakeNode.node.FindValue(args, reply)
	return nil
}

// FindNode is a stub function that exposes the FINDNODE RPC
func (fakeNode *NodeRPC) FindNode(args FindNodeArgs, reply *FindNodeReply) error {
	fakeNode.node.FindNode(args, reply)
	return nil
}

// NodeRPC is a wrapper struct that is used to control which RPCs are exposed
type NodeRPC struct {
	node *Node
}
