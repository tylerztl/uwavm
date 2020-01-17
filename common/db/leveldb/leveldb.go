package leveldb

import (
	"fmt"
	"path"
	"sync"
	"syscall"

	"github.com/BeDreamCoder/uwavm/common/log"
	"github.com/BeDreamCoder/uwavm/common/util"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	dbutil "github.com/syndtr/goleveldb/leveldb/util"
)

var logger = log.New("uwavm", "leveldb")

type dbState int32

const (
	closed dbState = iota
	opened
)

// Conf configuration for `DB`
type Conf struct {
	DBPath string
}

// DB - a wrapper on an actual store
type DB struct {
	conf    *Conf
	db      *leveldb.DB
	dbState dbState
	mux     sync.Mutex

	readOpts        *opt.ReadOptions
	writeOptsNoSync *opt.WriteOptions
	writeOptsSync   *opt.WriteOptions
}

// CreateDB constructs a `DB`
func CreateDB() *DB {
	dbPath := path.Join(util.GoPath(), "src/github.com/BeDreamCoder/uwavm/output/dbdata")
	readOpts := &opt.ReadOptions{}
	writeOptsNoSync := &opt.WriteOptions{}
	writeOptsSync := &opt.WriteOptions{}
	writeOptsSync.Sync = true

	return &DB{
		conf:            &Conf{DBPath: dbPath},
		dbState:         closed,
		readOpts:        readOpts,
		writeOptsNoSync: writeOptsNoSync,
		writeOptsSync:   writeOptsSync}
}

// Open opens the underlying db
func (db *DB) Open() {
	db.mux.Lock()
	defer db.mux.Unlock()
	if db.dbState == opened {
		return
	}
	dbOpts := &opt.Options{}
	dbPath := db.conf.DBPath
	var err error
	var dirEmpty bool
	if dirEmpty, err = util.CreateDirIfMissing(dbPath); err != nil {
		panic(fmt.Sprintf("Error creating dir if missing: %s", err))
	}
	dbOpts.ErrorIfMissing = !dirEmpty
	if db.db, err = leveldb.OpenFile(dbPath, dbOpts); err != nil {
		panic(fmt.Sprintf("Error opening leveldb: %s", err))
	}
	db.dbState = opened
}

// Close closes the underlying db
func (db *DB) Close() {
	db.mux.Lock()
	defer db.mux.Unlock()
	if db.dbState == closed {
		return
	}
	if err := db.db.Close(); err != nil {
		logger.Error("Error closing leveldb", "error", err)
	}
	db.dbState = closed
}

// Get returns the value for the given key
func (db *DB) Get(key []byte) ([]byte, error) {
	value, err := db.db.Get(key, db.readOpts)
	if err == leveldb.ErrNotFound {
		value = nil
		err = nil
	}
	if err != nil {
		logger.Error("Error retrieving leveldb key", "key", key, "error", err)
		return nil, errors.Wrapf(err, "error retrieving leveldb key [%#v]", key)
	}
	return value, nil
}

// Put saves the key/value
func (db *DB) Put(key []byte, value []byte) error {
	err := db.db.Put(key, value, db.writeOptsSync)
	if err != nil {
		logger.Error("Error writing leveldb key", "key", key)
		return errors.Wrapf(err, "error writing leveldb key [%#v]", key)
	}
	return nil
}

// Delete deletes the given key
func (db *DB) Delete(key []byte) error {
	err := db.db.Delete(key, db.writeOptsSync)
	if err != nil {
		logger.Error("Error deleting leveldb key", "key", key)
		return errors.Wrapf(err, "error deleting leveldb key [%#v]", key)
	}
	return nil
}

// GetIterator returns an iterator over key-value store. The iterator should be released after the use.
// The resultset contains all the keys that are present in the db between the startKey (inclusive) and the endKey (exclusive).
// A nil startKey represents the first available key and a nil endKey represent a logical key after the last available key
func (db *DB) GetIterator(startKey []byte, endKey []byte) iterator.Iterator {
	return db.db.NewIterator(&dbutil.Range{Start: startKey, Limit: endKey}, db.readOpts)
}

// WriteBatch writes a batch
func (db *DB) WriteBatch(batch *leveldb.Batch, sync bool) error {
	wo := db.writeOptsNoSync
	if sync {
		wo = db.writeOptsSync
	}
	if err := db.db.Write(batch, wo); err != nil {
		return errors.Wrap(err, "error writing batch to leveldb")
	}
	return nil
}

// FileLock encapsulate the DB that holds the file lock.
// As the FileLock to be used by a single process/goroutine,
// there is no need for the semaphore to synchronize the
// FileLock usage.
type FileLock struct {
	db       *leveldb.DB
	filePath string
}

// NewFileLock returns a new file based lock manager.
func NewFileLock(filePath string) *FileLock {
	return &FileLock{
		filePath: filePath,
	}
}

// Lock acquire a file lock. We achieve this by opening
// a db for the given filePath. Internally, leveldb acquires a
// file lock while opening a db. If the db is opened again by the same or
// another process, error would be returned. When the db is closed
// or the owner process dies, the lock would be released and hence
// the other process can open the db. We exploit this leveldb
// functionality to acquire and release file lock as the leveldb
// supports this for Windows, Solaris, and Unix.
func (f *FileLock) Lock() error {
	dbOpts := &opt.Options{}
	var err error
	var dirEmpty bool
	if dirEmpty, err = util.CreateDirIfMissing(f.filePath); err != nil {
		panic(fmt.Sprintf("Error creating dir if missing: %s", err))
	}
	dbOpts.ErrorIfMissing = !dirEmpty
	f.db, err = leveldb.OpenFile(f.filePath, dbOpts)
	if err != nil && err == syscall.EAGAIN {
		return errors.Errorf("lock is already acquired on file %s", f.filePath)
	}
	if err != nil {
		panic(fmt.Sprintf("Error acquiring lock on file %s: %s", f.filePath, err))
	}
	return nil
}

// Unlock releases a previously acquired lock. We achieve this by closing
// the previously opened db. FileUnlock can be called multiple times.
func (f *FileLock) Unlock() {
	if f.db == nil {
		return
	}
	if err := f.db.Close(); err != nil {
		logger.Warn("unable to release the lock on file", "filePath", f.filePath, "error", err)
		return
	}
	f.db = nil
}
