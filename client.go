// Copyright (c) 2021 Alexey Khan
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
	"net/url"
)

// Provider is a wrapper for authorization and making requests.
type Client interface {
	AuthorizeURL(state, mode string) (*url.URL, error)
	TokenByCode(code string) (Token, error)
	SetToken(token Token) error
	SetDomain(domain string) error
	Accounts() Accounts
}

// Verify interface compliance.
var _ Client = (*amoCRM)(nil)

type amoCRM struct {
	api *api
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

// New allocates and returns a new amoCRM API Client.
func New(clientID, clientSecret, redirectURL string) Client {
	return &amoCRM{
		api: newAPI(clientID, clientSecret, redirectURL),
	}
}

// AuthorizeURL returns a URL of page to ask for permissions.
func (a *amoCRM) AuthorizeURL(state, mode string) (*url.URL, error) {
	if state == "" {
		return nil, oauth2Err("empty state")
	}
	if mode != PostMessageMode && mode != PopupMode {
		return nil, oauth2Err("unexpected mode")
	}

	query := url.Values{
		"mode":      []string{mode},
		"state":     []string{state},
		"client_id": []string{a.api.clientID},
	}.Encode()

	authURL := "https://www.amocrm.ru/oauth?" + query

	return url.Parse(authURL)
}

// SetToken stores given token to sign API requests.
func (a *amoCRM) SetToken(token Token) error {
	return a.api.setToken(token)
}

// SetToken stores given domain to build accounts-specific API endpoints.
func (a *amoCRM) SetDomain(domain string) error {
	return a.api.setDomain(domain)
}

// TokenByCode makes a handshake with amoCRM, exchanging given
// authorization code for a set of tokens.
func (a *amoCRM) TokenByCode(code string) (Token, error) {
	return a.api.getToken(authorizationCodeGrant, url.Values{
		"code":       []string{code},
		"grant_type": []string{"authorization_code"},
	}, nil)
}

// Accounts returns accounts repository.
func (a *amoCRM) Accounts() Accounts {
	return newAccounts(a.api)
}
