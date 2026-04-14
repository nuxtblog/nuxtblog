package storage

import (
	"fmt"
	"sync"

	"github.com/nuxtblog/nuxtblog/sdk"
)

// registeredAdapter holds a plugin-registered storage adapter and its type code.
type registeredAdapter struct {
	storageType int
	adapter     sdk.StorageAdapter
}

var (
	regMu       sync.RWMutex
	regAdapters = make(map[string]*registeredAdapter)
)

// RegisterAdapter registers a plugin-provided storage adapter under a unique name.
func RegisterAdapter(name string, storageType int, adapter sdk.StorageAdapter) error {
	regMu.Lock()
	defer regMu.Unlock()
	if _, exists := regAdapters[name]; exists {
		return fmt.Errorf("storage adapter %q already registered", name)
	}
	regAdapters[name] = &registeredAdapter{storageType: storageType, adapter: adapter}
	return nil
}

// UnregisterAdapter removes a previously registered adapter.
func UnregisterAdapter(name string) {
	regMu.Lock()
	defer regMu.Unlock()
	delete(regAdapters, name)
}

// GetRegisteredAdapter returns the adapter for the given name, or nil if not found.
func GetRegisteredAdapter(name string) (sdk.StorageAdapter, int, bool) {
	regMu.RLock()
	defer regMu.RUnlock()
	ra, ok := regAdapters[name]
	if !ok {
		return nil, 0, false
	}
	return ra.adapter, ra.storageType, true
}

// ListRegisteredAdapters returns a snapshot of all plugin-registered adapter names and their storage types.
func ListRegisteredAdapters() map[string]int {
	regMu.RLock()
	defer regMu.RUnlock()
	result := make(map[string]int, len(regAdapters))
	for name, ra := range regAdapters {
		result[name] = ra.storageType
	}
	return result
}
