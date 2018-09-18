package protocol

import (
	"testing"
	"crypto/rand"
	"srcd/account"
	"fmt"
	"srcd/crypto/ed25519/chainkd"
	"srcd/protocol/transaction"
	vm2 "srcd/protocol/vm"
	"encoding/json"
	"github.com/bytom/util"
	"github.com/bytom/protocol/bc/types"
)

func TestBuildUtxoTemplate(t *testing.T) {
	xpubs := []chainkd.XPub{}
	xprv, xpub, err := chainkd.NewXKeys(rand.Reader)
	xpubs = append(xpubs, xpub)
	fmt.Println(xprv)
	program, err := account.CreateP2PKH(xpub)

	if err != nil {
		fmt.Printf("Newxkeys is error :%x\n", err)
	}

	utxo := &transaction.UTXO{}
	utxo.SourceID = transaction.Hash{V0: 2}
	utxo.AssetID = *transaction.SRCAssetID
	utxo.Amount = 1000000000
	utxo.SourcePos = 0
	utxo.ControlProgram = program.ControlProgram
	utxo.Address = program.Address

	inst, err := transaction.UtxoInputs(xpubs, utxo)
	out := transaction.UtxoOutputs(*transaction.SRCAssetID, 100, []byte{byte(vm2.OP_FAIL)})

	utxoTpl, _, err := transaction.BuildUtxoTemplate([]transaction.InputAndSigInst{inst}, []*transaction.TxOutput{&out})
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("%x\n",txData.Inputs)
	if b, err := json.Marshal(utxoTpl); err == nil {
		fmt.Println(string(b))
	}

}

func TestTxSign(t *testing.T) {
	xpubs := []chainkd.XPub{}
	xprv, xpub, err := chainkd.NewXKeys(rand.Reader);
	if err != nil {
		t.Fatal(err)
	}

	xpubs = append(xpubs, xpub)

	program, err := account.CreateP2PKH(xpub)

	if err != nil {
		t.Fatal(err)
	}

	utxo := &transaction.UTXO{}
	utxo.SourceID = transaction.Hash{V0: 2}
	utxo.AssetID = *transaction.SRCAssetID
	utxo.Amount = 1000000000
	utxo.SourcePos = 0
	utxo.ControlProgram = program.ControlProgram
	utxo.Address = program.Address

	inst, err := transaction.UtxoInputs(xpubs, utxo)
	out := transaction.UtxoOutputs(*transaction.SRCAssetID, 100, []byte{byte(vm2.OP_FAIL)})

	//build tx
	utxoTpl, txData, err := transaction.BuildUtxoTemplate([]transaction.InputAndSigInst{inst}, []*transaction.TxOutput{&out})
	if err != nil {
		t.Fatal(err)
	}
	if b, err := json.Marshal(utxoTpl); err == nil {
		fmt.Println(string(b))
	}
	//sign tx
	//070100010160015e1e673900965623ec3305cead5a78dfb68a34599f8bc078460f3f202256c3dfa6ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8090dfc04a0101160014f0a08c419d3a19688bf4e3fa7e60dbc518db3b03010001013dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c8afa025011600145056532ecd3621c9ce8adde5505c058610b287cf00
	//070100010160015e1e673900965623ec3305cead5a78dfb68a34599f8bc078460f3f202256c3dfa6ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8090dfc04a0101160014f0a08c419d3a19688bf4e3fa7e60dbc518db3b03630240124455b8f85128e38e4aeca142b146e5506c356e8c689257d3570f57e83d813fd6ddc337b776c940f8a2d178210a523b2ac77122e2b885ebe5037d5980837e06200347d22eb3f60ff9eeed4b05b887e438384c8cace0f0a3128cbe28114c82740301013dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c8afa025011600145056532ecd3621c9ce8adde5505c058610b287cf00
	transaction.TxSign(utxoTpl, xprv, xpub)

	fmt.Printf("%x\n", txData.Inputs)
	if b, err := json.Marshal(utxoTpl); err == nil {
		fmt.Println(string(b))
	}
}

var cliJsonStr = `
{
	"raw_transaction":"070100010161015f1af627ecf5cd04ebf00478325c0bdd0b26360e0ea7a03bfc0789e659388f8cfeffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8099c4d59901000116001401c3c624dc4d01dd6fed6bad8ed61c2b82183376630240006e535cbabfc9562ca5ee8998c7e3ca9e29813c9343255cc450661034041ea7687ec3c4f6f25ab0ec70e91c7ac7afdee1470eabe787f051e11705995061ff08205c4e1bf420d9cf31122f12f1a692dc864d1b72e1d9250f69f4957db877f40a8402013dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08ef6b474011600141525f1f92c3aa251ddd6c2ae5b844b94b683146000013dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c8afa02501160014f8a276a21e1f7a91c45ea7365f30281aefd77d9a00"
}`
var jsonStr = `
	{
		"raw_transaction":"070100010160015e1e673900965623ec3305cead5a78dfb68a34599f8bc078460f3f202256c3dfa6ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8090dfc04a0101160014f0a08c419d3a19688bf4e3fa7e60dbc518db3b03630240124455b8f85128e38e4aeca142b146e5506c356e8c689257d3570f57e83d813fd6ddc337b776c940f8a2d178210a523b2ac77122e2b885ebe5037d5980837e06200347d22eb3f60ff9eeed4b05b887e438384c8cace0f0a3128cbe28114c82740301013dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c8afa025011600145056532ecd3621c9ce8adde5505c058610b287cf00"
	}`

//func TestSubmitTransaction(t *testing.T){
//	response, e := transaction.TxSubmit(cliJsonStr)
//	if e != nil{
//		t.Fatal(e)
//	}
//	fmt.Printf("TxID:%x\n",response.TxID)
//}

func TestSubmitClientTransaction(t *testing.T) {
	var ins = struct {
		Tx types.Tx `json:"raw_transaction"`
	}{}

	err := json.Unmarshal([]byte(cliJsonStr), &ins)
	fmt.Printf("submitTransaction:%s\n", &ins)

	if err != nil {
	}
	data, _ := util.ClientCall("/submit-transaction", &ins)
	fmt.Println(data)
}
