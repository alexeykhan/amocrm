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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// PostMessageMode and PopupMode are the only options for
// amoCRM OAuth2.0 "mode" request parameter.
const (
	PostMessageMode = "post_message"
	PopupMode       = "popup"
)

type GrantType struct {
	code   string
	fields []string
}

var (
	authorizationCodeGrant = GrantType{
		code:   "authorization_code",
		fields: []string{"code"},
	}
	refreshTokenGrant = GrantType{
		code:   "refresh_token",
		fields: []string{"refresh_token"},
	}
	// clientCredentialsGrant = GrantType{
	// 	code:   "client_credentials",
	// 	fields: []string{},
	// }
	// passwordGrant = GrantType{
	// 	code:   "password",
	// 	fields: []string{"username", "password"},
	// }
)

const (
	userAgent      = "AmoCRM-API-Golang-Client"
	apiVersion     = uint8(4)
	requestTimeout = 20 * time.Second
)

// api implements Client interface.
type api struct {
	clientID     string
	clientSecret string
	redirectURL  string

	domain string
	token  Token

	http *http.Client
}

func newAPI(clientID, clientSecret, redirectURL string) *api {
	return &api{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		http: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

func (a *api) get(ep endpoint, q url.Values, h http.Header) (*http.Response, error) {
	if a.token == nil {
		return nil, errors.New("invalid token")
	}

	if a.token.Expired() {
		if err := a.refreshToken(); err != nil {
			return nil, err
		}
	}

	header := a.header()
	for k, v := range h {
		if _, reserved := header[k]; !reserved {
			header[k] = v
		}
	}

	apiURL, err := a.url(ep.path(), q)
	if err != nil {
		return nil, err
	}

	return a.http.Do(&http.Request{
		Method: http.MethodGet,
		Header: header,
		URL:    apiURL,
	})
}

func (a *api) setToken(token Token) error {
	if token == nil {
		return errors.New("invalid token")
	}
	a.token = token
	return nil
}

func (a *api) setDomain(domain string) error {
	if !isValidDomain(domain) {
		return errors.New("invalid domain")
	}

	a.domain = domain
	return nil
}

func (a *api) authorizationURL(state, mode string) (*url.URL, error) {
	if state == "" {
		return nil, oauth2Err("empty state")
	}
	if mode != PostMessageMode && mode != PopupMode {
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

func (a *api) getToken(grant GrantType, options url.Values, header http.Header) (Token, error) {
	if !isValidDomain(a.domain) {
		return nil, oauth2Err("invalid accounts domain")
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

	token := &tokenSource{
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
	if a.token.RefreshToken() == "" {
		return oauth2Err("empty refresh token")
	}

	token, err := a.getToken(refreshTokenGrant, url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{a.token.RefreshToken()},
	}, nil)
	if err != nil {
		return err
	}

	a.token = token
	return nil
}

func (a *api) url(path string, q url.Values) (*url.URL, error) {
	if !isValidDomain(a.domain) {
		return nil, oauth2Err("invalid accounts domain")
	}

	endpointURL := "https://" + a.domain + path + "?" + q.Encode()

	return url.Parse(endpointURL)
}

func (a *api) header() http.Header {
	authHeader := a.token.TokenType() + " " + a.token.AccessToken()

	header := a.baseHeader()
	header["Authorization"] = []string{authHeader}

	return header
}

func (a *api) baseHeader() http.Header {
	return http.Header{
		"User-Agent": []string{userAgent},
	}
}

func isValidDomain(domain string) bool {
	if domain == "" {
		return false
	}

	parts := strings.Split(domain, ".")
	if len(parts) != 3 ||
		parts[0] == "" ||
		parts[0] == "www" ||
		len(parts[0]) > 63 ||
		parts[1] != "amocrm" ||
		parts[2] != "ru" && parts[2] != "com" {
		return false
	}

	return true
}

func oauth2Err(format string, args ...interface{}) error {
	return fmt.Errorf("oauth2: "+format, args...)
}
