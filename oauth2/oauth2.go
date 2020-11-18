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
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
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
	AccessTokenByCode(code string) (*Token, error)

	SetDomain(domain string) Client
}

// Verify interface compliance.
var _ Client = (*AmoCRMClient)(nil)

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

func (c *AmoCRMClient) SetDomain(domain string) Client {
	tokenURL := protocol + domain + "/oauth2/access_token"
	c.tokenURL = &tokenURL
	return c
}

// AuthorizeURL returns a URL of consent page to ask for permissions.
func (c *AmoCRMClient) AuthorizeURL(state string, mode RedirectMode) (string, error) {
	if state == "" {
		return "", oauth2Err("state must not be empty")
	}

	data := url.Values{
		"state":     []string{state},
		"client_id": []string{c.clientID},
		"mode":      []string{mode.String()},
	}.Encode()

	var buf bytes.Buffer
	buf.WriteString("https://www.amocrm.ru/oauth?")
	buf.WriteString(data)
	return buf.String(), nil
}

func (c *AmoCRMClient) AccessTokenByCode(code string) (*Token, error) {
	if c.tokenURL == nil {
		return nil, oauth2Err("account domain is not set")
	}

	data := url.Values{
		"code":          []string{code},
		"client_id":     []string{c.clientID},
		"client_secret": []string{c.clientSecret},
		"redirect_uri":  []string{c.redirectURL},
		"grant_type":    []string{"authorization_code"},
	}

	reqBody := strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", *c.tokenURL, reqBody)
	if err != nil {
		return nil, oauth2Err("build request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{
		Timeout: defaultTimeout,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, oauth2Err("send request")
	}

	respBody, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if closeBodyErr := resp.Body.Close(); closeBodyErr != nil {
		return nil, oauth2Err("close response body")
	}
	if err != nil {
		return nil, oauth2Err("fetch response body")
	}

	if statusCode := resp.StatusCode; statusCode < 200 || statusCode > 299 {
		return nil, oauth2Err("fetch token: response: %v - %s", resp.Status, respBody)
	}

	var jsonToken tokenJSON
	if err = json.Unmarshal(respBody, &jsonToken); err != nil {
		return nil, oauth2Err("parse token from json")
	}

	token := &Token{
		accessToken:  jsonToken.AccessToken,
		tokenType:    jsonToken.TokenType,
		refreshToken: jsonToken.RefreshToken,
		expiresAt:    time.Now().Add(time.Duration(jsonToken.ExpiresIn) * time.Second),
	}

	if token.accessToken == "" {
		return nil, oauth2Err("server response missing access_token")
	}

	return token, nil
}

func oauth2Err(format string, args ...interface{}) error {
	return fmt.Errorf("oauth2: "+format, args...)
}
