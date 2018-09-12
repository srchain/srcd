package accounts

import (
	"errors"
)

// ErrUnknownAccount is returned for any requested operation for which no backend
// provides the specified account.
var ErrUnknownAccount = errors.New("unknown account")
