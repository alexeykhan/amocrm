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
	PostMessageMode = "post_message"
	PopupMode       = "popup"
)

type api struct {
	clientID     string
	clientSecret string
	redirectURL  string
	domain       string

	token *TokenSource
	http  *http.Client
}

// authorizeURL returns a URL of page to ask for permissions.
func (a *api) authorizeURL(state, mode string) (*url.URL, error) {
	if state == "" {
		return nil, oauth2Err("state must not be empty")
	}

	if !isValidMode(mode) {
		return nil, oauth2Err("unexpected mode")
	}

	query := url.Values{
		"mode":      []string{mode},
		"state":     []string{state},
		"client_id": []string{a.clientID},
	}.Encode()

	authURL := "https://www.amocrm.ru/oauth?" + query

	return url.Parse(authURL)
}

func (a *api) accessTokenByCode(code string) (*TokenSource, error) {
	return a.accessToken(authorizationCodeGrant(), url.Values{
		"code":       []string{code},
		"grant_type": []string{"authorization_code"},
	}, nil)
}

func (a *api) accessTokenByRefreshToken(refreshToken string) (*TokenSource, error) {
	return a.accessToken(refreshTokenGrant(), url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{refreshToken},
	}, nil)
}

func (a *api) accessToken(grant grantType, options url.Values, header http.Header) (*TokenSource, error) {
	if err := a.domainIsSet(); err != nil {
		return nil, err
	}

	// Validate required grantType-specific fields
	for _, key := range grant.fields {
		if values, ok := options[key]; len(values) == 0 || !ok {
			return nil, oauth2Err("missing required %s grant parameter %s", grant.code, key)
		}
	}

	// Default request parameters
	data := url.Values{
		"client_id":     []string{a.clientID},
		"client_secret": []string{a.clientSecret},
		"redirect_uri":  []string{a.redirectURL},
		"grant_type":    []string{grant.code},
	}

	// Merge options with default parameters
	for k, v := range options {
		if _, reserved := data[k]; !reserved {
			data[k] = v
		}
	}

	// Set request URL
	tokenURL, err := a.url("/oauth2/access_token", nil)
	if err != nil {
		return nil, oauth2Err("build request url")
	}

	// Set request headers
	reqHeader := a.baseHeader()
	reqHeader["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	for k, v := range header {
		if _, reserved := reqHeader[k]; !reserved {
			reqHeader[k] = v
		}
	}

	// Create request body
	reqBody := ioutil.NopCloser(strings.NewReader(data.Encode()))

	// Build request
	req := &http.Request{
		Method: http.MethodPost,
		Header: reqHeader,
		URL:    tokenURL,
		Body:   reqBody,
	}

	resp, err := a.http.Do(req)
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
		return nil, oauth2Err("fetch token: response: %v - %s, request: %+v", resp.Status, respBody, req)
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

func (a *api) refreshToken() error {
	if a.token == nil {
		return oauth2Err("token is not set")
	}

	if a.token.RefreshToken() == "" {
		return oauth2Err("empty refresh token")
	}

	token, err := a.accessTokenByRefreshToken(a.token.RefreshToken())
	if err != nil {
		return err
	}

	a.token = token
	return nil
}

func (a *api) domainIsSet() error {
	if a.domain == "" {
		return oauth2Err("account domain is not set")
	}

	return nil
}

func oauth2Err(format string, args ...interface{}) error {
	return fmt.Errorf("oauth2: "+format, args...)
}
