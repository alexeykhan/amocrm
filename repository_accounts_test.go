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

package amocrm_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alexeykhan/amocrm"
)

func TestAccounts_Current(t *testing.T) {
	noTokenClient := amocrm.New(clientID, clientSecret, redirectURL)

	almostValidClient := amocrm.New(clientID, clientSecret, redirectURL)
	_ = almostValidClient.SetToken(amocrm.NewToken(accessToken, refreshToken, tokenType, time.Time{}))
	_ = almostValidClient.SetDomain("example.amocrm.ru")

	relations := []string{
		amocrm.WithUUID,
		amocrm.WithVersion,
		amocrm.WithAmojoID,
		amocrm.WithTaskTypes,
		amocrm.WithUserGroups,
		amocrm.WithAmojoRights,
		amocrm.WithDatetimeSettings,
	}

	cases := []struct {
		client amocrm.Client
		config amocrm.AccountsConfig
		wanted *amocrm.Account
		error  error
	}{
		{
			client: noTokenClient,
			error:  errors.New("unexpected account relation: example"),
			wanted: (*amocrm.Account)(nil),
			config: amocrm.AccountsConfig{
				Relations: []string{"example"},
			},
		},
		{
			client: noTokenClient,
			error:  errors.New("get accounts: invalid token"),
			wanted: (*amocrm.Account)(nil),
			config: amocrm.AccountsConfig{Relations: relations},
		},
		{
			client: almostValidClient,
			error:  nil,
			wanted: &amocrm.Account{},
			config: amocrm.AccountsConfig{Relations: relations},
		},
	}

	for _, tc := range cases {
		got, err := tc.client.Accounts().Current(tc.config)
		require.Exactly(t, tc.wanted, got)

		if tc.error == nil {
			require.NoError(t, err)
		} else {
			require.EqualError(t, err, tc.error.Error())
		}
	}
}
