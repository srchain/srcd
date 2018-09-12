package discover5

import (
	"sync"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"encoding/binary"
	"bytes"
	"os"
	"github.com/srchain/srcd/rlp"
	"fmt"
	"github.com/srchain/srcd/log"
)

var (
	nodeDBNilNodeID      = NodeID{}
	nodeDBVersionKey = []byte("version") // Version of the database to flush if changes
	nodeDBItemPrefix = []byte("n:")      // Identifier to prefix node entries with

)

// nodeDB stores all nodes we know about.
type nodeDB struct {
	lvl    *leveldb.DB   // Interface to the database itself
	self   NodeID        // Own node id to prevent adding it into the database
	runner sync.Once     // Ensures we can start at most one expirer
	quit   chan struct{} // Channel to signal the expiring thread to stop
}

func newNodeDB(path string, version int, self NodeID) (*nodeDB, error) {
	if path == "" {
		return newMemoryNodeDB(self)
	}
	return newPersistentNodeDB(path,version,self)
}
func newMemoryNodeDB(self NodeID) (*nodeDB, error) {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return nil, err
	}
	return &nodeDB{
		lvl: db,
		self: self,
		quit: make(chan struct{}),
	}, nil
}

func newPersistentNodeDB(path string, version int , self NodeID) (*nodeDB, error) {
	opts := &opt.Options{OpenFilesCacheCapacity: 5}
	db, err := leveldb.OpenFile(path, opts)
	if _, iscorrupted := err.(*errors.ErrCorrupted); iscorrupted {
		db, err = leveldb.RecoverFile(path,nil)
	}
	if err != nil {
		return nil, err
	}

	currentVer := make([]byte, binary.MaxVarintLen64)
	currentVer = currentVer[:binary.PutVarint(currentVer, int64(version))]

	blob, err := db.Get(nodeDBVersionKey,nil)
	switch err {
	case leveldb.ErrNotFound:
		if err := db.Put(nodeDBVersionKey, currentVer,nil); err != nil {
			db.Close()
			return nil, err
		}
	case nil:
		if !bytes.Equal(blob, currentVer) {
			db.Close()
			if err = os.RemoveAll(path); err != nil {
				return nil, err
			}
			return newPersistentNodeDB(path, version, self)
		}
	}

	return &nodeDB{
		lvl: db,
		self: self,
		quit: make(chan struct{}),
	}, nil

}

func makeKey(id NodeID, field string) []byte {
	if bytes.Equal(id[:], nodeDBNilNodeID[:]) {
		return []byte(field)
	}
}

func splitKey(key []byte) (id NodeID, field string) {
	if !bytes.HasPrefix(key, nodeDBItemPrefix) {
		return  NodeID{}, string(key)
	}

	item := key[len(nodeDBItemPrefix):]
	copy(id[:],item[:len(id)])
	field = string(item[len(id):])
	return id, field
}

func (db *nodeDB) fetchInt64(key []byte) int64 {
	blob, err := db.lvl.Get(key,nil)
	if err != nil {
		return 0
	}
	val, read := binary.Varint(blob)
	if read <= 0 {
		return 0
	}
	return val
}

func (db *nodeDB) storeInt64(key []byte, n int64) error {
	blob := make([]byte, binary.MaxVarintLen64)
	blob = blob[:binary.PutVarint(blob, n )]
	return db.lvl.Put(key, blob,nil)
}


func (db *nodeDB) storeRLP(key []byte, val interface{}) error {
	blob, err := rlp.EncodeToBytes(val)
	if err != nil {
		return err
	}
	return db.lvl.Put(key, blob, nil)
}

func (db *nodeDB) fetchRLP(key []byte, val interface{}) error {
	blob, err := db.lvl.Get(key,nil)
	if err != nil {
		return err
	}
	err = rlp.DecodeBytes(blob, val)
	if err != nil {
		log.Warn(fmt.Sprintf("key &x (%T) %v",key,val, err))
	}
	return err
}

func (db *nodeDB) node(id NodeID) *Node {
	var node Node
}