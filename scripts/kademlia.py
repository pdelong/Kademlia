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

    def store_here(self, key, value):
        requests.post("http://{}/store_here/{}".format(self.address, key), data=value)

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

        r = requests.get(url)
        contacts = json.loads(r.text)
        for entry in contacts:
            key = entry['Id']
            addr = entry['Addr']
            print("%x %s"%(key, addr))

    def findvalue(self, key, oneshot=True):
        if oneshot:
            method = "oneshot"
        else:
            method = "iterative"

        url = "http://{}/{}/findvalue/{}".format(self.address, method, key)

        requests.get(url)

    def table(self):
        r = requests.get("http://{}/table".format(self.address))
        table = json.loads(r.text)

        new_table = {}
        for entry in table:
            key = entry['key']
            value = base64.b64decode(entry['value'])
            isOrigin = entry['isOrigin']

            new_table[key] = {'value': value, 'isOrigin': isOrigin}

        return new_table
