package transaction

import (
	"srcd/crypto/sha3pool"
	"bytes"
	"encoding/hex"
	"io"
	"srcd/protocol/transaction/extend"
	"fmt"
	"srcd/errors"
)

// Tx is a wrapper for the entries-based representation of a transaction.
type TxWrap struct {
	*TxHeader
	ID       Hash
	Entries  map[Hash]Entry
	InputIDs []Hash // 1:1 correspondence with TxData.Inputs

	SpentOutputIDs []Hash
	GasInputIDs    []Hash
}
// TxData encodes a transaction in the blockchain.
type TxData struct {
	Version        uint64
	SerializedSize uint64
	TimeRange      uint64
	Inputs         []*TxInput
	Outputs        []*TxOutput
}

// Output try to get the output entry by given hash
func (tx *TxWrap) Output(id Hash) (*Output, error) {
	e, ok := tx.Entries[id]
	if !ok || e == nil {
		return nil, errors.New("")
	}
	o, ok := e.(*Output)
	if !ok {
		return nil, errors.New("")
	}
	return o, nil
}

// Tx holds a transaction along with its hash.
type Tx struct {
	TxData
	TxWrap `json:"-"`
}

func NewTx(data TxData) Tx{
	return Tx{
		data,
		MapTxWrap(data),
	}
}

// UnmarshalText fulfills the encoding.TextUnmarshaler interface.
func (tx *Tx) UnmarshalText(p []byte) error {
	if err := tx.TxData.UnmarshalText(p); err != nil {
		return err
	}

	tx.TxWrap = MapTxWrap(tx.TxData)
	return nil
}

// UnmarshalText fulfills the encoding.TextUnmarshaler interface.
func (tx *TxData) UnmarshalText(p []byte) error {
	b := make([]byte, hex.DecodedLen(len(p)))
	if _, err := hex.Decode(b, p); err != nil {
		return err
	}

	r := extend.NewReader(b)
	if err := tx.readFrom(r); err != nil {
		return err
	}

	if trailing := r.Len(); trailing > 0 {
		return fmt.Errorf("trailing garbage (%d bytes)", trailing)
	}
	return nil
}
func (tx *TxData) readFrom(r *extend.Reader) (err error) {
	startSerializedSize := r.Len()
	var serflags [1]byte
	if _, err = io.ReadFull(r, serflags[:]); err != nil {
		return errors.New("reading serialization flags")
	}
	if serflags[0] != serRequired {
		return errors.New("unsupported serflags")
	}

	if tx.Version, err = extend.ReadVarint63(r); err != nil {
		return errors.New( "reading transaction version")
	}
	if tx.TimeRange, err = extend.ReadVarint63(r); err != nil {
		return err
	}

	n, err := extend.ReadVarint31(r)
	if err != nil {
		return errors.New("reading number of transaction inputs")
	}

	for ; n > 0; n-- {
		ti := new(TxInput)
		if err = ti.readFrom(r); err != nil {
			return errors.New("reading input ")
		}
		tx.Inputs = append(tx.Inputs, ti)
	}

	n, err = extend.ReadVarint31(r)
	if err != nil {
		return errors.New("reading number of transaction outputs")
	}

	for ; n > 0; n-- {
		to := new(TxOutput)
		if err = to.readFrom(r); err != nil {
			//return errors.Wrapf(err, "reading output %d", len(tx.Outputs))
		}
		tx.Outputs = append(tx.Outputs, to)
	}
	tx.SerializedSize = uint64(startSerializedSize - r.Len())
	return nil
}

// SigHash ...
func (tx *TxWrap) SigHash(n uint32) (hash Hash) {
	hasher := sha3pool.Get256()
	defer sha3pool.Put256(hasher)

	tx.InputIDs[n].WriteTo(hasher)
	tx.ID.WriteTo(hasher)
	hash.ReadFrom(hasher)
	return hash
}
// SetInputArguments sets the Arguments field in input n.
func (tx *Tx) SetInputArguments(n uint32, args [][]byte) {
	tx.Inputs[n].SetArguments(args)
	id := tx.TxWrap.InputIDs[n]
	e := tx.Entries[id]
	switch e := e.(type) {
	case *Spend:
		e.WitnessArguments = args
	}
}
// SetArguments set the args for the input
func (t *TxInput) SetArguments(args [][]byte) {
	switch inp := t.TypedInput.(type) {
	case *SpendInput:
		inp.Arguments = args
	}
}

//MarshalText fulfills the json.Marshaler interface.
func (tx *TxData) MarshalText() ([]byte, error) {
	var buf bytes.Buffer
	if _, err := tx.WriteTo(&buf); err != nil {
		return nil, nil
	}

	b := make([]byte, hex.EncodedLen(buf.Len()))
	hex.Encode(b, buf.Bytes())
	return b, nil
}
const serRequired = 0x7 // Bit mask accepted serialization flag.
// WriteTo writes tx to w.
func (tx *TxData) WriteTo(w io.Writer) (int64, error) {
	ew := extend.NewWriter(w)
	if err := tx.writeTo(ew, serRequired); err != nil {
		return 0, err
	}
	return ew.Written(), ew.Err()
}

func (tx *TxData) writeTo(w io.Writer, serflags byte) error {
	if _, err := w.Write([]byte{serflags}); err != nil {
		return errors.New("write byte error")
	}
	if _, err := extend.WriteVarint63(w, tx.Version); err != nil {
		return errors.New( "writing transaction version")
	}
	if _, err := extend.WriteVarint63(w, tx.TimeRange); err != nil {
		return errors.New("writing transaction maxtime")
	}

	if _, err := extend.WriteVarint31(w, uint64(len(tx.Inputs))); err != nil {
		return errors.New( "writing tx input count")
	}

	for _, ti := range tx.Inputs {
		if err := ti.writeTo(w); err != nil {
			return err
		}
	}

	if _, err := extend.WriteVarint31(w, uint64(len(tx.Outputs))); err != nil {
		return errors.New("writing tx output count")
	}

	for _, to := range tx.Outputs {
		if err := to.writeTo(w); err != nil {
			return err
		}
	}
	return nil
}
