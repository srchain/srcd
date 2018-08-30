package transaction

import (
	"srcd/rpc"
	"github.com/btcsuite/btcd/btcjson"
	"srcd/params/chainhash"
	"fmt"
)

// TransactionInput represents the inputs to a transaction.  Specifically a
// transaction hash and output number pair.
type TransactionInput struct {
	Txid string `json:"txid"`
	Vout uint32 `json:"vout"`
}

// CreateRawTransactionCmd defines the createrawtransaction.
type TxCmd struct {
	Inputs   []TransactionInput
	Amounts  map[string]float64
	LockTime *int64
}


// NewOutPoint returns a new src transaction outpoint point with the
// provided hash and index.
func NewOutPoint(hash *chainhash.Hash, index uint32) *OutPoint {
	return &OutPoint{
		Hash:  *hash,
		Index: index,
	}
}
// NewTxIn returns a new bitcoin transaction input with the provided
// previous outpoint point and signature script with a default sequence of
// MaxTxInSequenceNum.
func NewTxIn(prevOut *OutPoint, signatureScript []byte, witness [][]byte) *TxIn {
	return &TxIn{
		PreviousOutPoint: *prevOut,
		SignatureScript:  signatureScript,
		Witness:          witness,
		Sequence:         MaxTxInSequenceNum,
	}
}

func CreateRawTransaction(c *TxCmd)(interface{},error)  {

	// Validate the locktime, if given.
	if c.LockTime != nil &&
		(*c.LockTime < 0 || *c.LockTime > int64(MaxTxInSequenceNum)) {
		return nil, &rpc.RPCError{
			Code:    rpc.ErrRPCInvalidParameter,
			Message: "Locktime out of range",
		}
	}

	tx := NewMsgTx(TxVersion)
	for _, input := range c.Inputs {
		txHash, err := chainhash.NewHashFromStr(input.Txid)
		if err!=nil{
			return nil,rpcDecodeHexError(input.Txid)
		}

		prevOut := NewOutPoint(txHash, input.Vout)
		txIn := NewTxIn(prevOut, []byte{}, nil)
		if c.LockTime != nil && *c.LockTime != 0 {
			txIn.Sequence = MaxTxInSequenceNum - 1
		}
		tx.AddTxIn(txIn)
	}

	//for encodeAddr, amount:= range c.Amounts{
	//
	//	//Decode the provided address
	//
	//}

}


// rpcDecodeHexError is a convenience function for returning a nicely formatted
// RPC error which indicates the provided hex string failed to decode.
func rpcDecodeHexError(gotHex string) *btcjson.RPCError {
	return btcjson.NewRPCError(btcjson.ErrRPCDecodeHexString,
		fmt.Sprintf("Argument must be hexadecimal string (not %q)",
			gotHex))
}
