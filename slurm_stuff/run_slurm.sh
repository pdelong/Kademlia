#!/bin/bash
if [ "$#" -ne 1 ] ; then
  echo "Usage: $0 [current node's IP address]"
    exit 1
fi

# rebuild the project and save the bootstrap node's address
cd /home/pdelong/go/src/github.com/peterdelong/kademlia/cmd/kademlia_node
echo $1:8000 > bootstrap_nodes
go build

# start up the bootstrap node and save the id for later
./kademlia_node $1:8000 b > /home/pdelong/go/src/github.com/peterdelong/kademlia/logs/bootstrap_log &
jobid=$!

# start up the other nodes and wait until they finish, so we can kill
# the bootstrap node
cd /home/pdelong/go/src/github.com/peterdelong/kademlia/slurm_stuff
sbatch --wait start_nodes.cmd

# kill the bootstrap node
kill $jobid
