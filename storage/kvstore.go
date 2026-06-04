package storage

import (
	"sync"
)

// KVStore represents thread safe in-memory key-value store
type KVStore struct {
	mu   sync.RWMutex      // The lock that protects the map with best performance
	data map[string]string // The actual Data
	wal  *WAL              // The Write-Ahead Log for durability and crash recovery
}

// NewKVStore is a constructor function to initialize our KVStore
// If we dont initialize the map with "make()", it cause a panic.
func NewKVStore(wal *WAL) *KVStore {
	return &KVStore{
		data: make(map[string]string),
		wal:  wal,
	}
}

// Put safely adds or updates a key-value pair in the store
func (store *KVStore) Put(key, value string) error {

	store.mu.Lock() // Grabbing the exclusive write lock so no one can read or write to the store while we are updating it.

	defer store.mu.Unlock() // Making sure it unlocks it.

	// If the disk write fails, we return an error and DO NOT update RAM.
	// Saving in disk first.
	if err := store.wal.Append(key, value); err != nil {
		return err
	}

	// IF the disk write is successful, we can update the in-memory map. This way we ensure durability.
	store.data[key] = value
	return nil
}

// Getting the data associated with the key.
func (store *KVStore) Get(key string) (string, bool) {
	store.mu.RLock() // Grabbing the shared read lock so multiple readers can access the store concurrently, but no writers can modify it.
	defer store.mu.RUnlock()

	value, exists := store.data[key] // Map optionally returns a bool if the key exists or not.

	return value, exists
}
