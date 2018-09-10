package transaction

import (
	 "srcd/crypto/ed25519/chainkd"
)

type UTXO struct {
	SourceID       Hash
	AssetID        AssetID
	Amount         uint64
	SourcePos      uint64	//utxo sourece index
	ControlProgram []byte  //receipt program
	Address        string  //receipt address
}

// UtxoToInputs convert an utxo to the txinput
func UtxoInputs(xpubs []chainkd.XPub, u *UTXO) (InputAndSigInst, error) {
	txInput := NewSpendInput(nil, u.SourceID, u.AssetID, u.Amount, u.SourcePos, u.ControlProgram)
	sigInst := &SigningInstruction{}

	if u.Address == "" {
		return InputAndSigInst{}, nil
	}

	//address, err := address2.DecodeAddress(u.Address,params.TestNetParams)
	//if err != nil {
	//	return nil, nil, err
	//}

	derivedPK := xpubs[0].PublicKey()
	sigInst.WitnessComponents = append(sigInst.WitnessComponents, DataWitness([]byte(derivedPK)))

	return InputAndSigInst{txInput,sigInst},nil
	//return txInput, sigInst, nil
}

//convert an utxo to th txoutput
func UtxoOutputs(assetID AssetID,amount uint64,controlProgram []byte)TxOutput  {

	return TxOutput{
		AssetVersion: 1,
		OutputCommitment: OutputCommitment{
			AssetAmount: AssetAmount{
				AssetId: &assetID,
				Amount:  amount,
			},
			VMVersion:      1,
			ControlProgram: controlProgram,
		},
	}
}

