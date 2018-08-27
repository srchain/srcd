package srcdb

import (
	"io/ioutil"
	"os"

	"testing"
	"bytes"
)

func newTestLDB() (*LDBDatabase,func()) {
	dirname, err := ioutil.TempDir(os.TempDir(),"ethdb_test_")
	if err != nil {
		panic("failed to create test file" + err.Error())
	}

	db, err := NewLDBDatabase(dirname,0,0)
	if err != nil {
		panic("failed to create test database" + err.Error())
	}

	return db, func() {
		db.Close()
		os.RemoveAll(dirname)
	}
}

func TestLDB_PutGet(t *testing.T) {
	db, remove := newTestLDB()
	defer  remove()
	testPutGet(db,t)
}

var test_values = []string{"","a","1251","\x00123\x00"}

func testPutGet(db *LDBDatabase, t *testing.T) {
	t.Parallel()


	for _, k := range test_values {
		err := db.Put([]byte(k), nil)
		if err != nil {
			t.Fatalf("put failed: %v", err)
		}
	}

	for _, k := range test_values {
		data, err := db.Get([]byte(k))
		if err != nil {
			t.Fatalf("get failed: %v",err)
		}
		if len(data) != 0 {
			t.Fatalf("get returned wrong result, got %q expected nil", string(data))
		}
	}

	_, err := db.Get([]byte("non-exist-key"))
	if err == nil {
		t.Fatalf("expect to return a not found error")
	}

	for _, v := range test_values {
		err := db.Put([]byte(v),[]byte(v))
		if err != nil {
			t.Fatalf("put failed: %v", err)
		}
	}

	for _, v := range test_values {
		data, err := db.Get([]byte(v))
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if !bytes.Equal(data,[]byte(v)) {
			t.Fatalf("get returned wrong result, get %q expected %q", string(data),v)
		}
	}

	for _, v := range test_values {
		err := db.Put([]byte(v),[]byte("?"))
		if err != nil {
			t.Fatalf("put override failed: %v",err)
		}
	}

	for _, v := range  test_values {
		data, err := db.Get([]byte(v))
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if !bytes.Equal(data, []byte("?")) {
			t.Fatalf("get returned wrong result, got %q expected ?", string(data))
		}
	}



}
