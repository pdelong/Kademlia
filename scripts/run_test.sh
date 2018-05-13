#!/bin/sh

~/go/src/github.com/peterdelong/kademlia/scripts/test.py test --zipf 1.0625 --keys 1024 --retrieveonly --size 16 --times 61440 172.16.133.246:8001 $(for f in $(ls logs); do if [[ ${f%.log} != 'bootstrap' ]]; then echo -n "${f%.log}:8001 "; fi; done) > data/$1
