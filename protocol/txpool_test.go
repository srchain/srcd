package protocol

import (
	"testing"
	"srcd/protocol/transaction"
	"fmt"
	"srcd/core/blockchain"
)

func TestAddTransaction(t *testing.T) {
	tp := transaction.NewTxPool()

	chain := blockchain.BlockChain{}

	response, e := tp.TxSubmit(cliJsonStr,chain)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Printf("TxID:%x\n", response.TxID)

	fmt.Println(tp)

}