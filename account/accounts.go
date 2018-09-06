package account

import (
	"srcd/crypto/ed25519/chainkd"
	"srcd/crypto/ripemd160"
	address2 "srcd/wallet/address"
	"srcd/params"
	"srcd/protocol/vm"
)


type CtrlProgram struct {
	AccountID      string
	Address        string
	KeyIndex       uint64
	ControlProgram []byte
	Change         bool // Mark whether this control program is for UTXO change
}

func  CreateP2PKH(xpub chainkd.XPub) (*CtrlProgram, error) {
	pubKey := xpub.PublicKey()
	pubHash := ripemd160.Ripemd160(pubKey)

	// TODO: pass different params due to config
	address, err := address2.NewAddressWitnessPubKeyHash(pubHash, &params.TestNetParams)
	if err != nil {
		return nil, err
	}

	control, err := vm.P2WSHProgram([]byte(pubHash))
	if err != nil {
		return nil, err
	}

	return &CtrlProgram{
		Address:        address.EncodeAddress(),
		ControlProgram: control,
	}, nil
}