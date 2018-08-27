package srcdb

import (
	"github.com/btcsuite/goleveldb/leveldb"
	"sync"
	log "github.com/inconshreveable/log15"

	"github.com/btcsuite/goleveldb/leveldb/opt"
	"github.com/btcsuite/goleveldb/leveldb/filter"
	"github.com/btcsuite/goleveldb/leveldb/errors"
	"github.com/btcsuite/goleveldb/leveldb/iterator"
	"github.com/btcsuite/goleveldb/leveldb/util"
)

type LDBDatabase struct {
	filename string
	db *leveldb.DB

	quitLock sync.Mutex
	quitChan chan chan error

	log log.Logger
}

func NewLDBDatabase(file string, cache int , handles int) (*LDBDatabase, error) {
	logger := log.New("database",file)
	if cache < 16 {
		cache = 16
	}

	if handles < 16 {
		handles = 16
	}
	logger.Info("Allocated cache and file handles","cache",cache,"handles",handles)

	db, err := leveldb.OpenFile(file,&opt.Options{
		OpenFilesCacheCapacity: handles,
		BlockCacheCapacity: cache / 2 * opt.MiB,
		WriteBuffer:		cache / 4 * opt.MiB,
		Filter:				filter.NewBloomFilter(10),
	})
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(file,nil)
	}
	if err != nil {
		return nil, err
	}
	return  &LDBDatabase{
		filename: file,
		db:db,
		log:logger,
	} , nil


}

func (db *LDBDatabase) Path() string {
	return db.filename
}

func (db *LDBDatabase) Put(key []byte, value []byte) error {
	return db.db.Put(key,value,nil)
}

func (db *LDBDatabase) Has(key []byte) (bool, error) {
	return db.db.Has(key,nil)
}

func (db *LDBDatabase) Get(key []byte) ([]byte,error) {
	dat, err := db.db.Get(key,nil)
	if err != nil {
		return nil,err
	}
	return dat, nil
}

func (db *LDBDatabase) Delete(key []byte) error {
	return db.db.Delete(key,nil)
}

func (db *LDBDatabase) NewIterator() iterator.Iterator {
	return db.db.NewIterator(nil,nil)
}

func (db *LDBDatabase) NewIteratorWithPrefix(prefix []byte) iterator.Iterator {
	return db.db.NewIterator(util.BytesPrefix(prefix),nil)
}

func (db *LDBDatabase) Close() {
	db.quitLock.Lock()
	defer db.quitLock.Unlock()
	err := db.db.Close()
	if err == nil {
		db.log.Info("Database closed")
	} else {
		db.log.Error("Failed to close database","err",err)
	}
}

func (db *LDBDatabase) LDB() *leveldb.DB {
	return db.db
}