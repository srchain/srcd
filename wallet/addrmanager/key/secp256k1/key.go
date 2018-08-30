package secp256k1

import (
	"math/big"
	"srcd/wallet/addrmanager/key/elliptic"
)

// PublicKey represents a Bitcoin public key.
type PublicKey struct {
	elliptic.Point
}

// PrivateKey represents a Bitcoin private key.
type PrivateKey struct {
	PublicKey
	D *big.Int
}
