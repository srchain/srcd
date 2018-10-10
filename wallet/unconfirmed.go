package wallet

import (
	"srcd/protocol/transaction"
	"srcd/log"
	"github.com/gin-gonic/gin/json"
)

type Utxo struct {
	OutputID            transaction.Hash
	SourceID            transaction.Hash
	AssetID             transaction.AssetID
	Amount              uint64
	SourcePos           uint64
	ControlProgram      []byte
	AccountID           string
	Address             string
	ControlProgramIndex uint64
	Change              bool
} 

// AddUnconfirmedTx handle wallet status update when tx add into txpool
func (w *Wallet) AddUnconfirmedTx(tx *transaction.TxPoolMsg) {
	//db
	if err := w.saveUnconfirmedTx(tx.Tx); err != nil {
		log.Error("err",err," fail on saveUnconfirmedTx ")
	}
	//buffer
	utxos := txOutToUtxos(tx.Tx)
	w.utxokeeper.AddUnconfirmedTx(utxos)
}

func (w *Wallet) saveUnconfirmedTx(tx transaction.Tx) error {
	bytes, e := json.Marshal(tx)
	if e != nil{
		return e
	}
	w.db.Put([]byte("UTXS:"+tx.ID.String()),bytes)
}

func txOutToUtxos(tx transaction.Tx)[]*Utxo{
	utxos := []*Utxo{}
	for i, out := range tx.Outputs {
		bcOut, err := tx.Output(*tx.ResultIds[i])
		if err != nil {
			continue
		}

		utxos = append(utxos, &Utxo{
			OutputID:       *tx.OutputID(i),
			AssetID:        *out.AssetAmount.AssetId,
			Amount:         out.Amount,
			ControlProgram: out.ControlProgram,
			SourceID:       *bcOut.Source.Ref,
			SourcePos:      bcOut.Source.Position,
		})
	}
	return utxos
}
