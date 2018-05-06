"""Kademlia Test Script

Usage:
    test.py <addr> store <key>
    test.py <addr> ping (id | ip) <target> 
    test.py <addr> findnode (iterative | oneshot) <target>
    test.py <addr> findvalue (iterative | oneshot) <target>
    test.py <addr> shutdown
    test.py <addr> table
    test.py test [--zipf <alpha> | --uniform | --exponential <k> | --linear <m>] [-n <times>] [--size <size>] <nodes>...
    test.py (-h | --help)
    test.py --version

    Options:
        --zipf <alpha>              Use zipfian distribution
        --uniform                   Use uniform distribution (--linear 1)
        --linear <m>                Use linear distribution
        --size <size>               How large to make the values (in bytes) [default: 128]
        -n <times>                  Number of times to test [default: 100]
        -h --help                   Show this screen
        --version                   Show version
        -v                          Show extra debug information
"""

from kademlia import KademliaNode
from distributions import Zipf
from docopt import docopt
import sys

def get_distribution(arguments):
    pass


if __name__ == '__main__':
    arguments = docopt(__doc__, version='Kademlia Test Script 0.1')
    print(arguments)

    if not arguments['test']:
        node = KademliaNode(arguments['<addr>'])
        print("Created Kademlia node with address: {}".format(arguments['<addr>']))

    if arguments['store']:
        print("Attempting store")

        key = arguments['<key>']
        value = sys.stdin.buffer.read()

        node.store(key, value)
    elif arguments['ping']:
        pass
    elif arguments['findnode']:
        pass
    elif arguments['findvalue']:
        pass
    elif arguments['shutdown']:
        print("Attempting to shut down node")
        node.shutdown()
    elif arguments['table']:
        print("Retrieving list of all keys on node")
        node.table()
    elif arguments['test']:
        nodes = [KademliaNode(addr) for addr in arguments['<nodes>']]
        distribution = get_distribution(arguments)
        for i in range(int(arguments['-n'])):
            print(i)
            # Query nodes repeatedly with the given key distribution
        pass
