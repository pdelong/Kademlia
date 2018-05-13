#!/bin/bash
if [ "$#" -ne 1 ] ; then
  echo "Usage: $0 [current node's IP address]"
    exit 1
fi

ROOT=/home/pdelong/go/src/github.com/peterdelong/kademlia
PORT=8001

#remove all logs
cd $ROOT/logs
rm *

# rebuild the project and save the bootstrap node's address
cd $ROOT/cmd/kademlia_node
rm bootstrap_nodes
go build

# start up the bootstrap node and save the id for later
./kademlia_node $1:$PORT b > $ROOT/logs/bootstrap.log &
jobid=$!

# start up the other nodes and wait until they finish, so we can kill
# the bootstrap node
cd $ROOT/slurm_stuff
sbatch --wait start_nodes.cmd

# kill the bootstrap node
kill $jobid
