package sniffer

import (
	"sync"
)

//template type ConcurrentMap(Key, Value)

type cmapSubjectHandlersMap map[string]*HandlerMap

// A "thread" safe map of type Key:Value
type SubjectHandlersMap struct {
	items map[string]*HandlerMap
	mutex sync.RWMutex
}

// Creates a new concurrent map.
func NewSubjectHandlersMap() *SubjectHandlersMap {
	return &SubjectHandlersMap{
		items: make(map[string]*HandlerMap),
		mutex: sync.RWMutex{},
	}
}

// Sets the given value under the specified key.
func (this *SubjectHandlersMap) Set(key string, value *HandlerMap) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.items[key] = value
}

// Removes an element from the map.
func (this *SubjectHandlersMap) Remove(key string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	delete(this.items, key)
}

// Retrieves an element from map under given key.
func (this *SubjectHandlersMap) Get(key string) (*HandlerMap, bool) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	v, ok := this.items[key]
	return v, ok
}

// Retrieves an element from map under given key.
// If it exists, removes it from map.
func (this *SubjectHandlersMap) GetAndRemove(key string) (*HandlerMap, bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	v, ok := this.items[key]
	if ok {
		delete(this.items, key)
	}
	return v, ok
}

// Removes an element from the map by key and value.
func (this *SubjectHandlersMap) RemoveWithValue(key string, value *HandlerMap) (*HandlerMap, bool) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	v, ok := this.items[key]
	if ok && v == value {
		delete(this.items, key)
	}
	return v, ok
}

// Iteratively removes all elements from the map that contain a value.
func (this *SubjectHandlersMap) IterRemoveWithValue(value *HandlerMap) uint {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var counter uint = 0
	for k, v := range this.items {
		if v == value {
			delete(this.items, k)
			counter++
		}
	}

	return counter
}

// Returns the number of elements within the map.
func (this *SubjectHandlersMap) Count() int {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	return len(this.items)
}

// Looks up an item under specified key
func (this *SubjectHandlersMap) Has(key string) bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	// See if element is within shard.
	_, ok := this.items[key]
	return ok
}

// Checks if map is empty.
func (this *SubjectHandlersMap) IsEmpty() bool {
	return this.Count() == 0
}

// Returns a slice containing all map keys
func (this *SubjectHandlersMap) Keys() []string {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	keys := make([]string, 0, len(this.items))
	for k, _ := range this.items {
		keys = append(keys, k)
	}
	return keys
}

// Returns a slice containing all map values
func (this *SubjectHandlersMap) Values() []*HandlerMap {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	values := make([]*HandlerMap, 0, len(this.items))
	for _, v := range this.items {
		values = append(values, v)
	}
	return values
}

// Returns a <strong>snapshot</strong> (copy) of current map items which could be used in a for range loop.
// One <strong>CANNOT</strong> change the contents of this map by means of this method, since it returns only a copy.
func (this *SubjectHandlersMap) Iter() map[string]*HandlerMap {
	this.mutex.RLock()
	defer this.mutex.RUnlock()

	results := make(map[string]*HandlerMap, len(this.items))
	for k, v := range this.items {
		results[k] = v
	}

	return results
}
