package main

import (
	"fmt"
	"github.com/peterdelong/kademlia"
)

func main() {
	node := kademlia.NewNode("140.180.129.24", 7878)
	fmt.Println(node)
}
