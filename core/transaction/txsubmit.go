package transaction

import (
	"encoding/json"

	"github.com/srchain/srcd/errors"
)

type TxSubmitResponse struct {
	TxID []byte `json:"tx_id"`
	Status string `json:"status"`
}
const (
	SUCCESS = "success"
	FAIL = "fail"
)


func (tp *TxPool)TxSubmit(raw_transaction string,height uint64)(TxSubmitResponse,error)  {

	var entity= struct {
		Tx Tx `json:"raw_transaction"`
	}{}

	err := json.Unmarshal([]byte(raw_transaction), &entity)
	if err != nil {}


	err = tp.AddTransaction(entity.Tx, height, 2)
	if err != nil {
		return TxSubmitResponse{nil,FAIL},errors.New("add tx to pool fail")
	}
	return TxSubmitResponse{entity.Tx.ID.Bytes(),SUCCESS},nil
}
