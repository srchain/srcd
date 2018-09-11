package discover5

import (
	"bytes"
	"os"
	"encoding/binary"
	"sync"
	"github.com/btcsuite/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/srchain/srcd/rlp"
)

var (
	nodeDBNilNodeID = NodeID{}
	nodeDBItemPrefix = []byte("n:")
)


var (
	nodeDBVersionKey = []byte("version")
)

type nodeDB struct {
	lvl    *leveldb.DB   // Interface to the database itself
	self   NodeID        // Own node id to prevent adding it into the database
	runner sync.Once     // Ensures we can start at most one expirer
	quit   chan struct{} // Channel to signal the expiring thread to stop
}

// newNodeDB creates a new node database for storing and retrieving infos about
// known peers in the network. If no path is given, an in-memory, temporary
// database is constructed.
func newNodeDB(path string, version int, self NodeID) (*nodeDB, error) {
	if path == "" {
		return newMemoryNodeDB(self)
	}
	return newPersistentNodeDB(path, version, self)
}

// newMemoryNodeDB creates a new in-memory node database without a persistent
// backend.
func newMemoryNodeDB(self NodeID) (*nodeDB, error) {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return nil, err
	}
	return &nodeDB{
		lvl:  db,
		self: self,
		quit: make(chan struct{}),
	}, nil
}

func newPersistentNodeDB(path string, version int, self NodeID) (*nodeDB, error) {
	opts := &opt.Options{OpenFilesCacheCapacity: 5}
	db, err := leveldb.OpenFile(path,opts)
	if _, iscorrupted := err.(* errors.ErrCorrupted); iscorrupted {
		db, err = leveldb.RecoverFile(path,nil)
	}
	if err != nil {
		return nil , err
	}
	currentVer := make([]byte,binary.MaxVarintLen64)
	currentVer = currentVer[:binary.PutVarint(currentVer,int64(version))]

	blob, err := db.Get(nodeDBVersionKey,nil)
	switch err {
	case leveldb.ErrNotFound:

		// 找不到版本，进行插入
		if err := db.Put(nodeDBVersionKey, currentVer,nil); err != nil {
			db.Close()
			return nil, err
		}
	case nil:

		// 已经有版本信息了且并不相同 ，flush 掉
		if !bytes.Equal(blob,currentVer) {
			db.Close()
			if err = os.RemoveAll(path); err != nil {
				return nil, err
			}
			return newPersistentNodeDB(path, version,self)
		}

	}

	return &nodeDB{
		lvl: db,
		self: self,
		quit: make(chan struct{}),
	}, nil
}

// 生成一个leveldb node 缓存key
func makeKey(id NodeID, field string) []byte {
	if bytes.Equal(id[:],nodeDBNilNodeID[:]) {
		return []byte(field)
	}
	return append(nodeDBItemPrefix, append(id[:],field...)...)
}

// 将leveldb node 缓存key的 field 分离
func splitKey(key []byte) (id NodeID, field string) {
	if !bytes.HasPrefix(key,nodeDBItemPrefix) {
		return NodeID{}, string(key)
	}
	item := key[len(nodeDBItemPrefix):]
	copy(id[:], item[:len(id)])
	field = string((item[len(id):]))
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
	blob = blob[:binary.PutVarint(blob,n)]
	return db.lvl.Put(key,blob,nil)
}

func (db *nodeDB) storeRLP(key []byte, val interface{}) error {
	blob, err := rlp.EncodeToBytes(val)
	if err != nil {
		return err
	}
	return db.lvl.Put(key,blob,nil)
}


