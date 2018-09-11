package address

import (
	"fmt"
	"testing"
	"strings"
)

const (
	addr = "sr1q0p8attf4l0vg6lqjfxf0rdenmpdqruhug6pd2t"
)
func TestAddrPrefix(t *testing.T)  {

	oneIndex := strings.LastIndexByte(addr, '1')
	prefix := addr[:oneIndex+1]
	hrp := prefix[:len(prefix)-1]

	fmt.Println(prefix)
	fmt.Println(hrp)


}

func TestDecodeSegWitAddress(t *testing.T)  {
	// Decode the bech32 encoded address.
	_, data, err := Bech32Decode(addr)
	if err != nil {
		fmt.Errorf("no witness version")
	}

	// The first byte of the decoded address is the witness version, it must
	// exist.
	if len(data) < 1 {
		fmt.Errorf("no witness version")
	}

	// ...and be <= 16.
	version := data[0]
	if version > 16 {
		fmt.Errorf("invalid witness version: %v", version)
	}

	// The remaining characters of the address returned are grouped into
	// words of 5 bits. In order to restore the original witness program
	// bytes, we'll need to regroup into 8 bit words.
	regrouped, err := ConvertBits(data[1:], 5, 8, false)
	if err != nil {
		fmt.Errorf("no witness version")
	}

	// The regrouped data must be between 2 and 40 bytes.
	if len(regrouped) < 2 || len(regrouped) > 40 {
		fmt.Errorf("no witness version")
	}

	// For witness version 0, address MUST be exactly 20 or 32 bytes.
	if version == 0 && len(regrouped) != 20 && len(regrouped) != 32 {
		fmt.Errorf("invalid data length for witness "+
			"version 0: %v", len(regrouped))
	}

	fmt.Printf("%x,%x\n",version,regrouped)

}