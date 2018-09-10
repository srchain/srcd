package transaction

import (
	"srcd/crypto/sha3pool"
	"bytes"
	"encoding/hex"
	"io"
	"github.com/bytom/errors"
	"github.com/bytom/encoding/blockchain"
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

// Tx holds a transaction along with its hash.
type Tx struct {
	TxData
	TxWrap `json:"-"`
}

func NewTx(data TxData) Tx{
	return Tx{
		data,
		MapTx(data),
	}
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
	ew := errors.NewWriter(w)
	if err := tx.writeTo(ew, serRequired); err != nil {
		return 0, err
	}
	return ew.Written(), ew.Err()
}

func (tx *TxData) writeTo(w io.Writer, serflags byte) error {
	if _, err := w.Write([]byte{serflags}); err != nil {
		return errors.New("write byte error")
	}
	if _, err := blockchain.WriteVarint63(w, tx.Version); err != nil {
		return errors.Wrap(err, "writing transaction version")
	}
	if _, err := blockchain.WriteVarint63(w, tx.TimeRange); err != nil {
		return errors.Wrap(err, "writing transaction maxtime")
	}

	if _, err := blockchain.WriteVarint31(w, uint64(len(tx.Inputs))); err != nil {
		return errors.Wrap(err, "writing tx input count")
	}

	//for i, ti := range tx.Inputs {
	//	if err := ti.writeTo(w); err != nil {
	//		return errors.Wrapf(err, "writing tx input %d", i)
	//	}
	//}

	if _, err := blockchain.WriteVarint31(w, uint64(len(tx.Outputs))); err != nil {
		return errors.Wrap(err, "writing tx output count")
	}

	//for i, to := range tx.Outputs {
	//	if err := to.writeTo(w); err != nil {
	//		return errors.Wrapf(err, "writing tx output %d", i)
	//	}
	//}
	return nil
}
