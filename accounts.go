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
	"net/url"

	"github.com/alexeykhan/amocrm/entity"
)

// AccountsRepository describes methods available for Accounts entity.
type AccountsRepository interface {
	Current() (*entity.Account, error)
}

type accounts struct {
	api *api
}

// Current returns an Accounts entity for current authorized user.
func (r accounts) Current() (res *entity.Account, err error) {
	res = &entity.Account{}

	query := url.Values{}
	for _, with := range res.Relations() {
		query.Add("with", with)
	}

	resp, rErr := r.api.get(accountEndpoint, query, nil)
	if rErr != nil {
		return res, fmt.Errorf("get accounts: %w", rErr)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			err = fmt.Errorf("close response body: %w", err)
		}
	}()

	if dErr := json.NewDecoder(resp.Body).Decode(res); dErr != nil {
		return res, fmt.Errorf("decode response json: %w", dErr)
	}

	return
}
