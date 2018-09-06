package transaction

import (
	vm2 "srcd/protocol/vm"
	"encoding/binary"
)

var SRCAssetID = &AssetID{
	V0: binary.BigEndian.Uint64([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
	V1: binary.BigEndian.Uint64([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
	V2: binary.BigEndian.Uint64([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
	V3: binary.BigEndian.Uint64([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
}


func MapTx(old TxData)TxWrap{
	txID, txHeader, entries := mapTx(old)
	tx := TxWrap{
		TxHeader: txHeader,
		ID:       txID,
		Entries:  entries,
		InputIDs: make([]Hash, len(old.Inputs)),
	}

	spentOutputIDs := make(map[Hash]bool)
	for id, e := range entries {
		var ord uint64
		switch e := e.(type) {

		case *Spend:
			ord = e.Ordinal
			spentOutputIDs[*e.SpentOutputId] = true
			if *e.WitnessDestination.Value.AssetId == *SRCAssetID {
				tx.GasInputIDs = append(tx.GasInputIDs, id)
			}

		default:
			continue
		}

		if ord >= uint64(len(old.Inputs)) {
			continue
		}
		tx.InputIDs[ord] = id
	}

	for id := range spentOutputIDs {
		tx.SpentOutputIDs = append(tx.SpentOutputIDs, id)
	}
	return tx
}

func mapTx(tx TxData) (headerID Hash, hdr *TxHeader, entryMap map[Hash]Entry) {
	entryMap = make(map[Hash]Entry)
	addEntry := func(e Entry) Hash {
		id := EntryID(e)
		entryMap[id] = e
		return id
	}

	var (
		spends    []*Spend
	)

	muxSources := make([]*ValueSource, len(tx.Inputs))
	for i, input := range tx.Inputs {
		switch inp := input.TypedInput.(type) {

		case *SpendInput:
			// create entry for prevout
			prog := &Program{VmVersion: inp.VMVersion, Code: inp.ControlProgram}
			src := &ValueSource{
				Ref:      &inp.SourceID,
				Value:    &inp.AssetAmount,
				Position: inp.SourcePosition,
			}
			prevout := NewOutput(src, prog, 0) // ordinal doesn't matter for prevouts, only for result outputs

			prevoutID := addEntry(prevout)
			// create entry for spend
			spend := NewSpend(&prevoutID, uint64(i))
			spend.WitnessArguments = inp.Arguments
			spendID := addEntry(spend)
			// setup mux
			muxSources[i] = &ValueSource{
				Ref:   &spendID,
				Value: &inp.AssetAmount,
			}
			spends = append(spends, spend)

		}
	}

	mux := NewMux(muxSources, &Program{VmVersion: 1, Code: []byte{byte(vm2.OP_TRUE)}})
	muxID := addEntry(mux)

	// connect the inputs to the mux
	for _, spend := range spends {
		spentOutput := entryMap[*spend.SpentOutputId].(*Output)
		spend.SetDestination(&muxID, spentOutput.Source.Value, spend.Ordinal)
	}

	// convert types.outputs to the bc.output
	var resultIDs []*Hash
	for i, out := range tx.Outputs {
		src := &ValueSource{
			Ref:      &muxID,
			Value:    &out.AssetAmount,
			Position: uint64(i),
		}
		var resultID Hash
		//if vmutil.IsUnspendable(out.ControlProgram) {
		//	// retirement
		//	r := bc.NewRetirement(src, uint64(i))
		//	resultID = addEntry(r)
		//} else {
		//
		//}
		// non-retirement
		prog := &Program{out.VMVersion, out.ControlProgram}
		o := NewOutput(src, prog, uint64(i))
		resultID = addEntry(o)

		dest := &ValueDestination{
			Value:    src.Value,
			Ref:      &resultID,
			Position: 0,
		}
		resultIDs = append(resultIDs, &resultID)
		mux.WitnessDestinations = append(mux.WitnessDestinations, dest)
	}

	h := NewTxHeader(tx.Version, tx.SerializedSize, tx.TimeRange, resultIDs)
	return addEntry(h), h, entryMap
}