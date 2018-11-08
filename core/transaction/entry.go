package transaction

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"

	"github.com/srchain/srcd/core/transaction/extend"
	"github.com/srchain/srcd/crypto/sha3pool"
	"github.com/srchain/srcd/errors"
	"github.com/golang/protobuf/proto"
)

// Entry is the interface implemented by each addressable unit in a
// blockchain: transaction components such as spends, issuances,
// outputs, and retirements (among others), plus blockheaders.
type Entry interface {
	proto.Message

	// type produces a short human-readable string uniquely identifying
	// the type of this entry.
	typ() string

	// writeForHash writes the entry's body for hashing.
	writeForHash(w io.Writer)
}

func EntryID(e Entry) (hash Hash) {
	if e == nil {
		return hash
	}

	// Nil pointer; not the same as nil interface above. (See
	// https://golang.org/doc/faq#nil_error.)
	if v := reflect.ValueOf(e); v.Kind() == reflect.Ptr && v.IsNil() {
		return hash
	}

	hasher := sha3pool.Get256()
	defer sha3pool.Put256(hasher)

	hasher.Write([]byte("entryid:"))
	hasher.Write([]byte(e.typ()))
	hasher.Write([]byte{':'})

	bh := sha3pool.Get256()
	defer sha3pool.Put256(bh)

	e.writeForHash(bh)

	var innerHash [32]byte
	bh.Read(innerHash[:])

	hasher.Write(innerHash[:])

	hash.ReadFrom(hasher)
	return hash
}

var byte32zero [32]byte

func mustWriteForHash(w io.Writer, c interface{}) {
	if err := writeForHash(w, c); err != nil {
		panic(err)
	}
}

func writeForHash(w io.Writer, c interface{}) error {
	switch v := c.(type) {
	case byte:
		_, err := w.Write([]byte{v})
		return err
	case uint64:
		buf := [8]byte{}
		binary.LittleEndian.PutUint64(buf[:], v)
		_, err := w.Write(buf[:])
		return err
	case []byte:
		_, err := extend.WriteVarstr31(w, v)
		return err
	case [][]byte:
		_, err := extend.WriteVarstrList(w, v)
		return err
	case string:
		_, err := extend.WriteVarstr31(w, []byte(v))
		return err
	case *Hash:
		if v == nil {
			_, err := w.Write(byte32zero[:])
			return err
		}
		_, err := w.Write(v.Bytes())
		return err
	case *AssetID:
		if v == nil {
			_, err := w.Write(byte32zero[:])
			return err
		}
		_, err := w.Write(v.Bytes())
		return err
	case Hash:
		_, err := v.WriteTo(w)
		return err
	case AssetID:
		_, err := v.WriteTo(w)
		return err
	}

	// The two container types in the spec (List and Struct)
	// correspond to slices and structs in Go. They can't be
	// handled with type assertions, so we must use reflect.
	switch v := reflect.ValueOf(c); v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		elem := v.Elem()
		return writeForHash(w, elem.Interface())
	case reflect.Slice:
		l := v.Len()
		if _, err := extend.WriteVarint31(w, uint64(l)); err != nil {
			return err
		}
		for i := 0; i < l; i++ {
			c := v.Index(i)
			if !c.CanInterface() {
			}
			if err := writeForHash(w, c.Interface()); err != nil {
				return err
			}
		}
		return nil

	case reflect.Struct:
		typ := v.Type()
		for i := 0; i < typ.NumField(); i++ {
			c := v.Field(i)
			if !c.CanInterface() {
			}
			if err := writeForHash(w, c.Interface()); err != nil {
				t := v.Type()
				f := t.Field(i)
				fmt.Printf("writing struct field %d (%s.%s) for hash", i, t.Name(), f.Name)
				return err
			}
		}
		return nil
	}

	return errors.New("bad type ")
}
