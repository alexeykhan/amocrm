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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout = 15 * time.Second
	protocol       = "https://"
)

type Client interface {
	AccountURL() string
	AuthorizeURL(state string) (string, error)

	AccessTokenByCode(code string) (*TokenSource, error)
	AccessTokenByRefreshToken(refreshToken string) (*TokenSource, error)

	SetRedirectMode(mode RedirectMode) Client
	SetDomain(domain Domain) Client
	SetToken(token Token) Client
}

// Verify interface compliance.
var _ Client = (*AuthClient)(nil)

type AuthClient struct {
	clientID     string
	clientSecret string
	redirectURL  string

	domain       Domain
	redirectMode RedirectMode

	tokenURL *string
	token    *TokenSource
}

func New(clientID, clientSecret, redirectURL string) *AuthClient {
	return &AuthClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		redirectMode: RedirectMode{
			value: PostMessageMode,
		},
		domain: Domain{
			account: "www",
			zone:    "ru",
		},
	}
}

func (a *AuthClient) SetRedirectMode(mode RedirectMode) Client {
	a.redirectMode = mode
	return a
}

func (a *AuthClient) SetDomain(domain Domain) Client {
	tokenURL := protocol + domain.String() + "/oauth2/access_token"
	a.tokenURL = &tokenURL
	a.domain = domain
	return a
}

func (a *AuthClient) SetToken(token Token) Client {
	a.token = &TokenSource{
		tokenType:    token.Type(),
		accessToken:  token.AccessToken(),
		refreshToken: token.RefreshToken(),
		expiresAt:    token.ExpiresAt(),
	}
	return a
}

func (a *AuthClient) AccountURL() string {
	return protocol + a.domain.String() + "/"
}

// AuthorizeURL returns a URL of consent page to ask for permissions.
func (a *AuthClient) AuthorizeURL(state string) (string, error) {
	if state == "" {
		return "", oauth2Err("state must not be empty")
	}

	query := url.Values{
		"state":     []string{state},
		"client_id": []string{a.clientID},
		"mode":      []string{a.redirectMode.value},
	}.Encode()

	authURL := protocol + a.domain.String() + "/oauth?" + query

	return authURL, nil
}

func (a *AuthClient) AccessTokenByCode(code string) (*TokenSource, error) {
	return a.accessToken(AuthorizationCodeGrant(), url.Values{
		"code":       []string{code},
		"grant_type": []string{"authorization_code"},
	})
}

func (a *AuthClient) AccessTokenByRefreshToken(refreshToken string) (*TokenSource, error) {
	return a.accessToken(RefreshTokenGrant(), url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{refreshToken},
	})
}

func (a *AuthClient) accessToken(grant GrantType, options url.Values) (*TokenSource, error) {
	if a.tokenURL == nil {
		return nil, oauth2Err("account domain is not set")
	}

	if err := verifyGrantParameters(grant, options); err != nil {
		return nil, oauth2Err("verify grant parameters: %w", err)
	}

	data := url.Values{
		"client_id":     []string{a.clientID},
		"client_secret": []string{a.clientSecret},
		"redirect_uri":  []string{a.redirectURL},
		"grant_type":    []string{grant.Name()},
	}

	for k, v := range options {
		if _, reserved := data[k]; !reserved {
			data[k] = v
		}
	}

	reqBody := strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", *a.tokenURL, reqBody)
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

	token := &TokenSource{
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

func verifyGrantParameters(grant GrantType, options url.Values) error {
	for _, key := range grant.Parameters() {
		if values, ok := options[key]; len(values) == 0 || !ok {
			return fmt.Errorf("missing required %s grant parameter %s", grant.Name(), key)
		}
	}

	return nil
}

func oauth2Err(format string, args ...interface{}) error {
	return fmt.Errorf("oauth2: "+format, args...)
}
