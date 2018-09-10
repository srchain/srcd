package transaction

type SigningInstruction struct {
	Position          uint32             `json:"position"`
	WitnessComponents []witnessComponent `json:"witness_components,omitempty"`
}
type witnessComponent interface {
	materialize(*[][]byte) error
}

type RawWitness struct {
	Quorum int                  `json:"quorum"`
	Sigs   string `json:"signatures"`
}

func (sw RawWitness) materialize(args *[][]byte) error {
	return nil
}

