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

// type Client interface {
// 	Account() (*account.Resource, error)
// 	AuthorizeURL(state, mode string) (string, error)
// 	AccessTokenByCode(domain, code string) (oauth2.Token, error)
// }

// Verify interface compliance.
// var _ Client = (*AmoCRM)(nil)

var (
	ErrAccountDomainNotSet = errors.New("account domain is not set")
)

type AmoCRM struct {
	api *api
}

func New(clientID, clientSecret, redirectURL string) *AmoCRM {
	return &AmoCRM{
		api: &api{
			clientID:     clientID,
			clientSecret: clientSecret,
			redirectURL:  redirectURL,
			http: &http.Client{
				Timeout: RequestTimeout,
			},
		},
	}
}

func RandomState() string {
	// Converting bytes to hex will always double length. Hence, we can reduce
	// the amount of bytes by half to produce the correct length of 32 characters.
	key := make([]byte, 16)

	// https://golang.org/pkg/math/rand/#Rand.Read
	// Ignore errors as it always returns a nil error.
	_, _ = rand.Read(key)

	return hex.EncodeToString(key)
}

func (a *AmoCRM) AuthorizeURL(state, mode string) (*url.URL, error) {
	return a.api.authorizeURL(state, mode)
}

func (a *AmoCRM) SetToken(token *TokenSource) error {
	if token == nil {
		return errors.New("invalid token")
	}

	a.api.token = token
	return nil
}

func (a *AmoCRM) SetDomain(domain string) error {
	if !isValidDomain(domain) {
		return errors.New("invalid domain")
	}

	a.api.domain = domain
	return nil
}

func (a *AmoCRM) TokenByCode(code string) (*TokenSource, error) {
	if code == "" {
		return nil, errors.New("empty authorization code")
	}

	return a.api.accessTokenByCode(code)
}

func (a *AmoCRM) Account() AccountResource {
	return account{
		api: a.api,
	}
}
