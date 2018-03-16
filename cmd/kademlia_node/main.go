package main

import (
	"fmt"
	"github.com/peterdelong/kademlia"
	"os"
)

// usage: kademlia_node <node_addr> [bootstrap_addr]
func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("usage: kademlia_node <node_addr> [bootstrap_addr]")
		return
	}

	addr := args[0]
	//bootstrapAddr := args[1]

	node := kademlia.NewNode(addr)
	fmt.Println(node)

	node.Run()
}
