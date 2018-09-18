package transaction

import (
	"io"
)

const (
	IssuanceInputType uint8 = iota
	SpendInputType
	CoinbaseInputType
)

// SpendInput satisfies the TypedInput interface and represents a spend transaction.
type SpendInput struct {
	SpendCommitmentSuffix []byte   // The unconsumed suffix of the output commitment
	Arguments             [][]byte // Witness
	SpendCommitment
}

// NewSpendInput create a new SpendInput struct.
func NewSpendInput(arguments [][]byte, sourceID Hash, assetID AssetID, amount, sourcePos uint64, controlProgram []byte) *TxInput {
	sc := SpendCommitment{
		AssetAmount: AssetAmount{
			AssetId: &assetID,
			Amount:  amount,
		},
		SourceID:       sourceID,
		SourcePosition: sourcePos,
		VMVersion:      1,
		ControlProgram: controlProgram,
	}
	return &TxInput{
		AssetVersion: 1,
		TypedInput: &SpendInput{
			SpendCommitment: sc,
			Arguments:       arguments,
		},
	}
}

// InputType is the interface function for return the input type.
func (si *SpendInput) InputType() uint8 { return SpendInputType }

// SpendCommitment contains the commitment data for a transaction output.
type SpendCommitment struct {
	AssetAmount
	SourceID       Hash
	SourcePosition uint64
	VMVersion      uint64
	ControlProgram []byte
}


func (sc *SpendCommitment) writeExtensibleString(w io.Writer, suffix []byte, assetVersion uint64) error {
	//_, err := blockchain.WriteExtensibleString(w, suffix, func(w io.Writer) error {
	//	return sc.writeContents(w, suffix, assetVersion)
	//})
	return nil
}
