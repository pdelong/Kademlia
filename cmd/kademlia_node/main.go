package main

import (
	"bufio"
	"fmt"
	"github.com/peterdelong/kademlia"
	"io"
	"log"
	"os"
)

func checkIOError(e error) {
	if e != nil && e != io.EOF {
		log.Fatal(e)
	}
}

// usage: kademlia_node <node_addr> [bootstrap_addr]
func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("usage: kademlia_node <node_addr> <b/nb> [bootstrap_addr]")
		return
	}

	addr := args[0]

	//if len(args) >= 2 {
	//	bootstrapAddr = args[1]
	//}
	// if this isn't a bootstrap node, read from a list of nodes in the system
	bootstrapAddr := ""
	if args[1] == "nb" {
		file, err := os.Open("bootstrap_nodes")
		checkIOError(err)
		reader := bufio.NewReader(file)
		bootstrapAddrBytes, _, err := reader.ReadLine()
		checkIOError(err)
		bootstrapAddr = string(bootstrapAddrBytes)
		fmt.Println("Contacting ", bootstrapAddr)
	}

	node := kademlia.NewNode(addr)

	fmt.Println(node)

	node.Run(bootstrapAddr)
}
