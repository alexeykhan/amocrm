// Copyright (c) 2020 Alexey Khan
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package amocrm

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"net/url"
)

// Client describes an interface for interacting with amoCRM API.
type Client interface {
	AuthorizeURL(state, mode string) (*url.URL, error)
	TokenByCode(code string) (*TokenSource, error)

	SetToken(token *TokenSource) error
	SetDomain(domain string) error

	Accounts() Account
}

// Verify interface compliance.
var _ Client = (*api)(nil)

// api implements Client interface.
type api struct {
	clientID     string
	clientSecret string
	redirectURL  string
	domain       string

	token *TokenSource
	http  *http.Client
}

// New allocates and returns a new amoCRM API Client.
func New(clientID, clientSecret, redirectURL string) Client {
	return &api{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		http: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

// RandomState generates a new random state.
func RandomState() string {
	// Converting bytes to hex will always double length. Hence, we can reduce
	// the amount of bytes by half to produce the correct length of 32 characters.
	key := make([]byte, 16)

	// https://golang.org/pkg/math/rand/#Rand.Read
	// Ignore errors as it always returns a nil error.
	_, _ = rand.Read(key)

	return hex.EncodeToString(key)
}

// AuthorizeURL returns a URL of page to ask for permissions.
func (a *api) AuthorizeURL(state, mode string) (*url.URL, error) {
	return a.authorizeURL(state, mode)
}

// SetToken stores given token to sign API requests.
func (a *api) SetToken(token *TokenSource) error {
	if token == nil {
		return errors.New("invalid token")
	}

	a.token = token
	return nil
}

// SetToken stores given domain to build account-specific API endpoints.
func (a *api) SetDomain(domain string) error {
	if !isValidDomain(domain) {
		return errors.New("invalid domain")
	}

	a.domain = domain
	return nil
}

// TokenByCode makes a handshake with amoCRM, exchanging given
// authorization code for a set of tokens.
func (a *api) TokenByCode(code string) (*TokenSource, error) {
	if code == "" {
		return nil, errors.New("empty authorization code")
	}

	return a.accessTokenByCode(code)
}

// Accounts returns an Account.
func (a *api) Accounts() Account {
	return account{
		api: a,
	}
}
