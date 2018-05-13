#!/bin/bash 

#SBATCH --nodes=75  # node count 
#SBATCH --ntasks-per-node=1 
#SBATCH --time=20:00


srun --exclusive ./run_kademlia.sh

