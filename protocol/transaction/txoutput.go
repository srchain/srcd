package transaction

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
