package accounts

import (
	"srcd/common/common"
	"srcd/core/types"
)

// Wallet represents a software or hardware wallet that might contain one or more
// accounts (derived from the same seed).
type Wallet interface {
	// URL retrieves the canonical path under which this wallet is reachable. It is
	// user by upper layers to define a sorting order over all wallets from multiple
	// backends.
	// URL() URL

	// Status returns a textual status to aid the user in the current state of the
	// wallet. It also returns an error indicating any failure the wallet might have
	// encountered.
	Status() (string, error)

	// Open initializes access to a wallet instance. It is not meant to unlock or
	// decrypt account keys, rather simply to establish a connection to hardware
	// wallets and/or to access derivation seeds.
	//
	// The passphrase parameter may or may not be used by the implementation of a
	// particular wallet instance. The reason there is no passwordless open method
	// is to strive towards a uniform wallet handling, oblivious to the different
	// backend providers.
	//
	// Please note, if you open a wallet, you must close it to release any allocated
	// resources (especially important when working with hardware wallets).
	Open(passphrase string) error

	// Close releases any resources held by an open wallet instance.
	Close() error

	// Address retrieves the list of signing address the wallet is currently aware
	// of. For hierarchical deterministic wallets, the list will not be exhaustive,
	// rather only contain the accounts explicitly pinned during account derivation.
	Address() []common.Address

	// Contains returns whether an address is part of this particular wallet or not.
	Contains(address common.Address) bool

	// Derive attempts to explicitly derive a hierarchical deterministic account at
	// the specified derivation path. If requested, the derived account will be added
	// to the wallet's tracked account list.
	// Derive(path DerivationPath, pin bool) (Account, error)

	// SelfDerive sets a base account derivation path from which the wallet attempts
	// to discover non zero accounts and automatically add them to list of tracked
	// accounts.
	//
	// Note, self derivaton will increment the last component of the specified path
	// opposed to decending into a child path to allow discovering accounts starting
	// from non zero components.
	//
	// You can disable automatic account discovery by calling SelfDerive with a nil
	// chain state reader.
	// SelfDerive(base DerivationPath, chain ethereum.ChainStateReader)

	// SignHash requests the wallet to sign the given hash.
	SignHash(address common.Address, hash []byte) ([]byte, error)

	// SignTx requests the wallet to sign the given transaction.
	SignTx(address common.Address, tx *types.Transaction) (*types.Transaction, error)

	// SignHashWithPassphrase requests the wallet to sign the given hash with the
	// given passphrase as extra authentication information.
	SignHashWithPassphrase(address common.Address, passphrase string, hash []byte) ([]byte, error)

	// SignTxWithPassphrase requests the wallet to sign the given transaction, with the
	// given passphrase as extra authentication information.
	SignTxWithPassphrase(address common.Address, passphrase string, tx *types.Transaction) (*types.Transaction, error)
}

// Backend is a "wallet provider" that may contain a batch of accounts they can
// sign transactions with and upon request, do so.
type Backend interface {
	// Wallets retrieves the list of wallets the backend is currently aware of.
	// The returned wallets are not opened by default. For software HD wallets this
	// means that no base seeds are decrypted, and for hardware wallets that no actual
	// connection is established.
	Wallets() []Wallet

	// Subscribe creates an async subscription to receive notifications when the
	// backend detects the arrival or departure of a wallet.
	// Subscribe(sink chan<- WalletEvent) event.Subscription
}
