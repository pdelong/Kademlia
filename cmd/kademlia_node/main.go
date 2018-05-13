package main

import (
	"bufio"
	"fmt"
	"github.com/peterdelong/kademlia"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

// path to bootstrap_nodes

func checkIOError(e error) {
	if e != nil && e != io.EOF {
		log.Fatal(e)
	}
}

// usage: kademlia_node <node_addr> [bootstrap_addr]
func main() {
	fmt.Println("Started")
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("usage: kademlia_node <node_addr> <b/nb> [bootstrap_addr]")
		return
	}

	addr := args[0]
	nodes := make([]string, 0, 100)

	// if this isn't a bootstrap node, read from a list of nodes in the system
	bootstrapAddr := ""
	if args[1] == "nb" {
		num := r1.Intn(7000)
		time.Sleep(time.Duration(num) * time.Millisecond)
		if len(args) < 3 {
			file, err := os.Open(kademlia.Bootstrap_node_path)
			checkIOError(err)
			reader := bufio.NewReader(file)
			for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
				// need to cut off the delimiter
				address := line[0:len(line)-1]
				nodes = append(nodes, address)
			}
		} else {
			bootstrapAddr = args[2]
		}
		index := r1.Intn(len(nodes))
		bootstrapAddr = nodes[index]

		fmt.Println("Contacting ", bootstrapAddr)
	}

	node := kademlia.NewNode(addr)

	fmt.Println(node)

	node.Run(bootstrapAddr)
}
