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

type Reader struct {
	buf []byte
}
// ReadByte reads and returns the next byte from the input.
//
// It implements the io.ByteReader interface.
func (r *Reader) ReadByte() (byte, error) {
	if len(r.buf) == 0 {
		return 0, io.EOF
	}

	b := r.buf[0]
	r.buf = r.buf[1:]
	return b, nil
}
// Read reads up to len(p) bytes into p. It implements
// the io.Reader interface.
func (r *Reader) Read(p []byte) (n int, err error) {
	n = copy(p, r.buf)
	r.buf = r.buf[n:]
	if len(r.buf) == 0 {
		err = io.EOF
	}
	return
}
// NewReader constructs a new reader with the provided bytes. It
// does not create a copy of the bytes, so the caller is responsible
// for copying the bytes if necessary.
func NewReader(b []byte) *Reader {
	return &Reader{buf: b}
}

// Len returns the number of unread bytes.
func (r *Reader) Len() int {
	return len(r.buf)
}

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

func ReadVarint63(r *Reader) (uint64, error) {
	val, err := binary.ReadUvarint(r)
	if err != nil {
		return 0, err
	}
	if val > math.MaxInt64 {
		return 0, ErrRange
	}
	return val, nil
}
func ReadVarstr31(r *Reader) ([]byte, error) {
	l, err := ReadVarint31(r)
	if err != nil {
		return nil, err
	}
	if l == 0 {
		return nil, nil
	}
	if int(l) > len(r.buf) {
		return nil, io.ErrUnexpectedEOF
	}
	str := r.buf[:l]
	r.buf = r.buf[l:]
	return str, nil
}
func ReadVarint31(r *Reader) (uint32, error) {
	val, err := binary.ReadUvarint(r)
	if err != nil {
		return 0, err
	}
	if val > math.MaxInt32 {
		return 0, ErrRange
	}
	return uint32(val), nil
}

func WriteVarstr31(w io.Writer, str []byte) (int, error) {
	n, err := WriteVarint31(w, uint64(len(str)))
	if err != nil {
		return n, err
	}
	n2, err := w.Write(str)
	return n + n2, err
}

// ReadVarstrList reads a varint31 length prefix followed by
// that many varstrs.
func ReadVarstrList(r *Reader) (result [][]byte, err error) {
	nelts, err := ReadVarint31(r)
	if err != nil {
		return nil, err
	}
	if nelts == 0 {
		return nil, nil
	}

	for ; nelts > 0 && err == nil; nelts-- {
		var s []byte
		s, err = ReadVarstr31(r)
		result = append(result, s)
	}
	if len(result) < int(nelts) {
		err = io.ErrUnexpectedEOF
	}
	return result, err
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
// ReadExtensibleString reads a varint31 length prefix and that many
// bytes from r. It then calls the given function to consume those
// bytes, returning any unconsumed suffix.
func ReadExtensibleString(r *Reader, f func(*Reader) error) (suffix []byte, err error) {
	s, err := ReadVarstr31(r)
	if err != nil {
		return nil, err
	}

	sr := NewReader(s)
	//err = f(sr)
	//if err != nil {
	//	return nil, err
	//}
	return sr.buf, nil
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
