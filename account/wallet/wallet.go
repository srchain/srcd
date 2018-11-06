package wallet

import "srcd/database"

type Wallet struct {
	db         database.Database
	utxokeeper utxoKeeper
}
