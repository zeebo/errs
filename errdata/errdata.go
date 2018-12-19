// Package errdata helps with associating some data to error classes
package errdata

import (
	"sync"

	"github.com/zeebo/errs"
)

// registry is a concurrent map[key]interface{}. we use this because it is
// expected to be frequently read, with a one time initial set of writes.
var registry sync.Map

// key is the type of keys for the registry map.
type key struct {
	class *errs.Class
	key   interface{}
}

// makeKey is a helper to create a key for the registry map.
func makeKey(class *errs.Class, k interface{}) key {
	return key{
		class: class,
		key:   k,
	}
}

// Set associates the value for the given key and class. Errors wrapped by the
// class will return the value in the call to Get for the key.
func Set(class *errs.Class, key interface{}, value interface{}) {
	registry.Store(makeKey(class, key), value)
}

// Get returns the value associated to the key for the error if any of the
// classes the error is part of have a value associated for that key.
func Get(err error, key interface{}) interface{} {
	for _, class := range errs.Classes(err) {
		value, ok := registry.Load(makeKey(class, key))
		if ok {
			return value
		}
	}
	return nil
}
