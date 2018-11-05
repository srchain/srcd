package account

import (
	"srcd/account/wallet/address"
	"srcd/core/vm"
	"srcd/crypto/ed25519/chainkd"
	"srcd/crypto/ripemd160"
	"srcd/params"
)

type CtrlProgram struct {
	AccountID      string
	Address        string
	KeyIndex       uint64
	ControlProgram []byte
	Change         bool    // Mark whether this control program is for UTXO change
}

func CreateP2PKH(xpub chainkd.XPub) (*CtrlProgram, []byte, error) {
	pubKey := xpub.PublicKey()
	pubHash := ripemd160.Ripemd160(pubKey)

	address, err := address.NewAddressWitnessPubKeyHash(pubHash, &params.TestNetParams)
	if err != nil {
		return nil, nil, err
	}

	control, err := vm.P2WSHProgram([]byte(pubHash))
	if err != nil {
		return nil, nil, err
	}

	return &CtrlProgram{
		Address:        address.EncodeAddress(),
		ControlProgram: control,
	}, pubHash, nil
}
