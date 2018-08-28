package secp256k1

import (
	"testing"
	"crypto/rand"
	"log"
	"fmt"
)

func TestGenerateKey(t *testing.T) {
	priv, err := GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Generating keypair: %s\n", err)
	}
	fmt.Println(priv)
}

func TestPublicKey_Toddress(t *testing.T) {
	priv, err := GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Generating keypair: %s\n", err)
	}
	address := priv.Toddress()
	fmt.Println(address)
}