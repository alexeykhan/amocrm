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
	"errors"

	"github.com/alexeykhan/amocrm/oauth2"
	"github.com/alexeykhan/amocrm/repository/account"
)

type Client interface {
	Account() (*account.Resource, error)
	SetToken(token oauth2.Token) Client
}

// Verify interface compliance.
var _ Client = (*API)(nil)

type API struct {
	auth  *oauth2.AuthClient
	token *oauth2.TokenSource
}

func New(clientID, clientSecret, redirectURL string) *API {
	return &API{
		auth: oauth2.New(clientID, clientSecret, redirectURL),
	}
}

func (a *API) AuthorizeURL(state, mode string) (string, error) {
	redirectMode, err := oauth2.NewRedirectMode(mode)
	if err != nil {
		return "", err
	}

	a.auth.SetRedirectMode(*redirectMode)

	return a.auth.AuthorizeURL(state)
}

func (a *API) AccessTokenByCode(domain, code string) (oauth2.Token, error) {
	accountDomain, err := oauth2.NewDomain(domain)
	if err != nil {
		return nil, err
	}

	a.auth.SetDomain(*accountDomain)

	return a.auth.AccessTokenByCode(code)
}

func (a *API) SetToken(token oauth2.Token) Client {
	a.auth.SetToken(token)
}

func (a *API) Account() (*account.Resource, error) {
	return nil, nil
}

func (a *API) refreshAccessToken() error {
	if a.token == nil {
		return errors.New("token is not set")
	}
	if a.token.RefreshToken() == "" {
		return errors.New("empty refresh token")
	}

	token, err := a.auth.AccessTokenByRefreshToken(a.token.RefreshToken())
	if err != nil {
		return err
	}

	a.SetToken(token)
	return nil
}
