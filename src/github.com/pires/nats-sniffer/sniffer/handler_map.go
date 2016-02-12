package sniffer

import (
	"sync"
)

//template type ConcurrentMap(Key, Value)

type cmapHandlerMap map[string]SniffedMessageHandler

// A "thread" safe map of type Key:Value
type HandlerMap struct {
	items map[string]SniffedMessageHandler
	mutex sync.RWMutex
}

// Creates a new concurrent map.
func NewHandlerMap() *HandlerMap {
	return &HandlerMap{
		items: make(map[string]SniffedMessageHandler),
		mutex: sync.RWMutex{},
	}
}

// Sets the given value under the specified key.
func (this *HandlerMap) Set(key string, value SniffedMessageHandler) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.items[key] = value
}

// Removes an element from the map.
func (this *HandlerMap) Remove(key string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	delete(this.items, key)
}

// Retrieves an element from map under given key.
func (this *HandlerMap) Get(key string) (SniffedMessageHandler, bool) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	v, ok := this.items[key]
	return v, ok
}

// Retrieves an element from map under given key.
// If it exists, removes it from map.
func (this *HandlerMap) GetAndRemove(key string) (SniffedMessageHandler, bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	v, ok := this.items[key]
	if ok {
		delete(this.items, key)
	}
	return v, ok
}

// Returns the number of elements within the map.
func (this *HandlerMap) Count() int {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	return len(this.items)
}

// Looks up an item under specified key
func (this *HandlerMap) Has(key string) bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	// See if element is within shard.
	_, ok := this.items[key]
	return ok
}

// Checks if map is empty.
func (this *HandlerMap) IsEmpty() bool {
	return this.Count() == 0
}

// Returns a slice containing all map keys
func (this *HandlerMap) Keys() []string {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	keys := make([]string, 0, len(this.items))
	for k, _ := range this.items {
		keys = append(keys, k)
	}
	return keys
}

// Returns a slice containing all map values
func (this *HandlerMap) Values() []SniffedMessageHandler {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	values := make([]SniffedMessageHandler, 0, len(this.items))
	for _, v := range this.items {
		values = append(values, v)
	}
	return values
}

// Returns a <strong>snapshot</strong> (copy) of current map items which could be used in a for range loop.
// One <strong>CANNOT</strong> change the contents of this map by means of this method, since it returns only a copy.
func (this *HandlerMap) Iter() map[string]SniffedMessageHandler {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	results := make(map[string]SniffedMessageHandler, len(this.items))
	for k, v := range this.items {
		results[k] = v
	}

	return results
}
