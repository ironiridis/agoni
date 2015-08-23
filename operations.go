package agoni

// OperationStatus indicates the state of an Operation.
type OperationStatus int

const (
	// OperationUninitialized is the zero value of an Operation.
	OperationUninitialized OperationStatus = iota
	// OperationReady indicates the Operation has been initialized.
	OperationReady
	// OperationCommitted indicates the data modification was applied to the KeyValueStore.
	OperationCommitted
	// OperationFailed indicates the data modification was not applied.
	OperationFailed
)

// An Operation is the representation of a single change to the KeyValueStore.
type Operation interface {
	Commit()
	Key() *Key
	Result(error)
}

// A DeleteOperation represents the removal of a Key from the KeyValueStore.
type DeleteOperation struct {
	KVS  *KeyValueStore
	K    *Key
	OldV Value
}

// A NewOperation represents the addition of a Key & Value to the KeyValueStore.
type NewOperation struct {
	KVS *KeyValueStore
	K   *Key
	V   Value
}

// An UpdateOperation represents the modification of a Key to a new Value in the KeyValueStore.
type UpdateOperation struct {
	KVS  *KeyValueStore
	K    *Key
	V    Value
	OldV Value
}

// Commit deletes the requested Key from the KeyValueStore.
func (o *DeleteOperation) Commit() {

}

// Key returns the Key that will be deleted.
func (o *DeleteOperation) Key() *Key {
	return o.K
}

func (o *DeleteOperation) Result(e error) {

}

// Commit creates the new Key and sets it to the Value specified earlier (TODO).
func (o *NewOperation) Commit() {

}

// Key returns the Key that will be created.
func (o *NewOperation) Key() *Key {
	return o.K
}

func (o *NewOperation) Result(e error) {

}

// Commit updates the existing Key to the Value specified earlier (TODO).
func (o *UpdateOperation) Commit() {

}

// Key returns the Key that will be updated.
func (o *UpdateOperation) Key() *Key {
	return o.K
}

func (o *UpdateOperation) Result(e error) {

}

// New returns a NewOperation, ready to Commit.
func (kvs *KeyValueStore) New(k *Key, v *Value) *NewOperation {
	return &NewOperation{KVS: kvs, K: k, V: *v}
}

// Update returns an UpdateOperation, ready to Commit.
func (kvs *KeyValueStore) Update(k *Key, v *Value) *UpdateOperation {
	return &UpdateOperation{KVS: kvs, K: k, V: *v}
}

// Delete returns a DeleteOperation, ready to Commit.
func (kvs *KeyValueStore) Delete(k *Key) *DeleteOperation {
	return &DeleteOperation{KVS: kvs, K: k}
}
