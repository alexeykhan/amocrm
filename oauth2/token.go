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

package oauth2

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

type TokenSource interface {
	AccessToken() string
	RefreshToken() string
	ExpiresAt() time.Time
	Type() string
	Valid() bool
}

type Token struct {
	accessToken  string
	refreshToken string
	tokenType    string
	expiresAt    time.Time
}

// AccessToken returns the token that authorizes and
// authenticates the requests.
func (t Token) AccessToken() string {
	return t.accessToken
}

// RefreshToken returns a token that's used by the application
// (as opposed to the user) to refresh the access token
// if it expires.
func (t Token) RefreshToken() string {
	return t.refreshToken
}

// ExpiresAt returns the optional expiration time of the access token.
//
// If zero, TokenSource implementations will reuse the same
// token forever and RefreshToken or equivalent
// mechanisms for that TokenSource will not be used.
func (t Token) ExpiresAt() time.Time {
	return t.expiresAt
}

// TokenType returns either this or "Bearer", the default.
func (t Token) Type() string {
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

// Valid reports whether t has an AccessToken, and is not expired.
func (t Token) Valid() bool {
	return t.accessToken != "" && !t.expired()
}

func (t Token) expired() bool {
	if t.expiresAt.IsZero() {
		return false
	}

	return t.expiresAt.Round(0).Add(-expiryDelta).Before(timeNow())
}
