package main

import (
	"bufio"
	"fmt"
	"github.com/peterdelong/kademlia"
	"io"
	"log"
	"math/rand"
	"os"
	"io/ioutil"
	"time"
)

// path to bootstrap_nodes
const bootstrap_node_path = "/home/pdelong/go/src/github.com/peterdelong/kademlia/cmd/kademlia_node/bootstrap_nodes"

func checkIOError(e error) {
	if e != nil && e != io.EOF {
		log.Fatal(e)
	}
}

// usage: kademlia_node <node_addr> [bootstrap_addr]
func main() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	num := r1.Intn(1000)
	time.Sleep(time.Duration(num) * time.Millisecond)

	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("usage: kademlia_node <node_addr> <b/nb> [bootstrap_addr]")
		return
	}

	addr := args[0]

	// if this isn't a bootstrap node, read from a list of nodes in the system
	bootstrapAddr := ""
	if args[1] == "nb" {
		if len(args) < 3 {
			file, err := os.Open(bootstrap_node_path)
			checkIOError(err)
			reader := bufio.NewReader(file)
			bootstrapAddrBytes, _, err := reader.ReadLine()
			checkIOError(err)
			bootstrapAddr = string(bootstrapAddrBytes)
		} else {
			bootstrapAddr = args[2]
		}

		fmt.Println("Contacting ", bootstrapAddr)
	}

	// Disable logging if necessary (see option in globals.go)
	if !loggingEnable {
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)
	}

	node := kademlia.NewNode(addr)

	fmt.Println(node)

	node.Run(bootstrapAddr)
}
