package transaction

import (
	"encoding/json"

	"github.com/srchain/srcd/crypto/ed25519/chainkd"
)

// RawTxSigWitness is like SignatureWitness but doesn't involve
// signature programs.
type RawTxSigWitness struct {
	Quorum int                  `json:"quorum"`
	Keys   []keyID              `json:"keys"`
	Sigs   []HexBytes `json:"signatures"`
}
type keyID struct {
	XPub           chainkd.XPub         `json:"xpub"`
	DerivationPath []HexBytes `json:"derivation_path"`
}

func (sw RawTxSigWitness) materialize(args *[][]byte) error {
	var nsigs int
	for i := 0; i < len(sw.Sigs) && nsigs < sw.Quorum; i++ {
		if len(sw.Sigs[i]) > 0 {
			*args = append(*args, sw.Sigs[i])
			nsigs++
		}
	}
	return nil
}

// MarshalJSON convert struct to json
func (sw RawTxSigWitness) MarshalJSON() ([]byte, error) {
	obj := struct {
		Type   string               `json:"type"`
		Quorum int                  `json:"quorum"`
		Keys   []keyID              `json:"keys"`
		Sigs   []HexBytes `json:"signatures"`
	}{
		Type:   "raw_tx_signature",
		Quorum: sw.Quorum,
		Keys:   sw.Keys,
		Sigs:   sw.Sigs,
	}
	return json.Marshal(obj)
}
