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
	"net/url"

	"github.com/alexeykhan/amocrm/api"
	"github.com/alexeykhan/amocrm/api/accounts"
)

// Provider describes an interface for interacting with amoCRM API.
type Provider interface {
	AuthorizeURL(state, mode string) (*url.URL, error)
	TokenByCode(code string) (api.Token, error)
	SetToken(token api.Token) error
	SetDomain(domain string) error
	Accounts() accounts.Repository
}

// Verify interface compliance.
var _ Provider = (*amoCRM)(nil)

type amoCRM struct {
	api api.Client
}

// New allocates and returns a new amoCRM API Client.
func New(clientID, clientSecret, redirectURL string) Provider {
	return &amoCRM{
		api: api.New(clientID, clientSecret, redirectURL),
	}
}

// AuthorizeURL returns a URL of page to ask for permissions.
func (a *amoCRM) AuthorizeURL(state, mode string) (*url.URL, error) {
	return a.api.AuthorizationURL(state, mode)
}

// SetToken stores given token to sign API requests.
func (a *amoCRM) SetToken(token api.Token) error {
	return a.api.SetToken(token)
}

// SetToken stores given domain to build accounts-specific API endpoints.
func (a *amoCRM) SetDomain(domain string) error {
	return a.api.SetDomain(domain)
}

// TokenByCode makes a handshake with amoCRM, exchanging given
// authorization code for a set of tokens.
func (a *amoCRM) TokenByCode(code string) (api.Token, error) {
	return a.api.GetToken(api.AuthorizationCodeGrant, url.Values{
		"code":       []string{code},
		"grant_type": []string{"authorization_code"},
	}, nil)
}

// Accounts returns accounts repository.
func (a *amoCRM) Accounts() accounts.Repository {
	return accounts.New(a.api)
}
