package couchdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	base *url.URL
	user string
	pass string
	db   string
	hc   *http.Client
}

func New(baseURL, user, pass, db string) (*Client, error) {
	u, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return nil, fmt.Errorf("parse couchdb url: %w", err)
	}
	if db == "" {
		return nil, fmt.Errorf("missing couchdb db name")
	}
	return &Client{
		base: u,
		user: user,
		pass: pass,
		db:   db,
		hc: &http.Client{
			Timeout: 15 * time.Second,
		},
	}, nil
}

func (c *Client) EnsureDB(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.dbURL().String(), nil)
	if err != nil {
		return err
	}
	c.setAuth(req)
	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("couchdb ensure db: %w", err)
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 201, 202:
		return nil
	case 412: // already exists
		return nil
	default:
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return fmt.Errorf("couchdb ensure db status=%d body=%s", resp.StatusCode, string(b))
	}
}

func (c *Client) Get(ctx context.Context, id string, out any) (found bool, rev string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.docURL(id).String(), nil)
	if err != nil {
		return false, "", err
	}
	c.setAuth(req)
	resp, err := c.hc.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("couchdb get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return false, "", nil
	}
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return false, "", fmt.Errorf("couchdb get status=%d body=%s", resp.StatusCode, string(b))
	}
	var m map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return false, "", fmt.Errorf("decode couchdb doc: %w", err)
	}
	if r, ok := m["_rev"].(string); ok {
		rev = r
	}
	// Re-marshal into out to preserve json semantics and allow struct targets.
	b, _ := json.Marshal(m)
	if err := json.Unmarshal(b, out); err != nil {
		return false, "", fmt.Errorf("unmarshal couchdb doc: %w", err)
	}
	return true, rev, nil
}

func (c *Client) Put(ctx context.Context, id string, doc any) (newRev string, err error) {
	body, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.docURL(id).String(), bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	c.setAuth(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.hc.Do(req)
	if err != nil {
		return "", fmt.Errorf("couchdb put: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 && resp.StatusCode != 202 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return "", fmt.Errorf("couchdb put status=%d body=%s", resp.StatusCode, string(b))
	}
	var res struct {
		Ok  bool   `json:"ok"`
		Rev string `json:"rev"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	return res.Rev, nil
}

func (c *Client) AllDocs(ctx context.Context, query url.Values, out any) error {
	u := c.dbURL()
	u.Path = strings.TrimRight(u.Path, "/") + "/_all_docs"
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	c.setAuth(req)
	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("couchdb all_docs: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return fmt.Errorf("couchdb all_docs status=%d body=%s", resp.StatusCode, string(b))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) setAuth(req *http.Request) {
	if c.user != "" || c.pass != "" {
		req.SetBasicAuth(c.user, c.pass)
	}
}

func (c *Client) dbURL() *url.URL {
	u := *c.base
	u.Path = strings.TrimRight(u.Path, "/") + "/" + c.db
	return &u
}

func (c *Client) docURL(id string) *url.URL {
	u := c.dbURL()
	u.Path = strings.TrimRight(u.Path, "/") + "/" + url.PathEscape(id)
	return u
}

