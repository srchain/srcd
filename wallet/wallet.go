package wallet

import (
	"srcd/wallet/addrmanager/key/secp256k1"
	"crypto/rand"
	"fmt"
	"crypto/sha256"
	"srcd/wallet/addrmanager/key/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	PublicKey  secp256k1.PublicKey
	PrivateKey secp256k1.PrivateKey
}

//NewWallet create new wallet
func NewWallet() *Wallet {

	priv, err := secp256k1.GenerateKey(rand.Reader)


	if(err != nil){
		fmt.Errorf("GenerateKey: %v", err)
	}
	return &Wallet{priv.PublicKey,priv}

}

// GetAddress converts a SRC public key to a compressed SRC address string.
func (w *Wallet)GetAddress() (address string){

	pub_bytes := w.PublicKey.ToBytes()

	/* SHA256 Hash */
	sha256_h := sha256.New()
	sha256_h.Reset()
	sha256_h.Write(pub_bytes)
	pub_hash_1 := sha256_h.Sum(nil)

	/* RIPEMD-160 Hash */
	ripemd160_h := ripemd160.New()
	ripemd160_h.Reset()
	ripemd160_h.Write(pub_hash_1)
	pub_hash_2 := ripemd160_h.Sum(nil)

	/* Convert hash bytes to base58 check encoded sequence */
	address = base58.B58checkencode(0x00, pub_hash_2)
	return address
}