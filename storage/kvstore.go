package storage

import (
	"sync"
	"bufio"
	"strings"
	"os"
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

func (store *KVStore) Restore(filepath string) error {

	file, err := os.Open(filepath)
	if err != nil {
		// if the file is doesnt exists yet (first run), that is fine.
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")

		if len(parts) == 3 && parts[0] == "PUT" {
			key := strings.TrimSpace(parts[1])
			value := strings.TrimSpace(parts[2])
			/* 
			Not using Mutex lock or WAL append, as the server hasnt started yet, 
			No cleints are connected, 
			and we are just restoring the state from disk to RAM.
			*/
			store.data[key] = value 
		}
	}

	return scanner.Err()
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
