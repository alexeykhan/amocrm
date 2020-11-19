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

package repository

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/alexeykhan/amocrm/oauth2"
)

type Endpoint string

func (e Endpoint) Path() string {
	return fmt.Sprintf("/api/v%d/%s", Version, e)
}

const (
	Version        uint8 = 4
	UserAgent            = "github.com/alexeykhan/amocrm"
	RequestTimeout       = 20 * time.Second
	With                 = "with"
)

type Repository struct {
	token oauth2.Token
}

func New(token oauth2.Token) *Repository {
	return &Repository{
		token: token,
	}
}

func (r *Repository) Get(ep Endpoint, q url.Values, h http.Header) (*http.Response, error) {
	if r.token.Expired() {
		if err := r.refreshAccessToken(); err != nil {
			return nil, err
		}
	}

	header := baseHeaders()
	for k, v := range h {
		if _, reserved := header[k]; !reserved {
			header[k] = v
		}
	}

	apiURL, err := a.url(ep, q)
	if err != nil {
		return nil, err
	}

	return client().Do(&http.Request{
		Method: http.MethodGet,
		Header: header,
		URL:    apiURL,
	})
}

func client() *http.Client {
	return &http.Client{
		Timeout: RequestTimeout,
	}
}

func (a *API) url(endpoint Endpoint, query url.Values) (*url.URL, error) {
	path := a.auth.AccountURL() + endpoint.Path() + "?" + query.Encode()
	return url.Parse(path)
}

func baseHeaders() http.Header {
	authHeader := a.token.Type() + " "
	a.token.AccessToken()
	return http.Header{
		"Authorization": []string{authHeader},
		"User-Agent":    []string{UserAgent},
	}
}
