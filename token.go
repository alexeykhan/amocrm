// The MIT License (MIT)
//
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
	"strings"
	"time"
)

// expiryDelta determines how earlier a token should be considered
// expired than its actual expiration time. It is used to avoid late
// expirations due to client-server time mismatches.
const expiryDelta = 10 * time.Second

// timeNow is time.Now but pulled out as a variable for tests.
var timeNow = time.Now

var (
	mac    = "MAC"
	bearer = "Bearer"
	basic  = "Basic"
)

// tokenJSON is the struct representing the HTTP response from OAuth2
// providers returning a token in JSON form.
type tokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
}

type TokenSource struct {
	accessToken  string
	refreshToken string
	tokenType    string
	expiresAt    time.Time
}

func NewToken(accessToken, refreshToken, tokenType string, expiresAt time.Time) *TokenSource {
	return &TokenSource{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		tokenType:    tokenType,
		expiresAt:    expiresAt,
	}
}

// AccessToken returns the token that authorizes and
// authenticates the requests.
func (t *TokenSource) AccessToken() string {
	if t == nil {
		return ""
	}

	return t.accessToken
}

// RefreshToken returns a token that's used by the application
// (as opposed to the user) to refresh the access token
// if it expires.
func (t *TokenSource) RefreshToken() string {
	if t == nil {
		return ""
	}

	return t.refreshToken
}

// ExpiresAt returns the optional expiration time of the access token.
//
// If zero, TokenSource implementations will reuse the same
// token forever and RefreshToken or equivalent
// mechanisms for that TokenSource will not be used.
func (t *TokenSource) ExpiresAt() time.Time {
	if t == nil {
		return time.Now().Add(-expiryDelta)
	}

	return t.expiresAt
}

// TokenType returns either this or "Bearer", the default.
func (t *TokenSource) Type() string {
	if t == nil {
		return ""
	}

	if strings.EqualFold(t.tokenType, "bearer") {
		return bearer
	}
	if strings.EqualFold(t.tokenType, "mac") {
		return mac
	}
	if strings.EqualFold(t.tokenType, "basic") {
		return basic
	}
	if t.tokenType != "" {
		return t.tokenType
	}
	return bearer
}

// Expired reports whether t has no AccessToken or is expired.
func (t *TokenSource) Expired() bool {
	if t.expiresAt.IsZero() {
		return false
	}

	if t.accessToken == "" {
		return true
	}

	return t.expiresAt.Round(0).Add(-expiryDelta).After(timeNow())
}
