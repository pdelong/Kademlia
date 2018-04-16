#!/bin/bash

ips=($(hostname -I))
cd /home/pdelong/go/src/github.com/peterdelong/kademlia/cmd/kademlia_node

./kademlia_node ${ips[0]}:8000 nb > /home/pdelong/go/src/github.com/peterdelong/kademlia/logs/${ips[0]} 
