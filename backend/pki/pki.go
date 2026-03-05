package pki

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

type KeyPair struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

// Generate RSA keypair
func GenerateKeyPair() (*KeyPair, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &KeyPair{PublicKey: &privKey.PublicKey, PrivateKey: privKey}, nil
}

// Sign data
func Sign(privKey *rsa.PrivateKey, data string) (string, error) {
	hash := sha256.Sum256([]byte(data))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// Verify signature
func Verify(pubKey *rsa.PublicKey, data string, signature string) bool {
	sigBytes, _ := base64.StdEncoding.DecodeString(signature)
	hash := sha256.Sum256([]byte(data))
	err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], sigBytes)
	return err == nil
}