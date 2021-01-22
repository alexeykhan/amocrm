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
	"fmt"
	"net/url"
)

// Account relations.
const (
	WithUUID             = "uuid"
	WithVersion          = "version"
	WithAmojoID          = "amojo_id"
	WithTaskTypes        = "task_types"
	WithUserGroups       = "users_groups"
	WithAmojoRights      = "amojo_rights"
	WithDatetimeSettings = "datetime_settings"
)

// Accounts describes methods available for Accounts entity.
type Accounts interface {
	Current(cfg AccountsConfig) (*Account, error)
}

// Verify interface compliance.
var _ Accounts = accounts{}

type accounts struct {
	api *api
}

// Use AccountsConfig to set account parameters.
type AccountsConfig struct {
	Relations []string
}

func newAccounts(api *api) Accounts {
	return accounts{api: api}
}

// Current returns an Accounts entity for current authorized user.
func (a accounts) Current(cfg AccountsConfig) (dto *Account, err error) {
	query := url.Values{}
	for _, relation := range cfg.Relations {
		switch relation {
		case WithUUID, WithVersion, WithAmojoID, WithTaskTypes, WithUserGroups, WithAmojoRights, WithDatetimeSettings:
			query.Add("with", relation)
		default:
			return dto, fmt.Errorf("unexpected account relation: %s", relation)
		}
	}

	resp, rErr := a.api.get(accountsEndpoint, query, nil)
	if rErr != nil {
		return dto, fmt.Errorf("get accounts: %w", rErr)
	}
	defer func() {
		if clErr := resp.Body.Close(); clErr != nil {
			if err != nil {
				err = fmt.Errorf("close response body: %v: %v", clErr, err)
			} else {
				err = fmt.Errorf("close response body: %w", clErr)
			}
		}
	}()

	dto = &Account{}
	if dErr := json.NewDecoder(resp.Body).Decode(dto); dErr != nil {
		return dto, fmt.Errorf("decode json response: %w", dErr)
	}

	return dto, err
}
