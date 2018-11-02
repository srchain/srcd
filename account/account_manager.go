package account

import (
	"srcd/database"
	"srcd/crypto/ed25519/chainkd"
	"crypto/rand"
	"fmt"
	"srcd/log"
)

var (
	IDPrefix = []byte("ACID")
)

type AccountManager struct {
	db database.Database
	accounts []Account
}

func NewAccountManager(db database.Database) *AccountManager {
	return &AccountManager{db: db}
}

func (am AccountManager) CreateAccount() (Account, error) {
	xpubs := []chainkd.XPub{}
	_, xpub, _ := chainkd.NewXKeys(rand.Reader)
	xpubs = append(xpubs, xpub)
	program, pubhash, err := CreateP2PKH(xpub)
	if err != nil {
		return Account{}, fmt.Errorf("create account fail:%x\n", err)
	}
	//ID-Pubhash
	keyID := append(IDPrefix, pubhash[:]...)
	am.db.Put(keyID, pubhash)
	//Pubhash-Account
	am.db.Put(pubhash, []byte(program.Address))
	return Account{"", program.Address}, nil
}

func (am AccountManager) GetCurrentNodeAccounts(id []byte)([]Account,error){
	accounts := []Account{}
	iter := am.db.NewIteratorWithPrefix(IDPrefix)
	defer iter.Release()
	for iter.Next() {
		acc := Account{}
		account_address, err := am.db.Get(iter.Value())
		if err != nil{
			log.Debug("Failed to iterator value")
		}
		acc.Address = string(account_address)
		accounts = append([]Account{acc}, accounts...)
	}
	return accounts,nil
}