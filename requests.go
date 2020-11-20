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
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	userAgent      = "AmoCRM-API-Golang-Client"
	apiVersion     = uint8(4)
	requestTimeout = 20 * time.Second

	accountEndpoint endpoint = "accounts"
)

type endpoint string

func (e endpoint) path() string {
	return fmt.Sprintf("/api/v%d/%s", apiVersion, e)
}

func (a *api) get(ep endpoint, q url.Values, h http.Header) (*http.Response, error) {
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
