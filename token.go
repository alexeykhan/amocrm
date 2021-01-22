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

package amocrm

import (
	"strings"
	"time"
)

// GetToken stores a set of GetToken, RefreshToken and meta data.
type Token interface {
	AccessToken() string
	RefreshToken() string
	ExpiresAt() time.Time
	TokenType() string
	Expired() bool
}

// expiryDelta determines how earlier a token should be considered
// expired than its actual expiration time. It is used to avoid late
// expiration due to client-server time mismatches.
const expiryDelta = 10 * time.Second

// tokenJSON is the struct representing the HTTP response from OAuth2
// providers returning a token in JSON form.
type tokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
}

// tokenSource implements GetToken interface.
type tokenSource struct {
	accessToken  string
	refreshToken string
	tokenType    string
	expiresAt    time.Time
}

// Verify interface compliance.
var _ Token = tokenSource{}

// NewToken allocates and returns a new TokenSource.
func NewToken(accessToken, refreshToken, tokenType string, expiresAt time.Time) Token {
	return tokenSource{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		tokenType:    tokenType,
		expiresAt:    expiresAt,
	}
}

// GetToken returns the token that authorizes and
// authenticates the requests.
func (t tokenSource) AccessToken() string {
	return t.accessToken
}

// RefreshToken returns a token that's used by the application
// (as opposed to the user) to refresh the access token
// if it expires.
func (t tokenSource) RefreshToken() string {
	return t.refreshToken
}

// ExpiresAt returns the optional expiration time of the access token.
//
// If zero, TokenSource implementations will reuse the same token forever
// and RefreshToken or equivalent mechanisms for that TokenSource will
// not be used.
func (t tokenSource) ExpiresAt() time.Time {
	return t.expiresAt
}

// TokenType returns token type or "Bearer" by default.
func (t tokenSource) TokenType() string {
	switch {
	case strings.EqualFold(t.tokenType, "bearer"), t.tokenType == "":
		return "Bearer"
	case strings.EqualFold(t.tokenType, "mac"):
		return "MAC"
	case strings.EqualFold(t.tokenType, "basic"):
		return "Basic"
	default:
		return t.tokenType
	}
}

// Expired reports whether t has no GetToken or is expired.
func (t tokenSource) Expired() bool {
	if t.expiresAt.IsZero() {
		return false
	}

	if t.accessToken == "" {
		return true
	}

	return t.expiresAt.Round(0).Add(-expiryDelta).Before(time.Now())
}

