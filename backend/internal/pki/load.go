package pki

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func LoadRootCA(certPEM, privateKeyPEM string) (*RootCA, error) {
	certBlock, _ := pem.Decode([]byte(certPEM))
	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("invalid root cert pem")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse root cert: %w", err)
	}

	keyBlock, _ := pem.Decode([]byte(privateKeyPEM))
	if keyBlock == nil {
		return nil, fmt.Errorf("invalid root private key pem")
	}
	key, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse root private key: %w", err)
	}

	return &RootCA{Cert: cert, Key: key}, nil
}

