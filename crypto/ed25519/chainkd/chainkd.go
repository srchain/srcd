package chainkd

import (
	"io"
	"crypto/rand"
	"crypto/hmac"
	"crypto/sha512"
	"srcd/crypto/ed25519/ecmath"
	"srcd/crypto/ed25519"
)

type (
	//XPrv external private key
	XPrv [64]byte
	//XPub external public key
	XPub [64]byte
)

// NewXPrv takes a source of random bytes and produces a new XPrv.
// If r is nil, crypto/rand.Reader is used.
func NewXPrv(r io.Reader) (xprv XPrv, err error) {
	if r == nil {
		r = rand.Reader
	}
	var entropy [32]byte
	_, err = io.ReadFull(r, entropy[:])
	if err != nil {
		return xprv, err
	}
	return RootXPrv(entropy[:]), nil
}

// RootXPrv takes a seed binary string and produces a new xprv.
func RootXPrv(seed []byte) (xprv XPrv) {
	h := hmac.New(sha512.New, []byte{'R', 'o', 'o', 't'})
	h.Write(seed)
	h.Sum(xprv[:0])
	pruneRootScalar(xprv[:32])
	return
}

// s must be >= 32 bytes long and gets rewritten in place.
// This is NOT the same pruning as in Ed25519: it additionally clears the third
// highest bit to ensure subkeys do not overflow the second highest bit.
func pruneRootScalar(s []byte) {
	s[0] &= 248
	s[31] &= 31 // clear top 3 bits
	s[31] |= 64 // set second highest bit
}


// XPub derives an extended public key from a given xprv.
func (xprv XPrv) XPub() (xpub XPub) {
	var scalar ecmath.Scalar
	copy(scalar[:], xprv[:32])

	var P ecmath.Point
	P.ScMulBase(&scalar)
	buf := P.Encode()

	copy(xpub[:32], buf[:])
	copy(xpub[32:], xprv[32:])

	return
}

// PublicKey extracts the ed25519 public key from an xpub.
func (xpub XPub) PublicKey() ed25519.PublicKey {
	return ed25519.PublicKey(xpub[:32])
}

func (xprv XPrv)Sign(msg []byte) []byte {
	return Ed25519InnerSign(xprv.ExpandedPrivateKey(), msg)
	//return nil
}
// ExpandedPrivateKey generates a 64-byte key where
// the first half is the scalar copied from xprv,
// and the second half is the `prefix` is generated via PRF
// from the xprv.
func (xprv XPrv) ExpandedPrivateKey() ExpandedPrivateKey {
	var res [64]byte
	h := hmac.New(sha512.New, []byte{'E', 'x', 'p', 'a', 'n', 'd'})
	h.Write(xprv[:])
	h.Sum(res[:0])
	copy(res[:32], xprv[:32])
	return res[:]
}