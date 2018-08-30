package secp256k1

import "bytes"

// ToBytes converts a Bitcoin public key to a 33-byte byte slice with point compression.
func (pub *PublicKey) ToBytes() (b []byte) {
	/* See Certicom SEC1 2.3.3, pg. 10 */

	x := pub.X.Bytes()

	/* Pad X to 32-bytes */
	padded_x := append(bytes.Repeat([]byte{0x00}, 32-len(x)), x...)

	/* Add prefix 0x02 or 0x03 depending on ylsb */
	if pub.Y.Bit(0) == 0 {
		return append([]byte{0x02}, padded_x...)
	}

	return append([]byte{0x03}, padded_x...)
}