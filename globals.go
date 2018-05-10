package kademlia

import "time"

// tExpire is the time after which a key/value pair expires
// TTL from original publication date
const tExpire = 864000 * time.Second

// tRefresh is the time after which an unaccessed bucket must be refreshed
const tRefresh = 3600 * time.Second

// tReplicate is the interval between replication events, when a node is
// required to publish its entire database
const tReplicate = 3600 * time.Second

// tRepublish is the time after which the original publisher must republish a
// key/value pair
const tRepublish = 86400 * time.Second

// Alpha is the degree of parallelism in network calls
const alpha = 3

// k is the maximum number of contacts stored in a bucket
const k = 4

// keys should be stored as hex when in string form
const keyBase = 16
