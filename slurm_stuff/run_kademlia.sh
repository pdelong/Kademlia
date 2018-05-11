#!/bin/bash

ips=($(hostname -I))
PORT=8001
ROOT=/home/pdelong/go/src/github.com/peterdelong/kademlia

$ROOT/cmd/kademlia_node/kademlia_node ${ips[0]}:$PORT nb > $ROOT/logs/${ips[0]}.log 

#sleep 10
#
#$ROOT/scripts/test.py ${ips[0]}:$PORT findnode iterative 4eabc738192049c89d9e03319384 > $ROOT/logs/${ips[0]}.test.log
#
#
## Make sure it's dead
#kill $! 2> /dev/null
#if [[ $? == '0' ]]; then
#    echo "Node had to be killed (shutdown probably failed)"
#    exit 1
#fi
