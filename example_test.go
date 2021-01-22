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
	"fmt"
	"time"

	"github.com/alexeykhan/amocrm"
)

var (
	env = struct {
		clientID     string
		clientSecret string
		redirectURL  string
	}{
		clientID:     "CLIENT_ID",
		clientSecret: "CLIENT_SECRET",
		redirectURL:  "REDIRECT_URI",
	}

	storage = struct {
		domain       string
		accessToken  string
		refreshToken string
		tokenType    string
		expiresAt    time.Time
	}{
		domain:       "example.amocrm.ru",
		accessToken:  "access_token",
		refreshToken: "refresh_token",
		tokenType:    "bearer",
		expiresAt:    time.Now(),
	}
)

func Example_getAuthURL() {
	// Initialize amoCRM API Client.
	amoCRM := amocrm.New(env.clientID, env.clientSecret, env.redirectURL)

	// Save this random state as a session identifier to verify
	// user identity when they are redirected back with code.
	// Set required mode parameter: "post_message" or "popup".
	state := amocrm.RandomState()
	mode := amocrm.PostMessageMode

	// Redirect user to authorization URL.
	authURL, err := amoCRM.AuthorizeURL(state, mode)
	if err != nil {
		fmt.Println("Failed to Get auth url:", err)
		return
	}

	fmt.Println("Redirect user to this URL:")
	fmt.Println(authURL)
}

func Example_getTokenByCode() {
	// Initialize amoCRM API Client.
	amoCRM := amocrm.New(env.clientID, env.clientSecret, env.redirectURL)

	// Use the accounts domain and authorization code that are
	// pushed to the redirect URL as "referer" and "code GET
	// parameters respectively. AccessTokenByCode will do the
	// handshake to retrieve tokens.
	domain := "example.amocrm.ru"
	authCode := "def502000ba3e1724cac79...92146f93b70fd4ca31"

	// Set amoCRM API accounts domain.
	if err := amoCRM.SetDomain(domain); err != nil {
		fmt.Println("set domain:", err)
		return
	}

	// Exchange authorization code for token.
	token, err := amoCRM.TokenByCode(authCode)
	if err != nil {
		fmt.Println("Get token by code:", err)
		return
	}

	// Store received token.
	fmt.Println("access_token:", token.AccessToken())
	fmt.Println("refresh_token:", token.RefreshToken())
	fmt.Println("token_type:", token.TokenType())
	fmt.Println("expires_at:", token.ExpiresAt().Unix())
}

func Example_getCurrentAccount() {
	// Initialize amoCRM API Client.
	amoCRM := amocrm.New(env.clientID, env.clientSecret, env.redirectURL)

	// Retrieve domain from storage.
	if err := amoCRM.SetDomain(storage.domain); err != nil {
		fmt.Println("set domain:", err)
		return
	}

	// Retrieve token from storage.
	token := amocrm.NewToken(storage.accessToken, storage.refreshToken, storage.tokenType, storage.expiresAt)
	if err := amoCRM.SetToken(token); err != nil {
		fmt.Println("set token:", err)
		return
	}

	// Set up accounts request config.
	cfg := amocrm.AccountsConfig{
		Relations: []string{
			amocrm.WithUUID,
			amocrm.WithVersion,
			amocrm.WithAmojoID,
			amocrm.WithTaskTypes,
			amocrm.WithUserGroups,
			amocrm.WithAmojoRights,
			amocrm.WithDatetimeSettings,
		},
	}

	// Fetch current accounts with AccountsRepository.
	account, err := amoCRM.Accounts().Current(cfg)
	if err != nil {
		fmt.Println("fetch current accounts:", err)
		return
	}

	fmt.Println("current accounts:", account)
}
