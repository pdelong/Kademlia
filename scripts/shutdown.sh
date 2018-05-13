#!/bin/sh

echo $(for f in $(ls logs); do if [[ ${f%.log} != 'bootstrap' ]]; then echo "${f%.log}:8001 "; fi; done) | tr ' ' "\n" | xargs -n1 -P100 -I{} ./scripts/test.py {} shutdown
./scripts/test.py 172.16.133.246:8001 shutdown
