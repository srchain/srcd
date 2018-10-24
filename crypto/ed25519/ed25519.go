package ed25519


const (
	// SignatureSize is the size, in bytes, of signatures generated and verified by this package.
	SignatureSize = 64
)

// PublicKey is the type of Ed25519 public keys.
type PublicKey []byte