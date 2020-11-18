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

package oauth2_test

import (
	"fmt"

	"github.com/alexeykhan/amocrm/oauth2"
)

var (
	clientID     = "CLIENT_ID"
	clientSecret = "CLIENT_SECRET"
	redirectURL  = "REDIRECT_URI"
)

func ExampleOAuth2_AuthorizeURL() {
	// Create amoCRM-specific OAuth2 client.
	client := oauth2.New(clientID, clientSecret, redirectURL)

	// Save this random state as a session identifier to verify
	// user identity when they are redirected back with code.
	state := oauth2.GenerateState()

	// Redirect user to authorization URL.
	// Example: https://www.amocrm.ru/oauth?client_id=CLIENT_ID&state=GENERATED_STATE&mode=post_message
	url, _ := client.AuthorizeURL(state, oauth2.PostMessageMode)
	fmt.Println("Redirect user to this URL:", url)
}

func ExampleOAuth2_AccessTokenByCode() {
	// Use the authorization code that is pushed to the redirect
	// URL. AccessTokenByCode will do the handshake to retrieve
	// initial access and refresh tokens.
	code := "AUTHORIZATION_CODE"

	// Use the account domain that is pushed to the redirect URL
	// as a "referer" GET-parameter. AccessTokenByCode will do
	// the handshake to retrieve tokens.
	// https://www.amocrm.ru/developers/content/oauth/step-by-step#context:~:text=referer
	domain := "example.amocrm.ru"

	// Create OAuth2 client with specific amoCRM account domain.
	client := oauth2.New(clientID, clientSecret, redirectURL).SetDomain(domain)

	token, err := client.AccessTokenByCode(code)
	if err != nil {
		fmt.Println("error:", err)

		// Here you'll get an error as we entered invalid credentials:
		// error: oauth2: fetch token: response: 400 Bad Request - {"hint":"Cannot decrypt the authorization code",
		// "title":"Некорректный запрос","type":"https://developers.amocrm.ru/v3/errors/OAuthProblemJson","status":400,
		// "detail":"В запросе отсутствует ряд параметров или параметры невалидны"}
		return
	}

	fmt.Printf("access_token: %s", (*token).AccessToken())
	fmt.Printf("refresh_token: %s", (*token).RefreshToken())
	fmt.Printf("expires_at: %v", (*token).ExpiresAt())
	fmt.Printf("is_valid: %t", (*token).Valid())
}
