#!/bin/bash 

#SBATCH --nodes=4  # node count 
#SBATCH --ntasks-per-node=1 
#SBATCH --time=5:00


srun --exclusive ./run_kademlia.sh

