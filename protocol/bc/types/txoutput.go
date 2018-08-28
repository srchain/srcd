package types

// TxOutput is the top level struct of tx output.
type TxOutput struct {
	AssetVersion uint64
	OutputCommitment
	// Unconsumed suffixes of the commitment and witness extensible strings.
	CommitmentSuffix []byte
}

// OutputCommitment contains the commitment data for a transaction output.
type OutputCommitment struct {
	//bc.AssetAmount
	VMVersion      uint64
	ControlProgram []byte
}