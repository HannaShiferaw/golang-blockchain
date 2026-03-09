package pki

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"coffee-consortium/backend/internal/domain"
)

type RootCA struct {
	Cert *x509.Certificate
	Key  *ecdsa.PrivateKey
}

func NewRootCA(commonName string) (*RootCA, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate root key: %w", err)
	}

	serial, err := randSerial()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	tpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"Ethiopia Coffee Export Consortium"},
		},
		NotBefore:             now.Add(-5 * time.Minute),
		NotAfter:              now.AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}

	der, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &key.PublicKey, key)
	if err != nil {
		return nil, fmt.Errorf("create root cert: %w", err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, fmt.Errorf("parse root cert: %w", err)
	}

	return &RootCA{Cert: cert, Key: key}, nil
}

type IssuedIdentity struct {
	ID          domain.Identity
	Cert        *x509.Certificate
	PrivateKey  *ecdsa.PrivateKey
	CertPEM     string
	PrivatePEM  string
	Fingerprint string
}

func (ca *RootCA) IssueIdentity(id, name string, role domain.Role) (*IssuedIdentity, error) {
	if ca == nil || ca.Cert == nil || ca.Key == nil {
		return nil, fmt.Errorf("root CA not initialized")
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate identity key: %w", err)
	}

	serial, err := randSerial()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	tpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:         name,
			OrganizationalUnit: []string{string(role)},
			Organization:       []string{"Ethiopia Coffee Export Consortium"},
		},
		NotBefore:    now.Add(-5 * time.Minute),
		NotAfter:     now.AddDate(3, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		IsCA:         false,
		SubjectKeyId: []byte(id),
	}

	der, err := x509.CreateCertificate(rand.Reader, tpl, ca.Cert, &key.PublicKey, ca.Key)
	if err != nil {
		return nil, fmt.Errorf("issue identity cert: %w", err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, fmt.Errorf("parse identity cert: %w", err)
	}

	certPEM, err := certToPEM(der)
	if err != nil {
		return nil, err
	}
	privPEM, err := ecPrivToPEM(key)
	if err != nil {
		return nil, err
	}

	fp, err := fingerprintSHA256(cert.Raw)
	if err != nil {
		return nil, err
	}

	return &IssuedIdentity{
		ID: domain.Identity{
			ID:            id,
			Name:          name,
			Role:          role,
			CertPEM:       certPEM,
			PrivateKeyPEM: privPEM,
		},
		Cert:        cert,
		PrivateKey:  key,
		CertPEM:     certPEM,
		PrivatePEM:  privPEM,
		Fingerprint: fp,
	}, nil
}

func (ca *RootCA) CertPEM() (string, error) {
	if ca == nil || ca.Cert == nil {
		return "", fmt.Errorf("root CA not initialized")
	}
	return certToPEM(ca.Cert.Raw)
}

func (ca *RootCA) PrivateKeyPEM() (string, error) {
	if ca == nil || ca.Key == nil {
		return "", fmt.Errorf("root CA not initialized")
	}
	return ecPrivToPEM(ca.Key)
}

func certToPEM(der []byte) (string, error) {
	b := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	if b == nil {
		return "", fmt.Errorf("encode cert pem")
	}
	return string(b), nil
}

func ecPrivToPEM(k *ecdsa.PrivateKey) (string, error) {
	der, err := x509.MarshalECPrivateKey(k)
	if err != nil {
		return "", fmt.Errorf("marshal ec private key: %w", err)
	}
	b := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der})
	if b == nil {
		return "", fmt.Errorf("encode private key pem")
	}
	return string(b), nil
}

func randSerial() (*big.Int, error) {
	serialLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serial, err := rand.Int(rand.Reader, serialLimit)
	if err != nil {
		return nil, fmt.Errorf("serial: %w", err)
	}
	return serial, nil
}

