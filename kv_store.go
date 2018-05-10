package kademlia

import ()

// KVStore holds mappings from keys to values and keeps track if a given node is
// the owner of the value
type KVStore struct {
	//owner    *Node
	ht       map[string][]byte
	isOrigin map[string]bool
}

// NewKVStore returns a newly initialized KVStore
func NewKVStore() *KVStore {
	kvStore := new(KVStore)
	kvStore.ht = make(map[string][]byte)
	kvStore.isOrigin = make(map[string]bool)

	//kvStore.owner = owner
	return kvStore
}

func (store *KVStore) get(key string) ([]byte, bool) {
	if val, ok := store.ht[key]; ok {
		return val, true
	}
	return nil, false
}

// Will overwrite existing value
func (store *KVStore) add(key string, val []byte, isOrigin bool) {
	store.ht[key] = val
	store.isOrigin[key] = isOrigin
}

// KV contains all the information we have for a key
type KV struct {
	key      string
	val      []byte
	isOrigin bool
}

// Iterator returns a channel that iterates over all the keys that we've stored
func (store *KVStore) Iterator() chan *KV {
	ch := make(chan *KV)
	go func() {
		for k, v := range store.ht {
			kv := new(KV)
			kv.key = k
			kv.val = v
			kv.isOrigin = store.isOrigin[k]
			ch <- kv
		}
		close(ch)
	}()
	return ch
}
