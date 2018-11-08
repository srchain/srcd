package wallet

import "github.com/srchain/srcd/database"

type Wallet struct {
	db         database.Database
	utxokeeper utxoKeeper
}
