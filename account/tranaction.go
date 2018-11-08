package account

import (
	"github.com/srchain/srcd/account/wallet/address"
	"github.com/srchain/srcd/core/vm"
	"github.com/srchain/srcd/crypto/ed25519/chainkd"
	"github.com/srchain/srcd/crypto/ripemd160"
	"github.com/srchain/srcd/params"
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
