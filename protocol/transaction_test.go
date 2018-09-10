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
)

func TestBuildUtxoTemplate(t *testing.T)  {
	xpubs := []chainkd.XPub{}
	xprv, xpub, err := chainkd.NewXKeys(rand.Reader)
	xpubs = append(xpubs, xpub)
	fmt.Println(xprv)

	program, err := account.CreateP2PKH(xpub)

	fmt.Println(program.Address)

	if err != nil{
		fmt.Printf("Newxkeys is error :%x\n",err)
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

	utxoTpl, txData , err := transaction.BuildUtxoTemplate([]transaction.InputAndSigInst{inst}, []*transaction.TxOutput{&out})
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(utxoTpl)

	fmt.Printf("%x\n",txData.Inputs)
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

	//sign tx
	//070100010160015e1e673900965623ec3305cead5a78dfb68a34599f8bc078460f3f202256c3dfa6ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8090dfc04a0101160014f0a08c419d3a19688bf4e3fa7e60dbc518db3b03010001013dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c8afa025011600145056532ecd3621c9ce8adde5505c058610b287cf00
	//070100010160015e1e673900965623ec3305cead5a78dfb68a34599f8bc078460f3f202256c3dfa6ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8090dfc04a0101160014f0a08c419d3a19688bf4e3fa7e60dbc518db3b03630240124455b8f85128e38e4aeca142b146e5506c356e8c689257d3570f57e83d813fd6ddc337b776c940f8a2d178210a523b2ac77122e2b885ebe5037d5980837e06200347d22eb3f60ff9eeed4b05b887e438384c8cace0f0a3128cbe28114c82740301013dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c8afa025011600145056532ecd3621c9ce8adde5505c058610b287cf00
	transaction.TxSign(utxoTpl, xprv, xpub)

	fmt.Printf("%x\n", txData.Inputs)
	if b, err := json.Marshal(utxoTpl); err == nil {
		fmt.Println("================struct 到json str==")
		fmt.Println(string(b))
	}

}