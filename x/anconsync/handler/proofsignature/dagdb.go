package proofsignature

import (
	"context"
	"sync"

	"github.com/anconprotocol/node/x/anconsync"
	"github.com/pkg/errors"
	dbm "github.com/tendermint/tm-db"
)

var (
	// errBatchClosed is returned when a closed or written batch is used.
	errBatchClosed = errors.New("batch has been written or closed")

	// errKeyEmpty is returned when attempting to use an empty or nil key.
	errKeyEmpty = errors.New("key cannot be empty")

	// errValueNil is returned when attempting to set a nil value.
	errValueNil = errors.New("value cannot be nil")
)

func init() {

}

// item is a store.Item with byte slices as keys and values
type item struct {
	key   []byte
	value []byte
}

// newKey creates a new key item.
func newKey(key []byte) *item {
	return &item{key: key}
}

// newPair creates a new pair item.
func newPair(key, value []byte) *item {
	return &item{key: key, value: value}
}

// DagDB is an in-memory database backend using a B-tree for storage.
//
// For performance reasons, all given and returned keys and values are pointers to the in-memory
// database, so modifying them will cause the stored values to be modified as well. All DB methods
// already specify that keys and values should be considered read-only, but this is especially
// important with DagDB.
type DagDB struct {
	mtx   sync.RWMutex
	store *anconsync.Storage
}

var _ dbm.DB = (*DagDB)(nil)

// NewDagDB creates a new in-memory database.
func NewDagDB(s *anconsync.Storage) *DagDB {
	database := &DagDB{
			store: s,
	}
	return database
}

// Get implements DB.
func (db *DagDB) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errKeyEmpty
	}
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	i, err := db.store.DataStore.Get(context.Background(), string(key))
	if i != nil {
		return i, nil
	}
	return nil, err
}

// Has implements DB.
func (db *DagDB) Has(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, errKeyEmpty
	}
	db.mtx.RLock()
	defer db.mtx.RUnlock()

	return db.store.DataStore.Has(context.Background(), string(key))
}

// Set implements DB.
func (db *DagDB) Set(key []byte, value []byte) error {
	if len(key) == 0 {
		return errKeyEmpty
	}
	if value == nil {
		return errValueNil
	}
	db.mtx.Lock()
	defer db.mtx.Unlock()

	db.store.DataStore.Put(context.Background(), string(key), value)
	return nil
}

// SetSync implements DB.
func (db *DagDB) SetSync(key []byte, value []byte) error {
	return db.Set(key, value)
}

// Delete implements DB.
func (db *DagDB) Delete(key []byte) error {
	return errKeyEmpty
}

// DeleteSync implements DB.
func (db *DagDB) DeleteSync(key []byte) error {
	return db.Delete(key)
}

// Close implements DB.
func (db *DagDB) Close() error {
	// Close is a noop since for an in-memory database, we don't have a destination to flush
	// contents to nor do we want any data loss on invoking Close().
	// See the discussion in https://github.com/tendermint/tendermint/libs/pull/56
	return nil
}

// Print implements DB.
func (db *DagDB) Print() error {
	return nil
}

// Stats implements DB.
func (db *DagDB) Stats() map[string]string {
	return nil
}

// NewBatch implements DB.
func (db *DagDB) NewBatch() dbm.Batch {
	return nil
}

// dbm.Iterator implements DB.
// Takes out a read-lock on the database until the iterator is closed.
func (db *DagDB) Iterator(start, end []byte) (dbm.Iterator, error) {
	return nil, errKeyEmpty

}

// ReverseIterator implements DB.
// Takes out a read-lock on the database until the iterator is closed.
func (db *DagDB) ReverseIterator(start, end []byte) (dbm.Iterator, error) {
	return nil, errKeyEmpty

}

// IteratorNoMtx makes an iterator with no mutex.
func (db *DagDB) IteratorNoMtx(start, end []byte) (dbm.Iterator, error) {
	return nil, errKeyEmpty

}

// ReverseIteratorNoMtx makes an iterator with no mutex.
func (db *DagDB) ReverseIteratorNoMtx(start, end []byte) (dbm.Iterator, error) {
	return nil, errKeyEmpty
}
