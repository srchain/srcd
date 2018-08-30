package wallet

import (
	"testing"
)

func TestNewWallets(t *testing.T) {
	wallets := NewWallets()
	wallets.CreateNewWallet()
}