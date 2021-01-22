// The MIT License (MIT)
//
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

package amocrm_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alexeykhan/amocrm"
)

var (
	clientID     = "client_id"
	clientSecret = "client_secret"
	redirectURL  = "redirect_url"
)

func TestNew(t *testing.T) {
	cl := amocrm.New(clientID, clientSecret, redirectURL)
	require.Implements(t, (*amocrm.Client)(nil), cl)
}

func TestRandomState(t *testing.T) {
	state := amocrm.RandomState()
	require.IsType(t, "", state)
	require.Equal(t, 32, len(state))
}

func TestAmoCRM_Accounts(t *testing.T) {
	cl := amocrm.New(clientID, clientSecret, redirectURL)
	require.Implements(t, (*amocrm.Accounts)(nil), cl.Accounts())
}

func TestAmoCRM_AuthorizeURL(t *testing.T) {
	cases := []struct {
		state string
		mode  string
		err   error
	}{
		{state: "", mode: "", err: errors.New("oauth2: empty state")},
		{state: "state", mode: "", err: errors.New("oauth2: unexpected mode")},
		{state: "state", mode: amocrm.PopupMode, err: nil},
	}

	cl := amocrm.New(clientID, clientSecret, redirectURL)

	for _, tc := range cases {
		_, err := cl.AuthorizeURL(tc.state, tc.mode)
		require.Exactly(t, err, tc.err)
	}
}

func TestAmoCRM_SetToken(t *testing.T) {
	cl := amocrm.New(clientID, clientSecret, redirectURL)
	require.EqualError(t, cl.SetToken(nil), "invalid token")

	token := amocrm.NewToken(accessToken, refreshToken, tokenType, time.Now())
	require.NoError(t, cl.SetToken(token))
}

func TestAmoCRM_SetDomain(t *testing.T) {
	cases := []struct {
		domain  string
		isValid bool
	}{
		{domain: "", isValid: false},
		{domain: "domain", isValid: false},
		{domain: "domain.com", isValid: false},
		{domain: ".domain.com", isValid: false},
		{domain: strings.Repeat("w", 64) + ".domain.com", isValid: false},
		{domain: "www.domain.com", isValid: false},
		{domain: "www.amocrm.any", isValid: false},
		{domain: "www.amocrm.ru", isValid: false},
		{domain: ".amocrm.ru", isValid: false},
		{domain: ".amocrm.", isValid: false},
		{domain: "any.amocrm.", isValid: false},
		{domain: "any.amocrm.ru", isValid: true},
		{domain: "any.amocrm.com", isValid: true},
	}

	cl := amocrm.New(clientID, clientSecret, redirectURL)

	for _, tc := range cases {
		if tc.isValid {
			require.NoError(t, cl.SetDomain(tc.domain))
		} else {
			require.EqualError(t, cl.SetDomain(tc.domain), "invalid domain")
		}
	}
}
