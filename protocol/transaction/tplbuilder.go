package transaction

import (
)

type InputAndSigInst struct {
	input *TxInput
	sigInst *SigningInstruction
}

func NewInputAndSigInst(input *TxInput,sigInst *SigningInstruction) InputAndSigInst {
	return InputAndSigInst{
		input:input,
		sigInst:sigInst,
	}
}

// Build build transactions with template
func BuildUtxoTemplate(inputs []InputAndSigInst, outputs []*TxOutput) (*Template, TxData, error) {
	tpl := &Template{}
	tx := TxData{}
	// Add all the built outputs.
	tx.Outputs = append(tx.Outputs, outputs...)

	// Add all the built inputs and their corresponding signing instructions.
	for _, in := range inputs {
		// Empty signature arrays should be serialized as empty arrays, not null.
		in.sigInst.Position = uint32(len(inputs))
		if in.sigInst.WitnessComponents == nil {
			in.sigInst.WitnessComponents = []witnessComponent{}
		}
		tpl.SigningInstructions = append(tpl.SigningInstructions, in.sigInst)
		tx.Inputs = append(tx.Inputs, in.input)
	}

	tpl.Transaction = NewTx(tx)
	return tpl, tx, nil
}


func BuildTransactionFromTx(baseTx Tx, outputIndex uint64, outputAmount uint64, ctrlProgram []byte)(*Tx,error){
	//spendInput, err := CreateSpendInput(baseTx, outputIndex)
	//if err != nil {
	//	return nil, err
	//}
	//
	//txInput := &TxInput{
	//	AssetVersion: 1,
	//	TypedInput:   spendInput,
	//}
	//
	//output := UtxoOutputs(*SRCAssetID, outputAmount, ctrlProgram)
	//builder := txbuilder.NewBuilder(time.Now())
	//builder.AddInput(txInput, &SigningInstruction{})
	//builder.AddOutput(output)
	//
	//tpl, _, err := builder.Build()
	//if err != nil {
	//	return nil, err
	//}
	//
	//txSerialized, err := tpl.Transaction.MarshalText()
	//if err != nil {
	//	return nil, err
	//}
	//
	//tpl.Transaction.Tx.SerializedSize = uint64(len(txSerialized))
	//tpl.Transaction.TxData.SerializedSize = uint64(len(txSerialized))
	//return tpl.Transaction, nil

}
