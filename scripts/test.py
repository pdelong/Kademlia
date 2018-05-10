#!/usr/bin/env python
"""Kademlia Test Script

Usage:
    test.py <addr> store <key>
    test.py <addr> ping (id | ip) <target> 
    test.py <addr> findnode (iterative | oneshot) <target>
    test.py <addr> findvalue (iterative | oneshot) <target>
    test.py <addr> shutdown
    test.py <addr> table
    test.py test [--zipf <alpha> | --uniform | --linear <m>] [--times <times>] [--keys <keys>] [--size <size>] <nodes>...
    test.py (-h | --help)
    test.py --version

    Options:
        --zipf <alpha>              Use zipfian distribution
        --uniform                   Use uniform distribution (--linear 1)
        --linear <m>                Use linear distribution
        --size <size>               How large to make the values (in bytes) [default: 128]
        -t --times <times>          Number of times to test [default: 100]
        -n --keys <keys>            Number of keys to insert [default: 1000]
        -h --help                   Show this screen
        --version                   Show version
        -v                          Show extra debug information
"""

from kademlia import KademliaNode
from distributions import Zipf, Uniform, Linear
from docopt import docopt
import sys
import os
import random
import binascii
import hashlib

def get_distribution(arguments, keys):
    if arguments['--zipf'] is not None:
        distribution = Zipf(keys, int(arguments['--zipf']))
    elif arguments['--uniform']:
        distribution = Uniform(keys)
    elif arguments['--linear'] is not None:
        distribution = Linear(keys, int(arguments['--linear']))
    else:
        distribution = Zipf(keys, 1)

    return distribution

def key_to_id(key):
    return binascii.hexlify(hashlib.sha1(str(key).encode('utf-8')).digest())


if __name__ == '__main__':
    arguments = docopt(__doc__, version='Kademlia Test Script 0.1')
    print(arguments)

    if not arguments['test']:
        node = KademliaNode(arguments['<addr>'])
        print("Created Kademlia node with address: {}".format(arguments['<addr>']))

    if arguments['store']:
        print("Attempting store")

        key = key_to_id(arguments['<key>'])
        value = sys.stdin.read()

        print("Storing key as {}".format(key))

        node.store(key, value)
    elif arguments['ping']:
        pass
    elif arguments['findnode']:
        key = arguments['<target>']
        print("Going to find node with id: {}".format(key))
        node.findnode(key, arguments['oneshot'])
    elif arguments['findvalue']:
        key = arguments['<target>']
        print("Looking for value with id: {}".format(key))
        node.findvalue(key, arguments['oneshot'])
    elif arguments['shutdown']:
        print("Attempting to shut down node")
        node.shutdown()
    elif arguments['table']:
        print("Retrieving list of all keys on node")
        for key, value in node.table().items():
            print("Key: {}, Value: {}, isOrigin: {}".format(key, value['value'], value['isOrigin']))
    elif arguments['test']:
        nodes = [KademliaNode(addr) for addr in arguments['<nodes>']]
        num_keys = int(arguments['--keys'])
        keys = [key_to_id(key) for key in range(num_keys)]
        distribution = get_distribution(arguments, keys)
        size = int(arguments['--size'])

        values = {}
        for i in keys:
            values[i] = os.urandom(size)

        storage_node = nodes[0]
        print('Using {} as the main node for storage setup'.format(storage_node.address))

        # Store values on the appropriate nodes
        for key, value in values.items():
            print("Storing {}".format(key))
            storage_node.store(key, value)

        table = storage_node.table()
        print(table)
        for key, value in storage_node.table().items():
            print("{}: {}".format(key, value))

        for i in range(int(arguments['--times'])):
            key = distribution.next()
            node = random.choice(nodes)
            value = node.findvalue(key)

            if value != values[key]:
                print("Value returned for {} but did not equal expected value".format(key))
                print("{} vs {}".format(value, values[key]))
            else:
                print("Correct value returned!")

