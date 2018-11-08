package transaction

import (
	"fmt"

	"github.com/srchain/srcd/crypto/ed25519/chainkd"
)

func TxSign(tpl *Template,xprv chainkd.XPrv,xpub chainkd.XPub) error{
	h := tpl.Hash(0).Byte32()
	sig := xprv.Sign(h[:])
	pub := xpub.PublicKey()
	// Test with more signatures than required, in correct order
	tpl.SigningInstructions = []*SigningInstruction{{
		WitnessComponents: []witnessComponent{
			&RawTxSigWitness{
				Quorum: 1,
				Sigs:   []HexBytes{sig},
			},
			DataWitness([]byte(pub)),
		},
	}}
	//return nil
	return materializeWitnesses(tpl)
}

func materializeWitnesses(txTemplate *Template) error {
	msg := txTemplate.Transaction
	for i, sigInst := range txTemplate.SigningInstructions {
		var witness [][]byte
		for j, wc := range sigInst.WitnessComponents {
			err := wc.materialize(&witness)
			if err != nil {
				fmt.Printf("error in witness component %d of input %d", j, i)
			}
		}
		msg.SetInputArguments(sigInst.Position, witness)
	}

	return nil
}
