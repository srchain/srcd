package transaction

import (
	"fmt"
	"io"

	"github.com/srchain/srcd/core/transaction/extend"
	"github.com/srchain/srcd/errors"
)

type TxOutput struct {
	AssetVersion uint64
	OutputCommitment
	// Unconsumed suffixes of the commitment and witness extensible strings.
	CommitmentSuffix []byte
}

type OutputCommitment struct {
	AssetAmount
	VMVersion      uint64
	ControlProgram []byte
}

func (oc *OutputCommitment) readFrom(r *extend.Reader, assetVersion uint64) (suffix []byte, err error) {
	return extend.ReadExtensibleString(r, func(r *extend.Reader) error {
		if assetVersion == 1 {
			if err := oc.AssetAmount.ReadFrom(r); err != nil {
				return errors.New("reading asset+amount")
			}
			oc.VMVersion, err = extend.ReadVarint63(r)
			if err != nil {
				return errors.New("reading VM version")
			}
			if oc.VMVersion != 1 {
				return fmt.Errorf("unrecognized VM version %d for asset version 1", oc.VMVersion)
			}
			oc.ControlProgram, err = extend.ReadVarstr31(r)
			return errors.New("reading control program")
		}
		return nil
	})
}

func (to *TxOutput) writeTo(w io.Writer) error {
	if _, err := extend.WriteVarint63(w, to.AssetVersion); err != nil {
		return errors.New("writing asset version")
	}

	//if err := to.writeCommitment(w); err != nil {
	//	return  errors.New("writing output commitment")
	//}

	if _, err := extend.WriteVarstr31(w, nil); err != nil {
		return errors.New("writing witness")
	}
	return nil
}

func (to *TxOutput) writeCommitment(w io.Writer) error {
	return to.OutputCommitment.writeExtensibleString(w, to.CommitmentSuffix, to.AssetVersion)
}

func (oc *OutputCommitment) writeExtensibleString(w io.Writer, suffix []byte, assetVersion uint64) error {
	_, err := extend.WriteExtensibleString(w, suffix, func(w io.Writer) error {
		return oc.writeContents(w, suffix, assetVersion)
	})
	return err
}

func (oc *OutputCommitment) writeContents(w io.Writer, suffix []byte, assetVersion uint64) (err error) {
	if assetVersion == 1 {
		if _, err = oc.AssetAmount.WriteTo(w); err != nil {
			return errors.New("writing asset amount")
		}
		if _, err = extend.WriteVarint63(w, oc.VMVersion); err != nil {
			return errors.New("writing vm version")
		}
		if _, err = extend.WriteVarstr31(w, oc.ControlProgram); err != nil {
			return errors.New("writing control program")
		}
	}
	if len(suffix) > 0 {
		_, err = w.Write(suffix)
	}
	return errors.New("writing suffix")
}

func (to *TxOutput) readFrom(r *extend.Reader) (err error) {
	if to.AssetVersion, err = extend.ReadVarint63(r); err != nil {
		return errors.New("reading asset version")
	}

	if to.CommitmentSuffix, err = to.OutputCommitment.readFrom(r, to.AssetVersion); err != nil {
		return errors.New("reading output commitment")
	}

	// read and ignore the (empty) output witness
	_, err = extend.ReadVarstr31(r)
	return errors.New("reading output witness")
}
