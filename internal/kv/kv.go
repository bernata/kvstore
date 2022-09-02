package kv

import (
	"fmt"
	"sync"
)

type Store struct {
	// Simple map backing store with a read-write mutex to protect against
	// current writes and reads.
	// Consider using other strategies when the access pattern for the cache is
	// better understood.
	// sync.Map is optimized for read only caches [i.e. only grow] and writes on disjoint keys.
	// This store allows deletes and write pattern is unclear.
	// Other concurrent map models exist where a copy of the map is maintained for writes and
	// swapped with the read only version periodically.
	data map[string]string

	mutex sync.RWMutex
}

func NewStore(capacity int) *Store {
	if capacity <= 0 {
		capacity = 10
	}
	return &Store{
		data: make(map[string]string, capacity),
	}
}

func (kv *Store) Get(key string) (string, bool) {
	kv.mutex.RLock()
	defer kv.mutex.RUnlock()
	result, ok := kv.data[key]
	return result, ok
}

func (kv *Store) Delete(key string) {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	delete(kv.data, key)
}

func (kv *Store) Write(key, value string) error {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()

	const maxKeyBytes = 250
	const maxValueBytes = 1024 * 1024

	if len(key) > maxKeyBytes {
		return errMaxKeyBytes
	} else if len(value) > maxValueBytes {
		return errMaxValueBytes
	}

	kv.data[key] = value
	return nil
}

var errMaxKeyBytes = fmt.Errorf("max key size allowed is 255 bytes")
var errMaxValueBytes = fmt.Errorf("max value size allowed is 1 MB")
