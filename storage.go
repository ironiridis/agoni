package agoni

import "errors"

// Key is a string value that uniquely identifies an index in the KeyValueStore.
type Key string

// Value is an arbitrary type representing data stored in the KeyValueStore.
type Value string

// KeyValueStore provides the core data structure.
type KeyValueStore struct {
	x    map[Key]Value
	subs []Subscription
	ch   chan *Operation
}

func (k *Key) String() string {
	return string(*k)
}

// Compare returns true if c is the same key.
func (k *Key) Compare(c *Key) bool {
	return (string(*k) == string(*c))
}

// NewKeyValueStore returns an instance of KeyValueStore.
func NewKeyValueStore() *KeyValueStore {
	kvs := KeyValueStore{
		x:    map[Key]Value{},
		subs: []Subscription{},
		ch:   make(chan *Operation, 0),
	}
	go kvs.exec()
	return &kvs
}

func (kvs KeyValueStore) exec() {
	var ok bool
	for op := range kvs.ch {
		switch o := (*op).(type) {
		case *NewOperation:
			// Test whether the key already exists
			_, ok = kvs.x[*(o.K)]
			if ok {
				o.Fail(errors.New("Key exists"))
				break
			}
			kvs.x[*(o.K)] = o.V
		case *UpdateOperation:
			o.OldV, ok = kvs.x[*(o.K)]
			if !ok {
				o.Fail(errors.New("Key doesn't exist"))
				break
			}
			kvs.x[*(o.K)] = o.V
		case *DeleteOperation:
			o.OldV, ok = kvs.x[*(o.K)]
			if !ok {
				o.Fail(errors.New("Key doesn't exist"))
				break
			}
			delete(kvs.x, *(o.K))
		}
	}
}

// Close makes the KeyValueStore immutable by terminating its goroutine.
func (kvs *KeyValueStore) Close() {
	close(kvs.ch)
}
