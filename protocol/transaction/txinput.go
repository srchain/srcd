package transaction

import (
	"io"
	"srcd/errors"
	"srcd/protocol/transaction/extend"
)

// TxInput is the top level struct of tx input.
type (
	TxInput struct {
		AssetVersion     uint64
		CommitmentSuffix []byte
		WitnessSuffix    []byte
		TypedInput
	}
	// TypedInput return the txinput type.
	TypedInput interface {
		InputType() uint8
	}
)


func (t *TxInput) writeTo(w io.Writer) error {
	if _, err := extend.WriteVarint63(w, t.AssetVersion); err != nil {
		return errors.New("write byte error")
	}
	if _, err := extend.WriteExtensibleString(w, t.CommitmentSuffix, t.writeInputCommitment); err != nil {
		return errors.New("write byte error")
	}
	_, err := extend.WriteExtensibleString(w, t.WitnessSuffix, t.writeInputWitness)

	return err
}

func (t *TxInput)writeInputCommitment(w io.Writer)(err error){

	switch inp := t.TypedInput.(type){
	case *SpendInput:
		if _, err = w.Write([]byte{SpendInputType}); err != nil {
			return err
		}
		return inp.SpendCommitment.writeExtensibleString(w, inp.SpendCommitmentSuffix, t.AssetVersion)
	}

	return nil
}
func (t *TxInput) writeInputWitness(w io.Writer) error {
	switch inp := t.TypedInput.(type) {
	case *SpendInput:
		_, err := extend.WriteVarstrList(w, inp.Arguments)
		return err
	}
	return nil
}