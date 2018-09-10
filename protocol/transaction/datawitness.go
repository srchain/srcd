package transaction


import (
"encoding/json"
)

// DataWitness used sign transaction
type DataWitness HexBytes


func (dw DataWitness) materialize(args *[][]byte) error {
	*args = append(*args, dw)
	return nil
}

// MarshalJSON marshal DataWitness
func (dw DataWitness) MarshalJSON() ([]byte, error) {
	x := struct {
		Type  string             `json:"type"`
		Value HexBytes `json:"value"`
	}{
		Type:  "data",
		Value: HexBytes(dw),
	}
	return json.Marshal(x)
}

