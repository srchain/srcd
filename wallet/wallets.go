package wallet

import "fmt"

type  Wallets struct {
	Wallets map[string]*Wallet
}

//NewWallets create wallet list to store wallet address
func NewWallets() *Wallets  {
	wallets := &Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	return wallets
}

func (w *Wallets) CreateNewWallet()  {
	wallet := NewWallet()
	fmt.Printf("Address：%s\n",wallet.GetAddress())
	w.Wallets[string(wallet.GetAddress())] = wallet
}
