package kademlia

import ()

type KVStore struct {
	//owner    *Node
	ht       map[string][]byte
	isOrigin map[string]bool
}

func NewKVStore() *KVStore {
	kvStore := new(KVStore)
	kvStore.ht = make(map[string][]byte)
	kvStore.isOrigin = make(map[string]bool)

	//kvStore.owner = owner
	return kvStore
}

func (self *KVStore) get(key string) ([]byte, bool) {
	if val, ok := self.ht[key]; ok {
		return val, true
	}
	return nil, false
}

// Will overwrite existing value
func (self *KVStore) add(key string, val []byte, isOrigin bool) {
	self.ht[key] = val
	self.isOrigin[key] = isOrigin
}

type KV struct {
	key      string
	val      []byte
	isOrigin bool
}

//func (self *KVStore) getKeys
func (self *KVStore) Iterator() chan *KV {
	ch := make(chan *KV)
	go func() {
		for k, v := range self.ht {
			kv := new(KV)
			kv.key = k
			kv.val = v
			kv.isOrigin = self.isOrigin[k]
			ch <- kv
		}
		close(ch)
	}()
	return ch
}
