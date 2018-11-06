package address

import (
	"bytes"
	"fmt"
	"strings"

	"srcd/errors"
	"srcd/params"

	"github.com/bytom/common/bech32"
)

type Address interface {
	String() string
	EncodeAddress() string
}

// AddressWitnessPubKeyHash is an Address for a pay-to-witness-pubkey-hash
// (P2WPKH) output. See BIP 173 for further details regarding native segregated
// witness address encoding:
// https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki
type AddressWitnessPubKeyHash struct {
	hrp            string
	witnessVersion byte
	witnessProgram [20]byte
}

//NewAddressWitnessPubKeyHash returns a new AddressWitnessPubKeyHash.
func NewAddressWitnessPubKeyHash(witnessProg []byte, param *params.NetParams) (*AddressWitnessPubKeyHash, error) {
	return newAddressWitnessPubKeyHash(param.Bech32HRPSegwit, witnessProg)
}

// newAddressWitnessPubKeyHash is an internal helper function to create an
// AddressWitnessPubKeyHash with a known human-readable part, rather than
// looking it up through its parameters.
func newAddressWitnessPubKeyHash(hrp string, witnessProg []byte) (*AddressWitnessPubKeyHash, error) {
	// Check for valid program length for witness version 0, which is 20
	// for P2WPKH.
	if len(witnessProg) != 20 {
		return nil, errors.New("witness program must be 20 bytes for p2wpkh")
	}

	addr := &AddressWitnessPubKeyHash{
		hrp:            strings.ToLower(hrp),
		witnessVersion: 0x00,
	}

	copy(addr.witnessProgram[:], witnessProg)

	return addr, nil
}

// String returns a human-readable string for the AddressWitnessPubKeyHash.
// This is equivalent to calling EncodeAddress, but is provided so the type
// can be used as a fmt.Stringer.
// Part of the Address interface.
func (a *AddressWitnessPubKeyHash) String() string {
	return a.EncodeAddress()
}

// EncodeAddress returns the bech32 string encoding of an
// AddressWitnessPubKeyHash.
// Part of the Address interface.
func (a *AddressWitnessPubKeyHash) EncodeAddress() string {
	str, err := encodeSegWitAddress(a.hrp, a.witnessVersion, a.witnessProgram[:])
	if err != nil {
		return ""
	}
	return str
}

func DecodeAddress(addr string,param params.NetParams)(Address, error){
	oneIndex := strings.LastIndexByte(addr, '1')
	if oneIndex > 1 {
		prefix := addr[:oneIndex+1]
		if params.IsBech32SegwitPrefix(prefix, param) {
			witnessVer, witnessProg, err := decodeSegWitAddress(addr)
			if err != nil {
				return nil, err
			}

			// We currently only support P2WPKH and P2WSH, which is
			// witness version 0.
			if witnessVer != 0 {
				return nil, nil
			}

			// The HRP is everything before the found '1'.
			hrp := prefix[:len(prefix)-1]

			return newAddressWitnessPubKeyHash(hrp, witnessProg)
		}
	}
	return nil, nil
}

// encodeSegWitAddress creates a bech32 encoded address string representation
// from witness version and witness program.
func encodeSegWitAddress(hrp string, witnessVersion byte, witnessProgram []byte) (string, error) {
	// Group the address bytes into 5 bit groups, as this is what is used to
	// encode each character in the address string.
	converted, err := bech32.ConvertBits(witnessProgram, 8, 5, true)
	if err != nil {
		return "", err
	}

	// Concatenate the witness version and program, and encode the resulting
	// bytes using bech32 encoding.
	combined := make([]byte, len(converted)+1)
	combined[0] = witnessVersion
	copy(combined[1:], converted)
	bech, err := bech32.Bech32Encode(hrp, combined)
	if err != nil {
		return "", err
	}

	// Check validity by decoding the created address.
	version, program, err := decodeSegWitAddress(bech)
	if err != nil {
		return "", fmt.Errorf("invalid segwit address: %v", err)
	}

	if version != witnessVersion || !bytes.Equal(program, witnessProgram) {
		return "", fmt.Errorf("invalid segwit address")
	}

	return bech, nil
}

// decodeSegWitAddress parses a bech32 encoded segwit address string and
// returns the witness version and witness program byte representation.
func decodeSegWitAddress(address string) (byte, []byte, error) {
	// Decode the bech32 encoded address.
	_, data, err := bech32.Bech32Decode(address)
	if err != nil {
		return 0, nil, err
	}

	// The first byte of the decoded address is the witness version, it must
	// exist.
	if len(data) < 1 {
		return 0, nil, fmt.Errorf("no witness version")
	}

	// ...and be <= 16.
	version := data[0]
	if version > 16 {
		return 0, nil, fmt.Errorf("invalid witness version: %v", version)
	}

	// The remaining characters of the address returned are grouped into
	// words of 5 bits. In order to restore the original witness program
	// bytes, we'll need to regroup into 8 bit words.
	regrouped, err := bech32.ConvertBits(data[1:], 5, 8, false)
	if err != nil {
		return 0, nil, err
	}

	// The regrouped data must be between 2 and 40 bytes.
	if len(regrouped) < 2 || len(regrouped) > 40 {
		return 0, nil, fmt.Errorf("invalid data length")
	}

	// For witness version 0, address MUST be exactly 20 or 32 bytes.
	if version == 0 && len(regrouped) != 20 && len(regrouped) != 32 {
		return 0, nil, fmt.Errorf("invalid data length for witness "+
			"version 0: %v", len(regrouped))
	}

	return version, regrouped, nil
}
