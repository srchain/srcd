package crypto


import (
	"testing"
	"srcd/crypto/ed25519/chainkd"
	"crypto/rand"
	"fmt"
	"srcd/account"
)


func TestNewXKeys(t *testing.T){
	xprv, xpub, err := chainkd.NewXKeys(rand.Reader)

	program, err := account.CreateP2PKH(xpub)

	fmt.Println(program.Address)
	fmt.Printf("%x\n",program.ControlProgram)

	if err != nil{
		fmt.Printf("Newxkeys is error :%x\n",err)
	}
	fmt.Printf("%x\n,%x\n",xprv,xpub)
}
