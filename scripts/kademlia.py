#!/usr/bin/env python3

import requests
import sys
import time
import base64
import json


# with Stopwatch():
#     pass
class Stopwatch:
    def __enter__(self):
        self.start = time.clock()
        return self

    def __exit__(self, *args):
        interval = time.clock() - self.start
        print("timespan: {}".format(interval))


class KademliaNode:
    def __init__(self, address):
        self.address = address

    def store(self, key, value):
        requests.post("http://{}/store/{}".format(self.address, key), data=value)

    def ping(self, target, byid):
        if byid:
            method = "id"
        else:
            method = "ip"

        url = "http://{}/ping/{}/{}".format(self.address, method, target)

        requests.get(url)

    def shutdown(self):
            try:
                requests.get("http://{}/shutdown".format(self.address))
            except:
                pass

    def findnode(self, target, oneshot):
        if oneshot:
            method = "oneshot"
        else:
            method = "iterative"

        url = "http://{}/{}/findnode/{}".format(self.address, method, target)

        requests.get(url)

    def findvalue(self, key, oneshot):
        if oneshot:
            method = "oneshot"
        else:
            method = "iterative"

        url = "http://{}/{}/findvalue/{}".format(self.address, method, key)

        requests.get(url)

    def table(self):
        r = requests.get("http://{}/table".format(self.address))
        table = json.loads(r.text)

        for entry in table:
            key = entry['key']
            value = base64.b64decode(entry['value'])
            isOrigin = entry['isOrigin']

            print("key: {}, value: {}, isOrigin: {}".format(key, value, isOrigin))


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print('usage: kademlia.py SERVER [action] [args...]')
        sys.exit(1)

    node = KademliaNode(sys.argv[1])
    print("Created kademlia node with address: {}".format(sys.argv[1]))

    if len(sys.argv) < 3:
        sys.exit(0)

    command = sys.argv[2]

    if command == 'store':
        if len(sys.argv) < 4:
            print("key and value required for store")
            sys.exit(1)

        print("Attempting store")

        key = sys.argv[3]
        value = sys.stdin.buffer.read()

        node.store(key, value)
    elif command == 'ping':
        pass
    elif command == 'findnode':
        pass
    elif command == 'findvalue':
        pass
    elif command == 'shutdown':
        print("Attempting to shut down node")
        node.shutdown()
    elif command == 'table':
        node.table()
    else:
        print("Unknown command: {}".format(command))
        sys.exit(2)
