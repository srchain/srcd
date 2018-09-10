package extend

import (
	"io"
	"srcd/protocol/transaction/bufpool"
	"encoding/binary"
	"srcd/common/math"
	"srcd/errors"
	"bytes"
)
var ErrRange = errors.New("value out of range")

func WriteVarint31(w io.Writer, val uint64) (int, error) {
	if val > math.MaxInt32 {
		return 0, ErrRange
	}
	var buf = make([]byte, 9)
	n := binary.PutUvarint(buf[:], val)
	b, err := w.Write(buf[:n])
	bufpool.Put(bytes.NewBuffer(buf))
	return b, err
}

func WriteVarstr31(w io.Writer, str []byte) (int, error) {
	n, err := WriteVarint31(w, uint64(len(str)))
	if err != nil {
		return n, err
	}
	n2, err := w.Write(str)
	return n + n2, err
}
func WriteVarint63(w io.Writer, val uint64) (int, error) {
	if val > math.MaxInt64 {
		return 0, ErrRange
	}
	buf := make([]byte, 9)
	n := binary.PutUvarint(buf[:], val)
	b, err := w.Write(buf[:n])
	bufpool.Put(bytes.NewBuffer(buf))
	return b, err
}
func WriteExtensibleString(w io.Writer,suffix []byte,f func(writer io.Writer) error)(int ,error)  {
	buf := bufpool.Get()
	defer bufpool.Put(buf)
	err := f(buf)
	if err != nil {
		return 0, err
	}
	if len(suffix) > 0 {
		_, err := buf.Write(suffix)
		if err != nil {
			return 0, err
		}
	}
	return WriteVarstr31(w, buf.Bytes())
}

// WriteVarstrList writes a varint31 length prefix followed by the
// elements of l as varstrs.
func WriteVarstrList(w io.Writer, l [][]byte) (int, error) {
	n, err := WriteVarint31(w, uint64(len(l)))
	if err != nil {
		return n, err
	}
	for _, s := range l {
		n2, err := WriteVarstr31(w, s)
		n += n2
		if err != nil {
			return n, err
		}
	}
	return n, err
}
