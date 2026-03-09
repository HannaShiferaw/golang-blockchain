package pki

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
)

type ecdsaSig struct {
	R, S *big.Int
}

func SignPayload(privateKeyPEM string, payload []byte) (string, error) {
	k, err := parseECPrivateKeyPEM(privateKeyPEM)
	if err != nil {
		return "", err
	}

	h := sha256.Sum256(payload)
	r, s, err := ecdsa.Sign(rand.Reader, k, h[:])
	if err != nil {
		return "", fmt.Errorf("sign: %w", err)
	}

	der, err := asn1.Marshal(ecdsaSig{R: r, S: s})
	if err != nil {
		return "", fmt.Errorf("marshal sig: %w", err)
	}
	return base64.StdEncoding.EncodeToString(der), nil
}

func VerifyPayload(certPEM string, payload []byte, sigB64 string) error {
	cert, err := parseCertPEM(certPEM)
	if err != nil {
		return err
	}
	pub, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("cert public key is not ecdsa")
	}

	sigDER, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return fmt.Errorf("decode sig: %w", err)
	}
	var sig ecdsaSig
	if _, err := asn1.Unmarshal(sigDER, &sig); err != nil {
		return fmt.Errorf("unmarshal sig: %w", err)
	}

	h := sha256.Sum256(payload)
	if !ecdsa.Verify(pub, h[:], sig.R, sig.S) {
		return fmt.Errorf("invalid signature")
	}
	return nil
}

func fingerprintSHA256(derCert []byte) (string, error) {
	if len(derCert) == 0 {
		return "", fmt.Errorf("empty cert")
	}
	sum := sha256.Sum256(derCert)
	return hex.EncodeToString(sum[:]), nil
}

func parseCertPEM(certPEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("invalid cert pem")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse cert: %w", err)
	}
	return cert, nil
}

func parseECPrivateKeyPEM(privateKeyPEM string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("invalid private key pem")
	}
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse ec private key: %w", err)
	}
	return key, nil
}

