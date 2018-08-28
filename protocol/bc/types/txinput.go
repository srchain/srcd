package types

type (
	// TxInput is the top level struct of tx input.
	TxInput struct {
		AssetVersion uint64
		TypedInput
		CommitmentSuffix []byte
		WitnessSuffix    []byte
	}

	// TypedInput return the txinput type.
	TypedInput interface {
		InputType() uint8
	}
)
