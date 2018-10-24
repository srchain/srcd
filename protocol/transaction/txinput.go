package transaction

import (
	"io"
	"errors"
	"srcd/protocol/transaction/extend"
	"fmt"
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

func (t *TxInput) readFrom(r *extend.Reader) (err error) {
	if t.AssetVersion, err = extend.ReadVarint63(r); err != nil {
		return err
	}

	//var assetID AssetID
	t.CommitmentSuffix, err = extend.ReadExtensibleString(r, func(r *extend.Reader) error {
		if t.AssetVersion != 1 {
			return nil
		}
		var icType [1]byte
		if _, err = io.ReadFull(r, icType[:]); err != nil {
			return errors.New("aaa")
		}
		switch icType[0] {

		case SpendInputType:
			si := new(SpendInput)
			t.TypedInput = si
			//if si.SpendCommitmentSuffix, err = si.SpendCommitment.readFrom(r, 1); err != nil {
			//	return err
			//}

		default:
			return fmt.Errorf("unsupported input type %d", icType[0])
		}
		return nil
	})
	if err != nil {
		return err
	}

	t.WitnessSuffix, err = extend.ReadExtensibleString(r, func(r *extend.Reader) error {
		if t.AssetVersion != 1 {
			return nil
		}

		switch inp := t.TypedInput.(type) {
		case *SpendInput:
			if inp.Arguments, err = extend.ReadVarstrList(r); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
