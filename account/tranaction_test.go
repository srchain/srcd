package account

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/srchain/srcd/crypto/ed25519/chainkd"
	"github.com/srchain/srcd/core/transaction"
	"github.com/srchain/srcd/core/vm"
	"github.com/srchain/srcd/database"
	"github.com/srchain/srcd/trie"
)

func must(err error) {
	fmt.Println(err)
}

func mockAM() *AccountManager {

	diskdb, err := database.NewLDBDatabase("/Users/zhangrongxing/Downloads/srcd_wallet", 768, 1280)
	trie.NewDatabase(diskdb)
	if err != nil {
		panic(fmt.Sprintf("can't create temporary database: %v", err))
	}
	return NewAccountManager(diskdb)
}


func TestCreateAccount(t *testing.T) {
	am := mockAM()
	account, e := am.CreateAccount()
	must(e)
	fmt.Println(account)
}

func TestBuildUtxoTemplate(t *testing.T) {
	xpubs := []chainkd.XPub{}
	xprv, xpub, err := chainkd.NewXKeys(rand.Reader)
	//fmt.Println(err)
	xpubs = append(xpubs, xpub)
	fmt.Println(xprv, xpubs)
	program,_, err := CreateP2PKH(xpub)
	must(err)
	fmt.Printf("%x\n", program.Address)

	utxo := &transaction.UTXO{}
	utxo.SourceID = transaction.Hash{V0: 2}
	utxo.AssetID = *transaction.SRCAssetID
	utxo.Amount = 1000000000
	utxo.SourcePos = 0
	utxo.ControlProgram = program.ControlProgram
	utxo.Address = program.Address

	inst, err := transaction.UtxoInputs(xpubs, utxo)
	out := transaction.UtxoOutputs(*transaction.SRCAssetID, 100, []byte{byte(vm.OP_FAIL)})

	utxoTpl, _, err := transaction.BuildUtxoTemplate([]transaction.InputAndSigInst{inst}, []*transaction.TxOutput{&out})
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("%x\n",txData.Inputs)
	if b, err := json.Marshal(utxoTpl); err == nil {
		fmt.Println(string(b))
	}
}

func TestSignTx(t *testing.T) {
	xpubs := []chainkd.XPub{}
	xprv, xpub, err := chainkd.NewXKeys(rand.Reader)
	//fmt.Println(err)
	xpubs = append(xpubs, xpub)
	//fmt.Println(xprv, xpubs)
	program,_, err := CreateP2PKH(xpub)
	//must(err)
	//fmt.Printf("%x\n", program.Address)

	utxo := &transaction.UTXO{}
	utxo.SourceID = transaction.Hash{V0: 2}
	utxo.AssetID = *transaction.SRCAssetID
	utxo.Amount = 1000000000
	utxo.SourcePos = 0
	utxo.ControlProgram = program.ControlProgram
	utxo.Address = program.Address

	inst, err := transaction.UtxoInputs(xpubs, utxo)
	out := transaction.UtxoOutputs(*transaction.SRCAssetID, 100, []byte{byte(vm.OP_FAIL)})

	utxoTpl, _, err := transaction.BuildUtxoTemplate([]transaction.InputAndSigInst{inst}, []*transaction.TxOutput{&out})
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("%x\n",txData.Inputs)
	if b, err := json.Marshal(utxoTpl); err == nil {
		fmt.Println(string(b))
	}

	transaction.TxSign(utxoTpl, xprv, xpub)

	//fmt.Printf("%x\n", txData.Inputs)
	if b, err := json.Marshal(utxoTpl); err == nil {
		fmt.Println(string(b))
	}
}

var cliJsonStr = `
{
	"raw_transaction":"070100010161015f1af627ecf5cd04ebf00478325c0bdd0b26360e0ea7a03bfc0789e659388f8cfeffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8099c4d59901000116001401c3c624dc4d01dd6fed6bad8ed61c2b82183376630240006e535cbabfc9562ca5ee8998c7e3ca9e29813c9343255cc450661034041ea7687ec3c4f6f25ab0ec70e91c7ac7afdee1470eabe787f051e11705995061ff08205c4e1bf420d9cf31122f12f1a692dc864d1b72e1d9250f69f4957db877f40a8402013dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08ef6b474011600141525f1f92c3aa251ddd6c2ae5b844b94b683146000013dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c8afa02501160014f8a276a21e1f7a91c45ea7365f30281aefd77d9a00"
}`

func TestAddTransaction(t *testing.T) {
	tp := transaction.NewTxPool()

	//chain := blockchain.BlockChain{}

	//height := chain.CurrentBlock().Header().Number.Uint64()
	response, e := tp.TxSubmit(cliJsonStr,uint64(100))
	if e != nil {
		t.Fatal(e)
	}
	fmt.Printf("TxID:%x\n", response.TxID)
	//fmt.Println(tp)
}
