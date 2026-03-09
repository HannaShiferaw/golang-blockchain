package couchdb

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"coffee-consortium/backend/internal/domain"
)

type PKIRepo struct {
	c *Client
}

func NewPKIRepo(c *Client) *PKIRepo {
	return &PKIRepo{c: c}
}

func (r *PKIRepo) GetRoot(ctx context.Context) (certPEM, keyPEM string, found bool, err error) {
	var doc struct {
		CertPEM string `json:"certPem"`
		KeyPEM  string `json:"privateKeyPem"`
	}
	ok, _, err := r.c.Get(ctx, "pki:root", &doc)
	if err != nil || !ok {
		return "", "", ok, err
	}
	return doc.CertPEM, doc.KeyPEM, true, nil
}

func (r *PKIRepo) PutRoot(ctx context.Context, certPEM, keyPEM string) error {
	doc := map[string]any{
		"_id":          "pki:root",
		"certPem":      certPEM,
		"privateKeyPem": keyPEM,
	}
	// Merge _rev
	var existing map[string]any
	found, _, err := r.c.Get(ctx, "pki:root", &existing)
	if err != nil {
		return err
	}
	if found {
		if rev, ok := existing["_rev"].(string); ok {
			doc["_rev"] = rev
		}
	}
	_, err = r.c.Put(ctx, "pki:root", doc)
	return err
}

func (r *PKIRepo) PutIdentity(ctx context.Context, it domain.Identity) error {
	id := "pki:identity:" + it.ID

	doc := map[string]any{
		"_id":           id,
		"identityId":    it.ID,
		"name":          it.Name,
		"role":          string(it.Role),
		"certPem":       it.CertPEM,
		"privateKeyPem": it.PrivateKeyPEM,
	}

	var existing map[string]any
	found, _, err := r.c.Get(ctx, id, &existing)
	if err != nil {
		return err
	}
	if found {
		if rev, ok := existing["_rev"].(string); ok {
			doc["_rev"] = rev
		}
	}
	_, err = r.c.Put(ctx, id, doc)
	return err
}

func (r *PKIRepo) ListIdentities(ctx context.Context) ([]domain.Identity, error) {
	q := url.Values{}
	q.Set("include_docs", "true")
	q.Set("startkey", `"pki:identity:"`)
	q.Set("endkey", `"pki:identity:\ufff0"`)

	var res struct {
		Rows []struct {
			Doc json.RawMessage `json:"doc"`
		} `json:"rows"`
	}
	if err := r.c.AllDocs(ctx, q, &res); err != nil {
		return nil, err
	}
	out := make([]domain.Identity, 0, len(res.Rows))
	for _, row := range res.Rows {
		var doc struct {
			IdentityID    string `json:"identityId"`
			Name          string `json:"name"`
			Role          string `json:"role"`
			CertPEM       string `json:"certPem"`
			PrivateKeyPEM string `json:"privateKeyPem"`
		}
		if err := json.Unmarshal(row.Doc, &doc); err != nil {
			continue
		}
		role, err := domain.ParseRole(strings.ToUpper(doc.Role))
		if err != nil {
			continue
		}
		out = append(out, domain.Identity{
			ID:            doc.IdentityID,
			Name:          doc.Name,
			Role:          role,
			CertPEM:       doc.CertPEM,
			PrivateKeyPEM: doc.PrivateKeyPEM,
		})
	}
	return out, nil
}

