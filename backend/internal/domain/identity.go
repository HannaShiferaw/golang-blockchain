package domain

type Identity struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Role    Role   `json:"role"`
	CertPEM string `json:"certPem"`
	// For demo simplicity we store the private key server-side.
	// In production this would live in an HSM / client wallet.
	PrivateKeyPEM string `json:"privateKeyPem,omitempty"`
}

