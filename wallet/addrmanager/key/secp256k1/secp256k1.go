package secp256k1

import (
	"io"
	"fmt"
	"math/big"
	"srcd/wallet/addrmanager/key/elliptic"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"srcd/wallet/addrmanager/key/base58"
)

var secp256k1 elliptic.EllipticCurve

func init() {
	/* See Certicom's SEC2 2.7.1, pg.15 */
	/* secp256k1 elliptic curve parameters */
	secp256k1.P, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	secp256k1.A, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
	secp256k1.B, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	secp256k1.G.X, _ = new(big.Int).SetString("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798", 16)
	secp256k1.G.Y, _ = new(big.Int).SetString("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8", 16)
	secp256k1.N, _ = new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	secp256k1.H, _ = new(big.Int).SetString("01", 16)
}

//GenerateKey generates an ECDSA secp256k1 keypair
func GenerateKey(rand io.Reader) (priv PrivateKey, err error) {
	/* See Certicom's SEC1 3.2.1, pg.23 */
	/* See NSA's Suite B Implementer’s Guide to FIPS 186-3 (ECDSA) A.1.1, pg.18 */

	/* Select private key d randomly from [1, n) */

	/* Read N bit length random bytes + 64 extra bits  */
	b := make([]byte, secp256k1.N.BitLen()/8+8)
	_, err = io.ReadFull(rand, b)
	if err != nil {
		return priv, fmt.Errorf("Reading random reader: %v", err)
	}

	d := new(big.Int).SetBytes(b)

	/* Mod n-1 to shift d into [0, n-1) range */
	d.Mod(d, new(big.Int).Sub(secp256k1.N, big.NewInt(1)))
	/* Add one to shift d to [1, n) range */
	d.Add(d, big.NewInt(1))

	priv.D = d

	/* Derive public key from private key */
	priv.derive()

	return priv, nil
}

// derive derives a Bitcoin public key from a Bitcoin private key.
func (priv *PrivateKey) derive() (pub *PublicKey) {
	/* See Certicom's SEC1 3.2.1, pg.23 */

	/* Derive public key from Q = d*G */
	Q := secp256k1.ScalarBaseMult(priv.D)

	/* Check that Q is on the curve */
	if !secp256k1.IsOnCurve(Q) {
		panic("Catastrophic math logic failure in public key derivation.")
	}

	priv.X = Q.X
	priv.Y = Q.Y

	return &priv.PublicKey
}


// ToAddress converts a SRC public key to a compressed SRC address string.
func (pub *PublicKey) Toddress() (address string) {

	/* Convert the public key to bytes */
	pub_bytes := pub.ToBytes()

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