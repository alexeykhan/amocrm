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
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

type RedirectMode string

func (m RedirectMode) String() string {
	return string(m)
}

const (
	PostMessageMode RedirectMode = "post_message"
	PopupMode       RedirectMode = "popup"
)

const (
	defaultTimeout = 15 * time.Second
	protocol       = "https://"
)

var (
	ErrEmptyState          = errors.New("state must not be empty")
	ErrAccountDomainNotSet = errors.New("account domain is not set")
)

type RequestError struct {
	response *http.Response
	body     []byte
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("oauth2: cannot fetch token: %v", r.response.Status)
}

func (r *RequestError) Response() *http.Response {
	return r.response
}

func (r *RequestError) Body() []byte {
	return r.body
}

// tokenJSON is the struct representing the HTTP response from OAuth2
// providers returning a token in JSON form.
type tokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
}

func GenerateState() string {
	// Converting bytes to hex will always double length. Hence, we can reduce
	// the amount of bytes by half to produce the correct length of 32 characters.
	key := make([]byte, 16)

	// https://golang.org/pkg/math/rand/#Rand.Read
	// Ignore errors as it always returns a nil error.
	_, _ = rand.Read(key)

	return hex.EncodeToString(key)
}

type Client interface {
	AuthorizeURL(state string, mode RedirectMode) (string, error)
	AccessTokenByCode(ctx context.Context, code string) (*Token, error)

	SetAccountDomain(domain string) Client
}

type AmoCRMClient struct {
	clientID     string
	clientSecret string
	redirectURL  string

	tokenURL *string
}

func New(clientID, clientSecret, redirectURL string) Client {
	return &AmoCRMClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
	}
}

func (c *AmoCRMClient) SetAccountDomain(domain string) Client {
	tokenURL := protocol + domain + "/oauth2/access_token"
	c.tokenURL = &tokenURL
	return c
}

// AuthorizeURL returns a URL of consent page to ask for permissions.
func (c *AmoCRMClient) AuthorizeURL(state string, mode RedirectMode) (string, error) {
	if state == "" {
		return "", ErrEmptyState
	}

	v := url.Values{
		"state":     {state},
		"client_id": {c.clientID},
		"mode":      {mode.String()},
	}

	var buf bytes.Buffer
	buf.WriteString("https://www.amocrm.ru/oauth?")
	buf.WriteString(v.Encode())
	return buf.String(), nil
}

func (c *AmoCRMClient) AccessTokenByCode(ctx context.Context, code string) (*Token, error) {
	if c.tokenURL == nil {
		return nil, ErrAccountDomainNotSet
	}

	v := url.Values{
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"redirect_uri":  {c.redirectURL},
		"grant_type":    {"authorization_code"},
		"code":          {code},
	}

	req, err := http.NewRequest("POST", *c.tokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{Timeout: defaultTimeout}

	r, err := ctxhttp.Do(ctx, httpClient, req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	if closeBodyErr := r.Body.Close(); closeBodyErr != nil {
		return nil, closeBodyErr
	}

	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}

	if statusCode := r.StatusCode; statusCode < 200 || statusCode > 299 {
		return nil, &RequestError{
			response: r,
			body:     body,
		}
	}

	var tj tokenJSON
	if err = json.Unmarshal(body, &tj); err != nil {
		return nil, err
	}

	token := &Token{
		accessToken:  tj.AccessToken,
		tokenType:    tj.TokenType,
		refreshToken: tj.RefreshToken,
		expiresAt:    time.Now().Add(time.Duration(tj.ExpiresIn) * time.Second),
	}

	if token.accessToken == "" {
		return nil, errors.New("oauth2: server response missing access_token")
	}

	return token, nil
}
