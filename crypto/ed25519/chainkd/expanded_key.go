package chainkd

import (
	"strconv"
	"crypto/sha512"
	"srcd/crypto/ed25519/internal/edwards25519"
	"crypto"
	"srcd/crypto/ed25519"
)
const (
	// ExpandedPrivateKeySize is the size, in bytes, of a "secret key" as defined in NaCl.
	ExpandedPrivateKeySize = 64
)

type ExpandedPrivateKey []byte
// Public returns the PublicKey corresponding to secret key.
func (priv ExpandedPrivateKey) Public() crypto.PublicKey {
	var A edwards25519.ExtendedGroupElement
	var scalar [32]byte
	copy(scalar[:], priv[:32])
	edwards25519.GeScalarMultBase(&A, &scalar)
	var publicKeyBytes [32]byte
	A.ToBytes(&publicKeyBytes)
	return ed25519.PublicKey(publicKeyBytes[:])
}

func Ed25519InnerSign(privateKey ExpandedPrivateKey, message []byte) []byte{
	if l := len(privateKey); l != ExpandedPrivateKeySize {
		panic("ed25519: bad private key length: " + strconv.Itoa(l))
	}

	var messageDigest, hramDigest [64]byte

	h := sha512.New()
	h.Write(privateKey[32:])
	h.Write(message)
	h.Sum(messageDigest[:0])

	var messageDigestReduced [32]byte
	edwards25519.ScReduce(&messageDigestReduced, &messageDigest)
	var R edwards25519.ExtendedGroupElement
	edwards25519.GeScalarMultBase(&R, &messageDigestReduced)

	var encodedR [32]byte
	R.ToBytes(&encodedR)

	publicKey := privateKey.Public().(ed25519.PublicKey)
	h.Reset()
	h.Write(encodedR[:])
	h.Write(publicKey[:])
	h.Write(message)
	h.Sum(hramDigest[:0])
	var hramDigestReduced [32]byte
	edwards25519.ScReduce(&hramDigestReduced, &hramDigest)

	var sk [32]byte
	copy(sk[:], privateKey[:32])
	var s [32]byte
	edwards25519.ScMulAdd(&s, &hramDigestReduced, &sk, &messageDigestReduced)

	signature := make([]byte, ed25519.SignatureSize)
	copy(signature[:], encodedR[:])
	copy(signature[32:], s[:])

	return signature
}